"use client"

import React from "react"
import { format, startOfMonth, endOfMonth } from "date-fns"
import { LucideCalendar } from "lucide-react"

import { EventForm } from "./components/EventForm"
import { Calendar as EventCalendar } from "./components/Calendar"
import { Filters } from "./components/Filters"
import { DeleteDialog } from "./components/DeleteDialog"
import { EditEventDialog } from "./components/EditEventDialog"
import { EventActionsPopover } from "./components/EventActionsPopover"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { ThemeToggle } from "@/components/theme-toggle"
import { useToast } from "@/hooks/use-toast"
import { Toaster } from "@/components/ui/toaster"

import { eventsApi } from "./services/api"
import { useDebounce } from "./hooks/useDebounce"
import type { EventDTO, OccurrenceDTO, CreateEventPayload, UpdateEventPayload, ApiError } from "./types"

function App() {
  // State
  const [currentMonth, setCurrentMonth] = React.useState(new Date())
  const [nameFilter, setNameFilter] = React.useState("")
  const [occurrences, setOccurrences] = React.useState<OccurrenceDTO[]>([])
  const [events, setEvents] = React.useState<EventDTO[]>([])

  // Loading states
  const [loadingOccurrences, setLoadingOccurrences] = React.useState(false)
  const [loadingCreate, setLoadingCreate] = React.useState(false)
  const [loadingDelete, setLoadingDelete] = React.useState(false)
  const [loadingEdit, setLoadingEdit] = React.useState(false)

  // Error states
  const [createError, setCreateError] = React.useState<string | null>(null)
  const [occurrencesError, setOccurrencesError] = React.useState<string | null>(null)
  const [editError, setEditError] = React.useState<string | null>(null)

  // Dialog states
  const [deleteDialog, setDeleteDialog] = React.useState<{
    open: boolean
    occurrences: OccurrenceDTO[]
    selectedDate: Date | null
  }>({
    open: false,
    occurrences: [],
    selectedDate: null,
  })

  const [editDialog, setEditDialog] = React.useState<{
    open: boolean
    event: EventDTO | null
  }>({
    open: false,
    event: null,
  })

  // Event actions popover state
  const [eventActionsPopover, setEventActionsPopover] = React.useState<{
    open: boolean
    occurrences: OccurrenceDTO[]
    selectedDate: Date | null
  }>({
    open: false,
    occurrences: [],
    selectedDate: null,
  })

  // Snackbar state
  const { toast } = useToast()

  // Debounced filter
  const debouncedNameFilter = useDebounce(nameFilter, 300)

  // Abort controller for API calls
  const abortControllerRef = React.useRef<AbortController | null>(null)

  // Load occurrences
  const loadOccurrences = React.useCallback(async (month: Date, filter?: string) => {
    // Cancel previous request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    abortControllerRef.current = new AbortController()

    setLoadingOccurrences(true)
    setOccurrencesError(null)

    try {
      const start = format(startOfMonth(month), "yyyy-MM-dd")
      const end = format(endOfMonth(month), "yyyy-MM-dd")

      const response = await eventsApi.listOccurrences(
        { start, end, name: filter || undefined },
      )
    
      setOccurrences(response.data?.occurrences ?? [])
    } catch (error: any) {
      if (error.name !== "AbortError") {
        const apiError = error.response?.data as ApiError
        setOccurrencesError(apiError?.error || apiError?.message || "Error loading occurrences")
      }
    } finally {
      setLoadingOccurrences(false)
    }
  }, [])

  // Load events
  const loadEvents = React.useCallback(async () => {
    try {
      const response = await eventsApi.listEvents()
      setEvents(response.data?.events ?? []) 
      
    } catch (error: any) {
      console.error("Error loading events:", error)
    }
  }, [])

  // Effects
  React.useEffect(() => {
    loadOccurrences(currentMonth, debouncedNameFilter)
  }, [currentMonth, debouncedNameFilter, loadOccurrences])

  React.useEffect(() => {
    loadEvents()
  }, [loadEvents])

  // Handlers
  const handleCreateEvent = async (data: CreateEventPayload) => {
    setLoadingCreate(true)
    setCreateError(null)

    try {
      await eventsApi.createEvent(data)
      toast({
        title: "Success",
        description: "Event created successfully!",
      })

      // Reload data
      await Promise.all([loadOccurrences(currentMonth, debouncedNameFilter), loadEvents()])
    } catch (error: any) {
      const apiError = error.response?.data as ApiError
      setCreateError(apiError?.error || apiError?.message || "Error creating event")
      throw error
    } finally {
      setLoadingCreate(false)
    }
  }

  const handleMonthChange = (month: Date) => {
    setCurrentMonth(month)
  }

  const handleDayClick = (date: Date, dayOccurrences: OccurrenceDTO[]) => {
    setEventActionsPopover({
      open: true,
      occurrences: dayOccurrences,
      selectedDate: date,
    })
  }

  const handleDeleteConfirm = async (deleteType: "all" | "from") => {
    if (!deleteDialog.selectedDate || deleteDialog.occurrences.length === 0) return

    setLoadingDelete(true)

    try {
      const eventId = deleteDialog.occurrences[0].eventId

      if (deleteType === "all") {
        await eventsApi.deleteEvent(eventId)
      } else {
        const fromDate = format(deleteDialog.selectedDate, "yyyy-MM-dd")
        await eventsApi.deleteEventFrom(eventId, fromDate)
      }

      toast({
        title: "Success",
        description: "Event deleted successfully!",
      })

      setDeleteDialog({ open: false, occurrences: [], selectedDate: null })

      // Reload data
      await Promise.all([loadOccurrences(currentMonth, debouncedNameFilter), loadEvents()])
    } catch (error: any) {
      const apiError = error.response?.data as ApiError
      toast({
        title: "Error",
        description: apiError?.error || apiError?.message || "Error deleting event",
        variant: "destructive",
      })
    } finally {
      setLoadingDelete(false)
    }
  }

  const handleEditSubmit = async (data: UpdateEventPayload) => {
    if (!editDialog.event) return

    setLoadingEdit(true)
    setEditError(null)

    try {
      await eventsApi.updateEvent(editDialog.event.id, data)

      toast({
        title: "Success",
        description: "Event updated successfully!",
      })

      // Reload data
      await Promise.all([loadOccurrences(currentMonth, debouncedNameFilter), loadEvents()])
    } catch (error: any) {
      const apiError = error.response?.data as ApiError
      setEditError(apiError?.error || apiError?.message || "Error updating event")
      throw error
    } finally {
      setLoadingEdit(false)
    }
  }

  const handleDeleteFromPopover = (eventId: string) => {
    // Find the specific occurrence for this event
    const specificOccurrence = eventActionsPopover.occurrences.find((occ) => occ.eventId === eventId)
    if (specificOccurrence) {
      setDeleteDialog({
        open: true,
        occurrences: [specificOccurrence], // Only pass the specific occurrence
        selectedDate: eventActionsPopover.selectedDate,
      })
    }
  }

  const handleEditFromPopover = async (eventId: string) => {    
    const event = events.find((e) => e.id === eventId)
    
    if (event) {
      setEditDialog({
        open: true,
        event: event,
      })
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      {/* Modern Header */}
      <header className="glass-strong border-b sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-primary/10">
                <LucideCalendar className="h-6 w-6 text-primary" />
              </div>
              <div>
                <h1 className="text-2xl font-bold font-space-grotesk">Event Manager</h1>
                <p className="text-sm text-muted-foreground">Modern management of recurring events</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <ThemeToggle />
            </div>
          </div>
        </div>
      </header>

      {/* Loading Bar */}
      {(loadingCreate || loadingDelete || loadingEdit) && (
        <div className="h-1 bg-primary/20 overflow-hidden">
          <div className="h-full bg-primary animate-pulse" />
        </div>
      )}

      <main className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Event Form */}
          <div className="lg:col-span-4">
            <EventForm onSubmit={handleCreateEvent} loading={loadingCreate} error={createError} />
          </div>

          {/* Calendar and Filters */}
          <div className="lg:col-span-8">
            <div className="space-y-6">
              <Filters nameFilter={nameFilter} onNameFilterChange={setNameFilter} disabled={loadingOccurrences} />

              <EventCalendar
                currentMonth={currentMonth}
                onMonthChange={handleMonthChange}
                occurrences={occurrences}
                onDayClick={handleDayClick}
                loading={loadingOccurrences}
              />

              {occurrencesError && (
                <Alert variant="destructive">
                  <AlertDescription>{occurrencesError}</AlertDescription>
                </Alert>
              )}
            </div>
          </div>
        </div>
      </main>

      {/* Dialogs */}
      <EventActionsPopover
        open={eventActionsPopover.open}
        onClose={() => setEventActionsPopover({ open: false, occurrences: [], selectedDate: null })}
        occurrences={eventActionsPopover.occurrences}
        selectedDate={eventActionsPopover.selectedDate}
        onEdit={handleEditFromPopover}
        onDelete={handleDeleteFromPopover}
      />

      <DeleteDialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false, occurrences: [], selectedDate: null })}
        onConfirm={handleDeleteConfirm}
        occurrences={deleteDialog.occurrences}
        selectedDate={deleteDialog.selectedDate}
        loading={loadingDelete}
      />

      <EditEventDialog
        open={editDialog.open}
        onClose={() => setEditDialog({ open: false, event: null })}
        onSubmit={handleEditSubmit}
        event={editDialog.event}
        loading={loadingEdit}
        error={editError}
      />

      <Toaster />
    </div>
  )
}

export default App
