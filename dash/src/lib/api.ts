import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
})

export interface AuthStatus {
    authenticated: boolean
    backend_state: string
    auth_url?: string
    node_key?: string
    hostname?: string
    dns_name?: string
    ips?: string[]
    user_display_name?: string
    user_email?: string
    user_profile_pic?: string
    tailnet_name?: string
    tailnet_dns_suffix?: string
    is_admin?: boolean
    is_owner?: boolean
    created_at?: string
    key_expiry?: string
    peer_count?: number
}

export interface ApiResponse<T = any> {
    success: boolean
    message?: string
    data?: T
}

export interface ServiceInfo {
    port: number
    type: 'serve' | 'funnel'
    local_url: string
    public_url?: string
}

export interface ServeStatus {
    services: ServiceInfo[]
}

export interface ProxyInfo {
    port: number
    ip: string
    protocol: string
}

export interface ProxyStartRequest {
    mode: string // 'single', 'multi', 'all'
    port?: number
    ports?: number[]
    exclude_ports?: number[]
    scan_interval?: number
}

export const authApi = {
    getStatus: async (): Promise<AuthStatus> => {
        const { data } = await api.get<ApiResponse<AuthStatus>>('/api/v1/auth/status')
        return data.data!
    },
    logout: async (): Promise<void> => {
        await api.post('/api/v1/auth/logout')
    },
}

export const statusApi = {
    getStatus: async (): Promise<string> => {
        const { data } = await api.get<ApiResponse<string>>('/api/v1/status')
        return data.data!
    },
}

export const serveApi = {
    start: async (port: number, background: boolean = false) => {
        const { data } = await api.post<ApiResponse>('/api/v1/serve', { port, background })
        return data
    },
    getStatus: async (): Promise<ServeStatus> => {
        const { data } = await api.get<ApiResponse<ServeStatus>>('/api/v1/serve/status')
        return data.data!
    },
    reset: async () => {
        const { data } = await api.delete<ApiResponse>('/api/v1/serve')
        return data
    },
}

export const funnelApi = {
    start: async (port: number, background: boolean = false) => {
        const { data } = await api.post<ApiResponse>('/api/v1/funnel', { port, background })
        return data
    },
    getStatus: async (): Promise<ServeStatus> => {
        const { data } = await api.get<ApiResponse<ServeStatus>>('/api/v1/funnel/status')
        return data.data!
    },
    reset: async () => {
        const { data } = await api.delete<ApiResponse>('/api/v1/funnel')
        return data
    },
}

export const sshApi = {
    enable: async () => {
        const { data } = await api.post<ApiResponse>('/api/v1/ssh/enable')
        return data
    },
}

export const proxyApi = {
    start: async (req: ProxyStartRequest) => {
        const { data } = await api.post<ApiResponse>('/api/v1/proxy/start', req)
        return data
    },
    stop: async (ports: number[]) => {
        const { data } = await api.post<ApiResponse>('/api/v1/proxy/stop', { ports })
        return data
    },
    stopAll: async () => {
        const { data } = await api.delete<ApiResponse>('/api/v1/proxy')
        return data
    },
    getStatus: async (): Promise<ProxyInfo[]> => {
        const { data } = await api.get<ApiResponse<ProxyInfo[]>>('/api/v1/proxy/status')
        return data.data!
    },
}

export const logsApi = {
    getAppLogs: async (lines: string = '100') => {
        const { data } = await api.get<ApiResponse<string[]>>(`/api/v1/logs/app?lines=${lines}`)
        return data
    },
    clearLogs: async () => {
        const { data } = await api.delete<ApiResponse>('/api/v1/logs/app')
        return data
    },
}
