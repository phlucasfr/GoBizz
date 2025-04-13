"use client"

import Link from "next/link"
import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { deleteCookie, getCookie, useAuth } from "@/context/auth-context"
import { ThemeToggle } from "@/components/theme-toggle"
import { usePathname } from "next/navigation"
import { ConfirmDialog } from "@/components/confirm-dialog"
import { useState, useEffect, useRef } from "react"
import {
  LayoutDashboard,
  Link2,
  LogOut,
  ChevronLeft,
  ChevronRight,
} from "lucide-react"

const sidebarLinks = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    name: "Links",
    href: "/dashboard/links",
    icon: Link2,
  },
]

export default function DashboardSidebar() {
  const pathname = usePathname()
  const isFirstLoad = useRef(true)

  const { logout, refreshSession } = useAuth()
  const [sessionTime, setSessionTime] = useState(1800)
  const [isCollapsed, setIsCollapsed] = useState(false)
  const [showLogoutDialog, setShowLogoutDialog] = useState(false)

  const refreshToken = async () => {
    try {
      const userData = getCookie("user-data")
      const user = JSON.parse(userData)

      if (!user!.id) throw new Error("User id not found")
      if (!user!.email) throw new Error("User email not found")

      await refreshSession(user!.id, user!.email)
    } catch (error) {
      console.error("Failed to refresh session:", error)
      deleteCookie('user-data')

      window.location.href = "/login"
    }
  }

  useEffect(() => {
    const checkSession = () => {
      const token = getCookie("auth-token")
      if (token) {
        const { exp } = JSON.parse(atob(token.split(".")[1]))
        const timeLeft = exp - Math.floor(Date.now() / 1000)
        setSessionTime(Math.max(timeLeft, 0))
        if (timeLeft <= 0) {
          deleteCookie('user-data')
          window.location.href = "/login"
        }
      } else {
        deleteCookie('user-data')
        window.location.href = "/login"
      }
    }

    if (isFirstLoad.current) {
      isFirstLoad.current = false
      checkSession()
    }
  }, [refreshSession])

  useEffect(() => {
    const interval = setInterval(() => {
      setSessionTime((prevTime) => Math.max(prevTime - 1, 0))
    }, 1000)

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    const handleUserInteraction = () => {
      if (sessionTime > 0 && sessionTime < 300) {
        refreshToken()
      }
    }

    window.addEventListener("click", handleUserInteraction)

    return () => {
      window.removeEventListener("click", handleUserInteraction)
    }
  }, [sessionTime, refreshSession])

  const handleLogout = () => {
    setShowLogoutDialog(true)
  }

  const handleConfirmLogout = () => {
    logout()
  }

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${minutes}:${secs.toString().padStart(2, "0")}`
  }

  return (
    <>
      <motion.div
        className="bg-card border-r h-screen sticky top-0 flex flex-col transition-all duration-300"
        initial={{ width: isCollapsed ? 80 : 256 }}
        animate={{ width: isCollapsed ? 80 : 256 }}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
              <img src="/logo.png" alt="GoBizz Logo" className="w-8 h-8" />
            </div>
            {!isCollapsed && <span className="font-bold text-lg">GoBizz</span>}
          </div>
          <button
            onClick={() => setIsCollapsed((prev) => !prev)}
            className="p-2 rounded hover:bg-secondary"
          >
            {isCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 space-y-1">
          {sidebarLinks.map((link) => {
            const isActive = pathname === link.href
            const Icon = link.icon

            return (
              <Link
                key={link.href}
                href={link.href}
                className={`flex items-center space-x-2 px-3 py-2 rounded-md transition-colors relative ${isActive ? "bg-primary text-primary-foreground" : "hover:bg-secondary"
                  }`}
              >
                <Icon size={20} />
                {!isCollapsed && <span>{link.name}</span>}
                {isActive && (
                  <motion.div
                    className="absolute left-0 w-1 h-6 bg-primary rounded-r-full"
                    layoutId="sidebar-indicator"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.2 }}
                  />
                )}
              </Link>
            )
          })}
        </nav>

        {/* Footer */}
        <div className="mt-auto pt-8 border-t">
          {!isCollapsed && (
            <div className="flex items-center justify-between px-2 mb-4">
              <span className="text-sm text-muted-foreground">Theme</span>
              <ThemeToggle />
            </div>
          )}

          {!isCollapsed && (
            <div className="flex items-center justify-between px-2 mb-4">
              <span className="text-sm text-muted-foreground">
                Session Time: {formatTime(sessionTime)}
              </span>
            </div>
          )}

          <Button
            variant="outline"
            className={`w-full flex items-center space-x-2 justify-center text-destructive hover:text-destructive ${isCollapsed ? "justify-center" : ""
              }`}
            onClick={handleLogout}
          >
            <LogOut size={18} />
            {!isCollapsed && <span>Sign out</span>}
          </Button>
        </div>
      </motion.div>

      <ConfirmDialog
        open={showLogoutDialog}
        onOpenChange={setShowLogoutDialog}
        onConfirm={handleConfirmLogout}
        title="Sign out"
        description="Are you sure you want to sign out? You will need to sign in again to access your account."
      />
    </>
  )
}