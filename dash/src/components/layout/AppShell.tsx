import { type ReactNode } from 'react'
import Navbar from './Navbar'

interface AppShellProps {
    children: ReactNode
}

export default function AppShell({ children }: AppShellProps) {
    return (
        <div className="min-h-screen bg-background font-sans antialiased text-foreground selection:bg-primary/20 selection:text-primary">
            <Navbar />
            <main className="pt-24 px-6 pb-20 w-full">
                {children}
            </main>
        </div>
    )
}
