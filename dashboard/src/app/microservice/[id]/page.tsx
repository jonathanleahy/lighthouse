"use client"

import Link from 'next/link'
import {ArrowLeft, BarChart, AlertCircle, RefreshCw} from 'lucide-react'
import { ThemeToggle } from '@/components/theme-toggle'
import { Deployments } from '@/components/deployments'
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { useDashboardData } from "@/hooks/useDashboardData";
import React, { use } from "react";
import { useState, useEffect } from 'react';
import { Globe, RotateCw, Search } from 'lucide-react';
import Image from "next/image";
import GithubIcon from 'public/icons/github.svg';
import CodeFreshIcon from 'public/icons/codefresh.svg';
import ArgoIcon from 'public/icons/argo.svg';
import Grafana from 'public/icons/grafana.svg';

function ServiceOverview({ service, refreshData }: { service: { id: string, name: string, description: string, squad: string, stableTag: string }, refreshData: () => void }) {
  return (
      <Card className="relative">
        <CardContent className="pt-6">
          <div className="absolute top-2 right-2">
            <Button variant="ghost" size="icon" onClick={refreshData} aria-label="Refresh service status">
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
          <div className="space-y-2">
            <h2 className="text-2xl font-bold">{service.name}</h2>
            <p className="text-sm text-muted-foreground">{service.description}</p>
            <div className="flex items-center space-x-2">
              <span className="font-semibold">Squad:</span>
              <span>{service.squad}</span>
            </div>
            {service.stableTag && (
                <div className="flex items-center space-x-2">
                  <span className="font-semibold">Stable Tag:</span>
                  <span>{service.stableTag}</span>
                </div>
            )}
          </div>
        </CardContent>
      </Card>
  );
}

interface DashboardData {
  repoDesc?: string;
  repoSquad?: string;
  repoBitUrl?: string;
  repoNamespace?: string;
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

interface DeploymentData {
  repoName: string;
  repoBitUrl: string;
  repoNamespace: string;
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
  stableTag: string;
}

export default function MicroserviceDetail({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const selectedOption = resolvedParams.id;

  const { data: mockData, loading, error, refetch } = useDashboardData(
      selectedOption ? `http://localhost:8083/?repo=${selectedOption}` : null,
      true
  );

  const refreshData = (force = false) => {
    refetch(force);
  };

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

    const stableTag = data.tags?.find(tag => tag.status === 'stable')?.tag || '';

    return {
      repoName: selectedOption || '',
      repoBitUrl: data.repoBitUrl || '',
      repoNamespace: data.repoNamespace || '',
      apps: data.apps || [],
      stableTag: stableTag
    };
  };

  // const transformedData = transformToDeploymentData(mockData as DashboardData);

  const dashboardData = mockData as DashboardData;

  const transformedData = transformToDeploymentData(mockData as DashboardData);

  const service = {
    id: '12',
    name: selectedOption,
    description: dashboardData?.repoDesc || '',
    status: 'Warning',
    version: 'v1.9.2',
    requests: '500K',
    uptime: '99.5%',
    errors: 45,
    infos: 89,
    squad: dashboardData?.repoSquad || '',
    lastMainUpdate: '2023-11-18',
    lastBranchUpdate: '2023-11-21',
    lastDeploy: '2023-11-22',
    stableTag: transformedData?.stableTag || ''
  };

  return (
      <div className="min-h-screen bg-background">
        <header className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="container mx-auto px-4 sm:px-4 lg:px-4 flex h-12 items-center justify-between">
            <div className="flex items-center space-x-4">
              <Link href="/" className="flex items-center space-x-2">
                <ArrowLeft className="h-5 w-5"/>
                <span className="font-medium">Dashboard</span>
              </Link>
              <span className="text-muted-foreground">/</span>
              <h1 className="text-xl font-bold">{service.name}</h1>
            </div>
            <div className="flex items-center space-x-4">
              <ThemeToggle/>
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
                  <ServiceOverview key={service.id} service={{...service, id: service.id.toString()}} refreshData={refreshData} />
                </div>
                <Card key="quick-links">
                  <CardContent className="grid grid-cols-2 gap-4 p-4">
                    <Button variant="outline" className="w-full m-1" asChild>
                      <div>
                        <Image src={GithubIcon} alt="Github" width={16} height={16} className="text-black"/>
                        <Link href={dashboardData?.repoBitUrl || "#"} target="_blank" rel="noopener noreferrer">
                          {dashboardData?.repoBitUrl ? "View Github" : "No Github URL"}
                        </Link>
                      </div>
                    </Button>
                    <Button variant="outline" className="w-full m-1" asChild>
                      <div>
                        <Image src={CodeFreshIcon} alt="Codefresh" width={16} height={16} className="text-black"/>
                        <Link href={dashboardData?.repoCodefresh || "#"} target="_blank" rel="noopener noreferrer">
                          {dashboardData?.repoCodefresh ? "View Codefresh" : "No Codefresh URL"}
                        </Link>
                      </div>
                    </Button>
                    <Button variant="outline" className="w-full m-1" asChild>
                      <div>
                        <Image src={ArgoIcon} alt="Github" width={16} height={16} className="text-black"/>
                        <Link href={dashboardData?.argocd?.url || "#"} target="_blank" rel="noopener noreferrer">
                          {dashboardData?.argocd?.url ? "View ArgoCD" : "No ArgoCD URL"}
                        </Link>
                      </div>
                    </Button>
                    <Button variant="outline" className="w-full m-1" asChild>
                      <div>
                        <Image src={Grafana} alt="Grafana" width={16} height={16} className="text-black"/>
                        <a href="#" target="_blank" rel="noopener noreferrer">
                          Grafana :-(
                        </a>
                      </div>
                    </Button>
                  </CardContent>
                </Card>
              </div>
              <Deployments isLoading={false} mockData={transformedData} service={service} />
            </main>
        )}
      </div>
  )
}