// src/components/job/JobUtils.ts

// Type definitions
export interface StepResult {
    status: string;
    message: string;
    start_time: string;
    end_time: string;
}

export interface Step {
    status: string;
    last_updated: string;
    start_time: string;
    end_time: string;
    result?: StepResult;
}

export interface ActiveJob {
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

export interface JobHistory {
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

// Utility functions
export const isValidDate = (dateString: string): boolean => {
    const date = new Date(dateString);
    return date.toString() !== 'Invalid Date' && date.getTime() > new Date('2000-01-01').getTime();
};

export const formatDuration = (start: string, end: string): string => {
    if (!isValidDate(start)) return 'Not started';
    if (!isValidDate(end)) return 'In progress';

    const startTime = new Date(start).getTime();
    const endTime = new Date(end).getTime();
    const duration = endTime - startTime;

    if (duration < 0) return 'Invalid duration';
    if (duration < 1000) return `${duration}ms`;
    return `${(duration / 1000).toFixed(1)}s`;
};

export const getStepTiming = (step: Step): {
    start: string | null;
    end: string | null;
    message: string | null;
} => {
    // Prefer result timing if available and valid
    if (step.result) {
        const resultStart = step.result.start_time;
        const resultEnd = step.result.end_time;

        if (isValidDate(resultStart) && isValidDate(resultEnd)) {
            return {
                start: resultStart,
                end: resultEnd,
                message: step.result.message
            };
        }
    }

    // Fall back to step timing
    return {
        start: isValidDate(step.start_time) ? step.start_time : null,
        end: isValidDate(step.end_time) ? step.end_time : null,
        message: null
    };
};

export const calculateProgress = (job: ActiveJob): number => {
    const completedSteps = Object.values(job.steps).filter(
        step => step.status === 'completed'
    ).length;
    return (completedSteps / job.total_steps) * 100;
};

export const getJobStatus = (job: ActiveJob): 'pending' | 'processing' | 'completed' | 'failed' => {
    const steps = Object.values(job.steps);
    const hasFailedSteps = steps.some(step => step.status === 'failed');
    const hasCompletedSteps = steps.every(step => step.status === 'completed');
    const hasPendingSteps = steps.some(step => step.status === 'pending');

    if (hasFailedSteps) return 'failed';
    if (hasCompletedSteps) return 'completed';
    if (hasPendingSteps) return 'pending';
    return 'processing';
};

export const sortJobsByDate = (jobs: JobHistory[]): JobHistory[] => {
    return [...jobs].sort((a, b) => {
        const dateA = new Date(a.start_time);
        const dateB = new Date(b.start_time);
        return dateB.getTime() - dateA.getTime();
    });
};

export const getLatestJobResult = (jobs: JobHistory[]): JobHistory | null => {
    const sortedJobs = sortJobsByDate(jobs);
    return sortedJobs.length > 0 ? sortedJobs[0] : null;
};

// Polling utilities
export const DEFAULT_POLL_INTERVAL = 2000; // 2 seconds

export const shouldContinuePolling = (job: ActiveJob): boolean => {
    const status = getJobStatus(job);
    return status === 'processing' || status === 'pending';
};

export const hasJobChanged = (
    prevJob: ActiveJob | null,
    currentJob: ActiveJob
): boolean => {
    if (!prevJob) return true;

    return (
        prevJob.completed_steps !== currentJob.completed_steps ||
        JSON.stringify(prevJob.steps) !== JSON.stringify(currentJob.steps)
    );
};

// Add this function at the end of the file
export const formatWaitTime = (waitTimeString: string): string => {
    // Remove the 's' at the end if present
    const durationStr = waitTimeString.replace(/s$/, '');

    // Parse the duration
    const duration = parseFloat(durationStr);

    // Less than a minute
    if (duration < 60) {
        return `${duration.toFixed(1)} seconds`;
    }

    // Less than an hour
    if (duration < 3600) {
        const minutes = Math.floor(duration / 60);
        const seconds = Math.floor(duration % 60);
        return `${minutes} min ${seconds} sec`;
    }

    // Hours and beyond
    const hours = Math.floor(duration / 3600);
    const minutes = Math.floor((duration % 3600) / 60);

    if (hours > 0 && minutes > 0) {
        return `${hours} hr ${minutes} min`;
    } else if (hours > 0) {
        return `${hours} hr`;
    } else {
        return `${minutes} min`;
    }
};