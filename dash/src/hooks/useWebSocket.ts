import { useEffect, useState } from 'react'
import useWebSocket, { ReadyState } from 'react-use-websocket'
import { type AuthStatus } from '../lib/api'

interface WSEvent {
    type: string
    data: any
}

export function useTailscaleWS() {
    const [authStatus, setAuthStatus] = useState<AuthStatus | null>(null)
    const [serveStatus, setServeStatus] = useState<string | null>(null)
    const [connectionState, setConnectionState] = useState<string>('connecting')

    const wsUrl = import.meta.env.VITE_WS_URL

    const { lastJsonMessage, readyState } = useWebSocket(wsUrl, {
        shouldReconnect: () => true,
        reconnectInterval: 3000,
    }
    )

    useEffect(() => {
        if (lastJsonMessage) {
            const event = lastJsonMessage as WSEvent

            switch (event.type) {
                case 'auth_status_changed':
                    setAuthStatus(event.data)
                    break
                case 'serve_status_changed':
                    setServeStatus(event.data)
                    break
            }
        }
    }, [lastJsonMessage])

    useEffect(() => {
        const states = {
            [ReadyState.CONNECTING]: 'connecting',
            [ReadyState.OPEN]: 'connected',
            [ReadyState.CLOSING]: 'closing',
            [ReadyState.CLOSED]: 'disconnected',
            [ReadyState.UNINSTANTIATED]: 'uninstantiated',
        }
        setConnectionState(states[readyState])
    }, [readyState])

    useEffect(() => {
        const fetchInitialStatus = async () => {
            try {
                const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
                const response = await fetch(`${apiBaseUrl}/api/v1/auth/status`)
                const data = await response.json()
                if (data.success && data.data) {
                    setAuthStatus(data.data)
                }
            } catch (error) {
                console.error('Failed to fetch initial auth status:', error)
            }
        }
        fetchInitialStatus()
    }, [])

    return {
        authStatus,
        serveStatus,
        connectionState,
        isConnected: readyState === ReadyState.OPEN,
    }
}
