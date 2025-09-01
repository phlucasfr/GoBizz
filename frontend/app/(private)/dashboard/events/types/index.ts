export interface EventDTO {
  id: string
  name: string
  start_date: string
  interval_days: number
  stop_at: string | null
}

export interface OccurrenceDTO {
  eventId: string
  name: string
  date: string
}

export interface CreateEventPayload {
  name: string
  startDate: string
  intervalDays: number
  stopAt?: string | null
}

export interface UpdateEventPayload {
  name?: string
  startDate?: string
  intervalDays?: number
  stopAt?: string | null
}

export interface OccurrencesResponse {
  period: {
    start: string
    end: string
  }
  occurrences: OccurrenceDTO[]
}

export interface ApiError {
  error?: string
  message?: string
}
