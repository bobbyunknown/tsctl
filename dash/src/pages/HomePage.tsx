import { useState, useEffect } from 'react'
import { serveApi, funnelApi, sshApi } from '../lib/api'
import { useTailscaleWS } from '../hooks/useWebSocket'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { toast } from 'sonner'
import {
    ArrowRight,
    ShieldAlert,
    Terminal,
    Zap,
    Globe,
    Server,
    Network,
    Cpu,
    Lock,
    Wifi,
    Calendar,
    Users,
    ExternalLink,
    Copy,
} from 'lucide-react'
import { cn } from '@/lib/utils'


const InfoCard = ({
    label,
    value,
    subtext,
    icon: Icon,
    active = false,
    variant = 'default',
    accentColor = "text-violet-500"
}: {
    label: string
    value: string
    subtext: string
    icon: any
    active?: boolean
    variant?: 'default' | 'mono'
    accentColor?: string
}) => (
    <Card className="relative p-6 bg-card/50 backdrop-blur-sm border-border/60 hover:border-primary/30 card-hover-effect overflow-hidden group">
        {/* Inner Top Highlight */}
        <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-white/10 to-transparent opacity-50" />

        <div className="flex justify-between items-start mb-4">
            <span className="text-[10px] font-bold tracking-widest text-muted-foreground uppercase">{label}</span>
            <Icon className={cn("w-4 h-4 opacity-50 group-hover:opacity-100 transition-opacity", accentColor)} />
        </div>

        <div className="space-y-1">
            <div className={cn(
                "text-2xl font-semibold tracking-tight text-foreground break-all",
                variant === 'mono' && "font-mono text-xl"
            )}>
                {value}
            </div>
            <div className="text-xs font-medium text-muted-foreground flex items-center gap-2">
                {subtext}
                {active && (
                    <span className="flex h-1.5 w-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
                )}
            </div>
        </div>
    </Card>
)

export default function HomePage() {
    const { authStatus, serveStatus } = useTailscaleWS()
    const [servePort, setServePort] = useState('8080')
    const [funnelPort, setFunnelPort] = useState('443')
    const [loading, setLoading] = useState(false)
    const [activeServices, setActiveServices] = useState<any>(null)

    useEffect(() => {
        const fetchInitial = async () => {
            try {
                const data = await serveApi.getStatus()
                setActiveServices(data)
            } catch (error) {
                console.error('Failed to fetch initial services:', error)
            }
        }
        fetchInitial()
    }, [])

    useEffect(() => {
        if (serveStatus) {
            setActiveServices(serveStatus)
        }
    }, [serveStatus])


    const handleServeStart = async () => {
        setLoading(true)
        try {
            await serveApi.start(Number(servePort), false)
            toast.success('Serve started')
        } catch (error: any) {
            toast.error(error.message)
        }
        setLoading(false)
    }

    const handleFunnelStart = async () => {
        setLoading(true)
        try {
            await funnelApi.start(Number(funnelPort), false)
            toast.success('Funnel started')
        } catch (error: any) {
            toast.error(error.message)
        }
        setLoading(false)
    }

    const handleSSHEnable = async () => {
        setLoading(true)
        try {
            await sshApi.enable()
            toast.success('SSH access enabled')
        } catch (error: any) {
            toast.error(error.message)
        }
        setLoading(false)
    }

    return (
        <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in py-8">

            <header className="flex flex-col gap-2">
                <div className="flex items-center gap-3">
                    <h1 className="text-3xl font-semibold tracking-tight text-foreground">Network Overview</h1>
                    <Badge variant="outline" className="text-[10px] uppercase font-bold tracking-wider py-0.5 h-6 bg-background/50 backdrop-blur border-border/50 text-muted-foreground">
                        {authStatus?.authenticated ? 'Live' : 'Offline'}
                    </Badge>
                </div>
                <p className="text-sm text-muted-foreground max-w-2xl leading-relaxed">
                    Manage local node connectivity configuration, traffic routing, and system access controls.
                </p>
            </header>

            <section className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <InfoCard
                    label="Context"
                    value={authStatus?.authenticated ? 'Connected' : 'Disconnected'}
                    subtext="NODE_STATE"
                    icon={Wifi}
                    active={authStatus?.authenticated}
                    accentColor={authStatus?.authenticated ? 'text-emerald-500' : 'text-amber-500'}
                />
                <InfoCard
                    label="Backend"
                    value={authStatus?.backend_state || 'Idle'}
                    subtext="DAEMON_PROCESS"
                    icon={Cpu}
                    active={authStatus?.backend_state === 'Running'}
                    accentColor="text-violet-500"
                />
                <InfoCard
                    label="Address"
                    value={authStatus?.ips?.[0] || '---.---.---.---'}
                    subtext="TS_IPV4"
                    icon={Network}
                    variant="mono"
                    accentColor="text-blue-500"
                />
            </section>

            {authStatus && !authStatus.authenticated && authStatus.auth_url && (
                <div className="relative overflow-hidden rounded-xl border border-amber-500/20 bg-amber-500/5 p-8 flex flex-col sm:flex-row items-center justify-between gap-6 card-hover-effect">
                    <div className="absolute inset-0 bg-gradient-to-r from-amber-500/10 to-transparent opacity-50" />

                    <div className="relative z-10 flex items-start gap-4">
                        <div className="p-3 rounded-xl bg-amber-500/10 text-amber-500 shrink-0 border border-amber-500/20 shadow-lg shadow-amber-500/10">
                            <ShieldAlert className="w-6 h-6" />
                        </div>
                        <div className="space-y-1">
                            <h3 className="text-lg font-semibold text-amber-500">Authentication Required</h3>
                            <p className="text-sm text-muted-foreground max-w-md">
                                This node must be authenticated to join the tailnet and enable routing features.
                            </p>
                        </div>
                    </div>
                    <Button
                        onClick={() => window.open(authStatus.auth_url, '_blank')}
                        className="relative z-10 bg-amber-500 hover:bg-amber-600 text-white border-0 shadow-lg shadow-amber-500/20 h-10 px-6 font-medium"
                    >
                        Authenticate <ArrowRight className="ml-2 w-4 h-4" />
                    </Button>
                </div>
            )}

            {authStatus?.authenticated && (
                <section className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <Card className="p-6 bg-card/50 backdrop-blur-sm border-border/60 space-y-4">
                        <h3 className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Tailnet Information</h3>
                        <div className="space-y-3">
                            <div className="flex items-start justify-between">
                                <span className="text-xs text-muted-foreground">Network</span>
                                <span className="text-sm font-mono text-foreground">{authStatus.tailnet_name || 'Unknown'}</span>
                            </div>
                            <div className="flex items-start justify-between">
                                <span className="text-xs text-muted-foreground">DNS Suffix</span>
                                <span className="text-sm font-mono text-foreground">{authStatus.tailnet_dns_suffix || 'N/A'}</span>
                            </div>
                            <div className="flex items-start justify-between">
                                <span className="text-xs text-muted-foreground">Hostname</span>
                                <span className="text-sm font-mono text-foreground">{authStatus.hostname || 'N/A'}</span>
                            </div>
                        </div>
                    </Card>

                    <Card className="p-6 bg-card/50 backdrop-blur-sm border-border/60 space-y-4">
                        <h3 className="text-xs font-bold tracking-widest text-muted-foreground uppercase">Network Stats</h3>
                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                    <Users className="w-3.5 h-3.5" />
                                    Peers
                                </div>
                                <span className="text-lg font-bold text-foreground">{authStatus.peer_count || 0}</span>
                            </div>
                            {authStatus.created_at && (
                                <div className="flex items-start justify-between">
                                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                        <Calendar className="w-3.5 h-3.5" />
                                        Created
                                    </div>
                                    <span className="text-xs font-mono text-foreground">
                                        {new Date(authStatus.created_at).toLocaleDateString()}
                                    </span>
                                </div>
                            )}
                            {authStatus.key_expiry && (
                                <div className="flex items-start justify-between">
                                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                        <Lock className="w-3.5 h-3.5" />
                                        Key Expiry
                                    </div>
                                    <span className="text-xs font-mono text-foreground">
                                        {new Date(authStatus.key_expiry).toLocaleDateString()}
                                    </span>
                                </div>
                            )}
                        </div>
                    </Card>
                </section>
            )}

            {authStatus?.authenticated && activeServices?.services && activeServices.services.length > 0 && (
                <Card className="p-6 bg-card/50 backdrop-blur-sm border-border/60 space-y-4">
                    <h3 className="text-xs font-bold tracking-widest text-muted-foreground uppercase flex items-center gap-2">
                        <Server className="w-3.5 h-3.5" />
                        Active Services
                    </h3>
                    <div className="space-y-2">
                        {activeServices.services.map((service: any) => {
                            const cleanUrl = service.public_url?.replace(/\.+$/, '')
                            return (
                                <div key={service.port} className="flex items-center justify-between p-3 rounded-lg bg-secondary/40 border border-border/50">
                                    <div className="flex flex-col gap-1 flex-1 min-w-0">
                                        <div className="flex items-center gap-3">
                                            <Badge variant={service.type === 'funnel' ? "default" : "secondary"} className="text-[10px] font-mono">
                                                {service.type.toUpperCase()}
                                            </Badge>
                                            <span className="text-sm font-mono text-foreground">Port {service.port}</span>
                                        </div>
                                        {cleanUrl && (
                                            <div className="flex items-center gap-2 ml-16">
                                                <span className="text-xs font-mono text-muted-foreground truncate">{cleanUrl}</span>
                                                <Button
                                                    size="sm"
                                                    variant="ghost"
                                                    className="h-5 w-5 p-0"
                                                    onClick={() => {
                                                        navigator.clipboard.writeText(cleanUrl)
                                                        toast.success('URL copied')
                                                    }}
                                                >
                                                    <Copy className="w-3 h-3" />
                                                </Button>
                                            </div>
                                        )}
                                    </div>
                                    {cleanUrl && (
                                        <Button
                                            size="sm"
                                            variant="ghost"
                                            className="h-7 gap-1.5 text-xs"
                                            onClick={() => window.open(cleanUrl, '_blank')}
                                        >
                                            Open
                                            <ExternalLink className="w-3 h-3" />
                                        </Button>
                                    )}
                                </div>
                            )
                        })}
                    </div>
                </Card>
            )}

            <div className="h-px bg-gradient-to-r from-transparent via-border to-transparent" />

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">

                <section className="space-y-6 h-full flex flex-col">
                    <div className="flex items-center gap-2 mb-2">
                        <Globe className="w-4 h-4 text-violet-500" />
                        <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground">Exposure Configuration</h2>
                    </div>

                    <div className="space-y-4">
                        <Card className="group p-5 bg-card/50 backdrop-blur-sm border-border/60 hover:border-violet-500/50 hover:shadow-[0_0_20px_-10px_rgba(139,92,246,0.3)]">
                            <div className="flex items-center gap-4">
                                <div className="p-2.5 rounded-lg bg-violet-500/10 text-violet-500 group-hover:scale-105 transition-transform">
                                    <Zap className="w-5 h-5" />
                                </div>
                                <div className="flex-1">
                                    <h3 className="font-semibold text-sm text-violet-500 group-hover:text-violet-400">Serve</h3>
                                    <p className="text-xs text-muted-foreground mt-0.5">Local traffic exposure</p>
                                </div>
                                <div className="flex items-center gap-3">
                                    <Input
                                        placeholder="8080"
                                        value={servePort}
                                        onChange={(e) => setServePort(e.target.value)}
                                        className="w-20 font-mono text-center h-9 bg-background/50 focus-visible:border-violet-500/50 focus-visible:ring-violet-500/30"
                                    />
                                    <Button onClick={handleServeStart} disabled={loading} className="h-9 bg-violet-600 hover:bg-violet-500 text-white shadow-lg shadow-violet-500/20">
                                        Start
                                    </Button>
                                </div>
                            </div>
                        </Card>

                        <Card className="group p-5 bg-card/50 backdrop-blur-sm border-border/60 hover:border-violet-500/50 hover:shadow-[0_0_20px_-10px_rgba(139,92,246,0.3)]">
                            <div className="flex items-center gap-4">
                                <div className="p-2.5 rounded-lg bg-violet-500/10 text-violet-500 group-hover:scale-105 transition-transform">
                                    <Globe className="w-5 h-5" />
                                </div>
                                <div className="flex-1">
                                    <h3 className="font-semibold text-sm text-violet-500 group-hover:text-violet-400">Funnel</h3>
                                    <p className="text-xs text-muted-foreground mt-0.5">Public internet access</p>
                                </div>
                                <div className="flex items-center gap-3">
                                    <Input
                                        placeholder="443"
                                        value={funnelPort}
                                        onChange={(e) => setFunnelPort(e.target.value)}
                                        className="w-20 font-mono text-center h-9 bg-background/50 focus-visible:border-violet-500/50 focus-visible:ring-violet-500/30"
                                    />
                                    <Button onClick={handleFunnelStart} disabled={loading} className="h-9 bg-violet-600 hover:bg-violet-500 text-white shadow-lg shadow-violet-500/20">
                                        Start
                                    </Button>
                                </div>
                            </div>
                        </Card>
                    </div>
                </section>

                <section className="space-y-6 h-full flex flex-col">
                    <div className="flex items-center gap-2 mb-2">
                        <Lock className="w-4 h-4 text-emerald-500" />
                        <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground">System Access</h2>
                    </div>

                    <div className="space-y-4">
                        <Card className="group p-5 bg-card/50 backdrop-blur-sm border-border/60 hover:border-emerald-500/30">
                            <div className="flex items-center gap-4">
                                <div className="p-2.5 rounded-lg bg-secondary text-foreground group-hover:bg-emerald-500/10 group-hover:text-emerald-500 transition-colors">
                                    <Terminal className="w-5 h-5" />
                                </div>
                                <div className="flex-1">
                                    <h3 className="font-semibold text-sm">SSH Access</h3>
                                    <p className="text-xs text-muted-foreground mt-0.5">Remote shell connections</p>
                                </div>
                                <div className="flex items-center gap-3">
                                    <Input
                                        placeholder="22"
                                        disabled
                                        className="w-20 font-mono text-center h-9 bg-background/50 opacity-50"
                                    />
                                    <Button onClick={handleSSHEnable} disabled={loading} className="bg-emerald-600 hover:bg-emerald-500 text-white shadow-lg shadow-emerald-500/20">
                                        Enable
                                    </Button>
                                </div>
                            </div>
                        </Card>
                    </div>
                </section>

            </div>
        </div >
    )
}
