import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

const fakeSquadMembers = [
  { id: 1, name: "Alice Johnson", role: "Squad Lead", email: "alice@example.com" },
  { id: 2, name: "Bob Smith", role: "Senior Developer", email: "bob@example.com" },
  { id: 3, name: "Charlie Davis", role: "Backend Developer", email: "charlie@example.com" },
  { id: 4, name: "Diana Miller", role: "Frontend Developer", email: "diana@example.com" },
  { id: 5, name: "Ethan Brown", role: "QA Engineer", email: "ethan@example.com" },
]

export function SquadMembers({ className }: { className?: string }) {
  return (
    <div className={className}>
      <h3 className="text-lg font-semibold mb-4">Squad Members</h3>
      <div className="space-y-4">
        {fakeSquadMembers.map((member) => (
          <div key={member.id} className="flex items-center space-x-4">
            <Avatar>
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
    </div>
  )
}


