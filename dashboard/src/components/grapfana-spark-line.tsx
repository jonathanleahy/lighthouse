import React, { useState } from 'react';
import { LineChart, Line, ResponsiveContainer, ReferenceLine, ReferenceArea, XAxis, YAxis, CartesianGrid } from 'recharts';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';

interface ChartData {
    time: number;
    minutesAgo: number;
    errorRate: number;
    canaryErrorRate: number | null;
    deploymentStart: boolean;
    deploymentProgress: number;
    timestamp: string;
}

interface GrafanaSparklineProps {
    deploymentInProgress: boolean;
    data?: Array<ChartData>;
    canaryPercentage?: number;
    targetCanaryPercentage?: number;
}

const GrafanaSparkline = ({
                              data = [],
                              canaryPercentage = 10,
                              targetCanaryPercentage = 100
                          }: GrafanaSparklineProps) => {
    const [isModalOpen, setIsModalOpen] = useState(false);

    const mockData: ChartData[] = React.useMemo(() => {
        const numberOfPoints = 24;
        const deploymentStartIndex = 8;

        return Array.from({ length: numberOfPoints }, (_, i) => {
            const reversedIndex = numberOfPoints - 1 - i;

            const baseErrorRate = (Math.random() * 0.5) + 0.1;
            const spike = Math.random() > 0.8 ? Math.random() * 0.8 : 0;

            const canaryBaseRate = (Math.random() * 0.6) + 0.1;
            const canarySpike = Math.random() > 0.9 ? Math.random() * 1 : 0;

            const hasCanary = reversedIndex >= deploymentStartIndex;

            const deploymentProgress = hasCanary
                ? Math.min(100, Math.round((reversedIndex - deploymentStartIndex + 1) / 4 * canaryPercentage))
                : 0;

            return {
                time: i,
                minutesAgo: i,
                errorRate: baseErrorRate + spike,
                canaryErrorRate: hasCanary ? canaryBaseRate + canarySpike : null,
                deploymentStart: reversedIndex === deploymentStartIndex,
                deploymentProgress: deploymentProgress,
                timestamp: new Date(Date.now() - (i) * 60000).toISOString(),
            };
        }).reverse();
    }, [canaryPercentage]);

    const chartData = data.length > 0 ? data : mockData;
    const deploymentIndex = chartData.findIndex(d => d.deploymentStart);

    const metrics = React.useMemo(() => {
        if (deploymentIndex === -1) return null;

        const preDeploymentData = chartData.slice(0, deploymentIndex);
        const postDeploymentData = chartData.slice(deploymentIndex);

        const preDeploymentAvg = preDeploymentData.length
            ? preDeploymentData.reduce((sum, d) => sum + d.errorRate, 0) / preDeploymentData.length
            : 0;

        const postDeploymentAvg = postDeploymentData.length
            ? postDeploymentData.reduce((sum, d) => sum + d.errorRate, 0) / postDeploymentData.length
            : 0;

        const canaryAvg = postDeploymentData
                .filter(d => d.canaryErrorRate !== null)
                .reduce((sum, d) => sum + (d.canaryErrorRate ?? 0), 0) /
            postDeploymentData.filter(d => d.canaryErrorRate !== null).length || 0;

        const errorRateChange = canaryAvg - preDeploymentAvg;
        const percentageChange = ((canaryAvg - preDeploymentAvg) / preDeploymentAvg) * 100;
        const isSignificantIncrease = percentageChange > 20;

        return {
            preDeploymentAvg,
            postDeploymentAvg,
            canaryAvg,
            errorRateChange,
            percentageChange,
            isSignificantIncrease,
            sampleSize: postDeploymentData.length
        };
    }, [chartData, deploymentIndex]);

    const SparklineView = () => (
        <div className="h-8 w-full">
            <ResponsiveContainer width="100%" height="100%">
                <LineChart
                    data={chartData}
                    margin={{ top: 2, right: 0, left: 0, bottom: 2 }}
                >
                    {metrics && (
                        <>
                            <ReferenceArea
                                x1={0}
                                x2={deploymentIndex}
                                fill="#f3f4f6"
                                fillOpacity={0.5}
                            />
                            <ReferenceArea
                                x1={deploymentIndex}
                                x2={chartData.length - 1}
                                fill={metrics.isSignificantIncrease ? "#fee2e2" : "#f0fdf4"}
                                fillOpacity={0.3}
                            />
                        </>
                    )}

                    {deploymentIndex !== -1 && (
                        <ReferenceLine
                            x={deploymentIndex}
                            stroke="#6b7280"
                            strokeDasharray="3 3"
                        />
                    )}

                    <Line
                        type="monotone"
                        dataKey="errorRate"
                        stroke="#6b7280"
                        strokeWidth={1.5}
                        dot={false}
                        isAnimationActive={false}
                    />

                    <Line
                        type="monotone"
                        dataKey="canaryErrorRate"
                        stroke={metrics?.isSignificantIncrease ? "#ef4444" : "#22c55e"}
                        strokeWidth={1.5}
                        dot={false}
                        isAnimationActive={false}
                    />
                </LineChart>
            </ResponsiveContainer>
        </div>
    );

    const DetailedView = () => (
        <div className="space-y-6">
            <div className="h-64 w-full">
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart
                        data={chartData}
                        margin={{ top: 20, right: 40, left: 20, bottom: 20 }}
                    >
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                            dataKey="minutesAgo"
                            label={{ value: 'Minutes Ago', position: 'bottom' }}
                            tickFormatter={(value) => value}
                        />
                        <YAxis
                            yAxisId="error"
                            label={{ value: 'Error Rate (%)', angle: -90, position: 'left' }}
                        />
                        <YAxis
                            yAxisId="deployment"
                            orientation="right"
                            domain={[0, 100]}
                            label={{ value: 'Deployment %', angle: 90, position: 'right' }}
                        />

                        {metrics && (
                            <>
                                <ReferenceArea
                                    x1={0}
                                    x2={deploymentIndex}
                                    fill="#f3f4f6"
                                    fillOpacity={0.5}
                                    yAxisId="error"
                                />
                                <ReferenceArea
                                    x1={deploymentIndex}
                                    x2={chartData.length - 1}
                                    fill={metrics.isSignificantIncrease ? "#fee2e2" : "#f0fdf4"}
                                    fillOpacity={0.3}
                                    yAxisId="error"
                                />
                                <ReferenceLine
                                    y={metrics.preDeploymentAvg}
                                    yAxisId="error"
                                    stroke="#6b7280"
                                    strokeDasharray="3 3"
                                    label={{
                                        value: "Pre-deploy avg",
                                        position: "insideLeft"
                                    }}
                                />
                            </>
                        )}

                        {deploymentIndex !== -1 && (
                            <ReferenceLine
                                x={deploymentIndex}
                                stroke="#6b7280"
                                strokeDasharray="3 3"
                                label={{
                                    value: "Deployment",
                                    position: "top"
                                }}
                                yAxisId="error"
                            />
                        )}

                        <Line
                            yAxisId="error"
                            type="monotone"
                            dataKey="errorRate"
                            name="Production"
                            stroke="#6b7280"
                            strokeWidth={2}
                            dot={true}
                            isAnimationActive={false}
                        />

                        <Line
                            yAxisId="error"
                            type="monotone"
                            dataKey="canaryErrorRate"
                            name="Canary"
                            stroke={metrics?.isSignificantIncrease ? "#ef4444" : "#22c55e"}
                            strokeWidth={2}
                            dot={true}
                            isAnimationActive={false}
                        />

                        <Line
                            yAxisId="deployment"
                            type="stepAfter"
                            dataKey="deploymentProgress"
                            name="Deployment %"
                            stroke="#3b82f6"
                            strokeWidth={2}
                            strokeDasharray="3 3"
                            dot={false}
                            isAnimationActive={false}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>

            <div className="grid grid-cols-3 gap-6 p-4 bg-gray-50 rounded-lg">
                <div className="space-y-2">
                    <h3 className="font-medium text-gray-600">Pre-deployment</h3>
                    <div className="space-y-1">
                        <div className="text-sm">
                            Error Rate: {metrics?.preDeploymentAvg.toFixed(2)}%
                        </div>
                    </div>
                </div>
                <div className="space-y-2">
                    <h3 className="font-medium text-gray-600">Production</h3>
                    <div className="space-y-1">
                        <div className="text-sm">
                            Error Rate: {metrics?.postDeploymentAvg.toFixed(2)}%
                        </div>
                        <div className="text-sm">
                            Traffic: {100 - canaryPercentage}%
                        </div>
                    </div>
                </div>
                <div className="space-y-2">
                    <h3 className="font-medium text-gray-600">Canary</h3>
                    <div className="space-y-1">
                        <div className="text-sm">
                            Error Rate: {metrics?.canaryAvg.toFixed(2)}%
                        </div>
                        <div className="text-sm">
                            Traffic: {canaryPercentage}%
                        </div>
                        <div className="text-sm">
                            Target: {targetCanaryPercentage}%
                        </div>
                        <div className="text-sm">
                            Change: {metrics?.percentageChange.toFixed(1)}%
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    if (!metrics) return null;

    return (
        <>
            <div
                className="cursor-pointer"
                onClick={() => setIsModalOpen(true)}
            >
                <SparklineView />
            </div>

            <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
                <DialogContent className="max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>Canary Deployment Analysis</DialogTitle>
                    </DialogHeader>
                    <DetailedView />
                </DialogContent>
            </Dialog>
        </>
    );
};

export default GrafanaSparkline;