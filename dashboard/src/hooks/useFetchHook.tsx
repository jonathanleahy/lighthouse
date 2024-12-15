// hooks/useFetchHook.ts
import { useState, useEffect, useCallback } from 'react';

interface FetchState<T> {
    data: T | null;
    loading: boolean;
    error: Error | null;
    refetch: () => void;
}

export function useFetchData<T>(url: string): FetchState<T> {
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchData = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(url);
            const result = await response.json();
            setData(result);
        } catch (err) {
            setError(err as Error);
        } finally {
            setLoading(false);
        }
    }, [url]);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    return { data, loading, error, refetch: fetchData };
}