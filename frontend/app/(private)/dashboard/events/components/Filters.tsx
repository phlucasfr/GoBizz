"use client"

import type React from "react"
import { TextField, Box } from "@mui/material"
import { Search } from "@mui/icons-material"

interface FiltersProps {
  nameFilter: string
  onNameFilterChange: (value: string) => void
  disabled?: boolean
}

export const Filters: React.FC<FiltersProps> = ({ nameFilter, onNameFilterChange, disabled = false }) => {
  return (
    <Box sx={{ mb: 2 }}>
      <TextField
        value={nameFilter}
        onChange={(e) => onNameFilterChange(e.target.value)}
        label="Filter by event name"
        placeholder="Type to filter..."
        disabled={disabled}
        fullWidth
        InputProps={{
          startAdornment: <Search sx={{ mr: 1, color: "text.secondary" }} />,
        }}
        inputProps={{
          "aria-label": "Filter events by name",
        }}
      />
    </Box>
  )
}
