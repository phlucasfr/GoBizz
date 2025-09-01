"use client"

import React from "react"
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Box, Alert } from "@mui/material"
import { useForm, Controller } from "react-hook-form"
import type { EventDTO, UpdateEventPayload } from "../types"
import DeleteOutline from "@mui/icons-material/DeleteOutline"
import { InputAdornment, IconButton, Tooltip } from "@mui/material"

interface EditEventDialogProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: UpdateEventPayload) => Promise<void>
  event: EventDTO | null
  loading: boolean
  error: string | null
}

interface FormData {
  name: string
  startDate: string
  intervalDays: number
  stopAt: string
}

export const EditEventDialog: React.FC<EditEventDialogProps> = ({ open, onClose, onSubmit, event, loading, error }) => {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>()

  React.useEffect(() => {
    if (event && open) {
      reset({
        name: event.name,
        startDate: event.start_date,
        intervalDays: event.interval_days,
        stopAt: event.stop_at || "",
      })
    }
  }, [event, open, reset])

  const handleFormSubmit = async (data: FormData) => {
    const isEmpty = !data.stopAt || data.stopAt.trim() === ""

    const payload: UpdateEventPayload = {
      name: data.name,
      startDate: data.startDate,
      intervalDays: data.intervalDays,
      stopAt: isEmpty ? null : data.stopAt,
    }

    try {
      await onSubmit(payload)
      onClose()
    } catch {}
  }

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth aria-labelledby="edit-dialog-title">
      <DialogTitle id="edit-dialog-title">Edit Event</DialogTitle>

      <DialogContent>
        <Box
          component="form"
          onSubmit={handleSubmit(handleFormSubmit)}
          sx={{ display: "flex", flexDirection: "column", gap: 2, pt: 1 }}
          noValidate
        >
          {error && (
            <Alert severity="error" role="alert">
              {error}
            </Alert>
          )}

          <Controller
            name="name"
            control={control}
            rules={{
              required: "Event name is required",
              minLength: {
                value: 1,
                message: "Name cannot be empty",
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                label="Event Name"
                error={!!errors.name}
                helperText={errors.name?.message}
                disabled={loading}
                required
                fullWidth
              />
            )}
          />

          <Controller
            name="startDate"
            control={control}
            rules={{
              required: "Start date is required",
              pattern: {
                value: /^\d{4}-\d{2}-\d{2}$/,
                message: "Date must be in YYYY-MM-DD format",
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                label="Start Date"
                type="date"
                error={!!errors.startDate}
                helperText={errors.startDate?.message}
                disabled={loading}
                required
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            )}
          />

          <Controller
            name="intervalDays"
            control={control}
            rules={{
              required: "Interval in days is required",
              min: {
                value: 1,
                message: "Interval must be at least 1 day",
              },
              max: {
                value: 5000,
                message: "Interval cannot exceed 5000 days",
              },
            }}
            render={({ field: { onChange, value, ...field } }) => (
              <TextField
                {...field}
                value={value}
                onChange={(e) => onChange(Number.parseInt(e.target.value) || 0)}
                label="Recurrence in Days"
                type="number"
                error={!!errors.intervalDays}
                helperText={errors.intervalDays?.message}
                disabled={loading}
                required
                fullWidth
                inputProps={{ min: 1, max: 5000 }}
              />
            )}
          />

          <Controller
            name="stopAt"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                label="Stop Date (optional)"
                type="date"
                disabled={loading}
                fullWidth
                InputLabelProps={{ shrink: true }}
                helperText="Leave empty for endless event"
                InputProps={{
                  endAdornment: field.value ? (
                    <InputAdornment position="end">
                      <Tooltip title="Clear stop date">
                        <IconButton
                          aria-label="Clear stop date"
                          edge="end"
                          size="small"
                          onClick={() => {
                            field.onChange("")
                          }}
                        >
                          <DeleteOutline />
                        </IconButton>
                      </Tooltip>
                    </InputAdornment>
                  ) : null,
                }}
              />
            )}
          />

        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit(handleFormSubmit)}
          variant="contained"
          disabled={loading || isSubmitting}
          autoFocus
        >
          {loading ? "Saving..." : "Save"}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
