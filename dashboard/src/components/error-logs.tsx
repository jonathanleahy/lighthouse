import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

const fakeErrors = [
  { id: 1, message: "Invalid token", count: 23, lastOccurred: "2023-11-23 09:45 UTC" },
  { id: 2, message: "Database connection timeout", count: 5, lastOccurred: "2023-11-22 14:30 UTC" },
  { id: 3, message: "Rate limit exceeded", count: 150, lastOccurred: "2023-11-23 10:15 UTC" },
]

export function ErrorLogs() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Error Logs</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {fakeErrors.map((error) => (
            <div key={error.id} className="flex items-center justify-between">
              <div>
                <p className="font-medium">{error.message}</p>
                <p className="text-sm text-muted-foreground">Last occurred: {error.lastOccurred}</p>
              </div>
              <Badge variant="destructive">{error.count}</Badge>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}


