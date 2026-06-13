import { useEffect, useState } from 'react'
import { serveApi, funnelApi, proxyApi } from '../lib/api'
import type { ProxyInfo } from '../lib/api'
import { useTailscaleWS } from '../hooks/useWebSocket'
import { ExposureControls } from '../components/ExposureControls'
import { useConfirm } from '@/contexts/ConfirmContext'
import { toast } from 'sonner'
import {
    ArrowRight,
    ShieldAlert,
    Copy,
    ExternalLink,
    Wifi,
    Cpu,
    Network,
    Server,
    Activity,
    X
} from 'lucide-react'

export default function HomePage() {
    const { authStatus, serveStatus } = useTailscaleWS()
    const { confirm } = useConfirm()
    const [activeServices, setActiveServices] = useState<any>(null)
    const [proxies, setProxies] = useState<ProxyInfo[]>([])

    const fetchInitial = async () => {
        try {
            const data = await serveApi.getStatus()
            setActiveServices(data)
        } catch (error) {
            console.error('Failed to fetch initial services:', error)
        }
    }

    const fetchProxies = async () => {
        try {
            const proxyData = await proxyApi.getStatus()
            setProxies(proxyData || [])
        } catch (error) {
            console.error('Failed to fetch proxy status:', error)
        }
    }

    useEffect(() => {
        fetchInitial()
        fetchProxies()
        const interval = setInterval(fetchProxies, 5000)
        return () => clearInterval(interval)
    }, [])

    useEffect(() => {
        if (serveStatus) {
            setActiveServices(serveStatus)
        }
    }, [serveStatus])

    const handleResetServe = async () => {
        const isConfirmed = await confirm({
            title: 'Stop Serves',
            description: 'Are you sure you want to stop all serve routes?',
            confirmText: 'Stop All',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await serveApi.reset()
            toast.success('All serve routes stopped')
            const data = await serveApi.getStatus()
            setActiveServices(data)
        } catch (error: any) {
            toast.error(error.message || 'Failed to reset serve')
        }
    }

    const handleResetFunnel = async () => {
        const isConfirmed = await confirm({
            title: 'Stop Funnels',
            description: 'Are you sure you want to stop all funnel routes?',
            confirmText: 'Stop All',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await funnelApi.reset()
            toast.success('All funnel routes stopped')
            const data = await serveApi.getStatus()
            setActiveServices(data)
        } catch (error: any) {
            toast.error(error.message || 'Failed to reset funnel')
        }
    }

    const handleStopServe = async (port: number) => {
        const isConfirmed = await confirm({
            title: 'Stop Serve',
            description: `Are you sure you want to stop the serve route on port ${port}?`,
            confirmText: 'Stop',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await serveApi.stop(port)
            toast.success(`Serve on port ${port} stopped`)
            const data = await serveApi.getStatus()
            setActiveServices(data)
        } catch (error: any) {
            toast.error(error.message || 'Failed to stop serve')
        }
    }

    const handleStopFunnel = async (port: number) => {
        const isConfirmed = await confirm({
            title: 'Stop Funnel',
            description: `Are you sure you want to stop the funnel route on port ${port}?`,
            confirmText: 'Stop',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await funnelApi.stop(port)
            toast.success(`Funnel on port ${port} stopped`)
            const data = await serveApi.getStatus()
            setActiveServices(data)
        } catch (error: any) {
            toast.error(error.message || 'Failed to stop funnel')
        }
    }

    const handleProxyStop = async (port: number) => {
        const isConfirmed = await confirm({
            title: 'Stop Proxy',
            description: `Are you sure you want to stop the proxy on port ${port}?`,
            confirmText: 'Stop',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await proxyApi.stop([port])
            toast.success(`Proxy stopped for port ${port}`)
            fetchProxies()
        } catch (error: any) {
            toast.error(error.message || 'Failed to stop proxy')
        }
    }

    const handleProxyStopAll = async () => {
        const isConfirmed = await confirm({
            title: 'Stop Proxies',
            description: 'Are you sure you want to stop all active proxies?',
            confirmText: 'Stop All',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return
        try {
            await proxyApi.stopAll()
            toast.success('All proxies stopped')
            fetchProxies()
        } catch (error: any) {
            toast.error(error.message || 'Failed to stop proxies')
        }
    }

    if (!authStatus) {
        return (
            <div className="flex flex-col items-center justify-center min-h-[65vh] gap-6">
                <div className="relative w-14 h-14">
                    <div className="absolute inset-0 bg-[#02d7f2] rounded-full blur-md opacity-30 animate-pulse"></div>
                    <svg className="w-full h-full animate-[spin_2.5s_linear_infinite]" viewBox="0 0 100 100">
                        <circle 
                            cx="50" cy="50" r="40" 
                            stroke="rgba(2, 215, 242, 0.15)" 
                            strokeWidth="4" 
                            fill="none" 
                        />
                        <circle 
                            cx="50" cy="50" r="40" 
                            stroke="#02d7f2" 
                            strokeWidth="4" 
                            strokeDasharray="60 180" 
                            strokeLinecap="round"
                            fill="none" 
                            style={{ filter: 'drop-shadow(0 0 5px #02d7f2)' }}
                        />
                    </svg>
                </div>
                <div className="text-center font-mono tracking-[0.2em] text-[10px] text-cyan-glow uppercase animate-pulse">
                    Initializing control plane...
                </div>
            </div>
        )
    }

    const isConnected = authStatus.authenticated
    const serves = activeServices?.services?.filter((s: any) => s.type === 'serve') || []
    const funnels = activeServices?.services?.filter((s: any) => s.type === 'funnel') || []

    return (
        <div className="pb-10 text-foreground font-sans">
            <div className="max-w-6xl mx-auto space-y-6 px-4">
                
                {/* PAGE HEADER */}
                <header className="flex flex-col gap-4 border-b border-[#02d7f2]/15 pb-8 mb-8 mt-4">
                    <div className="flex items-center gap-4">
                        <Activity className="w-10 h-10 text-[#02d7f2]" style={{ filter: 'drop-shadow(0 0 10px #02d7f2)' }} />
                        <h1 className="text-5xl display-font font-bold tracking-[0.15em] text-white uppercase">DASHBOARD</h1>
                    </div>
                    <p className="font-mono text-zinc-400 text-sm max-w-3xl leading-relaxed tracking-wide opacity-80">
                        Control your Tailscale node configuration, expose local ports via proxy, and manage public internet access points (Funnel) directly from this terminal.
                    </p>
                </header>

                {/* TELEMETRY HEADER */}
                <div className="bg-card-defi border border-[#02d7f2]/15 rounded flex flex-col md:flex-row divide-y md:divide-y-0 md:divide-x divide-[rgba(0,255,255,0.15)] shadow-lg overflow-hidden">
                    <div className="flex-1 p-5 flex flex-col justify-between hover:bg-[rgba(2,215,242,0.02)] transition-colors">
                        <div className="flex justify-between items-center mb-2">
                            <span className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Connection</span>
                            <Wifi className={`w-4 h-4 ${isConnected ? "text-[#39ff14]" : "text-[#ff00ff]"}`} />
                        </div>
                        <div className={`text-lg font-bold display-font tracking-widest ${isConnected ? "text-[#39ff14]" : "text-muted-foreground"}`}>
                            {isConnected ? 'CONNECTED' : 'DISCONNECTED'}
                        </div>
                    </div>

                    <div className="flex-1 p-5 flex flex-col justify-between hover:bg-[rgba(2,215,242,0.02)] transition-colors">
                        <div className="flex justify-between items-center mb-2">
                            <span className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Daemon Status</span>
                            <Cpu className="w-4 h-4 text-[#007aff]" />
                        </div>
                        <div className="text-lg font-bold display-font tracking-widest text-[#007aff] uppercase">
                            {authStatus?.backend_state ? authStatus.backend_state : 'IDLE'}
                        </div>
                    </div>

                    <div className="flex-1 p-5 flex flex-col justify-between hover:bg-[rgba(2,215,242,0.02)] transition-colors">
                        <div className="flex justify-between items-center mb-2">
                            <span className="text-xs font-bold tracking-widest text-muted-foreground uppercase">IPv4 Address</span>
                            <Network className="w-4 h-4 text-[#02d7f2]" />
                        </div>
                        <div className="text-lg font-bold display-font tracking-widest text-[#02d7f2]">
                            {authStatus?.ips?.[0] || '---.---.---.---'}
                        </div>
                    </div>
                </div>

                {/* AUTHENTICATION REQUIRED BLOCK */}
                {authStatus && !authStatus.authenticated && authStatus.auth_url && (
                    <div className="bg-card-defi border border-[#ff00ff]/50 rounded p-6 flex flex-col sm:flex-row items-center justify-between gap-6 shadow-[0_0_20px_rgba(0,122,255,0.15)] relative overflow-hidden">
                        <div className="absolute inset-0 bg-gradient-to-r from-transparent via-[#ff00ff]/5 to-transparent animate-pulse"></div>
                        <div className="flex items-start gap-4 relative z-10">
                            <ShieldAlert className="w-8 h-8 text-[#ff00ff] shrink-0 mt-1" />
                            <div className="space-y-1">
                                <h3 className="text-lg font-bold display-font tracking-widest text-[#ff00ff]">AUTH REQUIRED</h3>
                                <p className="text-sm text-muted-foreground">
                                    Node is currently isolated. Connect your identity to join the tailnet.
                                </p>
                            </div>
                        </div>
                        <button
                            onClick={() => window.open(authStatus.auth_url, '_blank')}
                            className="gradient-btn px-6 py-3 rounded flex items-center justify-center gap-2 relative z-10 w-full sm:w-auto"
                        >
                            CONNECT NOW <ArrowRight className="w-5 h-5" />
                        </button>
                    </div>
                )}

                {/* MAIN GRID */}
                <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
                    
                    {/* LEFT COLUMN: Controls & Services */}
                    <div className="lg:col-span-8 space-y-6">
                        
                        <ExposureControls onSuccess={() => { fetchInitial(); fetchProxies(); }} />
                    </div>

                    {/* RIGHT COLUMN: Network Info */}
                    <div className="lg:col-span-4 space-y-6">
                        {isConnected && (
                            <div className="bg-card-defi border border-[#02d7f2]/15 rounded overflow-hidden shadow-lg h-full flex flex-col">
                                <div className="border-b border-[#02d7f2]/15 bg-black/40 px-5 py-4">
                                    <span className="text-sm font-bold tracking-widest text-cyan-glow display-font">NETWORK INFO</span>
                                </div>
                                <div className="p-5 space-y-6 flex-1 bg-black/20">
                                    
                                    <div className="space-y-1">
                                        <div className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Tailnet Domain</div>
                                        <div className="text-sm text-[#02d7f2] break-all">{authStatus.tailnet_name || 'UNKNOWN_NET'}</div>
                                    </div>

                                    <div className="space-y-1">
                                        <div className="text-xs font-bold tracking-widest text-muted-foreground uppercase">DNS Suffix</div>
                                        <div className="text-sm text-[#02d7f2] break-all">{authStatus.tailnet_dns_suffix || 'N/A'}</div>
                                    </div>

                                    <div className="space-y-1">
                                        <div className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Host ID</div>
                                        <div className="text-sm text-[#02d7f2]">{authStatus.hostname || 'N/A'}</div>
                                    </div>

                                    <div className="mt-6 pt-6 border-t border-[#02d7f2]/15 space-y-5">
                                        <div className="space-y-1">
                                            <div className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Peers Linked</div>
                                            <div className="text-lg font-bold display-font text-[#007aff]">{authStatus.peer_count || 0}</div>
                                        </div>
                                        {authStatus.key_expiry && (
                                            <div className="space-y-1">
                                                <div className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Key Expiry</div>
                                                <div className="text-sm font-bold text-[#39ff14] display-font">
                                                    {new Date(authStatus.key_expiry).toLocaleDateString()}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                    
                                </div>
                            </div>
                        )}
                    </div>

                </div>

                {/* Unified Active Services Card (FULL WIDTH) */}
                {isConnected && (serves.length > 0 || funnels.length > 0 || proxies.length > 0) && (
                    <div className="bg-card-defi border border-[#02d7f2]/15 rounded overflow-hidden shadow-lg w-full">
                        <div className="border-b border-[#02d7f2]/15 bg-black/40 px-5 py-4 flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <Server className="w-5 h-5 text-[#02d7f2]" />
                                <span className="text-sm font-bold tracking-widest text-cyan-glow display-font">ACTIVE SERVICES</span>
                            </div>
                            <div className="flex items-center gap-3">
                                {serves.length > 0 && (
                                    <button 
                                        onClick={handleResetServe}
                                        className="h-8 px-3 text-[10px] font-bold border border-[#ff1111] text-[#ff1111] hover:bg-[#ff1111] hover:text-black rounded transition-colors uppercase tracking-wider"
                                    >
                                        STOP ALL SERVES
                                    </button>
                                )}
                                {funnels.length > 0 && (
                                    <button 
                                        onClick={handleResetFunnel}
                                        className="h-8 px-3 text-[10px] font-bold border border-[#ff1111] text-[#ff1111] hover:bg-[#ff1111] hover:text-black rounded transition-colors uppercase tracking-wider"
                                    >
                                        STOP ALL FUNNELS
                                    </button>
                                )}
                                {proxies.length > 0 && (
                                    <button 
                                        onClick={handleProxyStopAll}
                                        className="h-8 px-3 text-[10px] font-bold border border-[#ff1111] text-[#ff1111] hover:bg-[#ff1111] hover:text-black rounded transition-colors uppercase tracking-wider"
                                    >
                                        STOP ALL PROXIES
                                    </button>
                                )}
                            </div>
                        </div>
                        <div className="p-5 space-y-4">
                            {/* Render Serves */}
                            {serves.map((service: any) => {
                                const tailscaleIp = authStatus?.ips?.[0]
                                const serveUrl = tailscaleIp ? `http://${tailscaleIp}:${service.port}` : service.local_url
                                return (
                                <div key={service.port} className="flex flex-col sm:flex-row sm:items-center justify-between p-4 border border-[#02d7f2]/15 rounded hover:border-[#02d7f2]/50 transition-all duration-200 bg-black/60 gap-4 shadow-[inset_0_0_10px_rgba(2,215,242,0.02)]">
                                    <div className="flex items-center gap-4">
                                        <div className="px-2.5 py-1 text-[10px] font-bold rounded border border-[#02d7f2]/30 text-[#02d7f2] bg-[#02d7f2]/5 tracking-widest font-mono uppercase">
                                            SERVE
                                        </div>
                                        <div className="flex flex-col sm:flex-row sm:items-baseline gap-2">
                                            <span className="text-lg font-bold display-font text-white tracking-wider">PORT {service.port}</span>
                                            {serveUrl && (
                                                <a 
                                                    href={serveUrl}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="text-xs text-[#02d7f2] font-mono select-all opacity-85 hover:opacity-100 hover:underline transition-opacity"
                                                >
                                                    ({tailscaleIp ? tailscaleIp : 'localhost'}:{service.port})
                                                </a>
                                            )}
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-2 sm:translate-y-1">
                                        {serveUrl && (
                                            <>
                                                <button
                                                    className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#02d7f2] hover:bg-[#02d7f2]/10 rounded border border-transparent hover:border-[#02d7f2]/20 transition-all duration-200"
                                                    onClick={() => {
                                                        navigator.clipboard.writeText(serveUrl)
                                                        toast.success('Copied serve URL')
                                                    }}
                                                    title="Copy URL"
                                                >
                                                    <Copy className="w-3.5 h-3.5" />
                                                </button>
                                                <button
                                                    className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#007aff] hover:bg-[#007aff]/10 rounded border border-transparent hover:border-[#007aff]/20 transition-all duration-200"
                                                    onClick={() => window.open(serveUrl, '_blank')}
                                                    title="Open Link"
                                                >
                                                    <ExternalLink className="w-3.5 h-3.5" />
                                                </button>
                                            </>
                                        )}
                                        <button 
                                            onClick={() => handleStopServe(service.port)}
                                            className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#ff1111] hover:bg-[#ff1111]/10 rounded border border-transparent hover:border-[#ff1111]/20 transition-all duration-200"
                                            title="Stop Serve"
                                        >                                          <X className="w-4 h-4" />
                                        </button>
                                    </div>
                                </div>
                            )})}

                            {/* Render Funnels */}
                            {funnels.map((service: any) => {
                                const cleanUrl = service.public_url?.replace(/\.+$/, '')
                                return (
                                    <div key={service.port} className="flex flex-col sm:flex-row sm:items-center justify-between p-4 border border-[#02d7f2]/15 rounded hover:border-[#007aff]/50 transition-all duration-200 bg-black/60 gap-4 shadow-[inset_0_0_10px_rgba(0,122,255,0.02)]">
                                        <div className="flex items-center gap-4">
                                            <div className="px-2.5 py-1 text-[10px] font-bold rounded border border-[#007aff]/30 text-[#007aff] bg-[#007aff]/5 tracking-widest font-mono uppercase">
                                                FUNNEL
                                            </div>
                                            <div className="flex flex-col sm:flex-row sm:items-baseline gap-2">
                                                <span className="text-lg font-bold display-font text-white tracking-wider">PORT {service.port}</span>
                                                {cleanUrl && (
                                                    <a 
                                                        href={cleanUrl}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="text-xs text-[#02d7f2] font-mono select-all opacity-85 hover:opacity-100 hover:underline transition-opacity"
                                                    >
                                                        ({cleanUrl})
                                                    </a>
                                                )}
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-2 sm:translate-y-1">
                                            {cleanUrl && (
                                                <>
                                                    <button
                                                        className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#02d7f2] hover:bg-[#02d7f2]/10 rounded border border-transparent hover:border-[#02d7f2]/20 transition-all duration-200"
                                                        onClick={() => {
                                                            navigator.clipboard.writeText(cleanUrl)
                                                            toast.success('Copied funnel URL')
                                                        }}
                                                        title="Copy URL"
                                                    >
                                                        <Copy className="w-3.5 h-3.5" />
                                                    </button>
                                                    <button
                                                        className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#007aff] hover:bg-[#007aff]/10 rounded border border-transparent hover:border-[#007aff]/20 transition-all duration-200"
                                                        onClick={() => window.open(cleanUrl, '_blank')}
                                                        title="Open Link"
                                                    >
                                                        <ExternalLink className="w-3.5 h-3.5" />
                                                    </button>
                                                </>
                                            )}
                                            <button 
                                                onClick={() => handleStopFunnel(service.port)}
                                                className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#ff1111] hover:bg-[#ff1111]/10 rounded border border-transparent hover:border-[#ff1111]/20 transition-all duration-200"
                                                title="Stop Funnel"
                                            >
                                                <X className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </div>
                                )
                            })}

                            {/* Render Proxies */}
                            {proxies.map((p: any) => {
                                const tailscaleIp = authStatus?.ips?.[0]
                                const proxyUrl = tailscaleIp ? `http://${tailscaleIp}:${p.port}` : ''
                                return (
                                    <div key={p.port} className="flex flex-col sm:flex-row sm:items-center justify-between p-4 border border-[#02d7f2]/15 rounded hover:border-[#f2e900]/50 transition-all duration-200 bg-black/60 gap-4 shadow-[inset_0_0_10px_rgba(242,233,0,0.02)]">
                                        <div className="flex items-center gap-4">
                                            <div className="px-2.5 py-1 text-[10px] font-bold rounded border border-[#f2e900]/30 text-[#f2e900] bg-[#f2e900]/5 tracking-widest font-mono uppercase">
                                                PROXY
                                            </div>
                                            <div className="flex flex-col sm:flex-row sm:items-baseline gap-2">
                                                <span className="text-lg font-bold display-font text-white tracking-wider">PORT {p.port}</span>
                                                {tailscaleIp && (
                                                    <a 
                                                        href={proxyUrl}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="text-xs text-[#02d7f2] font-mono select-all opacity-85 hover:opacity-100 hover:underline transition-opacity"
                                                    >
                                                        ({tailscaleIp}:{p.port})
                                                    </a>
                                                )}
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-2 sm:translate-y-1">
                                            {tailscaleIp && (
                                                <>
                                                    <button
                                                        className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#02d7f2] hover:bg-[#02d7f2]/10 rounded border border-transparent hover:border-[#02d7f2]/20 transition-all duration-200"
                                                        onClick={() => {
                                                            navigator.clipboard.writeText(proxyUrl)
                                                            toast.success('Copied proxy URL')
                                                        }}
                                                        title="Copy Proxy URL"
                                                    >
                                                        <Copy className="w-3.5 h-3.5" />
                                                    </button>
                                                    <button
                                                        className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#007aff] hover:bg-[#007aff]/10 rounded border border-transparent hover:border-[#007aff]/20 transition-all duration-200"
                                                        onClick={() => window.open(proxyUrl, '_blank')}
                                                        title="Open Link"
                                                    >
                                                        <ExternalLink className="w-3.5 h-3.5" />
                                                    </button>
                                                </>
                                            )}
                                            <button 
                                                onClick={() => handleProxyStop(p.port)}
                                                className="h-8 w-8 flex items-center justify-center text-muted-foreground hover:text-[#ff1111] hover:bg-[#ff1111]/10 rounded border border-transparent hover:border-[#ff1111]/20 transition-all duration-200"
                                                title="Stop Proxy"
                                            >
                                                <X className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </div>
                                )
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
