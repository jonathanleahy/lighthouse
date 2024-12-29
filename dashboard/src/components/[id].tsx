import { useRouter } from 'next/router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Separator } from "@/components/ui/separator"
import { ScrollArea } from "@/components/ui/scroll-area"
import { LineChart, BarChart, PieChart } from 'lucide-react'
import Link from 'next/link'

export default function MicroserviceDetails() {
  const router = useRouter()
  const { id } = router.query

  // In a real application, you would fetch the microservice details based on the ID
  const microserviceName = id === 'auth' ? 'User Authentication' :
                           id === 'payment' ? 'Payment Processing' :
                           id === 'inventory' ? 'Inventory Management' :
                           id === 'notification' ? 'Notification Service' : 'Unknown Microservice'

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-white">
      <header className="bg-gray-800 shadow-lg">
        <div className="container mx-auto px-4 py-6">
          <Link href="/" className="text-blue-500 hover:underline mb-2 inline-block">
            &larr; Back to Overview
          </Link>
          <h1 className="text-3xl font-bold">{microserviceName} Dashboard</h1>
        </div>
      </header>
      <main className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <StatusCard />
          <PerformanceCard />
          <ActivityCard />
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <DeploymentInfo />
          <ErrorLogs />
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <BuildInfo />
          <GitHubInfo />
        </div>
        <SquadMembers />
      </main>
    </div>
  )
}

function StatusCard() {
  return (
    <Card className="bg-gradient-to-br from-blue-600 to-blue-700">
      <CardHeader>
        <CardTitle className="text-white">Status</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <p className="text-5xl font-bold text-white">v2.3.1</p>
            <p className="text-lg text-blue-100">Current Version</p>
          </div>
          <Badge className="bg-green-500 hover:bg-green-600 text-white text-lg py-1 px-3">Healthy</Badge>
        </div>
      </CardContent>
    </Card>
  )
}

function PerformanceCard() {
  return (
    <Card className="bg-gradient-to-br from-purple-600 to-purple-700">
      <CardHeader>
        <CardTitle className="text-white">Performance</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <p className="text-5xl font-bold text-white">99.9%</p>
            <p className="text-lg text-purple-100">Uptime</p>
          </div>
          <LineChart className="w-16 h-16 text-purple-200" />
        </div>
      </CardContent>
    </Card>
  )
}

function ActivityCard() {
  return (
    <Card className="bg-gradient-to-br from-orange-600 to-orange-700">
      <CardHeader>
        <CardTitle className="text-white">Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <p className="text-5xl font-bold text-white">1.2M</p>
            <p className="text-lg text-orange-100">Requests / Day</p>
          </div>
          <BarChart className="w-16 h-16 text-orange-200" />
        </div>
      </CardContent>
    </Card>
  )
}

function DeploymentInfo() {
  const deployments = [
    { version: "v2.3.1", date: "2023-11-15 14:30 UTC", environment: "Production" },
    { version: "v2.3.0", date: "2023-11-10 09:15 UTC", environment: "Production" },
    { version: "v2.2.5", date: "2023-11-05 11:45 UTC", environment: "Production" },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <PieChart className="w-6 h-6 mr-2 text-blue-500" />
          Deployment Information
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-semibold mb-2">Current Deployment</h3>
            <div className="flex items-center space-x-2">
              <Badge variant="outline">v2.3.1</Badge>
              <span className="text-sm text-muted-foreground">Deployed on 2023-11-15 14:30 UTC</span>
            </div>
          </div>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold mb-2">Deployment History</h3>
            <ul className="space-y-2">
              {deployments.map((deployment, index) => (
                <li key={index} className="flex items-center justify-between">
                  <div>
                    <Badge variant="outline">{deployment.version}</Badge>
                    <span className="ml-2 text-sm text-muted-foreground">{deployment.environment}</span>
                  </div>
                  <span className="text-sm text-muted-foreground">{deployment.date}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function ErrorLogs() {
  const errors = [
    { id: 1, message: "Invalid token", count: 23, lastOccurred: "2023-11-23 09:45 UTC" },
    { id: 2, message: "Database connection timeout", count: 5, lastOccurred: "2023-11-22 14:30 UTC" },
    { id: 3, message: "Rate limit exceeded", count: 150, lastOccurred: "2023-11-23 10:15 UTC" },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <LineChart className="w-6 h-6 mr-2 text-red-500" />
          Error Logs
        </CardTitle>
        <CardDescription>Recent errors and exceptions</CardDescription>
      </CardHeader>
      <CardContent>
        <ul className="space-y-4">
          {errors.map((error) => (
            <li key={error.id} className="flex items-center justify-between">
              <div>
                <p className="font-medium">{error.message}</p>
                <p className="text-sm text-muted-foreground">Last occurred: {error.lastOccurred}</p>
              </div>
              <Badge variant="destructive">{error.count}</Badge>
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  )
}

function BuildInfo() {
  const builds = [
    { id: 1, version: "v2.3.1", status: "success", date: "2023-11-15 14:00 UTC", duration: "5m 23s" },
    { id: 2, version: "v2.3.0", status: "success", date: "2023-11-10 09:00 UTC", duration: "6m 12s" },
    { id: 3, version: "v2.2.5", status: "failed", date: "2023-11-05 11:30 UTC", duration: "4m 56s" },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <BarChart className="w-6 h-6 mr-2 text-green-500" />
          Build Information
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-semibold mb-2">Latest Build</h3>
            <div className="flex items-center space-x-2">
              <Badge variant="outline">v2.3.1</Badge>
              <Badge variant="default">Success</Badge>
              <span className="text-sm text-muted-foreground">Duration: 5m 23s</span>
            </div>
          </div>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold mb-2">Build History</h3>
            <ul className="space-y-2">
              {builds.map((build) => (
                <li key={build.id} className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline">{build.version}</Badge>
                    <Badge variant={build.status === "success" ? "default" : "destructive"}>
                      {build.status}
                    </Badge>
                  </div>
                  <div className="text-sm text-muted-foreground">
                    <span>{build.date}</span>
                    <span className="ml-2">({build.duration})</span>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function GitHubInfo() {
  const fakeReadme = `
# User Authentication Microservice

This microservice handles user authentication and authorization for our platform.

## Features

- User registration
- Login with JWT
- Password reset
- OAuth2 integration
- Role-based access control

## Getting Started

1. Clone the repository
2. Install dependencies: \`npm install\`
3. Set up environment variables
4. Run the server: \`npm start\`

## API Documentation

API documentation is available at \`/api-docs\` when running the server.

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
`

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <svg viewBox="0 0 24 24" className="w-6 h-6 mr-2 text-gray-400" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          GitHub Repository
        </CardTitle>
        <CardDescription>
          <a href="#" className="text-blue-500 hover:underline">https://github.com/ourcompany/user-auth-microservice</a>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex space-x-2">
            <Badge>main</Badge>
            <Badge variant="outline">v2.3.1</Badge>
            <Badge variant="outline">Node.js</Badge>
          </div>
          <div>
            <h3 className="text-lg font-semibold mb-2">Recent Commits</h3>
            <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
              <li>Update dependencies (2 hours ago)</li>
              <li>Fix OAuth token refresh bug (1 day ago)</li>
              <li>Add rate limiting to login endpoint (3 days ago)</li>
            </ul>
          </div>
          <Separator />
          <div>
            <h3 className="text-lg font-semibold mb-2">README.md</h3>
            <ScrollArea className="h-[200px] w-full rounded border p-4">
              <pre className="text-sm">{fakeReadme}</pre>
            </ScrollArea>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function SquadMembers() {
  const members = [
    { id: 1, name: "Alice Johnson", role: "Squad Lead", email: "alice@example.com" },
    { id: 2, name: "Bob Smith", role: "Senior Developer", email: "bob@example.com" },
    { id: 3, name: "Charlie Davis", role: "Backend Developer", email: "charlie@example.com" },
    { id: 4, name: "Diana Miller", role: "Frontend Developer", email: "diana@example.com" },
    { id: 5, name: "Ethan Brown", role: "QA Engineer", email: "ethan@example.com" },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6 mr-2 text-indigo-500">
            <path fillRule="evenodd" d="M8.25 6.75a3.75 3.75 0 117.5 0 3.75 3.75 0 01-7.5 0zM15.75 9.75a3 3 0 116 0 3 3 0 01-6 0zM2.25 9.75a3 3 0 116 0 3 3 0 01-6 0zM6.31 15.117A6.745 6.745 0 0112 12a6.745 6.745 0 016.709 3.115.75.75 0 01-.351.92 6.705 6.705 0 01-9.712 0 .75.75 0 01-.35-.92zm-.31 2.133a5.255 5.255 0 0110.62 0 .75.75 0 01-.351.92 6.705 6.705 0 01-9.712 0 .75.75 0 01-.351-.92zm9.06-1.484a3.742 3.742 0 012.43-1.493 3.742 3.742 0 012.43 1.493.75.75 0 01-.326.92 3.745 3.745 0 01-4.208 0 .75.75 0 01-.326-.92zm3.06 2.134a2.255 2.255 0 00-4.51 0 .75.75 0 01-.326.92 3.745 3.745 0 004.208 0 .75.75 0 01-.326-.92z" clipRule="evenodd" />
          </svg>
          Squad Members
        </CardTitle>
        <CardDescription>Team responsible for the User Authentication Microservice</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {members.map((member) => (
            <div key={member.id} className="flex items-center space-x-4 bg-muted/50 rounded-lg p-3">
              <Avatar className="w-12 h-12">
                <AvatarImage src={`https://api.dicebear.com/6.x/initials/svg?seed=${member.name}`} />
                <AvatarFallback>{member.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
              </Avatar>
              <div>
                <p className="font-medium">{member.name}</p>
                <p className="text-sm text-muted-foreground">{member.role}</p>
                <p className="text-sm text-muted-foreground">{member.email}</p>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}


