import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Clock } from 'lucide-react'

type DeploymentType = {
  version: string;
  date: string;
  environment: string;
}

export function DeploymentInfo({ deployment }: { deployment: DeploymentType }) {
  return (
    <Card>
      <CardContent className="pt-6">
        <div className="space-y-2">
          <h3 className="text-lg font-semibold">Latest Deployment</h3>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Version:</span>
            <Badge variant="secondary">{deployment.version}</Badge>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Deployed:</span>
            <div className="flex items-center">
              <Clock className="mr-1 h-4 w-4 text-muted-foreground" />
              <span>{deployment.date}</span>
            </div>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">Environment:</span>
            <Badge>{deployment.environment}</Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}


