"use client"

import Link from 'next/link'
import { ArrowLeft, BarChart, Users, AlertCircle } from 'lucide-react'
import { ThemeToggle } from '@/components/theme-toggle'
import { Deployments } from '@/components/deployments'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useDashboardData } from "@/hooks/useDashboardData";
import React, { use } from "react";
import { useState, useEffect } from 'react';
import { Globe, RotateCw, Search } from 'lucide-react';

function ServiceOverview({ service }: { service: { id: string, name: string, description: string, squad: string } }) {
  return (
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl font-bold">{service.name}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground mb-4">{service.description}</p>
          <div className="grid grid-cols-2 gap-4">
            <div className="flex items-center space-x-2">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span className="font-medium">Squad: {service.squad}</span>
            </div>
          </div>
        </CardContent>
      </Card>
  )
}

// Updated DashboardData interface to include all possible fields
interface DashboardData {
  repoDesc?: string;
  repoSquad?: string;
  repoBitUrl?: string;
  repoCodefresh?: string;
  argocd?: { url?: string };
  apps?: {
    appName: string;
    type: string;
    deployment?: {
      deployments: Array<{
        type: string;
        percentage: number;
        version: string;
      }>;
    };
    argocd?: {
      status: {
        weight: string;
        step: string[];
      };
      url: string;
    };
    grafana?: {
      url: string;
    };
  }[];
}

// Update the DeploymentData interface to match the one in the Deployments component
interface DeploymentData {
  repoName: string;
  repoBitUrl: string;
  apps: Array<{
    appName: string;
    type: string;
    deployment?: {
      deployments: Array<{
        type: string;
        percentage: number;
        version: string;
      }>;
    };
    argocd?: {
      status: {
        weight: string;
        step: string[];
      };
      url: string;
    };
    grafana?: {
      url: string;
    };
  }>;
}

export default function MicroserviceDetail({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const selectedOption = resolvedParams.id;

  const { data: mockData, loading, error } = useDashboardData(
      selectedOption ? `http://localhost:8083/?repo=${selectedOption}` : null,
      true
  );

  const LoadingIcons = ({ className }: { className?: string }) => {
    const [iconIndex, setIconIndex] = useState(0);
    const icons = [
      <Search key="search" className={`h-5 w-5 text-gray-800 ${className}`} />,
      <Globe key="globe" className={`h-5 w-5 text-gray-800 ${className}`} />,
      <RotateCw key="rotate" className={`h-5 w-5 text-gray-800 ${className}`} />
    ];

    useEffect(() => {
      const interval = setInterval(() => {
        setIconIndex((prevIndex) => (prevIndex + 1) % icons.length);
      }, 2000);

      return () => clearInterval(interval);
    }, [icons.length]);

    return icons[iconIndex];
  };

  const transformToDeploymentData = (data: DashboardData | null): DeploymentData | null => {
    if (!data) return null;

    return {
      repoName: selectedOption || '',
      repoBitUrl: data.repoBitUrl || '',
      apps: data.apps || []
    };
  };

  const transformedData = transformToDeploymentData(mockData);

  const service = {
    id: 12,
    name: selectedOption,
    description: (mockData as DashboardData)?.repoDesc || '',
    status: 'Warning',
    version: 'v1.9.2',
    requests: '500K',
    uptime: '99.5%',
    errors: 45,
    infos: 89,
    squad: (mockData as DashboardData)?.repoSquad || '',
    lastMainUpdate: '2023-11-18',
    lastBranchUpdate: '2023-11-21',
    lastDeploy: '2023-11-22'
  }

  return (
      <div className="min-h-screen bg-background">
        <header className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8 flex h-16 items-center justify-between">
            <div className="flex items-center space-x-4">
              <Link href="/" className="flex items-center space-x-2">
                <ArrowLeft className="h-5 w-5" />
                <span className="font-medium">Dashboard</span>
              </Link>
              <span className="text-muted-foreground">/</span>
              <h1 className="text-xl font-bold">{service.name}</h1>
            </div>
            <div className="flex items-center space-x-4">
              <ThemeToggle />
            </div>
          </div>
        </header>

        {loading ? (
            <div className="fixed inset-0 flex items-center justify-center bg-gray-500 bg-opacity-75 z-50">
              <div className="bg-white p-4 rounded shadow-lg">
                <div className="flex items-center space-x-2">
                  <LoadingIcons className="animate-spin"/>
                  <span>Loading...</span>
                </div>
              </div>
            </div>
        ) : null}

        {error ? (
            <div className="fixed inset-0 flex items-center justify-center z-10 pointer-events-none">
              <div className="bg-white p-4 rounded shadow-lg pointer-events-auto">
                <div className="flex items-center space-x-2">
                  <AlertCircle className="h-5 w-5 text-gray-800"/>
                  <span>No Data</span>
                </div>
              </div>
            </div>
        ) : (
            <main className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
              <div className="grid gap-6 md:grid-cols-3 mb-8">
                <div className="md:col-span-2">
                  <ServiceOverview key={service.id} service={{...service, id: service.id.toString()}}/>
                </div>
                <Card key="quick-links">
                  <CardHeader>
                    <CardTitle className="font-semibold">Quick Links</CardTitle>
                  </CardHeader>
                  <CardContent className="grid grid-cols-2 gap-4">
                    <Button variant="outline" className="w-full" asChild>
                      <Link href={(mockData as DashboardData)?.repoBitUrl || "#"} target="_blank" rel="noopener noreferrer">
                        {(mockData as DashboardData)?.repoBitUrl ? "View Github" : "No Github URL"}
                      </Link>
                    </Button>
                    <Button variant="outline" className="w-full" asChild>
                      <Link href={(mockData as DashboardData)?.repoCodefresh || "#"} target="_blank" rel="noopener noreferrer">
                        {(mockData as DashboardData)?.repoCodefresh ? "View Codefresh" : "No Codefresh URL"}
                      </Link>
                    </Button>
                    <Button variant="outline" className="w-full" asChild>
                      <Link href={(mockData as DashboardData)?.argocd?.url || "#"} target="_blank" rel="noopener noreferrer">
                        {(mockData as DashboardData)?.argocd?.url ? "View ArgoCD" : "No ArgoCD URL"}
                      </Link>
                    </Button>
                    <Button variant="outline" className="w-full" asChild>
                      <a href="#" target="_blank" rel="noopener noreferrer">
                        <BarChart className="mr-2 h-4 w-4" />
                        Grafana :-(
                      </a>
                    </Button>
                  </CardContent>
                </Card>
              </div>
              <Deployments isLoading={false} mockData={transformedData} />
            </main>
        )}
      </div>
  )
}