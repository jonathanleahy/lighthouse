import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { GitBranch, GitCommit, FileText, GitMerge } from 'lucide-react'

const fakeRepoInfo = {
    branch: 'main',
    version: 'v2.3.1',
    language: 'Node.js',
    lastUpdated: '2023-11-25',
    openIssues: 5,
    pullRequests: 2
}

const fakeCommits = [
    { id: '1', message: 'Update dependencies', author: 'Alice', date: '2 hours ago' },
    { id: '2', message: 'Fix OAuth token refresh bug', author: 'Bob', date: '1 day ago' },
    { id: '3', message: 'Add rate limiting to login endpoint', author: 'Charlie', date: '3 days ago' },
]

const fakeBranches = [
    { name: 'main', status: 'open', security: 'passed', build: 'passed', lastCommit: '1 hour ago', lastCommitAuthor: 'Alice Johnson' },
    { name: 'feature/user-profile', status: 'draft', security: 'passed', build: 'failed', lastCommit: '3 hours ago', lastCommitAuthor: 'Bob Smith' },
    { name: 'bugfix/login-issue', status: 'open', security: 'pending', build: 'passed', lastCommit: '1 day ago', lastCommitAuthor: 'Charlie Davis' },
]

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

function RepoOverview({ info }: { info: typeof fakeRepoInfo }) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                    <GitBranch className="h-5 w-5" />
                    <span>Repository Overview</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    <div className="flex space-x-2">
                        <Badge>{info.branch}</Badge>
                        <Badge variant="outline">{info.version}</Badge>
                        <Badge variant="outline">{info.language}</Badge>
                    </div>
                    <div className="space-y-2">
                        <p className="text-sm">
                            <span className="font-medium">Last Updated:</span> {info.lastUpdated}
                        </p>
                        <p className="text-sm">
                            <span className="font-medium">Open Issues:</span> {info.openIssues}
                        </p>
                        <p className="text-sm">
                            <span className="font-medium">Pull Requests:</span> {info.pullRequests}
                        </p>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}

function RecentCommits({ commits }: { commits: typeof fakeCommits }) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                    <GitCommit className="h-5 w-5" />
                    <span>Recent Commits</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <ul className="space-y-4">
                    {commits.map((commit) => (
                        <li key={commit.id} className="text-sm">
                            <p className="font-medium">{commit.message}</p>
                            <p className="text-muted-foreground">
                                {commit.author} - {commit.date}
                            </p>
                        </li>
                    ))}
                </ul>
            </CardContent>
        </Card>
    )
}

function CurrentBranches({ branches }: { branches: typeof fakeBranches }) {
    return (
        <Card className="col-span-full">
            <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                    <GitMerge className="h-5 w-5" />
                    <span>Current Branches</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <ul className="space-y-6">
                    {branches.map((branch) => (
                        <li key={branch.name} className="space-y-2">
                            <div className="flex items-center justify-between">
                                <p className="text-lg font-medium">{branch.name}</p>
                                <Badge variant={branch.status === 'open' ? 'default' : 'secondary'}>
                                    {branch.status}
                                </Badge>
                            </div>
                            <div className="flex space-x-4">
                                <div>
                                    <span className="text-sm text-muted-foreground">Security:</span>
                                    <Badge variant={branch.security === 'passed' ? 'default' : branch.security === 'pending' ? 'secondary' : 'destructive'}>
                                        {branch.security}
                                    </Badge>
                                </div>
                                <div>
                                    <span className="text-sm text-muted-foreground">Build:</span>
                                    <Badge variant={branch.build === 'passed' ? 'default' : 'destructive'}>
                                        {branch.build}
                                    </Badge>
                                </div>
                            </div>
                            <div className="text-sm text-muted-foreground">
                                Last commit: {branch.lastCommit} by {branch.lastCommitAuthor}
                            </div>
                        </li>
                    ))}
                </ul>
            </CardContent>
        </Card>
    )
}

function ReadmePreview({ readme }: { readme: string }) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                    <FileText className="h-5 w-5" />
                    <span>README.md</span>
                </CardTitle>
            </CardHeader>
            <CardContent>
        <pre className="bg-muted p-4 rounded-md overflow-auto text-sm whitespace-pre-wrap">
          {readme}
        </pre>
            </CardContent>
        </Card>
    )
}

export { RepoOverview, RecentCommits, CurrentBranches, ReadmePreview }
export { fakeRepoInfo, fakeCommits, fakeBranches, fakeReadme }