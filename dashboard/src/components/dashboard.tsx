"use client"

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Search, Activity, Clock, AlertTriangle, CheckCircle, XCircle, Info } from 'lucide-react'
import { ThemeToggle } from '@/components/theme-toggle'
import Link from 'next/link'

interface Microservice {
  id: string;
  name: string;
  status: string;
  version: string;
  requests: string;
  uptime: string;
  errors: number;
  infos: number;
}

const microservices: Microservice[] = [
  // { id: 'auth', name: 'User Authentication', status: 'Healthy', version: 'v2.3.1', requests: '1.2M', uptime: '99.9%', errors: 23, infos: 156 },
  // { id: 'payment', name: 'Payment Processing', status: 'Warning', version: 'v1.9.2', requests: '500K', uptime: '99.5%', errors: 45, infos: 89 },
  // { id: 'inventory', name: 'Inventory Management', status: 'Healthy', version: 'v3.0.0', requests: '800K', uptime: '99.8%', errors: 12, infos: 201 },
]

export default function Dashboard() {
  const [filter, setFilter] = useState('')

  const filteredMicroservices = microservices.filter(service =>
      service.name.toLowerCase().includes(filter.toLowerCase())
  )

  return (
      <div
          className="min-h-screen pb-16 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-900 dark:to-gray-800 text-gray-900 dark:text-white flex flex-col">
        <header className="bg-white dark:bg-gray-800 shadow-lg">
          <div className="container mx-auto px-4 py-6">
            <div className="flex justify-between items-center mb-4">
              <h1 className="text-3xl font-bold">Microservices Dashboard</h1>
              <ThemeToggle/>
            </div>
            <div className="flex items-center space-x-2">
              <Search className="text-gray-400"/>
              <Input
                  type="text"
                  placeholder="Filter microservices..."
                  className="flex-grow bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 border-gray-300 dark:border-gray-600"
                  value={filter}
                  onChange={(e) => setFilter(e.target.value)}
              />
            </div>
          </div>
        </header>
        <main className="container mx-auto px-4 py-8 pb-16 flex-grow">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            {filteredMicroservices.map((service) => (
                <Link href={`/microservice/${service.id}`} key={service.id}>
                  <Card
                      className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow duration-200">
                    <CardHeader>
                      <div className="flex justify-between items-center">
                        <CardTitle className="text-xl">{service.name}</CardTitle>
                        <StatusBadge status={service.status}/>
                      </div>
                      <CardDescription className="dark:text-gray-400">Version: {service.version}</CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="grid grid-cols-2 gap-4 mb-4">
                        <div className="flex items-center space-x-2">
                          <Activity className="text-blue-600 dark:text-blue-400"/>
                          <div>
                            <p className="text-sm text-gray-500 dark:text-gray-400">Requests (24h)</p>
                            <p className="text-xl font-bold">{service.requests}</p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Clock className="text-green-600 dark:text-green-400"/>
                          <div>
                            <p className="text-sm text-gray-500 dark:text-gray-400">Uptime</p>
                            <p className="text-xl font-bold">{service.uptime}</p>
                          </div>
                        </div>
                      </div>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="flex items-center space-x-2">
                          <AlertTriangle className="text-red-600 dark:text-red-400"/>
                          <div>
                            <p className="text-sm text-gray-500 dark:text-gray-400">Errors</p>
                            <p className="text-xl font-bold">{service.errors}</p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Info className="text-blue-600 dark:text-blue-400"/>
                          <div>
                            <p className="text-sm text-gray-500 dark:text-gray-400">Infos</p>
                            <p className="text-xl font-bold">{service.infos}</p>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </Link>
            ))}
          </div>
        </main>
        <footer className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 py-4">
          <div className="container mx-auto px-4">
            <p className="text-center text-gray-500 dark:text-gray-400">Microservices Dashboard © 2024</p>
          </div>
        </footer>
      </div>
  )
}

function StatusBadge({status}: { status: string }) {
  let color = "bg-gray-500 dark:bg-gray-400"
  let icon = null

  switch (status) {
    case 'Healthy':
      color = "bg-green-500 dark:bg-green-400"
      icon = <CheckCircle className="h-4 w-4" />
      break
    case 'Warning':
      color = "bg-yellow-500 dark:bg-yellow-400"
      icon = <AlertTriangle className="h-4 w-4" />
      break
    case 'Critical':
      color = "bg-red-500 dark:bg-red-400"
      icon = <XCircle className="h-4 w-4" />
      break
  }

  return (
      <Badge className={`${color} text-white dark:text-gray-900 flex items-center space-x-1 px-2 py-1`}>
        {icon}
        <span>{status}</span>
      </Badge>
  )
}