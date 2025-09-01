"use client"

import type React from "react"
import { useState } from "react"
import { Plus, Calendar, Clock, StopCircle } from "lucide-react"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Alert, AlertDescription } from "@/components/ui/alert"
import type { CreateEventPayload } from "../types"

interface EventFormProps {
  onSubmit: (data: CreateEventPayload) => Promise<void>
  loading: boolean
  error: string | null
}

export const EventForm: React.FC<EventFormProps> = ({ onSubmit, loading, error }) => {
  const [formData, setFormData] = useState({
    name: "",
    startDate: "",
    intervalDays: 1,
    stopAt: "",
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.name.trim()) {
      newErrors.name = "Event name is required"
    }

    if (!formData.startDate) {
      newErrors.startDate = "Start date is required"
    }

    if (formData.intervalDays < 1) {
      newErrors.intervalDays = "Interval must be at least 1 day"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) return

    try {
      const payload: CreateEventPayload = {
        name: formData.name,
        startDate: formData.startDate,
        intervalDays: formData.intervalDays,
        stopAt: formData.stopAt || null,
      }
      await onSubmit(payload)
      setFormData({ name: "", startDate: "", intervalDays: 1, stopAt: "" })
      setErrors({})
    } catch (err) {
      // Error is handled by parent component
    }
  }

  const handleChange = (field: string, value: string | number) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }))
    }
  }

  return (
    <Card className="glass h-fit">
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 font-space-grotesk">
          <div className="p-1.5 rounded-md bg-primary/10">
            <Plus className="h-4 w-4 text-primary" />
          </div>
          Create New Event
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4" noValidate>
          {error && (
            <Alert variant="destructive" role="alert">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <div className="space-y-2">
            <Label htmlFor="name" className="text-sm font-medium">
              Event Name *
            </Label>
            <Input
              id="name"
              value={formData.name}
              onChange={(e) => handleChange("name", e.target.value)}
              placeholder="E.g.: Weekly meeting"
              disabled={loading}
              className={errors.name ? "border-destructive" : ""}
              aria-describedby={errors.name ? "name-error" : undefined}
            />
            {errors.name && (
              <p id="name-error" className="text-sm text-destructive mt-1" role="alert">
                {errors.name}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="startDate" className="text-sm font-medium flex items-center gap-2">
              <Calendar className="h-3 w-3" />
              Start Date *
            </Label>
            <Input
              id="startDate"
              type="date"
              value={formData.startDate}
              onChange={(e) => handleChange("startDate", e.target.value)}
              disabled={loading}
              className={errors.startDate ? "border-destructive" : ""}
              aria-describedby={errors.startDate ? "startDate-error" : undefined}
            />
            {errors.startDate && (
              <p id="startDate-error" className="text-sm text-destructive mt-1" role="alert">
                {errors.startDate}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="intervalDays" className="text-sm font-medium flex items-center gap-2">
              <Clock className="h-3 w-3" />
              Recurrence in Days *
            </Label>
            <Input
              id="intervalDays"
              type="number"
              value={formData.intervalDays}
              onChange={(e) => handleChange("intervalDays", Number.parseInt(e.target.value) || 0)}
              placeholder="E.g.: 7 (weekly)"
              disabled={loading}
              min={1}
              max={5000}
              className={errors.intervalDays ? "border-destructive" : ""}
              aria-describedby={errors.intervalDays ? "intervalDays-error" : undefined}
            />
            {errors.intervalDays && (
              <p id="intervalDays-error" className="text-sm text-destructive mt-1" role="alert">
                {errors.intervalDays}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="stopAt" className="text-sm font-medium flex items-center gap-2">
              <StopCircle className="h-3 w-3" />
              End Date (optional)
            </Label>
            <Input
              id="stopAt"
              type="date"
              value={formData.stopAt}
              onChange={(e) => handleChange("stopAt", e.target.value)}
              disabled={loading}
              placeholder="Leave empty for endless event"
            />
          </div>

          <Button type="submit" disabled={loading} className="w-full mt-6 font-medium">
            {loading ? "Creating..." : "Create Event"}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
