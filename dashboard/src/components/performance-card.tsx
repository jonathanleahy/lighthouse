import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Activity } from 'lucide-react'

type ServiceType = {
  id: string;
  name: string;
  requests: string;
}

export function PerformanceCard({ service }: { service: ServiceType }) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Performance</CardTitle>
        <Activity className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="flex flex-col space-y-2">
          <span className="text-xs text-muted-foreground">Requests (24h)</span>
          <span className="text-2xl font-bold">{service.requests}</span>
        </div>
      </CardContent>
    </Card>
  )
}


