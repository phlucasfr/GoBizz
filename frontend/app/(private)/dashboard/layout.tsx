import type React from "react"

import DashboardSidebar from "@/components/dashboard-sidebar"

import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { verifyAuth } from "@/lib/auth"

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth-token")?.value

  if (!token || !(await verifyAuth(token))) {
    redirect("/login")
  }

  return (
    <div className="flex min-h-screen">
      <DashboardSidebar />
      <div className="flex-1 p-8">{children}</div>
    </div>
  )
}

