"use client"

import React, { useState, useEffect } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

const ChannelSummary = () => {
    const [channelData, setChannelData] = useState({});
    const [totalJobs, setTotalJobs] = useState(0);

    useEffect(() => {
        const fetchData = () => {
            fetch('http://localhost:8080/queue-length')
                .then(response => response.json())
                .then(data => {
                    setChannelData(data.channels);
                    const total = Object.values(data.channels)
                        .reduce((sum, channel) => sum + channel.length, 0);
                    setTotalJobs(total);
                })
                .catch(error => console.error('Error fetching data:', error));
        };

        fetchData();
        const interval = setInterval(fetchData, 5000);
        return () => clearInterval(interval);
    }, []);

    return (
        <Card className="w-full max-w-2xl mx-auto">
            <CardHeader>
                <CardTitle className="text-2xl font-bold text-center">Active Jobs Summary</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {Object.entries(channelData).map(([channel, jobs]) => (
                        <div key={channel} className="flex justify-between items-center bg-gray-100 p-4 rounded-lg">
                            <span className="font-medium">{channel}:</span>
                            <span className="bg-blue-500 text-white px-3 py-1 rounded-full">
                {jobs.length} jobs
              </span>
                        </div>
                    ))}

                    <div className="mt-6 border-t pt-4">
                        <div className="flex justify-between items-center bg-blue-100 p-4 rounded-lg">
                            <span className="font-bold text-lg">Total Active Jobs:</span>
                            <span className="bg-blue-600 text-white px-4 py-2 rounded-full font-bold">
                {totalJobs}
              </span>
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default ChannelSummary;