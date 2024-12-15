import { useState, useEffect } from 'react';

interface FetchState<T> {
    data: T | null;
    loading: boolean;
    error: Error | null;
}

export function useFetchData<T>(url: string) {
    const [state, setState] = useState<FetchState<T>>({
        data: null,
        loading: true,
        error: null,
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        const fetchData = async () => {
            try {
                setState(prev => ({ ...prev, loading: true }));

                const response = await fetch(url, {
                    signal,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();

                if (!signal.aborted) {
                    setState({
                        data,
                        loading: false,
                        error: null,
                    });
                }
            } catch (error) {
                if (!signal.aborted) {
                    setState({
                        data: null,
                        loading: false,
                        error: error instanceof Error ? error : new Error('An unknown error occurred'),
                    });
                }
            }
        };

        fetchData();

        return () => {
            abortController.abort();
        };
    }, [url]);

    return {
        data: state.data,
        loading: state.loading,
        error: state.error,
    };
}

// Optional: Add a version that automatically retries failed requests
export function useFetchDataWithRetry<T>(url: string, retryCount = 3, retryDelay = 1000) {
    const [state, setState] = useState<FetchState<T>>({
        data: null,
        loading: true,
        error: null,
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        const fetchWithRetry = async (retriesLeft: number): Promise<void> => {
            try {
                setState(prev => ({ ...prev, loading: true }));

                const response = await fetch(url, {
                    signal,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();

                if (!signal.aborted) {
                    setState({
                        data,
                        loading: false,
                        error: null,
                    });
                }
            } catch (error) {
                if (signal.aborted) return;

                if (retriesLeft > 0 && error instanceof Error && !error.message.includes('aborted')) {
                    // Wait for retryDelay milliseconds before retrying
                    await new Promise(resolve => setTimeout(resolve, retryDelay));
                    return fetchWithRetry(retriesLeft - 1);
                }

                setState({
                    data: null,
                    loading: false,
                    error: error instanceof Error ? error : new Error('An unknown error occurred'),
                });
            }
        };

        fetchWithRetry(retryCount);

        return () => {
            abortController.abort();
        };
    }, [url, retryCount, retryDelay]);

    return {
        data: state.data,
        loading: state.loading,
        error: state.error,
    };
}