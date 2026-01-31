import { Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, ScrollText, Zap, Sun, Moon, LogOut, LogIn, User } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useTailscaleWS } from '@/hooks/useWebSocket'
import { useTheme } from '@/contexts/ThemeContext'
import { authApi } from '@/lib/api'
import { toast } from 'sonner'
import { useState } from 'react'

export default function Navbar() {
    const location = useLocation()
    const { authStatus } = useTailscaleWS()
    const { isDark, toggleTheme } = useTheme()
    const [showLogout, setShowLogout] = useState(false)
    const [loading, setLoading] = useState(false)

    const nav = [
        { path: '/', icon: LayoutDashboard, label: 'Overview' },
        { path: '/logs', icon: ScrollText, label: 'System Logs' },
    ]

    const handleLogout = async () => {
        if (!confirm('Revoke node identity? This action cannot be undone.')) return
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
        <nav className="h-16 border-b border-[hsl(var(--navbar-border))] bg-[hsl(var(--navbar-bg))]/95 backdrop-blur-xl fixed top-0 w-full z-50 shadow-sm">
            <div className="max-w-7xl mx-auto h-full px-6 flex items-center justify-between gap-6">

                {/* Brand Section */}
                <Link to="/" className="flex items-center gap-3 group shrink-0">
                    <div className="relative">
                        <div className="absolute inset-0 bg-gradient-to-br from-indigo-500 to-violet-600 rounded-xl blur-sm opacity-40 group-hover:opacity-60 transition-opacity"></div>
                        <div className="relative w-9 h-9 bg-gradient-to-br from-indigo-500 via-violet-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg group-hover:scale-105 transition-transform duration-300">
                            <Zap className="w-5 h-5 text-white fill-white" />
                        </div>
                    </div>
                    <div className="flex flex-col">
                        <span className="font-bold text-base tracking-tight text-foreground leading-none">TS Control</span>
                        <span className="text-[10px] text-muted-foreground font-semibold uppercase tracking-widest mt-0.5">Plane</span>
                    </div>
                </Link>

                {/* Navigation Menu */}
                <div className="hidden md:flex items-center gap-1">
                    {nav.map((item) => {
                        const Icon = item.icon
                        const isActive = location.pathname === item.path

                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 relative",
                                    isActive
                                        ? "text-foreground"
                                        : "text-muted-foreground hover:text-foreground hover:bg-secondary/50"
                                )}
                            >
                                <Icon className="w-4 h-4" />
                                <span>{item.label}</span>
                                {isActive && (
                                    <span className="absolute bottom-0 left-2 right-2 h-0.5 bg-gradient-to-r from-indigo-600 to-violet-600 rounded-full" />
                                )}
                            </Link>
                        )
                    })}
                </div>

                {/* Right Section: Login + Auth Status + Theme Toggle + Avatar */}
                <div className="flex items-center gap-3 shrink-0">
                    {/* Login Button - shows when not authenticated */}
                    {!authStatus?.authenticated && authStatus?.auth_url && (
                        <a
                            href={authStatus.auth_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-2 px-4 py-2 rounded-xl border border-indigo-500/30 bg-indigo-500/10 text-indigo-500 hover:bg-indigo-500/15 text-xs font-bold uppercase tracking-wider transition-all duration-300 shadow-sm"
                        >
                            <LogIn className="w-4 h-4" />
                            <span className="leading-none">Login</span>
                        </a>
                    )}

                    {/* Auth Status Badge */}
                    <div className={cn(
                        "flex items-center gap-2.5 px-4 py-2 rounded-xl border text-xs font-bold uppercase tracking-wider transition-all duration-300 shadow-sm",
                        authStatus?.authenticated
                            ? "bg-emerald-500/10 border-emerald-500/30 text-emerald-500 hover:bg-emerald-500/15"
                            : "bg-amber-500/10 border-amber-500/30 text-amber-500 hover:bg-amber-500/15"
                    )}>
                        <div className="relative flex h-2.5 w-2.5">
                            {!authStatus?.authenticated && (
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-500 opacity-75"></span>
                            )}
                            <span className={cn(
                                "relative inline-flex rounded-full h-2.5 w-2.5 shadow-sm",
                                authStatus?.authenticated ? "bg-emerald-500" : "bg-amber-500"
                            )}></span>
                        </div>
                        <span className="leading-none">{authStatus?.authenticated ? 'Online' : 'Auth Required'}</span>
                    </div>

                    {/* Theme Toggle */}
                    <button
                        onClick={toggleTheme}
                        className="p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-secondary/50 transition-all duration-300"
                        title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
                    >
                        {isDark ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
                    </button>

                    {/* User Avatar Dropdown */}
                    {authStatus?.authenticated && (
                        <div className="relative">
                            <button
                                onClick={() => setShowLogout(!showLogout)}
                                className="w-9 h-9 rounded-full bg-gradient-to-br from-violet-500/20 to-blue-500/20 flex items-center justify-center hover:ring-2 hover:ring-violet-500/50 transition-all duration-200 overflow-hidden"
                                title="User Profile"
                            >
                                {authStatus.user_profile_pic ? (
                                    <img
                                        src={`${import.meta.env.VITE_API_BASE_URL}/api/v1/avatar`}
                                        alt={authStatus.user_display_name || 'User'}
                                        className="w-full h-full object-cover"
                                    />
                                ) : (
                                    <User className="w-4 h-4 text-violet-500" />
                                )}
                            </button>
                            {showLogout && (
                                <>
                                    <div className="fixed inset-0 z-40" onClick={() => setShowLogout(false)}></div>
                                    <div className="absolute right-0 top-full mt-2 w-64 bg-card border border-border rounded-lg shadow-lg z-50 p-4">
                                        <div className="flex items-start gap-3 mb-3 pb-3 border-b border-border">
                                            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-violet-500/20 to-blue-500/20 flex items-center justify-center shrink-0 overflow-hidden">
                                                {authStatus.user_profile_pic ? (
                                                    <img
                                                        src={`${import.meta.env.VITE_API_BASE_URL}/api/v1/avatar`}
                                                        alt={authStatus.user_display_name || 'User'}
                                                        className="w-full h-full object-cover"
                                                    />
                                                ) : (
                                                    <User className="w-5 h-5 text-violet-500" />
                                                )}
                                            </div>
                                            <div className="min-w-0 flex-1">
                                                <div className="font-semibold text-sm text-foreground mb-1">
                                                    {authStatus.user_display_name || 'Unknown User'}
                                                </div>
                                                <div className="flex items-center gap-1.5 mb-1">
                                                    {authStatus.is_owner && (
                                                        <span className="bg-violet-500 text-white text-[8px] px-1.5 py-0.5 rounded font-bold">OWNER</span>
                                                    )}
                                                    {authStatus.is_admin && (
                                                        <span className="bg-blue-500 text-white text-[8px] px-1.5 py-0.5 rounded font-bold">ADMIN</span>
                                                    )}
                                                </div>
                                                <div className="text-xs text-muted-foreground truncate">
                                                    {authStatus.user_email || 'No email'}
                                                </div>
                                            </div>
                                        </div>
                                        <div className="space-y-2">
                                            <div className="text-sm text-foreground font-semibold">Revoke Node Identity</div>
                                            <div className="text-xs text-muted-foreground">This action cannot be undone</div>
                                            <button
                                                onClick={handleLogout}
                                                disabled={loading}
                                                className="w-full px-3 py-2 bg-destructive text-destructive-foreground rounded-md hover:bg-destructive/90 transition-colors text-sm font-medium disabled:opacity-50 flex items-center justify-center gap-2"
                                            >
                                                <LogOut className="w-4 h-4" />
                                                {loading ? 'Logging out...' : 'Logout'}
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
