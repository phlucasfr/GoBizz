"use client"

import { toast } from "sonner"
import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/context/auth-context"
import { LinkCharts } from "./components/LinkCharts"
import { LinkReportPDF } from "./components/LinkReportPDF"
import { linksApi, Link } from "@/api/links"
import { useEffect, useState, useRef } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Activity, Link as LinkIcon, Loader2, RefreshCw, AlertCircle } from "lucide-react"

export default function Dashboard() {
  const { userSet } = useAuth()

  const [links, setLinks] = useState<Link[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const hasFetchedLinks = useRef(false)

  useEffect(() => {
    if (userSet && !hasFetchedLinks.current) {
      hasFetchedLinks.current = true
      fetchLinks()
    }
  }, [userSet])

  const fetchLinks = async () => {
    setIsLoading(true)
    try {
      const response = await linksApi.getCustomerLinks({ customerId: userSet!.id })
      if (response.success && response.data) {
        const linksData = Array.isArray(response.data.data) ? response.data.data : []
        setLinks(linksData)
      }
    } catch (error) {
      console.error("Error fetching links:", error)
      toast.error("Failed to fetch links")
      setLinks([])
    } finally {
      setIsLoading(false)
    }
  }

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  }

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 100,
        damping: 15,
      },
    },
  }

  const totalClicks = links.reduce((sum, link) => sum + (link.clicks || 0), 0)
  const mostClickedLinks = [...links]
    .sort((a, b) => (b.clicks || 0) - (a.clicks || 0))
    .slice(0, 5)
  const recentLinks = [...links]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)
  const expiredLinks = links.filter(link => link.expiration_date && new Date(link.expiration_date) < new Date())
  const activeLinks = links.filter(link => !link.expiration_date || new Date(link.expiration_date) >= new Date())

  return (
    <motion.div initial="hidden" animate="visible" variants={containerVariants} className="space-y-4 md:space-y-8 relative">
      {isLoading && (
        <div className="absolute inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center">
          <div className="flex flex-col items-center gap-2">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <span className="text-sm text-muted-foreground">Updating your links data...</span>
          </div>
        </div>
      )}
      <motion.div variants={itemVariants} className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold">Welcome, {userSet?.companyName || "User"}</h1>
          <p className="text-muted-foreground text-sm md:text-base">Here's an overview of your links performance</p>
        </div>
        <div className="flex gap-2">
          <LinkReportPDF links={links} />
          <Button
            variant="outline"
            size="icon"
            onClick={fetchLinks}
            disabled={isLoading}
            className="h-10 w-10"
          >
            <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </motion.div>

      {links.length === 0 ? (
        <div className="flex justify-center items-center py-10 md:py-20">
          <p className="text-muted-foreground">No links found. Create your first link to get started!</p>
        </div>
      ) : (
        <>
          <motion.div variants={itemVariants} className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Total Links</CardTitle>
                <LinkIcon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{links.length}</div>
                <p className="text-xs text-muted-foreground">All your shortened links</p>
              </CardContent>
            </Card>

            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Active Links</CardTitle>
                <LinkIcon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{activeLinks.length}</div>
                <p className="text-xs text-muted-foreground">Currently active links</p>
              </CardContent>
            </Card>

            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Expired Links</CardTitle>
                <AlertCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{expiredLinks.length}</div>
                <p className="text-xs text-muted-foreground">Links that have expired</p>
              </CardContent>
            </Card>

            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Total Clicks</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{totalClicks}</div>
                <p className="text-xs text-muted-foreground">Across all your links</p>
              </CardContent>
            </Card>
          </motion.div>

          <motion.div variants={itemVariants}>
            <LinkCharts links={links} />
          </motion.div>

          <motion.div variants={itemVariants} className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle>Most Clicked Links</CardTitle>
                <CardDescription>Your top performing links</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {mostClickedLinks.map((link) => (
                    <div key={link.id} className="flex items-center justify-between">
                      <div className="flex-1 truncate mr-4">
                        <a
                          href={link.short_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-sm hover:text-primary"
                        >
                          {link.short_url}
                        </a>
                      </div>
                      <div className="text-sm font-medium">{link.clicks || 0} clicks</div>
                    </div>
                  ))}
                  {mostClickedLinks.length === 0 && (
                    <p className="text-muted-foreground text-sm">No links with clicks yet</p>
                  )}
                </div>
              </CardContent>
            </Card>

            <Card className="bg-card hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle>Recent Links</CardTitle>
                <CardDescription>Your most recently created links</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {recentLinks.map((link) => (
                    <div key={link.id} className="flex items-center justify-between">
                      <div className="flex-1 truncate mr-4">
                        <a
                          href={link.short_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-sm hover:text-primary"
                        >
                          {link.short_url}
                        </a>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {new Date(link.created_at).toLocaleDateString()}
                      </div>
                    </div>
                  ))}
                  {recentLinks.length === 0 && (
                    <p className="text-muted-foreground text-sm">No links created yet</p>
                  )}
                </div>
              </CardContent>
            </Card>
          </motion.div>
        </>
      )}
    </motion.div>
  )
}