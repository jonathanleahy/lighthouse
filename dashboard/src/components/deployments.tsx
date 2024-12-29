"use client"

import { useState } from 'react';
import * as React from 'react';
import { ArrowUpDown, MoreVertical, RotateCw, PlayCircle, PauseCircle, CheckCircle2, XCircle } from 'lucide-react';
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHeader, TableHead, TableRow } from "@/components/ui/table";
import Image from 'next/image';
import ArgoIcon from 'public/icons/argo.svg';

// Type definitions
interface Deployment {
  type: string;
  percentage: number;
  version: string;
}

interface ArgocdStatus {
  weight: string;
  step: string[];
}

interface App {
  appName: string;
  type: string;
  deployment?: {
    deployments: Deployment[];
  };
  argocd?: {
    status: ArgocdStatus;
    url: string;
  };
  grafana?: {
    url: string;
  };
}

interface DeploymentData {
  repoName: string;
  repoBitUrl: string;
  repoNamespace: string;
  apps: App[];
}

const mockDataExample: DeploymentData = {
  repoName: 'example-repo',
  repoBitUrl: 'https://example.com/repo',
  repoNamespace: 'psm-namespace',
  apps: [
    {
      appName: 'example-repo-auth-service',
      type: 'primary',
      deployment: {
        deployments: [
          { type: 'stable', percentage: 80, version: 'v1.2.3' },
          { type: 'canary', percentage: 20, version: 'v1.2.4' }
        ]
      },
      argocd: {
        status: { weight: '20', step: ['{}'] },
        url: 'https://argocd.example.com/applications/auth-service'
      },
      grafana: { url: 'https://grafana.example.com/d/auth-service' }
    },
    {
      appName: 'example-repo-payment-service',
      type: 'secondary',
      deployment: {
        deployments: [
          { type: 'stable', percentage: 100, version: 'v2.0.1' }
        ]
      },
      argocd: {
        status: { weight: '0', step: ['{}'] },
        url: 'https://argocd.example.com/applications/payment-service'
      },
      grafana: { url: 'https://grafana.example.com/d/payment-service' }
    }
  ]
};

interface StatusIconProps {
  status: string;
  stableTab: string;
  argocdStep: string;
}

async function argocdResume(namespace: string, resourceName: string, region: string) {
  const url = `http://localhost:8083/argocd-resume?namespace=${namespace}&resourceName=${resourceName}&region=${region}`;

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ action: "resume" })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.error(`Error calling ${url}`, error);
  }
}

const StatusIcon: React.FC<StatusIconProps> = ({ status, argocdStep }) => {
  let parsedStep: { pause?: { duration?: string } } = {};
  let pausedFor: string | number = 0;

  try {
    parsedStep = JSON.parse(argocdStep);
    if (parsedStep.pause?.duration) {
      pausedFor = parsedStep.pause.duration;
    } else if (parsedStep.pause) {
      pausedFor = -1;
    }
  } catch {
    pausedFor = 0;
  }

  switch (status) {
    case 'In Progress': {
      let message = '';
      if (typeof pausedFor === 'string') {
        message = `Paused for ${pausedFor}`;
      } else if (pausedFor === -1) {
        message = 'Paused: Manual start required';
      } else if (pausedFor === 0) {
        message = 'ArgoCD is syncing';
      }
      return (
          <div className="flex items-center gap-2">
            {typeof pausedFor === 'string' ? (
                <PlayCircle className="text-blue-500 h-4 w-4" />
            ) : (
                <PlayCircle className={`text-blue-500 h-4 w-4 ${pausedFor === -1 ? 'animate-spin-back-and-forth' : 'animate-spin'}`} />
            )}
            <div>
              <div className="font-medium">In Progress</div>
              {message && <div className="text-xs text-muted-foreground">{message}</div>}
            </div>
          </div>
      );
    }
    case 'Up to date':
      return (
          <div className="flex items-center gap-2">
            <CheckCircle2 className="text-green-500 h-4 w-4" />
            <span className="font-medium">{status}</span>
          </div>
      );
    case 'Version Mismatch':
      return (
          <div className="flex items-center gap-2">
            <CheckCircle2 className="text-yellow-500 h-4 w-4" />
            <span className="font-medium">{status}</span>
          </div>
      );
    case 'No Deployment':
      return (
          <div className="flex items-center gap-2">
            <XCircle className="text-orange-500 h-4 w-4" />
            <span className="font-medium">{status}</span>
          </div>
      );
    default: {
      let message = 'Paused';
      if (typeof pausedFor === 'string') {
        message = `Paused for ${pausedFor}`;
      } else if (pausedFor === -1) {
        message = 'Waiting to kick off';
      }
      return (
          <div className="flex items-center gap-2">
            <PauseCircle className="text-blue-500 h-4 w-4" />
            <span className="font-medium">{message}</span>
          </div>
      );
    }
  }
};

interface DeploymentProgressProps {
  deployments: Deployment[];
  argocdWeight: number;
}

const DeploymentProgress: React.FC<DeploymentProgressProps> = ({ deployments, argocdWeight }) => {
  const stable = deployments.find(d => d.type === "stable");
  const canary = deployments.find(d => d.type === "canary");

  if (!stable || !canary) return null;

  const totalPercentage = stable.percentage + canary.percentage;
  if (totalPercentage !== 100) return null;

  const stablePercentage = isNaN(argocdWeight) ? 0 : 100 - argocdWeight;
  const canaryPercentage = isNaN(argocdWeight) ? 0 : argocdWeight;

  return (
      <div className="space-y-1 w-full">
        <div className="h-3 w-full bg-muted overflow-hidden rounded-full relative">
          <div
              className="h-full bg-gray-500 absolute left-0 top-0 transition-all duration-500"
              style={{ width: `${stablePercentage}%` }}
          />
          <div
              className="h-full bg-green-500 absolute left-0 top-0 transition-all duration-500"
              style={{ width: `${canaryPercentage}%`, marginLeft: `${stablePercentage}%` }}
          />
        </div>
        <div className="flex justify-between text-xs text-muted-foreground">
          <div>
            <span>stable {stablePercentage}%</span>
          </div>
          <div>
            <span>canary {canaryPercentage}%</span>
          </div>
        </div>
      </div>
  );
};

const getDeploymentStatus = (deployments: Deployment[], stableTag: string) => {
  if (deployments.length === 0) return { status: "No Deployment", color: "text-yellow-500" };

  if (deployments.length === 1 && deployments[0].percentage === 100) {
    if (deployments[0].version === stableTag) {
      return {status: "Up to date", color: "text-green-500"};
    } else {
      return {status: "Version Mismatch", color: "text-orange-500"};
    }
  }

  if (deployments.length > 1 || deployments[0].percentage < 100) {
    return { status: "In Progress", color: "text-blue-500" };
  }
  return { status: "Unknown", color: "text-gray-500" };
};

// After
interface DeploymentsProps {
  isLoading: boolean;
  mockData: DeploymentData | null;
  service: {
    id: string;
    name: string;
    description: string;
    status: string;
    version: string;
    requests: string;
    uptime: string;
    errors: number;
    infos: number;
    squad: string;
    lastMainUpdate: string;
    lastBranchUpdate: string;
    lastDeploy: string;
    stableTag: string;
  };
}

export const Deployments: React.FC<DeploymentsProps> = ({ isLoading, mockData, service }) => {
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("asc");
  const [sortColumn, setSortColumn] = useState("appName");
  const [searchTerm, setSearchTerm] = useState("");

  if (isLoading) {
    return (
        <div className="flex justify-center p-8">
          <RotateCw className="h-6 w-6 animate-spin text-gray-400" />
        </div>
    );
  }

  if (!mockData) return null;

  const data = mockData ?? mockDataExample;

  const filteredApps = (data.apps || []).filter(app =>
      app.appName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getPriority = (appName: string): number => {
    const lowerName = appName.toLowerCase();
    if (lowerName.startsWith('dev-')) return -2;
    if (lowerName.startsWith('integration-')) return -1;
    return 0;
  };

  const sortedApps = [...filteredApps].sort((a, b) => {
    const priorityA = getPriority(a.appName);
    const priorityB = getPriority(b.appName);

    if (priorityA !== priorityB) {
      return priorityA - priorityB;
    }

    const valueA = a[sortColumn as keyof App] ?? '';
    const valueB = b[sortColumn as keyof App] ?? '';

    if (typeof valueA === 'string' && typeof valueB === 'string') {
      return sortDirection === "asc"
          ? valueA.localeCompare(valueB)
          : valueB.localeCompare(valueA);
    }

    if (valueA < valueB) return sortDirection === "asc" ? -1 : 1;
    if (valueA > valueB) return sortDirection === "asc" ? 1 : -1;
    return 0;
  });

  const handleSort = (column: string) => {
    if (column === sortColumn) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortColumn(column);
      setSortDirection("asc");
    }
  };

  return (
      <div className="space-y-4">
        <div className="flex justify-between items-center">
          <h2 className="text-lg font-semibold">Github Terraform Deployments</h2>
          <div className="flex items-center gap-2">
            <Input
                placeholder="Filter deployments..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-64"
            />
          </div>
        </div>

        <div className="border rounded-md">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="pl-4 w-[250px]">
                  <Button
                      variant="ghost"
                      onClick={() => handleSort("appName")}
                      className="h-8 -ml-4 font-medium"
                  >
                    Region
                    <ArrowUpDown className="ml-2 h-4 w-4"/>
                  </Button>
                </TableHead>
                <TableHead className="">Version</TableHead>
                <TableHead className="">Progress</TableHead>
                <TableHead className="">Argo Status</TableHead>
                <TableHead className="w-8"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sortedApps.map((app) => {
                const deploymentStatus = app.deployment?.deployments ?
                    getDeploymentStatus(app.deployment.deployments, service.stableTag) :
                    {status: "No Deployment", color: "text-yellow-500"};
                const argocdWeight = app.argocd?.status ?
                    parseFloat(app.argocd.status.weight) : 0;

                const stableDeployment = app.deployment?.deployments?.find(d => d.type === 'stable');
                const canaryDeployment = app.deployment?.deployments?.find(d => d.type === 'canary');

                return (
                    <TableRow key={app.appName}>
                      <TableCell className="pl-4 pr-4">
                    <span className="font-medium">
                      {app.appName.replace(data?.repoName ? `${data.repoName}-` : '', '')}
                    </span>
                      </TableCell>
                      <TableCell className="pr-8">
                        <div className="space-y-1.5">
                          {stableDeployment && (
                              <div className="flex items-center gap-2">
                                <Badge variant="outline" className="w-16 justify-center">
                                  stable
                                </Badge>
                                <code className="text-xs bg-muted text-muted-foreground px-1.5 py-0.5 rounded">
                                  {stableDeployment.version}
                                </code>
                              </div>
                          )}
                          {canaryDeployment && (
                              <div className="flex items-center gap-2">
                                <Badge variant="outline" className="w-16 justify-center">
                                  canary
                                </Badge>
                                <code className="text-xs bg-muted text-muted-foreground px-1.5 py-0.5 rounded">
                                  {canaryDeployment.version}
                                </code>
                              </div>
                          )}
                        </div>
                      </TableCell>
                      <TableCell className="pr-8">
                        <DeploymentProgress
                            deployments={app.deployment?.deployments || []}
                            argocdWeight={argocdWeight}
                        />
                      </TableCell>
                      <TableCell className="pr-8">
                        <div className="flex items-center gap-2">
                          {deploymentStatus.status === 'In Progress' ? (
                              <>
                                <Button
                                    variant="default"
                                    className="m-1 bg-white hover:bg-gray-100"
                                    onClick={() => argocdResume(mockData.repoNamespace, service.name, app.appName.replace(data?.repoName ? `${data.repoName}-` : '', ''))}
                                >
                                  <Image src={ArgoIcon} alt="ArgoCD" width={16} height={16} className="text-black"/>
                                  <StatusIcon
                                      status={deploymentStatus.status}
                                      argocdStep={app.argocd?.status?.step?.[0] || '{}'}
                                  />
                                </Button>
                              </>
                          ) : (
                              <StatusIcon
                                  status={deploymentStatus.status}
                                  argocdStep={app.argocd?.status?.step?.[0] || '{}'}
                              />

                          )                         }
                        </div>
                      </TableCell>
                      <TableCell>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-8 w-8 p-0">
                              <MoreVertical className="h-4 w-4"/>
                              <span className="sr-only">Open menu</span>
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            {app.argocd?.url && (
                                <DropdownMenuItem>
                                  <a href={app.argocd.url} target="_blank" rel="noopener noreferrer">
                                    View in ArgoCD
                                  </a>
                                </DropdownMenuItem>
                            )}
                            {app.grafana?.url && (
                                <DropdownMenuItem>
                                  <a href={app.grafana.url} target="_blank" rel="noopener noreferrer">
                                    View in Grafana
                                  </a>
                                </DropdownMenuItem>
                            )}
                            <DropdownMenuItem>
                              View in CodeFresh
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      </div>
  );
};

export default Deployments;