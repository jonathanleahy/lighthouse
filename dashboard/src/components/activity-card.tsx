import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart } from 'lucide-react'

export function ActivityCard() {
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


