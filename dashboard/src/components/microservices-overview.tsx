import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import Link from 'next/link'
import { useEffect } from "react";
import { useCustomFields } from '@/lib/customFieldsContext';

const microservices = [
  { id: 'auth', name: 'User Authentication', status: 'Healthy', version: 'v2.3.1' },
  { id: 'payment', name: 'Payment Processing', status: 'Warning', version: 'v1.9.2' },
  { id: 'inventory', name: 'Inventory Management', status: 'Healthy', version: 'v3.0.0' },
  { id: 'notification', name: 'Notification Service', status: 'Critical', version: 'v1.5.4' },
]

export function MicroservicesOverview() {
  const { dispatch } = useCustomFields();

  useEffect(() => {
    const savedCustomFieldSets = localStorage.getItem('customFieldSets');
    const savedActiveSetId = localStorage.getItem('activeSetId');
    console.log('Loading data from localStorage');
    if (savedCustomFieldSets) {
      console.log('Custom Field Sets:', JSON.parse(savedCustomFieldSets));
      dispatch({ type: 'SET_CUSTOM_FIELD_SETS', payload: JSON.parse(savedCustomFieldSets) });
    }
    if (savedActiveSetId) {
      console.log('Active Set ID:', savedActiveSetId);
      dispatch({ type: 'SET_ACTIVE_SET_ID', payload: savedActiveSetId });
    }
  }, [dispatch]);

  return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {microservices.map((service) => (
            <Card key={service.id}>
              <CardHeader>
                <CardTitle>{service.name}</CardTitle>
                <CardDescription>Microservice</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex justify-between items-center">
                  <Badge
                      variant={service.status === 'Healthy' ? 'default' : service.status === 'Warning' ? 'outline' : 'destructive'}
                  >
                    {service.status}
                  </Badge>
                  <span className="text-sm text-muted-foreground">{service.version}</span>
                </div>
                <Link href={`/microservice/${service.id}`} className="mt-4 inline-block text-blue-500 hover:underline">
                  View Details
                </Link>
              </CardContent>
            </Card>
        ))}
      </div>
  )
}