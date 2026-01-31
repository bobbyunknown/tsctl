import { Routes, Route } from 'react-router-dom'
import HomePage from './pages/HomePage'
import LogsPage from './pages/LogsPage'
import AppShell from './components/layout/AppShell'

function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/logs" element={<LogsPage />} />
      </Routes>
    </AppShell>
  )
}

export default App
