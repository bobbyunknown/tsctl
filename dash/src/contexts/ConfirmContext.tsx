import React, { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

interface ConfirmOptions {
  title?: string;
  description?: string;
  confirmText?: string;
  cancelText?: string;
  isDestructive?: boolean;
}

interface ConfirmContextType {
  confirm: (options: ConfirmOptions) => Promise<boolean>;
}

const ConfirmContext = createContext<ConfirmContextType | undefined>(undefined);

export const useConfirm = () => {
  const context = useContext(ConfirmContext);
  if (!context) {
    throw new Error('useConfirm must be used within a ConfirmProvider');
  }
  return context;
};

export const ConfirmProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [options, setOptions] = useState<ConfirmOptions>({});
  const [resolvePromise, setResolvePromise] = useState<(value: boolean) => void>();

  const confirm = useCallback((opts: ConfirmOptions) => {
    setOptions(opts);
    setIsOpen(true);
    return new Promise<boolean>((resolve) => {
      setResolvePromise(() => resolve);
    });
  }, []);

  const handleConfirm = () => {
    setIsOpen(false);
    if (resolvePromise) resolvePromise(true);
  };

  const handleCancel = () => {
    setIsOpen(false);
    if (resolvePromise) resolvePromise(false);
  };

  const confirmClass = options.isDestructive
    ? "bg-transparent border border-[#ff1111] text-[#ff1111] hover:bg-[#ff1111] hover:text-black hover:shadow-[0_0_15px_rgba(255,17,17,0.4)] transition-all uppercase tracking-widest font-bold m-0"
    : "bg-transparent border border-[#02d7f2] text-[#02d7f2] hover:bg-[#02d7f2] hover:text-black hover:shadow-[0_0_15px_rgba(2,215,242,0.4)] transition-all uppercase tracking-widest font-bold m-0";

  return (
    <ConfirmContext.Provider value={{ confirm }}>
      {children}
      <AlertDialog open={isOpen} onOpenChange={(open) => {
          if (!open) {
              setIsOpen(false);
              if (resolvePromise) resolvePromise(false);
          }
      }}>
        <AlertDialogContent className="bg-card-defi border border-[#02d7f2]/15 shadow-[0_0_20px_rgba(2,215,242,0.1)] text-foreground max-w-md">
          <AlertDialogHeader>
            <AlertDialogTitle className={`display-font tracking-widest uppercase text-xl ${options.isDestructive ? 'text-[#ff1111]' : 'text-[#02d7f2]'}`}>
                {options.title || 'Confirm Action'}
            </AlertDialogTitle>
            <AlertDialogDescription className="text-muted-foreground mt-2">
              {options.description}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="mt-6 sm:justify-end gap-2">
            <AlertDialogCancel 
                onClick={handleCancel}
                className="border-[#02d7f2]/15 bg-black hover:bg-[#02d7f2]/10 hover:text-[#02d7f2] transition-colors uppercase tracking-widest font-bold m-0"
            >
              {options.cancelText || 'Cancel'}
            </AlertDialogCancel>
            <AlertDialogAction 
                onClick={handleConfirm}
                className={confirmClass}
            >
              {options.confirmText || 'Confirm'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </ConfirmContext.Provider>
  );
};
