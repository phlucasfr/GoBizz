"use client"

import React from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { format, startOfWeek, endOfWeek, eachDayOfInterval, isSameMonth, isSameDay } from "date-fns"
import { enUS } from "date-fns/locale"
import type { OccurrenceDTO } from "../types"

interface CalendarProps {
  currentMonth: Date
  onMonthChange: (month: Date) => void
  occurrences: OccurrenceDTO[]
  onDayClick: (date: Date, occurrences: OccurrenceDTO[]) => void
  loading: boolean
}

export const Calendar: React.FC<CalendarProps> = ({
  currentMonth,
  onMonthChange,
  occurrences,
  onDayClick,
  loading,
}) => {
  const occurrencesByDate = React.useMemo(() => {
    const map = new Map<string, OccurrenceDTO[]>()
    occurrences.forEach((occurrence) => {
      const dateKey = occurrence.date
      if (!map.has(dateKey)) {
        map.set(dateKey, [])
      }
      map.get(dateKey)!.push(occurrence)
    })
    return map
  }, [occurrences])

  const [selectedDate, setSelectedDate] = React.useState<Date | null>(null)

  const monthStart = new Date(currentMonth.getFullYear(), currentMonth.getMonth(), 1)
  const monthEnd = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 0)
  const calendarStart = startOfWeek(monthStart, { locale: enUS })
  const calendarEnd = endOfWeek(monthEnd, { locale: enUS })

  const calendarDays = eachDayOfInterval({
    start: calendarStart,
    end: calendarEnd,
  })

  const weekDays = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]

  const handleDayClick = (date: Date) => {
    const dateStr = format(date, "yyyy-MM-dd")
    const dayOccurrences = occurrencesByDate.get(dateStr) || []

    setSelectedDate(date)

    if (dayOccurrences.length > 0) {
      onDayClick(date, dayOccurrences)
    }
  }

  const handlePreviousMonth = () => {
    const newMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1, 1)
    onMonthChange(newMonth)
  }

  const handleNextMonth = () => {
    const newMonth = new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 1)
    onMonthChange(newMonth)
  }

  return (
    <Card className="glass">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <Button
            variant="ghost"
            size="icon"
            onClick={handlePreviousMonth}
            disabled={loading}
            className="h-8 w-8 text-foreground/70 hover:text-foreground hover:bg-muted/50"
            aria-label="Previous month"
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>

          <CardTitle className="text-xl font-space-grotesk text-foreground">
            {format(currentMonth, "MMMM yyyy", { locale: enUS })}
          </CardTitle>

          <Button
            variant="ghost"
            size="icon"
            onClick={handleNextMonth}
            disabled={loading}
            className="h-8 w-8 text-foreground/70 hover:text-foreground hover:bg-muted/50"
            aria-label="Next month"
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>

      <CardContent className={`transition-opacity duration-200 ${loading ? "opacity-50" : "opacity-100"}`}>
        <div className="grid grid-cols-7 gap-2 mb-4">
          {weekDays.map((day) => (
            <div
              key={day}
              className="h-8 flex items-center justify-center text-sm font-medium text-muted-foreground bg-muted/30 rounded-md"
            >
              {day}
            </div>
          ))}
        </div>

        <div className="grid grid-cols-7 gap-2">
          {calendarDays.map((date) => {
            const dateStr = format(date, "yyyy-MM-dd")
            const dayOccurrences = occurrencesByDate.get(dateStr) || []
            const hasEvents = dayOccurrences.length > 0
            const isCurrentMonth = isSameMonth(date, currentMonth)
            const isToday = isSameDay(date, new Date())
            const isSelected = selectedDate && isSameDay(date, selectedDate)

            return (
              <Button
                key={date.toISOString()}
                variant="ghost"
                className={`
                  h-12 w-full p-0 text-sm font-dm-sans relative rounded-lg border transition-all duration-200
                  ${
                    !isCurrentMonth
                      ? "text-muted-foreground/40 border-transparent bg-transparent hover:bg-muted/20"
                      : "text-foreground border-border/20 bg-card/50 hover:bg-card/80"
                  }
                  ${isToday ? "bg-primary/10 border-primary/30 text-primary font-semibold ring-1 ring-primary/20" : ""}
                  ${isSelected ? "bg-primary/20 border-primary/50 ring-2 ring-primary/30" : ""}
                  ${hasEvents ? "bg-cyan-50 dark:bg-cyan-950/30 border-cyan-200 dark:border-cyan-800/50 text-cyan-800 dark:text-cyan-200" : ""}
                  ${isToday && hasEvents ? "bg-gradient-to-br from-primary/15 to-cyan-50 dark:from-primary/20 dark:to-cyan-950/40" : ""}
                  hover:shadow-sm hover:scale-[1.02] active:scale-[0.98]
                  disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100
                `}
                onClick={() => handleDayClick(date)}
                disabled={loading}
                aria-label={`${format(date, "d 'of' MMMM", { locale: enUS })}${hasEvents ? ` - ${dayOccurrences.length} event(s)` : ""}`}
                aria-pressed={!!isSelected}
              >
                {format(date, "d")}
                {hasEvents && (
                  <div className="absolute bottom-1.5 left-1/2 transform -translate-x-1/2 flex gap-0.5">
                    {dayOccurrences.slice(0, 3).map((_, index) => (
                      <div key={index} className="w-1.5 h-1.5 bg-cyan-600 dark:bg-cyan-400 rounded-full shadow-sm" />
                    ))}
                    {dayOccurrences.length > 3 && (
                      <div className="w-1.5 h-1.5 bg-cyan-700 dark:bg-cyan-300 rounded-full opacity-80 shadow-sm" />
                    )}
                  </div>
                )}
              </Button>
            )
          })}
        </div>

        {loading && (
          <div className="text-center mt-6 text-sm text-muted-foreground bg-muted/30 rounded-lg py-3">
            Loading events...
          </div>
        )}
      </CardContent>
    </Card>
  )
}
