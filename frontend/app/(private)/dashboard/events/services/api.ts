import { apiRequest } from "@/api/api"
import type { EventDTO, CreateEventPayload, UpdateEventPayload, OccurrencesResponse } from "../types"


export const eventsApi = {
  createEvent: (payload: CreateEventPayload) =>
    apiRequest<EventDTO>({ method: "POST", endpoint: "/v1/events", body: payload }),

  listEvents: () =>
    apiRequest<{ events: EventDTO[] }>({ method: "GET", endpoint: "/v1/events" }),

  updateEvent: (id: string, payload: UpdateEventPayload) =>
    apiRequest<EventDTO>({ method: "PUT", endpoint: `/v1/events/${id}`, body: payload }),

  deleteEvent: (id: string) =>
    apiRequest<void>({ method: "DELETE", endpoint: `/v1/events/${id}` }),

  deleteEventFrom: (id: string, from: string) =>
    apiRequest<void>({
      method: "DELETE",
      endpoint: `/v1/events/${id}?from=${encodeURIComponent(from)}`,
    }),

  listOccurrences: (params: { start: string; end: string; name?: string }) => {
    const qs = new URLSearchParams()
    qs.set("start", params.start)
    qs.set("end", params.end)
    if (params.name && params.name.trim() !== "") qs.set("name", params.name.trim())

    return apiRequest<OccurrencesResponse>({
      method: "GET",
      endpoint: `/v1/events/occurrences?${qs.toString()}`,
    })
  },
}
