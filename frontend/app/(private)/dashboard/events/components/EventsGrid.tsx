"use client"

import * as React from "react"
import { format, parseISO } from "date-fns"
import { ptBR } from "date-fns/locale"
import { ArrowDownUp, Pencil, Trash2 } from "lucide-react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import type { EventDTO, OccurrenceDTO } from "../types"

type Row = {
    eventId: string
    name: string
    countInMonth: number
    firstDateInMonth: string
    lastDateInMonth: string
    startDate?: string
    intervalDays?: number
    stopAt?: string | null
}

type SortKey =
    | "firstDateInMonth"
    | "lastDateInMonth"
    | "name"
    | "countInMonth"
    | "startDate"
    | "intervalDays"
    | "stopAt"

interface EventsGridProps {
    month: Date
    occurrences: OccurrenceDTO[]
    events: EventDTO[]
    loading?: boolean
    pageSize?: number
    onEdit: (eventId: string) => void
    onDelete: (eventId: string, suggestedDate?: string) => void
}

export const EventsGrid: React.FC<EventsGridProps> = ({
    month,
    occurrences,
    events,
    loading = false,
    pageSize = 10,
    onEdit,
    onDelete,
}) => {
    const rows = React.useMemo<Row[]>(() => {
        const byEvent = new Map<string, { name: string; dates: string[] }>()
        for (const occ of occurrences) {
            const entry = byEvent.get(occ.eventId) ?? { name: occ.name, dates: [] }
            entry.dates.push(occ.date)
            byEvent.set(occ.eventId, entry)
        }

        const mapEvent = new Map(events.map((e) => [e.id, e]))
        const result: Row[] = []
        byEvent.forEach((v, eventId) => {
            v.dates.sort()
            const e = mapEvent.get(eventId)
            result.push({
                eventId,
                name: v.name,
                countInMonth: v.dates.length,
                firstDateInMonth: v.dates[0],
                lastDateInMonth: v.dates[v.dates.length - 1],
                startDate: e?.start_date,
                intervalDays: e?.interval_days,
                stopAt: e?.stop_at ?? null,
            })
        })

        result.sort((a, b) =>
            a.firstDateInMonth === b.firstDateInMonth
                ? a.name.localeCompare(b.name)
                : a.firstDateInMonth.localeCompare(b.firstDateInMonth),
        )
        return result
    }, [occurrences, events])

    const [sortKey, setSortKey] = React.useState<SortKey>("firstDateInMonth")
    const [sortDir, setSortDir] = React.useState<"asc" | "desc">("asc")

    const sortedRows = React.useMemo(() => {
        const arr = [...rows]

        const getVal = (r: Row) => {
            switch (sortKey) {
                case "name":
                    return r.name.toLowerCase()
                case "countInMonth":
                    return r.countInMonth
                case "intervalDays":
                    return r.intervalDays ?? Number.POSITIVE_INFINITY
                case "firstDateInMonth":
                case "lastDateInMonth":
                    return r[sortKey]
                case "startDate":
                    return r.startDate ?? "9999-12-31"
                case "stopAt":
                    return r.stopAt ?? "9999-12-31"
            }
        }

        arr.sort((a, b) => {
            const va = getVal(a)
            const vb = getVal(b)
            let cmp = 0
            if (typeof va === "number" && typeof vb === "number") {
                cmp = va - vb
            } else {
                cmp = String(va).localeCompare(String(vb))
            }
            return sortDir === "asc" ? cmp : -cmp
        })

        return arr
    }, [rows, sortKey, sortDir])

    const [page, setPage] = React.useState(1)
    const totalPages = Math.max(1, Math.ceil(sortedRows.length / pageSize))

    React.useEffect(() => {
        setPage(1)
    }, [month, sortedRows.length, pageSize, sortKey, sortDir])

    const start = (page - 1) * pageSize
    const end = start + pageSize
    const pageRows = sortedRows.slice(start, end)

    const fmt = (iso?: string | null) =>
        iso ? format(parseISO(iso), "dd/MM/yyyy", { locale: ptBR }) : "—"

    return (
        <Card className="glass">
            <CardHeader className="space-y-3">
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                    <CardTitle className="text-xl font-space-grotesk">Events of the Month</CardTitle>
                    <div className="text-sm text-muted-foreground">
                        {rows.length} {rows.length === 1 ? "event" : "events"}
                    </div>
                </div>

                <div className="flex flex-col sm:flex-row gap-2 sm:items-center">
                    <label className="text-sm text-muted-foreground" htmlFor="sortKey">
                        Sort by
                    </label>
                    <select
                        id="sortKey"
                        className="w-full sm:w-auto rounded-md border bg-background px-2 py-1 text-sm"
                        value={sortKey}
                        onChange={(e) => setSortKey(e.target.value as SortKey)}
                    >
                        <option value="firstDateInMonth">First in month</option>
                        <option value="lastDateInMonth">Last in month</option>
                        <option value="name">Name</option>
                        <option value="countInMonth">Occurrences in month</option>
                        <option value="startDate">Start</option>
                        <option value="intervalDays">Recurrence (days)</option>
                        <option value="stopAt">Ends at</option>
                    </select>

                    <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="w-full sm:w-auto"
                        onClick={() => setSortDir((d) => (d === "asc" ? "desc" : "asc"))}
                        aria-label={`Toggle order (${sortDir === "asc" ? "ascending" : "descending"})`}
                    >
                        <ArrowDownUp className="h-4 w-4 mr-2" />
                        {sortDir === "asc" ? "Ascending" : "Descending"}
                    </Button>
                </div>
            </CardHeader>

            <CardContent className="space-y-4">
                <div className="md:hidden space-y-3">
                    {loading ? (
                        <div className="py-6 text-center text-muted-foreground">Loading…</div>
                    ) : pageRows.length === 0 ? (
                        <div className="py-6 text-center text-muted-foreground">
                            No events with occurrences this month.
                        </div>
                    ) : (
                        pageRows.map((r) => (
                            <div
                                key={r.eventId}
                                className="rounded-lg border p-3 bg-card/60 backdrop-blur-sm"
                            >
                                <div className="flex items-start justify-between gap-2">
                                    <div>
                                        <div className="font-medium">{r.name}</div>
                                        <div className="text-xs text-muted-foreground">
                                            {r.countInMonth} occurrence(s) this month
                                        </div>
                                    </div>
                                    <div className="flex gap-2">
                                        <Button
                                            variant="outline"
                                            size="sm"
                                            onClick={() => onEdit(r.eventId)}
                                            aria-label={`Edit ${r.name}`}
                                        >
                                            <Pencil className="h-4 w-4" />
                                        </Button>
                                        <Button
                                            variant="destructive"
                                            size="sm"
                                            onClick={() => onDelete(r.eventId, r.firstDateInMonth)}
                                            aria-label={`Delete ${r.name}`}
                                        >
                                            <Trash2 className="h-4 w-4" />
                                        </Button>
                                    </div>
                                </div>

                                <div className="mt-3 grid grid-cols-2 gap-2 text-sm">
                                    <div className="rounded-md bg-muted/40 p-2">
                                        <div className="text-muted-foreground text-xs">First in month</div>
                                        <div>{fmt(r.firstDateInMonth)}</div>
                                    </div>
                                    <div className="rounded-md bg-muted/40 p-2">
                                        <div className="text-muted-foreground text-xs">Last in month</div>
                                        <div>{fmt(r.lastDateInMonth)}</div>
                                    </div>
                                    <div className="rounded-md bg-muted/40 p-2">
                                        <div className="text-muted-foreground text-xs">Start</div>
                                        <div>{fmt(r.startDate)}</div>
                                    </div>
                                    <div className="rounded-md bg-muted/40 p-2">
                                        <div className="text-muted-foreground text-xs">Recurrence</div>
                                        <div>{r.intervalDays ? `${r.intervalDays} day(s)` : "—"}</div>
                                    </div>
                                    <div className="rounded-md bg-muted/40 p-2 col-span-2">
                                        <div className="text-muted-foreground text-xs">Ends at</div>
                                        <div>{fmt(r.stopAt)}</div>
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>

                <div className="hidden md:block">
                    <div className="overflow-x-auto">
                        <table className="w-full border-collapse min-w-[1000px]">
                            <thead>
                                <tr className="text-left text-sm text-muted-foreground border-b">
                                    <th className="py-2 pr-3">Name</th>
                                    <th className="py-2 px-3">Occurrences</th>
                                    <th className="py-2 px-3 whitespace-nowrap">First in month</th>
                                    <th className="py-2 px-3 whitespace-nowrap">Last in month</th>
                                    <th className="py-2 px-3">Start</th>
                                    <th className="py-2 px-3">Recurrence</th>
                                    <th className="py-2 px-3">Ends at</th>
                                    <th className="py-2 pl-3 text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {loading ? (
                                    <tr>
                                        <td colSpan={8} className="py-8 text-center text-muted-foreground">
                                            Loading…
                                        </td>
                                    </tr>
                                ) : pageRows.length === 0 ? (
                                    <tr>
                                        <td colSpan={8} className="py-8 text-center text-muted-foreground">
                                            No events with occurrences this month.
                                        </td>
                                    </tr>
                                ) : (
                                    pageRows.map((r) => (
                                        <tr key={r.eventId} className="border-b hover:bg-muted/40">
                                            <td className="py-3 pr-3 font-medium">{r.name}</td>
                                            <td className="py-3 px-3">{r.countInMonth}</td>
                                            <td className="py-3 px-3 whitespace-nowrap">{fmt(r.firstDateInMonth)}</td>
                                            <td className="py-3 px-3 whitespace-nowrap">{fmt(r.lastDateInMonth)}</td>
                                            <td className="py-3 px-3 whitespace-nowrap">{fmt(r.startDate)}</td>
                                            <td className="py-3 px-3">
                                                {r.intervalDays ? `${r.intervalDays} day(s)` : "—"}
                                            </td>
                                            <td className="py-3 px-3 whitespace-nowrap">{fmt(r.stopAt)}</td>
                                            <td className="py-3 pl-3">
                                                <div className="flex justify-end gap-2">
                                                    <Button
                                                        variant="outline"
                                                        size="sm"
                                                        onClick={() => onEdit(r.eventId)}
                                                        aria-label={`Edit ${r.name}`}
                                                    >
                                                        <Pencil className="h-4 w-4 mr-1" />
                                                        Edit
                                                    </Button>
                                                    <Button
                                                        variant="destructive"
                                                        size="sm"
                                                        onClick={() => onDelete(r.eventId)}
                                                        aria-label={`Delete ${r.name}`}
                                                    >
                                                        <Trash2 className="h-4 w-4 mr-1" />
                                                        Delete
                                                    </Button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-2">
                    <div
                        className="text-sm text-muted-foreground"
                        aria-live="polite"
                        aria-atomic="true"
                    >
                        {sortedRows.length > 0
                            ? `Showing ${start + 1}–${Math.min(end, sortedRows.length)} of ${sortedRows.length}`
                            : "—"}
                    </div>

                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage((p) => Math.max(1, p - 1))}
                            disabled={page <= 1 || loading}
                        >
                            Previous
                        </Button>
                        <span className="text-sm text-muted-foreground">
                            Page {page} of {totalPages}
                        </span>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                            disabled={page >= totalPages || loading}
                        >
                            Next
                        </Button>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
