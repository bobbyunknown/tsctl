import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from 'sonner'
import { ThemeProvider } from './contexts/ThemeContext.tsx'
import { ConfirmProvider } from './contexts/ConfirmContext.tsx'
import App from './App.tsx'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider>
      <ConfirmProvider>
        <BrowserRouter>
          <App />
          <Toaster 
            position="top-right"
            theme="dark"
            toastOptions={{
              style: {
                background: 'rgba(0, 0, 0, 0.95)',
                backdropFilter: 'blur(16px)',
                border: '1px solid #f2e900',
                color: '#f2e900',
                fontFamily: '"JetBrains Mono", monospace',
                textTransform: 'uppercase',
                letterSpacing: '0.1em',
                borderRadius: '4px',
              },
              classNames: {
                toast: 'group flex items-center gap-3 p-4',
                title: 'font-bold text-sm',
                description: 'text-xs opacity-80',
                success: 'border-[#f2e900] text-[#f2e900] shadow-[0_0_15px_rgba(242,233,0,0.25)]',
                error: 'border-[#ff1111] text-[#ff1111] shadow-[0_0_15px_rgba(255,17,17,0.25)]',
                warning: 'border-[#eab308] text-[#eab308] shadow-[0_0_15px_rgba(234,179,8,0.2)]',
                info: 'border-[#02d7f2] text-[#02d7f2] shadow-[0_0_15px_rgba(2,215,242,0.2)]',
              }
            }}
          />
        </BrowserRouter>
      </ConfirmProvider>
    </ThemeProvider>
  </StrictMode>,
)
