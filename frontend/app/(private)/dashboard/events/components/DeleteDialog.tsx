"use client"

import React from "react"
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  RadioGroup,
  FormControlLabel,
  Radio,
  FormControl,
  FormLabel,
} from "@mui/material"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import type { OccurrenceDTO } from "../types"

interface DeleteDialogProps {
  open: boolean
  onClose: () => void
  onConfirm: (deleteType: "all" | "from") => void
  occurrences: OccurrenceDTO[]
  selectedDate: Date | null
  loading: boolean
}

export const DeleteDialog: React.FC<DeleteDialogProps> = ({
  open,
  onClose,
  onConfirm,
  occurrences,
  selectedDate,
  loading,
}) => {
  const [deleteType, setDeleteType] = React.useState<"all" | "from">("from")

  React.useEffect(() => {
    if (open) {
      setDeleteType("from")
    }
  }, [open])

  const handleConfirm = () => {
    onConfirm(deleteType)
  }

  const eventName = occurrences[0]?.name || "Event"
  const formattedDate = selectedDate ? format(selectedDate, "dd/MM/yyyy", { locale: ptBR }) : ""

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      aria-labelledby="delete-dialog-title"
      aria-describedby="delete-dialog-description"
    >
      <DialogTitle id="delete-dialog-title">What do you want to delete?</DialogTitle>

      <DialogContent>
        <Typography id="delete-dialog-description" sx={{ mb: 2 }}>
          Event: <strong>{eventName}</strong>
          <br />
          Selected date: <strong>{formattedDate}</strong>
        </Typography>

        <FormControl component="fieldset">
          <FormLabel component="legend">Delete options:</FormLabel>
          <RadioGroup
            value={deleteType}
            onChange={(e) => setDeleteType(e.target.value as "all" | "from")}
            aria-labelledby="delete-dialog-title"
          >
            <FormControlLabel
              value="from"
              control={<Radio />}
              label={`Delete all future occurrences from ${formattedDate} (inclusive)`}
            />
            <FormControlLabel value="all" control={<Radio />} label="Delete the entire event" />
          </RadioGroup>
        </FormControl>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button onClick={handleConfirm} variant="contained" color="error" disabled={loading} autoFocus>
          {loading ? "Deleting..." : "Confirm Deletion"}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
