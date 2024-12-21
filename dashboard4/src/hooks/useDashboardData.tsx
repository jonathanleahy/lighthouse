import { useState, useEffect, useCallback } from 'react';

interface DashboardData {
    repositories: Array<{
        id: string;
        name: string;
        [key: string]: string | number | boolean | null;
    }>;
}

interface DashboardDataResult {
    data: DashboardData | null;
    loading: boolean;
    error: Error | null;
    refetch: (force?: boolean) => void;
}

export function useDashboardData(url: string | null, initialFetch: boolean): DashboardDataResult {
    const [data, setData] = useState<DashboardData | null>(null);
    const [loading, setLoading] = useState<boolean>(initialFetch);
    const [error, setError] = useState<Error | null>(null);

    const fetchData = useCallback(async (force = false) => {
        if (!url) return;
        setLoading(true);
        setError(null);
        try {
            const fetchUrl = force ? `${url}&force=true` : url;
            const response = await fetch(fetchUrl);
            const result = await response.json();
            setData(result);
        } catch (err) {
            if (err instanceof Error) {
                setError(err);
            } else {
                setError(new Error('An unknown error occurred'));
            }
        } finally {
            setLoading(false);
        }
    }, [url]);

    useEffect(() => {
        if (initialFetch) {
            fetchData();
        }
    }, [fetchData, initialFetch]);

    return { data, loading, error, refetch: fetchData };
}