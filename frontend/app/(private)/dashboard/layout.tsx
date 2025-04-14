"use client"

import type React from "react"
import DashboardSidebar from "@/components/dashboard-sidebar"
import { redirect } from "next/navigation"
import { useEffect } from 'react'
import { verifyAuth } from "@/lib/auth"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  useEffect(() => {
    const token = localStorage.getItem("auth-token")
    if (!token) {
      redirect("/login")
    }

    verifyAuth(token).then((isValid) => {
      if (!isValid) redirect("/login")
    })
  }, [])

  return (
    <div className="flex min-h-screen">
      <DashboardSidebar />
      <div className="flex-1 p-8">{children}</div>
    </div>
  )
}
