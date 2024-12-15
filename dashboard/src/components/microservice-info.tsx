import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

type ServiceType = {
  id: string;
  name: string;
  description: string;
  status: string;
  version: string;
  squad: string;
  lastMainUpdate: string;
  lastBranchUpdate: string;
  lastDeploy: string;
}

export function MicroserviceInfo({ service }: { service: ServiceType }) {
  return (
      <Card>
        <CardContent className="pt-6">
          <div className="space-y-2">
            <h2 className="text-2xl font-bold">{service.name}</h2>
            <p className="text-sm text-muted-foreground">{service.description}</p>
            <div className="flex items-center space-x-2">
              <span className="font-semibold">Squad:</span>
              <span>{service.squad}</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="font-semibold">Status:</span>
              <Badge variant={service.status === 'Healthy' ? "default" : service.status === 'Warning' ? "secondary" : "destructive"}>{service.status}</Badge>
            </div>
            <div>
              <span className="font-semibold">Version:</span> {service.version}
            </div>
          </div>
        </CardContent>
      </Card>
  )
}