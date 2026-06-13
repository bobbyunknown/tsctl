import { Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, ScrollText, Sun, Moon, LogOut, LogIn, User } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useTailscaleWS } from '@/hooks/useWebSocket'
import { useTheme } from '@/contexts/ThemeContext'
import { authApi } from '@/lib/api'
import { toast } from 'sonner'
import { useState } from 'react'
import { useConfirm } from '@/contexts/ConfirmContext'

export default function Navbar() {
    const location = useLocation()
    const { authStatus } = useTailscaleWS()
    const { isDark, toggleTheme } = useTheme()
    const [showLogout, setShowLogout] = useState(false)
    const [loading, setLoading] = useState(false)
    const { confirm } = useConfirm()

    const nav = [
        { path: '/', icon: LayoutDashboard, label: 'OVERVIEW' },
        { path: '/logs', icon: ScrollText, label: 'SYSTEM LOGS' },
    ]

    const handleLogout = async () => {
        const isConfirmed = await confirm({
            title: 'Revoke Identity',
            description: 'Revoke node identity? This action cannot be undone.',
            confirmText: 'Revoke',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        setLoading(true)
        try {
            await authApi.logout()
            toast.success('Node identity revoked')
            setShowLogout(false)
        } catch (error: any) {
            toast.error(error.message)
        }
        setLoading(false)
    }

    return (
        <nav className="h-16 border-b border-[#02d7f2]/15 bg-[#0a0a0f]/90 backdrop-blur-xl fixed top-0 w-full z-50 shadow-[0_0_15px_rgba(2,215,242,0.1)]">
            <div className="max-w-7xl mx-auto h-full px-6 flex items-center justify-between gap-6">

                {/* Brand Section */}
                <Link to="/" className="flex items-center gap-3 group shrink-0">
                    <div className="relative">
                        <div className="absolute inset-0 bg-[#02d7f2] rounded-xl blur-md opacity-40 group-hover:opacity-70 transition-opacity"></div>
                        <div className="relative w-9 h-9 border border-[#02d7f2] bg-black rounded flex items-center justify-center shadow-lg group-hover:scale-105 transition-transform duration-300">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#02d7f2" strokeWidth="1.5" className="w-5 h-5">
                                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" className="origin-center animate-[spin_20s_linear_infinite]"></path>
                                <circle cx="12" cy="12" r="3" fill="#02d7f2" fillOpacity="0.2" stroke="#02d7f2" strokeWidth="1.5" className="animate-[pulse_3s_ease-in-out_infinite]"></circle>
                            </svg>
                        </div>
                    </div>
                    <div className="flex flex-col justify-center">
                        <span className="font-black text-xl tracking-[0.1em] text-nexus-glow logo-font leading-none uppercase">TSCTL</span>
                        <span className="text-[10px] tracking-[0.2em] text-muted-foreground uppercase mt-0.5 opacity-80">Control Plane</span>
                    </div>
                </Link>

                {/* Navigation Menu */}
                <div className="hidden md:flex flex-1 items-center justify-center gap-10">
                    {nav.map((item) => {
                        const isActive = location.pathname === item.path
                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "text-sm font-mono tracking-[0.15em] transition-colors duration-200 uppercase",
                                    isActive ? "text-[#02d7f2]" : "text-zinc-500 hover:text-[#02d7f2]"
                                )}
                            >
                                {item.label}
                            </Link>
                        )
                    })}
                </div>

                {/* Right Section */}
                <div className="flex items-center gap-4 shrink-0">
                    {/* Login Button */}
                    {!authStatus?.authenticated && authStatus?.auth_url && (
                        <a
                            href={authStatus.auth_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="gradient-btn px-4 py-2 rounded text-xs flex items-center gap-2"
                        >
                            <LogIn className="w-4 h-4" />
                            <span className="leading-none">CONNECT NODE</span>
                        </a>
                    )}

                    {/* Auth Status Badge */}
                    <div className={cn(
                        "flex items-center gap-2.5 px-4 py-2 rounded border text-xs font-bold uppercase tracking-wider transition-all duration-300 shadow-sm",
                        authStatus?.authenticated
                            ? "bg-[#39ff14]/10 border-[#39ff14]/30 text-[#39ff14]"
                            : "bg-[#ff00ff]/10 border-[#ff00ff]/30 text-[#ff00ff]"
                    )}>
                        <div className="relative flex h-2.5 w-2.5">
                            {authStatus?.authenticated ? (
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[#39ff14] opacity-75"></span>
                            ) : (
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[#ff00ff] opacity-75"></span>
                            )}
                            <span className={cn(
                                "relative inline-flex rounded-full h-2.5 w-2.5 shadow-sm",
                                authStatus?.authenticated ? "bg-[#39ff14]" : "bg-[#ff00ff]"
                            )}></span>
                        </div>
                        <span className="leading-none">{authStatus?.authenticated ? 'ONLINE' : 'AUTH REQUIRED'}</span>
                    </div>

                    {/* Theme Toggle */}
                    <button
                        onClick={toggleTheme}
                        className="p-2 rounded text-muted-foreground hover:text-[#02d7f2] transition-colors"
                        title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
                    >
                        {isDark ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
                    </button>

                    {/* User Avatar */}
                    {authStatus?.authenticated && (
                        <div className="relative">
                            <button
                                onClick={() => setShowLogout(!showLogout)}
                                className="w-9 h-9 rounded bg-[#0d1117] border border-[#02d7f2]/15 flex items-center justify-center hover:border-[#02d7f2] hover:shadow-[0_0_10px_rgba(2,215,242,0.3)] transition-all duration-200 overflow-hidden"
                            >
                                {authStatus.user_profile_pic ? (
                                    <img
                                        src={`${import.meta.env.VITE_API_BASE_URL}/api/v1/avatar`}
                                        alt={authStatus.user_display_name || 'User'}
                                        className="w-full h-full object-cover grayscale hover:grayscale-0 transition-all"
                                    />
                                ) : (
                                    <User className="w-4 h-4 text-[#02d7f2]" />
                                )}
                            </button>
                            {showLogout && (
                                <>
                                    <div className="fixed inset-0 z-40" onClick={() => setShowLogout(false)}></div>
                                    <div className="absolute right-0 top-full mt-2 w-64 bg-card-defi border border-[#02d7f2]/15 rounded shadow-lg z-50 p-4 backdrop-blur-xl">
                                        <div className="flex items-start gap-3 mb-3 pb-3 border-b border-[#02d7f2]/15">
                                            <div className="w-10 h-10 rounded bg-black border border-[#02d7f2]/15 flex items-center justify-center shrink-0 overflow-hidden">
                                                {authStatus.user_profile_pic ? (
                                                    <img
                                                        src={`${import.meta.env.VITE_API_BASE_URL}/api/v1/avatar`}
                                                        alt={authStatus.user_display_name || 'User'}
                                                        className="w-full h-full object-cover"
                                                    />
                                                ) : (
                                                    <User className="w-5 h-5 text-[#02d7f2]" />
                                                )}
                                            </div>
                                            <div className="min-w-0 flex-1">
                                                <div className="font-bold text-sm text-cyan-glow mb-1 truncate display-font">
                                                    {authStatus.user_display_name || 'UNKNOWN USER'}
                                                </div>
                                                <div className="flex items-center gap-1.5 mb-1">
                                                    {authStatus.is_owner && (
                                                        <span className="bg-[#007aff] text-black text-[10px] px-1.5 py-0.5 rounded font-bold uppercase">OWNER</span>
                                                    )}
                                                    {authStatus.is_admin && (
                                                        <span className="bg-[#02d7f2] text-black text-[10px] px-1.5 py-0.5 rounded font-bold uppercase">ADMIN</span>
                                                    )}
                                                </div>
                                                <div className="text-xs text-muted-foreground truncate">
                                                    {authStatus.user_email || 'NO_EMAIL'}
                                                </div>
                                            </div>
                                        </div>
                                        <div className="space-y-2">
                                            <button
                                                onClick={handleLogout}
                                                disabled={loading}
                                                className="w-full px-3 py-2 bg-black border border-[#ff1111] text-[#ff1111] rounded hover:bg-[#ff1111] hover:text-black hover:shadow-[0_0_10px_rgba(255,17,17,0.5)] transition-all text-xs font-bold disabled:opacity-50 flex items-center justify-center gap-2 uppercase tracking-widest"
                                            >
                                                <LogOut className="w-4 h-4" />
                                                {loading ? 'TERMINATING...' : 'REVOKE ACCESS'}
                                            </button>
                                        </div>
                                    </div>
                                </>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </nav>
    )
}
