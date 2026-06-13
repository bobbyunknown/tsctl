import { useState } from 'react'
import { serveApi, funnelApi, sshApi, proxyApi } from '../lib/api'
import { Input } from '@/components/ui/input'
import { toast } from 'sonner'
import { Zap, Globe, Terminal, Network } from 'lucide-react'

interface ExposureControlsProps {
    onSuccess?: () => void
}

export function ExposureControls({ onSuccess }: ExposureControlsProps) {
    const [loading, setLoading] = useState(false)
    const [selectedType, setSelectedType] = useState<'serve'|'funnel'|'ssh'|'proxy'|'proxy_all'>('serve')
    const [portInput, setPortInput] = useState('')

    const handleAddSubmit = async () => {
        setLoading(true)
        try {
            if (selectedType === 'serve') {
                if (!portInput) throw new Error("Port required")
                await serveApi.start(Number(portInput), false)
                toast.success('Serve started')
            } else if (selectedType === 'funnel') {
                if (!portInput) throw new Error("Port required")
                await funnelApi.start(Number(portInput), false)
                toast.success('Funnel started')
            } else if (selectedType === 'ssh') {
                await sshApi.enable()
                toast.success('SSH access enabled')
            } else if (selectedType === 'proxy') {
                if (!portInput) throw new Error("Port required")
                await proxyApi.start({ mode: 'single', port: Number(portInput) })
                toast.success('Proxy started')
            } else if (selectedType === 'proxy_all') {
                await proxyApi.start({ mode: 'all', scan_interval: 5 })
                toast.success('Auto-scan proxy started')
            }
            setPortInput('')
            if (onSuccess) onSuccess()
        } catch (error: any) {
            toast.error(error.message || 'Failed to apply configuration')
        }
        setLoading(false)
    }

    return (
        <div className="bg-card-defi border border-[#02d7f2]/15 rounded overflow-hidden shadow-lg uppercase h-full flex flex-col">
            <div className="border-b border-[#02d7f2]/15 bg-black/40 px-5 py-4">
                <h2 className="text-sm font-bold tracking-widest text-cyan-glow display-font">EXPOSURE CONTROLS</h2>
            </div>

            <div className="p-5 bg-black/20 flex-1 flex flex-col">
                <div className="grid grid-cols-2 sm:grid-cols-5 gap-4 mb-6">
                    {[
                        { id: 'serve', icon: Zap, label: 'SERVE' },
                        { id: 'funnel', icon: Globe, label: 'FUNNEL' },
                        { id: 'ssh', icon: Terminal, label: 'SSH' },
                        { id: 'proxy', icon: Network, label: 'PROXY PORT' },
                        { id: 'proxy_all', icon: Network, label: 'PROXY ALL' },
                    ].map(opt => {
                        const Icon = opt.icon
                        const isActive = selectedType === opt.id
                        return (
                            <button
                                key={opt.id}
                                onClick={() => setSelectedType(opt.id as any)}
                                className={`flex flex-col items-center justify-center p-4 border rounded text-xs gap-3 transition-all duration-200 display-font font-bold ${
                                    isActive 
                                    ? 'border-[#02d7f2] bg-[#02d7f2]/10 text-[#02d7f2] shadow-[0_0_15px_rgba(2,215,242,0.2)]' 
                                    : 'border-[#02d7f2]/15 hover:border-[#02d7f2]/50 text-muted-foreground hover:text-[#02d7f2] bg-black/50'
                                }`}
                            >
                                <Icon className="w-6 h-6" />
                                <span className="text-[10px] text-center tracking-widest uppercase">{opt.label}</span>
                            </button>
                        )
                    })}
                </div>

                <div className="min-h-[80px] flex items-center mb-6">
                    {['serve', 'funnel', 'proxy'].includes(selectedType) ? (
                        <div className="space-y-3 w-full">
                            <label className="text-xs font-bold text-[#02d7f2] block tracking-widest display-font">TARGET PORT</label>
                            <Input 
                                type="number" 
                                placeholder="8080" 
                                value={portInput}
                                onChange={e => setPortInput(e.target.value)}
                                className="input-defi max-w-[250px] font-bold text-lg h-12 rounded"
                                autoFocus
                            />
                        </div>
                    ) : selectedType === 'ssh' ? (
                        <div className="text-sm text-[#007aff] bg-black/50 p-4 border-l-2 border-[#007aff] w-full font-mono rounded-r">
                            ENABLE REMOTE SHELL OVER TAILNET.
                        </div>
                    ) : (
                        <div className="text-sm text-[#39ff14] bg-black/50 p-4 border-l-2 border-[#39ff14] w-full font-mono rounded-r">
                            AUTO-SCAN LOCAL PORTS -{'>'} TAILNET.
                        </div>
                    )}
                </div>

                <div className="flex justify-end pt-5 border-t border-[#02d7f2]/15 mt-auto">
                    <button 
                        onClick={handleAddSubmit} 
                        disabled={loading} 
                        className="px-8 h-12 border border-[#02d7f2] text-[#02d7f2] hover:bg-[#02d7f2] hover:text-black hover:shadow-[0_0_15px_rgba(2,215,242,0.4)] rounded font-bold tracking-widest uppercase transition-all duration-300 display-font"
                    >
                        EXECUTE_CMD
                    </button>
                </div>
            </div>
        </div>
    )
}
