"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface LinkChartsProps {
    links: Array<{
        id: string
        clicks: number
        short_url: string
        created_at: string
    }>
}

export function LinkCharts({ links }: LinkChartsProps) {
    const chartData = links.map(link => ({
        name: link.short_url,
        clicks: link.clicks,
    }))

    return (
        <Card className="bg-card hover:shadow-md transition-shadow">
            <CardHeader>
                <CardTitle>Link Performance</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="h-[300px]">
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="name" />
                            <YAxis />
                            <Tooltip />
                            <Bar dataKey="clicks" fill="#8884d8" />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </CardContent>
        </Card>
    )
} 