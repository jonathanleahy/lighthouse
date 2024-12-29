import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

type SummaryProps = {
  totalEnvironments: number
  latestTag: string
  environmentsOnLatest: number
}

export function DashboardSummary({ totalEnvironments, latestTag, environmentsOnLatest }: SummaryProps) {
  const percentage = (environmentsOnLatest / totalEnvironments) * 100

  return (
    <Card className="col-span-full">
      <CardHeader>
        <CardTitle>Dashboard Summary</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-4 md:grid-cols-3">
        <div>
          <p className="text-sm font-medium">Total Environments</p>
          <p className="text-2xl font-bold">{totalEnvironments}</p>
        </div>
        <div>
          <p className="text-sm font-medium">Latest Tag</p>
          <p className="text-2xl font-bold">{latestTag}</p>
        </div>
        <div>
          <p className="text-sm font-medium">Environments on Latest Tag</p>
          <p className="text-2xl font-bold">{environmentsOnLatest} out of {totalEnvironments}</p>
          <Progress value={percentage} className="mt-2" />
        </div>
      </CardContent>
    </Card>
  )
}


