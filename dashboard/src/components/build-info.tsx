import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

const fakeBuildHistory = [
  { id: 1, version: "v2.3.1", status: "success", date: "2023-11-15 14:00 UTC", duration: "5m 23s" },
  { id: 2, version: "v2.3.0", status: "success", date: "2023-11-10 09:00 UTC", duration: "6m 12s" },
  { id: 3, version: "v2.2.5", status: "failed", date: "2023-11-05 11:30 UTC", duration: "4m 56s" },
]

export function BuildInfo() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Build Information</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-semibold mb-2">Latest Build</h3>
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <span
className="font-medium">Version:</span>
                <Badge>v2.3.1</Badge>
              </div>
              <div className="flex items-center space-x-2">
                <span className="font-medium">Status:</span>
                <Badge variant="default">Success</Badge>
              </div>
              <div>
                <span className="font-medium">Date:</span> 2023-11-15 14:00 UTC
              </div>
              <div>
                <span className="font-medium">Duration:</span> 5m 23s
              </div>
            </div>
          </div>
          <div>
            <h3 className="text-lg font-semibold mb-2">Build History</h3>
            <table className="w-full">
              <thead>
                <tr className="text-left">
                  <th>Version</th>
                  <th>Status</th>
                  <th>Date</th>
                  <th>Duration</th>
                </tr>
              </thead>
              <tbody>
                {fakeBuildHistory.map((build) => (
                  <tr key={build.id}>
                    <td>{build.version}</td>
                    <td>
                      <Badge variant={build.status === "success" ? "default" : "destructive"}>
                        {build.status}
                      </Badge>
                    </td>
                    <td>{build.date}</td>
                    <td>{build.duration}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}


