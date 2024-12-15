import DashboardGrid from '@/components/dashboard-grid'

export default function Home() {
  return (
    <main className="min-h-screen">
      <div className="container mx-auto px-4">
        <h1 className="text-3xl font-bold tracking-tight mb-4 py-4">Microservices Dashboard</h1>
        <DashboardGrid />
      </div>
    </main>
  )
}

