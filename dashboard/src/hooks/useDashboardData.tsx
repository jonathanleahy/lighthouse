"use client";

import { useState, useEffect } from 'react';

interface DashboardData {
    repoDesc?: string;
    repoSquad?: string;
    repoBitUrl?: string;
    repoCodefresh?: string;
    argocd?: { url?: string };
}

interface DashboardDataResult {
    data: DashboardData | null;
    loading: boolean;
    error: Error | null;
}

export const useDashboardData = (url: string | null, fetchData: boolean): DashboardDataResult => {
    const [data, setData] = useState<DashboardData | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (!url) return;

        const fetchDashboardData = async () => {
            setLoading(true);
            try {
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error(`Error: ${response.statusText}`);
                }
                const result: DashboardData = await response.json();
                setData(result);
            } catch (err) {
                setError(err as Error);
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardData();
    }, [url, fetchData]);

    return { data, loading, error };
};