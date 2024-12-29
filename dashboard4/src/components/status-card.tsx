import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { GitBranch, Clock, CheckCircle, AlertTriangle, XCircle, Calendar } from 'lucide-react'

type ServiceType = {
  id: string;
  name: string;
  status: string;
  version: string;
  uptime: string;
  lastDeploy: string;
}

export function StatusCard({ service }: { service: ServiceType }) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">Status</CardTitle>
        <StatusBadge status={service.status} />
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4">
          <div className="flex flex-col space-y-2">
            <span className="text-xs text-muted-foreground">Version</span>
            <div className="flex items-center space-x-2">
              <GitBranch className="h-4 w-4 text-muted-foreground" />
              <span className="text-xl font-bold">{service.version}</span>
            </div>
          </div>
          <div className="flex flex-col space-y-2">
            <span className="text-xs text-muted-foreground">Uptime</span>
            <div className="flex items-center space-x-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span className="text-xl font-bold">{service.uptime}</span>
            </div>
          </div>
          <div className="flex flex-col space-y-2 col-span-2">
            <span className="text-xs text-muted-foreground">Last Deployment</span>
            <div className="flex items-center space-x-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <span className="text-xl font-bold">{service.lastDeploy}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function StatusBadge({ status }: { status: string }) {
  const statusConfig = {
    Healthy: { color: "bg-green-500", icon: CheckCircle },
    Warning: { color: "bg-yellow-500", icon: AlertTriangle },
    Critical: { color: "bg-red-500", icon: XCircle },
  }[status] || { color: "bg-gray-500", icon: AlertTriangle }

  const Icon = statusConfig.icon

  return (
    <Badge className={`${statusConfig.color} text-primary-foreground`}>
      <Icon className="mr-1 h-3 w-3" />
      {status}
    </Badge>
  )
}


