"use client"

import { toast } from "sonner"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/context/auth-context"
import { linksApi, type Link } from "@/api/links"
import { useState, useEffect, useRef } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Link2, Copy, Trash2, ExternalLink, CheckCircle, RefreshCw, Pencil, Search, ArrowUp, ArrowDown, ChevronLeft, ChevronRight } from "lucide-react"
import {
  Dialog,
  DialogTitle,
  DialogFooter,
  DialogHeader,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog"
import {
  Select,
  SelectItem,
  SelectValue,
  SelectTrigger,
  SelectContent,
} from "@/components/ui/select"

export default function LinksPage() {
  const { userSet } = useAuth()

  const isFetching = useRef(false)

  const [itemsPerPage] = useState(10)
  const [links, setLinks] = useState<Link[]>([])
  const [longUrl, setLongUrl] = useState("")
  const [copiedId, setCopiedId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [sortField, setSortField] = useState<"clicks" | "created_at" | "original_url" | "expiration_date" | null>(null)
  const [totalItems, setTotalItems] = useState(0)
  const [customSlug, setCustomSlug] = useState("")
  const [totalPages, setTotalPages] = useState(0)
  const [editingLink, setEditingLink] = useState<Link | null>(null)
  const [editLongUrl, setEditLongUrl] = useState("")
  const [searchQuery, setSearchQuery] = useState("")
  const [currentPage, setCurrentPage] = useState(1)
  const [isShortening, setIsShortening] = useState(false)
  const [filterStatus, setFilterStatus] = useState<"all" | "active" | "expired">("all")
  const [sortDirection, setSortDirection] = useState<"asc" | "desc">("desc")
  const [expirationDate, setExpirationDate] = useState("")
  const [editCustomSlug, setEditCustomSlug] = useState("")
  const [filterSlugType, setFilterSlugType] = useState<"all" | "custom" | "auto">("all")
  const [editExpirationDate, setEditExpirationDate] = useState("")

  const fetchLinks = async () => {
    if (!userSet?.id || isFetching.current) return

    isFetching.current = true
    try {
      setIsLoading(true)
      const response = await linksApi.getCustomerLinks({
        customerId: userSet.id,
        limit: itemsPerPage,
        offset: (currentPage - 1) * itemsPerPage,
        sort_by: sortField || undefined,
        sort_direction: sortDirection,
        search: searchQuery,
        status: filterStatus !== 'all' ? filterStatus : undefined,
        slug_type: filterSlugType !== 'all' ? filterSlugType : undefined
      })

      if (response?.success && response?.data) {
        const { data, total } = response.data
        setLinks(Array.isArray(data) ? data : [])
        setTotalItems(total || 0)
        setTotalPages(Math.ceil((total || 0) / itemsPerPage))
      } else {
        console.warn('Invalid response format:', response)
        setLinks([])
        setTotalItems(0)
        setTotalPages(0)
      }
    } catch (error) {
      console.error("Error fetching links:", error)
      toast.error("Failed to fetch links")
      setLinks([])
      setTotalItems(0)
      setTotalPages(0)
    } finally {
      setIsLoading(false)
      isFetching.current = false
    }
  }

  useEffect(() => {
    const debounceTimer = setTimeout(() => {
      if (userSet?.id && !isFetching.current) {
        fetchLinks()
      }
    }, 300)

    return () => clearTimeout(debounceTimer)
  }, [currentPage, sortField, sortDirection, searchQuery, filterStatus, filterSlugType])

  const validateUrl = (url: string) => {
    try {
      new URL(url)
      return true
    } catch {
      return false
    }
  }

  const validateCustomSlug = (slug: string) => {
    const regex = /^[a-zA-Z0-9_-]+$/
    return regex.test(slug)
  }

  const formatDateForBackend = (dateString: string) => {
    if (!dateString) return undefined
    const date = new Date(dateString)
    return date.toISOString()
  }

  const handleShortenLink = async () => {
    if (!longUrl || !userSet?.id) return

    if (!validateUrl(longUrl)) {
      toast.error("Please enter a valid URL")
      return
    }

    if (customSlug && !validateCustomSlug(customSlug)) {
      toast.error("Custom slug can only contain letters, numbers, underscores, and hyphens")
      return
    }

    if (expirationDate && new Date(expirationDate) < new Date()) {
      toast.error("Expiration date must be in the future")
      return
    }

    setIsShortening(true)
    try {
      const response = await linksApi.createLink({
        original_url: longUrl,
        custom_slug: customSlug || undefined,
        customer_id: userSet.id,
        expiration_date: formatDateForBackend(expirationDate),
      })

      if (response.success) {
        setLongUrl("")
        setCustomSlug("")
        setExpirationDate("")
        toast.success("Link created successfully")

        await fetchLinks()
      } else {
        const errorMessage = response.message || "Failed to create link. Please try again."
        toast.error(errorMessage)
      }
    } catch (error: any) {
      console.error("Error creating link:", error)

      if (error?.message?.includes("custom slug already exists") ||
        error?.message?.includes("rpc error: code = Unknown desc = custom slug already exists")) {
        toast.error("This custom URL is already taken. Please choose a different one.")
      } else {
        toast.error("Failed to create link. Please try again.")
      }
    } finally {
      setIsShortening(false)
    }
  }

  const handleCopyLink = (id: string, url: string) => {
    navigator.clipboard.writeText(url)
    setCopiedId(id)
    toast.success("Link copied to clipboard")

    setTimeout(() => {
      setCopiedId(null)
    }, 2000)
  }

  const handleDeleteLink = async (id: string) => {
    try {
      const response = await linksApi.deleteLink(id)
      if (response.success) {
        setLinks(prevLinks => prevLinks.filter((link) => link.id !== id))
        toast.success("Link deleted successfully")
      } else {
        toast.error(response.message || "Failed to delete link")
      }
    } catch (error) {
      console.error("Error deleting link:", error)
      toast.error("Failed to delete link")
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    }).format(date)
  }

  const truncateUrl = (url: string, maxLength = 50) => {
    return url.length > maxLength ? `${url.substring(0, maxLength)}...` : url
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

  const handleEditLink = async () => {
    if (!editingLink || !editLongUrl) return

    if (!validateUrl(editLongUrl)) {
      toast.error("Please enter a valid URL")
      return
    }

    if (editCustomSlug && !validateCustomSlug(editCustomSlug)) {
      toast.error("Custom slug can only contain letters, numbers, underscores, and hyphens")
      return
    }

    if (editExpirationDate && new Date(editExpirationDate) < new Date()) {
      toast.error("Expiration date must be in the future")
      return
    }

    try {
      const response = await linksApi.updateLink({
        id: editingLink.id,
        original_url: editLongUrl,
        custom_slug: editCustomSlug || undefined,
        expiration_date: editExpirationDate ? formatDateForBackend(editExpirationDate) : undefined,
      })

      if (response.success) {
        setEditingLink(null)
        setEditLongUrl("")
        setEditCustomSlug("")
        setEditExpirationDate("")
        toast.success("Link updated successfully")
        await fetchLinks()
      } else {
        toast.error(response.message || "Failed to update link")
      }
    } catch (error) {
      console.error("Error updating link:", error)
      toast.error("Failed to update link")
    }
  }

  const openEditModal = (link: Link) => {
    setEditingLink(link)
    setEditLongUrl(link.original_url)
    setEditCustomSlug(link.custom_slug || "")
    setEditExpirationDate(link.expiration_date ? new Date(link.expiration_date).toISOString().slice(0, 16) : "")
  }

  const isLinkExpired = (link: Link): boolean => {
    return Boolean(link.expiration_date && new Date(link.expiration_date) < new Date())
  }

  const handleSort = (field: "clicks" | "created_at" | "original_url" | "expiration_date") => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc")
    } else {
      setSortField(field)
      setSortDirection("desc")
    }
  }

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page)
    }
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={containerVariants} className="space-y-8">
      <motion.div variants={itemVariants}>
        <h1 className="text-3xl font-bold">Links</h1>
        <p className="text-muted-foreground">Create and manage your shortened links</p>
      </motion.div>

      <motion.div variants={itemVariants}>
        <Card>
          <CardHeader>
            <CardTitle>Create Link</CardTitle>
            <CardDescription>Create and track shortened links for your marketing campaigns</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4">
              <div>
                <Label htmlFor="long-url">Paste your long URL</Label>
                <Input
                  id="long-url"
                  placeholder="https://example.com/your-long-url"
                  value={longUrl}
                  onChange={(e) => setLongUrl(e.target.value)}
                  className={!longUrl || validateUrl(longUrl) ? "" : "border-red-500"}
                />
                {longUrl && !validateUrl(longUrl) && (
                  <p className="text-sm text-red-500 mt-1">Please enter a valid URL</p>
                )}
              </div>

              <div>
                <Label htmlFor="custom-slug">Customize your link (optional)</Label>
                <div className="flex items-center">
                  <div className="bg-muted px-3 py-2 rounded-l-md border-y border-l">gobizz.co/</div>
                  <Input
                    id="custom-slug"
                    className={`rounded-l-none ${customSlug && !validateCustomSlug(customSlug) ? "border-red-500" : ""}`}
                    placeholder="custom-name"
                    value={customSlug}
                    onChange={(e) => setCustomSlug(e.target.value)}
                  />
                </div>
                {customSlug && !validateCustomSlug(customSlug) && (
                  <p className="text-sm text-red-500 mt-1">Only letters, numbers, underscores, and hyphens are allowed</p>
                )}
              </div>

              <div>
                <Label htmlFor="expiration-date">Expiration Date (optional)</Label>
                <Input
                  id="expiration-date"
                  type="datetime-local"
                  value={expirationDate}
                  onChange={(e) => setExpirationDate(e.target.value)}
                  min={new Date(new Date().setHours(0, 0, 0, 0)).toISOString().slice(0, 16)}
                />
              </div>

              <Button
                onClick={handleShortenLink}
                disabled={!longUrl || isShortening || !validateUrl(longUrl) || (customSlug ? !validateCustomSlug(customSlug) : false)}
                className="w-full md:w-auto md:self-end"
              >
                {isShortening ? (
                  <>
                    <svg
                      className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Shortening...
                  </>
                ) : (
                  <>
                    <Link2 className="mr-2 h-4 w-4" />
                    Shorten URL
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>
      </motion.div>

      <motion.div variants={itemVariants}>
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Your Links</CardTitle>
                <CardDescription>Manage and track your shortened links</CardDescription>
              </div>
              <div className="flex items-center gap-2">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="Search links..."
                    value={searchQuery}
                    onChange={(e) => {
                      const value = e.target.value
                      setSearchQuery(value)
                      setCurrentPage(1)
                    }}
                    className="pl-9 w-[200px]"
                  />
                </div>
                <Select value={filterStatus} onValueChange={(value: "all" | "active" | "expired") => {
                  setFilterStatus(value)
                  setCurrentPage(1)
                }}>
                  <SelectTrigger className="w-[120px]">
                    <SelectValue placeholder="Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="active">Active</SelectItem>
                    <SelectItem value="expired">Expired</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={filterSlugType} onValueChange={(value: "all" | "custom" | "auto") => {
                  setFilterSlugType(value)
                  setCurrentPage(1)
                }}>
                  <SelectTrigger className="w-[120px]">
                    <SelectValue placeholder="Slug Type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Slugs</SelectItem>
                    <SelectItem value="custom">Custom</SelectItem>
                    <SelectItem value="auto">Auto</SelectItem>
                  </SelectContent>
                </Select>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={fetchLinks}
                  title="Refresh links"
                >
                  <RefreshCw className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead
                      className="cursor-pointer hover:text-primary"
                      onClick={() => {
                        handleSort("original_url")
                        setCurrentPage(1)
                      }}
                    >
                      <div className="flex items-center">
                        Original URL
                        {sortField === "original_url" && (
                          sortDirection === "asc" ?
                            <ArrowUp className="ml-1 h-4 w-4" /> :
                            <ArrowDown className="ml-1 h-4 w-4" />
                        )}
                      </div>
                    </TableHead>
                    <TableHead>Shortened URL</TableHead>
                    <TableHead
                      className="text-center cursor-pointer hover:text-primary"
                      onClick={() => {
                        handleSort("clicks")
                        setCurrentPage(1)
                      }}
                    >
                      <div className="flex items-center justify-center">
                        Clicks
                        {sortField === "clicks" && (
                          sortDirection === "asc" ?
                            <ArrowUp className="ml-1 h-4 w-4" /> :
                            <ArrowDown className="ml-1 h-4 w-4" />
                        )}
                      </div>
                    </TableHead>
                    <TableHead
                      className="cursor-pointer hover:text-primary"
                      onClick={() => {
                        handleSort("created_at")
                        setCurrentPage(1)
                      }}
                    >
                      <div className="flex items-center">
                        Created
                        {sortField === "created_at" && (
                          sortDirection === "asc" ?
                            <ArrowUp className="ml-1 h-4 w-4" /> :
                            <ArrowDown className="ml-1 h-4 w-4" />
                        )}
                      </div>
                    </TableHead>
                    <TableHead
                      className="cursor-pointer hover:text-primary"
                      onClick={() => {
                        handleSort("expiration_date")
                        setCurrentPage(1)
                      }}
                    >
                      <div className="flex items-center">
                        Expires
                        {sortField === "expiration_date" && (
                          sortDirection === "asc" ?
                            <ArrowUp className="ml-1 h-4 w-4" /> :
                            <ArrowDown className="ml-1 h-4 w-4" />
                        )}
                      </div>
                    </TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-8">
                        <div className="flex items-center justify-center">
                          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                          <span className="ml-2">Loading links...</span>
                        </div>
                      </TableCell>
                    </TableRow>
                  ) : links.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                        No links found matching your criteria.
                      </TableCell>
                    </TableRow>
                  ) : (
                    links.map((link) => (
                      <TableRow key={link.id} className={isLinkExpired(link) ? "opacity-50" : ""}>
                        <TableCell className="font-medium max-w-[200px] truncate" title={link.original_url}>
                          {truncateUrl(link.original_url)}
                        </TableCell>
                        <TableCell>
                          <a
                            href={link.short_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className={`flex items-center hover:text-primary ${isLinkExpired(link) ? "pointer-events-none" : ""}`}
                          >
                            {link.short_url}
                            <ExternalLink className="ml-1 h-3 w-3" />
                          </a>
                        </TableCell>
                        <TableCell className="text-center">{link.clicks}</TableCell>
                        <TableCell>{formatDate(link.created_at)}</TableCell>
                        <TableCell>
                          {link.expiration_date ? (
                            <div className="flex items-center gap-2">
                              {formatDate(link.expiration_date)}
                              {isLinkExpired(link) && (
                                <span className="text-xs text-red-500">(Expired)</span>
                              )}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">Never</span>
                          )}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex items-center justify-end space-x-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleCopyLink(link.id, link.short_url)}
                              title="Copy link"
                              disabled={isLinkExpired(link)}
                            >
                              {copiedId === link.id ? (
                                <CheckCircle className="h-4 w-4 text-green-500" />
                              ) : (
                                <Copy className="h-4 w-4" />
                              )}
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => openEditModal(link)}
                              title="Edit link"
                            >
                              <Pencil className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleDeleteLink(link.id)}
                              title="Delete link"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            {/* Pagination Controls */}
            {!isLoading && totalItems > 0 && (
              <div className="flex items-center justify-between mt-4">
                <div className="text-sm text-muted-foreground">
                  Showing {Math.min((currentPage - 1) * itemsPerPage + 1, totalItems)} to{" "}
                  {Math.min(currentPage * itemsPerPage, totalItems)} of{" "}
                  {totalItems} links
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <div className="flex items-center gap-1">
                    {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                      <Button
                        key={page}
                        variant={currentPage === page ? "default" : "outline"}
                        size="icon"
                        onClick={() => handlePageChange(page)}
                        className="w-8 h-8"
                      >
                        {page}
                      </Button>
                    ))}
                  </div>
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </motion.div>

      <Dialog open={!!editingLink} onOpenChange={() => setEditingLink(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Link</DialogTitle>
            <DialogDescription>
              Update your shortened link's destination URL and custom slug
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div>
              <Label htmlFor="edit-long-url">Destination URL</Label>
              <Input
                id="edit-long-url"
                placeholder="https://example.com/your-long-url"
                value={editLongUrl}
                onChange={(e) => setEditLongUrl(e.target.value)}
                className={!editLongUrl || validateUrl(editLongUrl) ? "" : "border-red-500"}
              />
              {editLongUrl && !validateUrl(editLongUrl) && (
                <p className="text-sm text-red-500 mt-1">Please enter a valid URL</p>
              )}
            </div>

            <div>
              <Label htmlFor="edit-custom-slug">Custom Slug (optional)</Label>
              <div className="flex items-center">
                <div className="bg-muted px-3 py-2 rounded-l-md border-y border-l">gobizz.co/</div>
                <Input
                  id="edit-custom-slug"
                  className={`rounded-l-none ${editCustomSlug && !validateCustomSlug(editCustomSlug) ? "border-red-500" : ""}`}
                  placeholder="custom-name"
                  value={editCustomSlug}
                  onChange={(e) => setEditCustomSlug(e.target.value)}
                />
              </div>
              {editCustomSlug && !validateCustomSlug(editCustomSlug) && (
                <p className="text-sm text-red-500 mt-1">Only letters, numbers, underscores, and hyphens are allowed</p>
              )}
            </div>

            <div>
              <Label htmlFor="edit-expiration-date">Expiration Date (optional)</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="edit-expiration-date"
                  type="datetime-local"
                  value={editExpirationDate}
                  onChange={(e) => setEditExpirationDate(e.target.value)}
                  min={new Date(new Date().setHours(0, 0, 0, 0)).toISOString().slice(0, 16)}
                  className="flex-1"
                />
                {editingLink?.expiration_date && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setEditExpirationDate("")}
                    className="shrink-0"
                  >
                    Remove expiration
                  </Button>
                )}
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditingLink(null)}>
              Cancel
            </Button>
            <Button
              onClick={handleEditLink}
              disabled={!editLongUrl || !validateUrl(editLongUrl) || (editCustomSlug ? !validateCustomSlug(editCustomSlug) : false)}
            >
              Save Changes
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </motion.div>
  )
}