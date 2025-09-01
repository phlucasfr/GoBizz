"use client"

import type React from "react"
import { Edit, Trash2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import type { OccurrenceDTO } from "../types"

interface EventActionsPopoverProps {
  open: boolean
  onClose: () => void
  occurrences: OccurrenceDTO[]
  selectedDate: Date | null
  onEdit: (eventId: string) => void
  onDelete: (eventId: string) => void
}

export const EventActionsPopover: React.FC<EventActionsPopoverProps> = ({
  open,
  onClose,
  occurrences,
  selectedDate,
  onEdit,
  onDelete,
}) => {
  if (!open || !selectedDate || occurrences.length === 0) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <Card className="glass w-80 max-w-[90vw]" onClick={(e) => e.stopPropagation()}>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-space-grotesk">
            {format(selectedDate, "d 'de' MMMM", { locale: ptBR })}
          </CardTitle>
          <p className="text-sm text-muted-foreground">{occurrences.length} event(s) on this day</p>
        </CardHeader>

        <CardContent className="space-y-3">
          {occurrences.map((occurrence) => (
            <div key={occurrence.eventId} className="p-3 rounded-lg bg-muted/30 border border-border/20">
              {/* <h4 className="font-semibold text-lg text-foreground mb-3 font-space-grotesk">{occurrence.eventName}</h4> */}
              <h4 className="font-medium text-foreground mb-2">{occurrence.name}</h4>

              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    onEdit(occurrence.eventId)
                    onClose()
                  }}
                  className="flex-1 gap-2 text-cyan-600 border-cyan-200 hover:bg-cyan-50 dark:text-cyan-400 dark:border-cyan-800 dark:hover:bg-cyan-950/30"
                >
                  <Edit className="h-4 w-4" />
                  Edit
                </Button>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    onDelete(occurrence.eventId)
                    onClose()
                  }}
                  className="flex-1 gap-2 text-red-600 border-red-200 hover:bg-red-50 dark:text-red-400 dark:border-red-800 dark:hover:bg-red-950/30"
                >
                  <Trash2 className="h-4 w-4" />
                  Delete
                </Button>
              </div>
            </div>
          ))}

          <Button variant="ghost" onClick={onClose} className="w-full mt-4">
            Cancel
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
