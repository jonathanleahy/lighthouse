"use client"

import React, {useState, useEffect} from 'react';
import {
    Pause, Play, StepForward, RefreshCw,
    AlertCircle, Server, List, Clock,
    RotateCw, History, XCircle, Activity, Layers, Briefcase, Database
} from 'lucide-react';
import {Alert, AlertDescription} from '@/components/ui/alert';
import {formatWaitTime} from "@/components/job/JobUtils";

// Type Definitions
interface ServiceType {
    description: string;
    queues: string[];
    handlers: {
        [key: string]: {
            name: string;
            cacheSeconds: number;
            description: string;
        }
    }
}

interface StepResult {
    status: string;
    message: string;
    start_time: string;
    end_time: string;
}

interface Step {
    status: string;
    last_updated: string;
    start_time: string;
    end_time: string;
    result?: StepResult;
}

interface ActiveJob {
    service_name: string;
    type: string;
    status: string;
    total_steps: number;
    completed_steps: number;
    steps: {
        [key: string]: Step;
    };
    start_time: string;
    last_updated: string;
}

interface JobHistory {
    service_name: string;
    type: string;
    start_time: string;
    end_time: string;
    duration: string;
    status: string;
    steps: {
        [key: string]: {
            status: string;
            result?: StepResult;
        };
    };
    tree_output: string;
}

interface QueuedJobInfo {
    service_name: string;
    type: string;
    queue_position: number;
    queue_time: string;
    wait_time: string;
    steps_to_run: string[];
}

interface QueueStats {
    queue_length: number;
    max_queue_size: number;
    active_checks: number;
    queued_services: string[];
    queued_jobs: QueuedJobInfo[];
}

interface SystemMetrics {
    total_cache_entries: number;
    active_processes: number;
    queue_stats: {
        queue_length: number;
        max_queue_size: number;
        active_checks: number;
        queued_services: string[];
        queued_jobs: QueuedJobInfo[];
    };
    system_state: {
        status: string;
        last_updated: string;
        step_mode: boolean;
        current_step?: string;
        message?: string;
    };
}


const API_BASE_URL = 'http://localhost:8080';

const ServiceManagementComponent: React.FC = () => {
    // State Management
    const [serverStatus, setServerStatus] = useState<'connecting' | 'connected' | 'error'>('connecting');
    const [serviceName, setServiceName] = useState('');
    const [serviceType, setServiceType] = useState('');
    const [serviceTypes, setServiceTypes] = useState<{ [key: string]: ServiceType }>({});
    const [checkResult, setCheckResult] = useState<any>(null);
    const [checkError, setCheckError] = useState<string | null>(null);
    const [ignoreCacheChecked, setIgnoreCacheChecked] = useState(false);
    const [systemStatus, setSystemStatus] = useState<any>(null);
    const [debugInfo, setDebugInfo] = useState<any>(null);
    const [debugFormat, setDebugFormat] = useState<'json' | 'text'>('json');
    const [activeJobs, setActiveJobs] = useState<{ [key: string]: ActiveJob }>({});
    const [jobHistory, setJobHistory] = useState<JobHistory[]>([]);
    const [selectedJob, setSelectedJob] = useState<JobHistory | null>(null);
    const [isPollingActive, setIsPollingActive] = useState(false);
    const [errorDetails, setErrorDetails] = useState<string>('');
    const [queuedJobs, setQueuedJobs] = useState<QueuedJobInfo[]>([]);
    const [systemMetrics, setSystemMetrics] = useState<SystemMetrics | null>(null);

    const fetchSystemMetrics = async () => {
        try {
            const response = await fetch(`${API_BASE_URL}/health`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            setSystemMetrics(data);
            console.log('System Metrics:', data);
        } catch (error) {
            console.error('Error fetching system metrics:', error);
        }
    };

    // Modify existing useEffect for periodic metrics fetch
    useEffect(() => {
        // Fetch metrics initially
        fetchSystemMetrics();

        // Set up periodic fetching
        const metricsInterval = setInterval(fetchSystemMetrics, 5000);

        // Cleanup interval on component unmount
        return () => clearInterval(metricsInterval);
    }, []);

    // Check server connectivity
    const checkServerConnection = async () => {

        try {
            const response = await fetch(`${API_BASE_URL}/health`, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
            });

            if (response.ok) {
                setServerStatus('connected');
                setErrorDetails('');
                return true;
            }
            throw new Error(`Server responded with status: ${response.status}`);
        } catch (error) {
            setServerStatus('error');
            setErrorDetails(error instanceof Error ?
                error.message :
                'Failed to connect to server');
            return false;
        }
    };

    // Enhanced fetch with retry logic
    const fetchWithRetry = async (url: string, options: RequestInit = {}, retries = 3) => {
        for (let i = 0; i < retries; i++) {
            try {
                const response = await fetch(url, {
                    ...options,
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json',
                        ...options.headers,
                    },
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                return await response.json();
            } catch (error) {
                if (i === retries - 1) throw error;
                await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
            }
        }
    };

    useEffect(() => {
        const fetchQueuedJobs = async () => {
            try {
                const data = await fetchWithRetry(`${API_BASE_URL}/jobs/queue`);
                console.log('Queued Jobs Fetched:', data.queued_jobs); // Add logging
                setQueuedJobs(data.queued_jobs || []);
            } catch (error) {
                console.error('Error fetching queued jobs:', error);
            }
        };

        // Add queued jobs fetch to your existing job polling logic
        if (isPollingActive) {
            fetchQueuedJobs();
        }
    }, [isPollingActive]);


    // Fetch Service Types on Component Mount
    useEffect(() => {
        const fetchServiceTypes = async () => {
            try {
                const data = await fetchWithRetry(`${API_BASE_URL}/service-types`);
                setServiceTypes(data);

                // Set default service type
                const types = Object.keys(data);
                if (types.length > 0) {
                    setServiceType(types[0]);
                }
            } catch (error) {
                console.error('Error fetching service types:', error);
                // Fallback service types
                setServiceTypes({
                    'check': {
                        description: 'Default Check Type',
                        handlers: {}
                    },
                    'report': {
                        description: 'Default Report Type',
                        handlers: {}
                    }
                });
                setServiceType('check');
            }
        };

        fetchServiceTypes();
    }, []);

    // Initialize connection check and polling
    useEffect(() => {
        checkServerConnection();
        const intervalId = setInterval(checkServerConnection, 30000);
        return () => clearInterval(intervalId);
    }, []);

    // Active Jobs Polling
    // Active Jobs Polling
    useEffect(() => {
        let interval: NodeJS.Timer | null = null;

        if (isPollingActive) {
            interval = setInterval(async () => {
                try {
                    // Fetch active jobs
                    const data = await fetchWithRetry(`${API_BASE_URL}/jobs/progress`);
                    setActiveJobs(data.active_jobs);

                    // Fetch queued jobs when active jobs are updated
                    const queueData = await fetchWithRetry(`${API_BASE_URL}/jobs/queue`);
                    setQueuedJobs(queueData.queued_jobs || []);

                    // If no active jobs, stop polling and refresh history
                    if (Object.keys(data.active_jobs).length === 0) {
                        setIsPollingActive(false);
                        fetchJobHistory();
                    }
                } catch (error) {
                    console.error('Error polling jobs:', error);
                    setIsPollingActive(false);
                }
            }, 2000);
        }

        return () => {
            if (interval) clearInterval(interval);
        };
    }, [isPollingActive]);

    useEffect(() => {
        // Add this to your existing useEffect or create a new one
        if (debugInfo) {
            console.log('Full Debug Info:', debugInfo);
            console.log('Queued Jobs:', debugInfo.queue_status?.queued_jobs);
        }
    }, [debugInfo]);

    useEffect(() => {
        const fetchQueuedJobs = async () => {
            try {
                const data = await fetchWithRetry(`${API_BASE_URL}/jobs/queue`);
                console.log('Fetched Queued Jobs:', data.queued_jobs);
                setQueuedJobs(data.queued_jobs || []);
            } catch (error) {
                console.error('Error fetching queued jobs:', error);
            }
        };

        // Fetch queued jobs when polling is active or when component mounts
        if (isPollingActive) {
            fetchQueuedJobs();
        }

        // Add a method to manually fetch queued jobs when button is clicked
        const handleManualQueueRefresh = () => {
            fetchQueuedJobs();
        };

        // Expose the method if needed
        return () => {
            // Cleanup if necessary
        };
    }, [isPollingActive]);

    // Fetch Job History
    const fetchJobHistory = async () => {
        try {
            const data = await fetchWithRetry(`${API_BASE_URL}/jobs/history`);
            setJobHistory(data.completed_jobs);
        } catch (error) {
            console.error('Error fetching job history:', error);
        }
    };

    // Submit Service Check
    const submitServiceCheck = async (e: React.FormEvent) => {
        e.preventDefault();
        setCheckError(null);
        setCheckResult(null);

        try {
            const queryParams = new URLSearchParams();
            if (ignoreCacheChecked) {
                queryParams.append('resetCache', 'true');
            }

            const data = await fetchWithRetry(`${API_BASE_URL}/check?${queryParams.toString()}`, {
                method: 'POST',
                body: JSON.stringify({
                    name: serviceName,
                    type: serviceType
                })
            });

            setCheckResult(data);
            setIsPollingActive(true);

        } catch (error) {
            setCheckError(error instanceof Error ? error.message : 'An unknown error occurred');
            console.error('Check submission error:', error);
        }
    };

    // System Control
    const controlSystem = async (command: string) => {
        try {
            const status = await fetchWithRetry(`${API_BASE_URL}/control?command=${command}`, {
                method: 'POST'
            });
            setSystemStatus(status);
            fetchDebugInfo(debugFormat);
        } catch (error) {
            console.error(`Error during ${command}:`, error);
        }
    };

    // Fetch Debug Information
    const fetchDebugInfo = async (format: 'json' | 'text' = 'json') => {
        try {
            const response = await fetch(`${API_BASE_URL}/health?format=${format}`);
            const data = await response.json(); // Always parse as JSON for the health endpoint
            console.log('Debug Info:', data); // Add this to help diagnose
            setDebugInfo(data);
            setDebugFormat(format);
        } catch (error) {
            console.error('Error fetching debug info:', error);
        }
    };

    // Cache Invalidation
    const invalidateCache = async (svcName: string, svcType: string) => {
        try {
            await fetchWithRetry(`${API_BASE_URL}/invalidate`, {
                method: 'POST',
                body: JSON.stringify({
                    service_name: svcName,
                    type: svcType,
                    reset_times: true
                })
            });

            fetchDebugInfo(debugFormat);
        } catch (error) {
            console.error('Error invalidating cache:', error);
        }
    };

    // Render loading state
    if (serverStatus === 'connecting') {
        return (
            <div className="p-4">
                <Alert>
                    <Clock className="h-4 w-4"/>
                    <AlertDescription>
                        Connecting to server...
                    </AlertDescription>
                </Alert>
            </div>
        );
    }

    // Render error state
    if (serverStatus === 'error') {
        return (
            <div className="p-4">
                <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4"/>
                    <AlertDescription>
                        Failed to connect to server: {errorDetails}
                        <br/>
                        Please ensure:
                        <ul className="list-disc ml-6 mt-2">
                            <li>The server is running on port 8080</li>
                            <li>No firewall is blocking the connection</li>
                            <li>CORS is properly configured on the server</li>
                        </ul>
                    </AlertDescription>
                </Alert>
            </div>
        );
    }

    // Main component render
    return (
        <div className="max-w-4xl mx-auto p-6 space-y-6">
            {/* Service Check Form */}
            <div className="bg-white shadow-md rounded-lg p-6">
                <h2 className="text-xl font-bold mb-4 flex items-center">
                    <Server className="mr-2 text-blue-600"/> Submit Service Check
                </h2>
                <form onSubmit={submitServiceCheck} className="space-y-4">
                    <div>
                        <label className="block mb-2">Service Name</label>
                        <input
                            type="text"
                            value={serviceName}
                            onChange={(e) => setServiceName(e.target.value)}
                            className="w-full px-3 py-2 border rounded"
                            placeholder="Enter service name"
                            required
                        />
                    </div>
                    <div>
                        <label className="block mb-2">Service Type</label>
                        <select
                            value={serviceType}
                            onChange={(e) => setServiceType(e.target.value)}
                            className="w-full px-3 py-2 border rounded"
                            required
                        >
                            {Object.entries(serviceTypes).map(([type, details]) => (
                                <option key={type} value={type}>
                                    {type} - {details.description}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div className="flex items-center">
                        <input
                            type="checkbox"
                            id="ignoreCache"
                            checked={ignoreCacheChecked}
                            onChange={() => setIgnoreCacheChecked(!ignoreCacheChecked)}
                            className="mr-2"
                        />
                        <label htmlFor="ignoreCache" className="text-sm">
                            Ignore/Reset Cache
                        </label>
                    </div>
                    <button
                        type="submit"
                        disabled={!serviceName || Object.keys(serviceTypes).length === 0}
                        className={`w-full py-2 rounded ${
                            serviceName && Object.keys(serviceTypes).length > 0
                                ? 'bg-blue-600 text-white hover:bg-blue-700'
                                : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                        }`}
                    >
                        Submit Service Check
                    </button>
                </form>

                {checkError && (
                    <div className="mt-4 bg-red-100 p-3 rounded flex items-center">
                        <AlertCircle className="mr-2 text-red-600"/>
                        <p className="text-red-800">{checkError}</p>
                    </div>
                )}

                {checkResult && (
                    <div className="mt-4 bg-green-50 p-4 rounded">
                        <h3 className="font-bold mb-2">Check Result</h3>
                        <pre className="bg-white p-2 rounded overflow-x-auto text-sm">
                            {JSON.stringify(checkResult, null, 2)}
                        </pre>
                    </div>
                )}
            </div>

            {/* Job Processing Section */}
            {(Object.keys(activeJobs).length > 0 ||
                (debugInfo?.queue_status?.queued_jobs && debugInfo.queue_status.queued_jobs.length > 0)) && (
                <div className="bg-white shadow-md rounded-lg p-6">
                    <h2 className="text-xl font-bold mb-4 flex items-center">
                        <RotateCw className="mr-2 animate-spin text-blue-600"/> Job Processing
                    </h2>

                    {/* Active Jobs Subsection */}
                    {Object.keys(activeJobs).length > 0 && (
                        <>
                            <h3 className="text-lg font-semibold mb-3 text-blue-700">Active Jobs</h3>
                            <div className="space-y-4 mb-6">
                                {Object.entries(activeJobs).map(([jobId, job]) => (
                                    <div key={jobId} className="border rounded p-4">
                                        <div className="flex justify-between items-center mb-2">
                                            <span className="font-medium">{job.service_name}</span>
                                            <span className="text-sm text-blue-600">
                                    {job.completed_steps} / {job.total_steps} steps
                                </span>
                                        </div>
                                        <div className="w-full bg-gray-200 rounded-full h-2 mb-4">
                                            <div
                                                className="bg-blue-600 rounded-full h-2 transition-all duration-500"
                                                style={{width: `${(job.completed_steps / job.total_steps) * 100}%`}}
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            {Object.entries(job.steps).map(([stepName, step]) => {
                                                const timing = step.result ? {
                                                    start: step.result.start_time,
                                                    end: step.result.end_time
                                                } : {
                                                    start: step.start_time,
                                                    end: step.end_time
                                                };

                                                const duration = step.status === 'completed' &&
                                                step.result?.start_time &&
                                                step.result?.end_time ?
                                                    new Date(step.result.end_time).getTime() -
                                                    new Date(step.result.start_time).getTime() : null;

                                                return (
                                                    <div key={stepName}
                                                         className="flex justify-between items-center text-sm border-b pb-2">
                                                        <div className="flex items-center space-x-2">
                                                            <span>{stepName}</span>
                                                            {duration && (
                                                                <span className="text-gray-500">
                                                        ({(duration / 1000).toFixed(1)}s)
                                                    </span>
                                                            )}
                                                        </div>
                                                        <span className={`${
                                                            step.status === 'completed' ? 'text-green-600' :
                                                                step.status === 'pending' ? 'text-gray-400' :
                                                                    'text-blue-600'
                                                        }`}>
                                                {step.status}
                                            </span>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </>
                    )}

                    {/* Queued Jobs Subsection */}
                    {debugInfo?.queue_status?.queued_jobs && debugInfo.queue_status.queued_jobs.length > 0 && (
                        <>
                            <h3 className="text-lg font-semibold mb-3 text-yellow-700">Queued Jobs</h3>
                            <div className="space-y-4">
                                {debugInfo.queue_status.queued_jobs.map((job, index) => (
                                    <div key={index} className="border rounded p-4 bg-yellow-50">
                                        <div className="flex justify-between items-center mb-2">
                                            <div>
                                                <span className="font-medium">{job.service_name}</span>
                                                <span className="text-sm text-gray-500 ml-2">({job.type})</span>
                                            </div>
                                            <span className="text-sm text-yellow-600">
                                    Position: {job.queue_position}
                                </span>
                                        </div>
                                        <div className="text-sm text-gray-600">
                                            <div className="flex justify-between">
                                                <span>Queued At: {new Date(job.queue_time).toLocaleTimeString()}</span>
                                                <span>Waiting: {job.wait_time}</span>
                                            </div>
                                            {job.steps_to_run && job.steps_to_run.length > 0 && (
                                                <div className="mt-2">
                                                    <span className="font-medium">Steps to Run:</span>
                                                    <span className="ml-2 text-gray-500">
                                            {job.steps_to_run.join(', ')}
                                        </span>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </>
                    )}
                </div>
            )}

            {queuedJobs.length > 0 && (
                <div className="bg-white shadow-md rounded-lg p-6">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-xl font-bold flex items-center">
                            <Clock className="mr-2 text-yellow-600"/> Queued Jobs
                            <span className="ml-2 text-sm text-gray-500">({queuedJobs.length} total)</span>
                        </h2>
                        <button
                            onClick={() => {
                                const fetchQueuedJobs = async () => {
                                    try {
                                        const data = await fetchWithRetry(`${API_BASE_URL}/jobs/queue`);
                                        setQueuedJobs(data.queued_jobs || []);
                                    } catch (error) {
                                        console.error('Error fetching queued jobs:', error);
                                    }
                                };
                                fetchQueuedJobs();
                            }}
                            className="text-blue-600 hover:text-blue-800"
                        >
                            <RefreshCw className="w-5 h-5"/>
                        </button>
                    </div>
                    <div className="space-y-4">
                        {queuedJobs.map((job, index) => (
                            <div key={index} className="border rounded p-4 bg-yellow-50">
                                <div className="flex justify-between items-center mb-2">
                                    <div>
                                        <span className="font-medium">{job.service_name}</span>
                                        <span className="text-sm text-gray-500 ml-2">({job.type})</span>
                                    </div>
                                    <span className="text-sm text-yellow-600">
                Position: {job.queue_position}
            </span>
                                </div>
                                <div className="text-sm text-gray-600">
                                    <div className="flex justify-between">
                                        <span>Queued At: {new Date(job.queue_time).toLocaleTimeString()}</span>
                                        <span>Waiting: {formatWaitTime(job.wait_time)}</span>
                                    </div>
                                    {job.steps_to_run && job.steps_to_run.length > 0 && (
                                        <div className="mt-2">
                                            <span className="font-medium">Steps to Run:</span>
                                            <span className="ml-2 text-gray-500">
                        {job.steps_to_run.join(', ')}
                    </span>
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* System Controls */}
            <div className="bg-white shadow-md rounded-lg p-6">
                <h2 className="text-xl font-bold mb-4 flex items-center">
                    <Clock className="mr-2 text-green-600"/> System Controls
                </h2>
                <div className="grid grid-cols-2 gap-4">
                    <button
                        onClick={() => controlSystem('pause')}
                        className="bg-yellow-500 text-white py-2 rounded flex items-center justify-center hover:bg-yellow-600"
                    >
                        <Pause className="mr-2"/> Pause
                    </button>
                    <button
                        onClick={() => controlSystem('step')}
                        className="bg-blue-500 text-white py-2 rounded flex items-center justify-center hover:bg-blue-600"
                    >
                        <StepForward className="mr-2"/> Step
                    </button>
                    <button
                        onClick={() => controlSystem('resume')}
                        className="bg-green-500 text-white py-2 rounded flex items-center justify-center hover:bg-green-600"
                    >
                        <Play className="mr-2"/> Resume
                    </button>
                    <button
                        onClick={() => controlSystem('reset')}
                        className="bg-red-500 text-white py-2 rounded flex items-center justify-center hover:bg-red-600"
                    >
                        <RefreshCw className="mr-2"/> Reset
                    </button>
                </div>

                {systemStatus && (
                    <div className="mt-4 bg-gray-50 p-4 rounded">
                        <h3 className="font-bold mb-2">System Status</h3>
                        <pre className="bg-white p-2 rounded overflow-x-auto text-sm">
                            {JSON.stringify(systemStatus, null, 2)}
                        </pre>
                    </div>
                )}
            </div>

            {/* Debug Information */}
            <div className="bg-white shadow-md rounded-lg p-6">
                <h2 className="text-xl font-bold mb-4 flex items-center">
                    <List className="mr-2 text-purple-600"/> System Debug Info
                </h2>
                <div className="flex space-x-4 mb-4">
                    <button
                        onClick={() => {
                            setDebugFormat('json');
                            fetchDebugInfo('json');
                        }}
                        className={`py-2 px-4 rounded ${
                            debugFormat === 'json'
                                ? 'bg-purple-600 text-white'
                                : 'bg-gray-200 text-black hover:bg-gray-300'
                        }`}
                    >
                        JSON
                    </button>
                    <button
                        onClick={() => {
                            setDebugFormat('text');
                            fetchDebugInfo('text');
                        }}
                        className={`py-2 px-4 rounded ${
                            debugFormat === 'text'
                                ? 'bg-purple-600 text-white'
                                : 'bg-gray-200 text-black hover:bg-gray-300'
                        }`}
                    >
                        Text
                    </button>
                    <button
                        onClick={() => fetchDebugInfo(debugFormat)}
                        className="py-2 px-4 rounded bg-gray-200 hover:bg-gray-300 flex items-center gap-2"
                    >
                        <RefreshCw className="w-4 h-4"/>
                        Refresh
                    </button>
                </div>

                {debugInfo && (
                    <div className="bg-gray-50 p-4 rounded">
                        <pre className="overflow-x-auto text-sm whitespace-pre-wrap">
                            {debugFormat === 'json'
                                ? JSON.stringify(debugInfo, null, 2)
                                : debugInfo
                            }
                        </pre>
                    </div>
                )}
            </div>

            {/* Job History */}
            <div className="bg-white shadow-md rounded-lg p-6">
                <div className="flex justify-between items-center mb-4">
                    <h2 className="text-xl font-bold flex items-center">
                        <History className="mr-2 text-indigo-600"/> Job History
                    </h2>
                    <button
                        onClick={fetchJobHistory}
                        className="text-indigo-600 hover:text-indigo-800"
                    >
                        <RefreshCw className="w-5 h-5"/>
                    </button>
                </div>
                <div className="space-y-4">
                    {jobHistory.map((job, index) => (
                        <div key={index} className="border rounded p-3 hover:bg-gray-50">
                            <div className="flex justify-between items-center">
                                <div>
                                    <span className="font-medium">{job.service_name}</span>
                                    <span className="text-sm text-gray-500 ml-2">({job.type})</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <span className="text-sm text-gray-500">{job.duration}</span>
                                    <button
                                        onClick={() => setSelectedJob(job)}
                                        className="text-blue-600 hover:text-blue-800 text-sm"
                                    >
                                        View Details
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Job Details Modal */}
            {selectedJob && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-lg max-w-2xl w-full max-h-[80vh] overflow-y-auto">
                        <div className="p-6">
                            <div className="flex justify-between items-center mb-4">
                                <div>
                                    <h3 className="text-lg font-bold">{selectedJob.service_name}</h3>
                                    <p className="text-sm text-gray-500">
                                        Duration: {selectedJob.duration}
                                    </p>
                                </div>
                                <button
                                    onClick={() => setSelectedJob(null)}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    <XCircle className="w-6 h-6"/>
                                </button>
                            </div>

                            <div className="mb-4">
                                <h4 className="font-medium mb-2">Tree Output</h4>
                                <pre className="bg-gray-50 p-4 rounded overflow-x-auto text-sm">
                                    {selectedJob.tree_output}
                                </pre>
                            </div>

                            <div className="mb-4">
                                <h4 className="font-medium mb-2">Step Details</h4>
                                {Object.entries(selectedJob.steps).map(([stepName, step]) => (
                                    <div key={stepName} className="mb-2 border-b pb-2">
                                        <div className="flex justify-between">
                                            <span className="font-medium">{stepName}</span>
                                            <span className={`text-sm ${
                                                step.status === 'completed' ? 'text-green-600' :
                                                    step.status === 'failed' ? 'text-red-600' :
                                                        'text-blue-600'
                                            }`}>
                                                {step.status}
                                            </span>
                                        </div>
                                        {step.result && (
                                            <p className="text-sm text-gray-600 mt-1">
                                                {step.result.message}
                                            </p>
                                        )}
                                    </div>
                                ))}
                            </div>

                            <div className="flex justify-end">
                                <button
                                    onClick={() => setSelectedJob(null)}
                                    className="bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded"
                                >
                                    Close
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {systemMetrics && (
                <div className="bg-white shadow-md rounded-lg p-6">
                    <h2 className="text-xl font-bold mb-4 flex items-center">
                        <Activity className="mr-2 text-green-600"/> System Metrics
                    </h2>

                    <div className="grid grid-cols-2 gap-4">
                        {/* System State Card */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center justify-between mb-3">
                                <h3 className="font-semibold flex items-center">
                                    <Layers className="mr-2 text-blue-600"/> System State
                                </h3>
                                <span className={`text-sm font-medium ${
                                    systemMetrics.system_state?.status === 'running'
                                        ? 'text-green-600'
                                        : 'text-yellow-600'
                                }`}>
                                    {systemMetrics.system_state?.status}
                                </span>
                            </div>
                            <div className="space-y-2 text-sm text-gray-600">
                                <div>
                                    <span className="font-medium">Last Updated:</span>{' '}
                                    {new Date(systemMetrics.system_state?.last_updated).toLocaleString()}
                                </div>
                                {systemMetrics.system_state?.step_mode && (
                                    <div className="text-blue-600 flex items-center">
                                        <StepForward className="mr-2 w-4 h-4"/> Step Mode Active
                                    </div>
                                )}
                                {systemMetrics.system_state?.current_step && (
                                    <div>
                                        <span className="font-medium">Current Step:</span>{' '}
                                        {systemMetrics.system_state.current_step}
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Queue Statistics Card */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center justify-between mb-3">
                                <h3 className="font-semibold flex items-center">
                                    <Briefcase className="mr-2 text-purple-600"/> Queue Statistics
                                </h3>
                                <span className="text-sm font-medium text-gray-600">
                                    {systemMetrics.queue_stats?.queue_length} / {systemMetrics.queue_stats?.max_queue_size}
                                </span>
                            </div>
                            <div className="space-y-2 text-sm text-gray-600">
                                <div>
                                    <span className="font-medium">Queued Jobs:</span>{' '}
                                    {systemMetrics.queue_stats?.queue_length}
                                </div>
                                <div>
                                    <span className="font-medium">Active Checks:</span>{' '}
                                    {systemMetrics.queue_stats?.active_checks}
                                </div>
                                {systemMetrics.queue_stats?.queued_services?.length > 0 && (
                                    <div>
                                        <span className="font-medium">Queued Services:</span>{' '}
                                        {systemMetrics.queue_stats?.queued_services?.join(', ')}
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* System Resources Card */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center justify-between mb-3">
                                <h3 className="font-semibold flex items-center">
                                    <Database className="mr-2 text-indigo-600"/> System Resources
                                </h3>
                            </div>
                            <div className="space-y-2 text-sm text-gray-600">
                                <div>
                                    <span className="font-medium">Total Cache Entries:</span>{' '}
                                    {systemMetrics.total_cache_entries}
                                </div>
                                <div>
                                    <span className="font-medium">Active Processes:</span>{' '}
                                    {systemMetrics.active_processes}
                                </div>
                            </div>
                        </div>

                        {/* Job Queue Details Card */}
                        <div className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-center justify-between mb-3">
                                <h3 className="font-semibold flex items-center">
                                    <Clock className="mr-2 text-orange-600"/> Queued Job Details
                                </h3>
                            </div>
                            <div className="space-y-2 text-sm">
                                {systemMetrics.queue_stats?.queued_jobs.length > 0 ? (
                                    systemMetrics.queue_stats?.queued_jobs.map((job, index) => (
                                        <div key={index} className="border-b pb-2 last:border-b-0">
                                            <div className="flex justify-between">
                                                <span className="font-medium">{job.service_name}</span>
                                                <span className="text-gray-500">
                                                    Position: {job.queue_position}
                                                </span>
                                            </div>
                                            <div className="text-gray-600">
                                                <span>Type: {job.type}</span>
                                                <span className="ml-2">
                                                    Wait Time: {formatWaitTime(job.wait_time)}
                                                </span>
                                            </div>
                                        </div>
                                    ))
                                ) : (
                                    <div className="text-gray-500 text-center">No jobs in queue</div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            )}

        </div>
    );
};

export default ServiceManagementComponent;