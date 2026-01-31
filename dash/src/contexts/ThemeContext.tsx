import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'

type ThemeMode = 'light' | 'dark' | 'system'

interface ThemeContextType {
    mode: ThemeMode
    isDark: boolean
    setMode: (mode: ThemeMode) => void
    toggleTheme: () => void
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [mode, setModeState] = useState<ThemeMode>(() => {
        const saved = localStorage.getItem('theme-mode')
        return (saved as ThemeMode) || 'system'
    })

    const [isDark, setIsDark] = useState(() => {
        if (mode === 'system') {
            return window.matchMedia('(prefers-color-scheme: dark)').matches
        }
        return mode === 'dark'
    })

    useEffect(() => {
        const updateTheme = () => {
            let dark: boolean
            if (mode === 'system') {
                dark = window.matchMedia('(prefers-color-scheme: dark)').matches
            } else {
                dark = mode === 'dark'
            }
            setIsDark(dark)

            if (dark) {
                document.documentElement.classList.add('dark')
            } else {
                document.documentElement.classList.remove('dark')
            }
        }

        updateTheme()

        if (mode === 'system') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
            mediaQuery.addEventListener('change', updateTheme)
            return () => mediaQuery.removeEventListener('change', updateTheme)
        }
    }, [mode])

    const setMode = (newMode: ThemeMode) => {
        setModeState(newMode)
        localStorage.setItem('theme-mode', newMode)
    }

    const toggleTheme = () => {
        if (mode === 'light') {
            setMode('dark')
        } else if (mode === 'dark') {
            setMode('system')
        } else {
            setMode('light')
        }
    }

    return (
        <ThemeContext.Provider value={{ mode, isDark, setMode, toggleTheme }}>
            {children}
        </ThemeContext.Provider>
    )
}

export function useTheme() {
    const context = useContext(ThemeContext)
    if (context === undefined) {
        throw new Error('useTheme must be used within a ThemeProvider')
    }
    return context
}
