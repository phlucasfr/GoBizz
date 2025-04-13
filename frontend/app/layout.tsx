import "./globals.css"

import type React from "react"
import type { Metadata } from "next/types"

import ProtectedNavigation from "@/components/protected-navigation"

import { Inter } from "next/font/google"
import { Toaster } from "@/components/ui/sonner"
import { AuthProvider } from "@/context/auth-context"
import { ThemeProvider } from "@/components/theme-provider"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "GoBizz - Business Management Platform",
  description: "Intelligent tools and complete resources to drive your business success.",
  generator: 'v0.dev'
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <link rel="icon" href="/gobizz.svg" type="image/svg+xml" />
      <body className={inter.className}>
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
          <AuthProvider>
            <ProtectedNavigation>
              {children}
            </ProtectedNavigation>
          </AuthProvider>
          <Toaster />
        </ThemeProvider>
      </body>
    </html>
  )
}


