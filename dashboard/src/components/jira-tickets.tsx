import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

type JiraTicket = {
  id: string
  title: string
  status: 'To Do' | 'In Progress' | 'In Review' | 'Done'
  priority: 'Low' | 'Medium' | 'High'
  assignee: {
    name: string
    avatar: string
  }
  dueDate: string
}

const mockJiraTickets: JiraTicket[] = [
  {
    id: 'PROJ-1',
    title: 'Implement OAuth 2.0 flow',
    status: 'In Progress',
    priority: 'High',
    assignee: {
      name: 'Alice Johnson',
      avatar: 'https://api.dicebear.com/6.x/initials/svg?seed=AJ'
    },
    dueDate: '2023-12-15'
  },
  {
    id: 'PROJ-2',
    title: 'Optimize database queries',
    status: 'To Do',
    priority: 'Medium',
    assignee: {
      name: 'Bob Smith',
      avatar: 'https://api.dicebear.com/6.x/initials/svg?seed=BS'
    },
    dueDate: '2023-12-20'
  },
  {
    id: 'PROJ-3',
    title: 'Update API documentation',
    status: 'In Review',
    priority: 'Low',
    assignee: {
      name: 'Charlie Davis',
      avatar: 'https://api.dicebear.com/6.x/initials/svg?seed=CD'
    },
    dueDate: '2023-12-10'
  },
  {
    id: 'PROJ-4',
    title: 'Fix login page responsiveness',
    status: 'Done',
    priority: 'High',
    assignee: {
      name: 'Diana Miller',
      avatar: 'https://api.dicebear.com/6.x/initials/svg?seed=DM'
    },
    dueDate: '2023-12-05'
  },
  {
    id: 'PROJ-5',
    title: 'Implement rate limiting',
    status: 'In Progress',
    priority: 'Medium',
    assignee: {
      name: 'Ethan Brown',
      avatar: 'https://api.dicebear.com/6.x/initials/svg?seed=EB'
    },
    dueDate: '2023-12-18'
  }
]

const getStatusColor = (status: JiraTicket['status']) => {
  switch (status) {
    case 'To Do':
      return 'bg-gray-500'
    case 'In Progress':
      return 'bg-blue-500'
    case 'In Review':
      return 'bg-yellow-500'
    case 'Done':
      return 'bg-green-500'
    default:
      return 'bg-gray-500'
  }
}

const getPriorityColor = (priority: JiraTicket['priority']) => {
  switch (priority) {
    case 'Low':
      return 'bg-green-500'
    case 'Medium':
      return 'bg-yellow-500'
    case 'High':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

export function JiraTickets() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Jira Tickets</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Title</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Priority</TableHead>
              <TableHead>Assignee</TableHead>
              <TableHead>Due Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {mockJiraTickets.map((ticket) => (
              <TableRow key={ticket.id}>
                <TableCell className="font-medium">{ticket.id}</TableCell>
                <TableCell>{ticket.title}</TableCell>
                <TableCell>
                  <Badge variant="secondary" className={`${getStatusColor(ticket.status)} text-white`}>
                    {ticket.status}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Badge variant="outline" className={`${getPriorityColor(ticket.priority)} text-white`}>
                    {ticket.priority}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div className="flex items-center space-x-2">
                    <Avatar className="h-6 w-6">
                      <AvatarImage src={ticket.assignee.avatar} alt={ticket.assignee.name} />
                      <AvatarFallback>{ticket.assignee.name.split(' ').map(n => n[0]).join('')}</AvatarFallback>
                    </Avatar>
                    <span>{ticket.assignee.name}</span>
                  </div>
                </TableCell>
                <TableCell>{ticket.dueDate}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}


