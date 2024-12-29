// src/components/job/JobProgressView.tsx
import React from 'react';
import { Clock, CheckCircle, AlertCircle, Loader2 } from 'lucide-react';
import { formatDuration, getStepTiming } from './JobUtils';
import type { ActiveJob } from './JobUtils';

const StepStatus: React.FC<{ status: string }> = ({ status }) => {
    switch (status.toLowerCase()) {
        case 'completed':
            return <CheckCircle className="w-4 h-4 text-green-500" />;
        case 'pending':
            return <Clock className="w-4 h-4 text-gray-400" />;
        case 'processing':
            return <Loader2 className="w-4 h-4 text-blue-500 animate-spin" />;
        default:
            return <AlertCircle className="w-4 h-4 text-yellow-500" />;
    }
};

interface JobProgressViewProps {
    job: ActiveJob;
    onViewDetails?: () => void;
}

const JobProgressView: React.FC<JobProgressViewProps> = ({ job, onViewDetails }) => {
    const completionPercent = (job.completed_steps / job.total_steps) * 100;

    return (
        <div className="bg-white rounded-lg shadow-sm border p-4">
            <div className="mb-4">
                <div className="flex justify-between items-center mb-2">
                    <h3 className="text-lg font-semibold">
                        {job.service_name}
                        <span className="text-sm text-gray-500 ml-2">({job.type})</span>
                    </h3>
                    <span className="text-sm text-gray-600">
                        {job.completed_steps} of {job.total_steps} steps completed
                    </span>
                </div>

                <div className="w-full bg-gray-100 rounded-full h-2">
                    <div
                        className="bg-blue-500 h-2 rounded-full transition-all duration-500"
                        style={{ width: `${completionPercent}%` }}
                    />
                </div>
            </div>

            <div className="space-y-3">
                {Object.entries(job.steps).map(([stepName, step]) => {
                    const timing = getStepTiming(step);
                    const duration = step.status === 'completed' && timing.start && timing.end
                        ? formatDuration(timing.start, timing.end)
                        : null;

                    return (
                        <div key={stepName} className="border rounded p-3">
                            <div className="flex items-center justify-between mb-1">
                                <div className="flex items-center gap-2">
                                    <StepStatus status={step.status} />
                                    <span className="font-medium capitalize">{stepName}</span>
                                </div>
                                <span className={`text-sm ${
                                    step.status === 'completed' ? 'text-green-600' :
                                        step.status === 'pending' ? 'text-gray-400' :
                                            'text-blue-600'
                                }`}>
                                    {step.status}
                                </span>
                            </div>

                            <div className="text-sm text-gray-600">
                                {step.status === 'completed' && (
                                    <div className="flex justify-between">
                                        <span>Duration: {duration || 'N/A'}</span>
                                        {timing.end && (
                                            <span>
                                                {new Date(timing.end).toLocaleTimeString()}
                                            </span>
                                        )}
                                    </div>
                                )}
                                {step.result?.message && (
                                    <div className="text-gray-500 mt-1">
                                        {step.result.message}
                                    </div>
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>

            {onViewDetails && (
                <div className="mt-4 flex justify-end">
                    <button
                        onClick={onViewDetails}
                        className="text-blue-600 hover:text-blue-800 text-sm"
                    >
                        View Details
                    </button>
                </div>
            )}
        </div>
    );
};

export default JobProgressView;