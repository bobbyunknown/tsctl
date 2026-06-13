import { useState, useEffect } from 'react'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { RefreshCcw, TerminalSquare, Trash2 } from 'lucide-react'
import { logsApi } from '@/lib/api'
import { useConfirm } from '@/contexts/ConfirmContext'

export default function LogsPage() {
    const [logs, setLogs] = useState<string[]>([])
    const [lines, setLines] = useState('100')
    const [loading, setLoading] = useState(false)
    const { confirm } = useConfirm()

    const fetchLogs = async () => {
        setLoading(true)
        try {
            const data = await logsApi.getAppLogs(lines)
            if (data.data) {
                setLogs(data.data)
            }
        } catch (error) {
            console.error(error)
        }
        setLoading(false)
    }

    const handleClearLogs = async () => {
        const isConfirmed = await confirm({
            title: 'Clear Logs',
            description: 'Are you sure you want to clear all logs? This action cannot be undone.',
            confirmText: 'Clear',
            cancelText: 'Cancel',
            isDestructive: true
        })
        if (!isConfirmed) return

        try {
            await logsApi.clearLogs()
            setLogs([])
            await fetchLogs()
        } catch (error) {
            console.error('Failed to clear logs:', error)
        }
    }

    useEffect(() => {
        fetchLogs()
    }, [lines])

    return (
        <div className="max-w-6xl mx-auto h-[calc(100vh-8rem)] flex flex-col gap-4 animate-in fade-in duration-500 py-8">
            <div className="flex items-center justify-between">
                <div className="flex flex-col gap-1">
                    <h1 className="text-2xl font-light tracking-tight text-foreground flex items-center gap-3">
                        <TerminalSquare className="w-6 h-6 text-muted-foreground" />
                        System Logs
                    </h1>
                    <p className="text-sm text-muted-foreground">Live stream from application output</p>
                </div>

                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-2 px-3 py-1.5 rounded bg-secondary/40 border border-border/50">
                        <span className="text-[10px] uppercase text-muted-foreground font-medium">Buffer Size</span>
                        <Input
                            value={lines}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setLines(e.target.value)}
                            className="w-16 h-5 text-xs font-mono bg-transparent border-none p-0 focus-visible:ring-0 text-right"
                        />
                    </div>

                    <Button onClick={fetchLogs} disabled={loading} size="sm" variant="ghost" className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground">
                        <RefreshCcw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>

                    <Button onClick={handleClearLogs} size="sm" variant="ghost" className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive">
                        <Trash2 className="w-4 h-4" />
                    </Button>
                </div>
            </div>

            <Card className="flex-1 bg-card/50 backdrop-blur-sm border-border/60 rounded-lg overflow-hidden flex flex-col relative">
                <div className="absolute top-0 left-0 right-0 h-6 bg-gradient-to-b from-background/20 to-transparent pointer-events-none z-10" />
                <div className="flex-1 overflow-auto p-4 font-mono text-xs space-y-1 scrollbar-thin scrollbar-thumb-border scrollbar-track-transparent">
                    {logs.length === 0 ? (
                        <div className="h-full flex flex-col items-center justify-center text-muted-foreground">
                            <TerminalSquare className="w-8 h-8 mb-3 opacity-20" />
                            <span>No active log output</span>
                        </div>
                    ) : (
                        logs.map((log, i) => (
                            <div key={i} className="flex gap-3 text-muted-foreground hover:bg-accent hover:text-foreground py-0.5 px-2 rounded -mx-2 transition-colors">
                                <span className="text-muted-foreground/50 select-none w-8 text-right shrink-0">{i + 1}</span>
                                <span className="break-all whitespace-pre-wrap">{log}</span>
                            </div>
                        ))
                    )}
                </div>
            </Card>
        </div>
    )
}
