"use client"

import { use } from "react"
import { toast } from "sonner"
import { linksApi } from "@/api/links"
import { useRouter } from "next/navigation"
import { useEffect, useRef } from "react"

export default function RedirectPage({ params }: { params: Promise<{ slug: string }> }) {
    const router = useRouter()
    const { slug } = use(params)
    const redirectAttempted = useRef(false)

    useEffect(() => {
        if (redirectAttempted.current) return
        redirectAttempted.current = true

        const redirect = async () => {
            try {
                const response = await linksApi.getLink(slug)

                if (response.success && response.data) {
                    const parsedData = typeof response.data === 'string' ? JSON.parse(response.data) : response.data

                    await linksApi.updateLinkClicks(parsedData.id)
                    window.location.href = parsedData.original_url
                } else {
                    toast.error(response.message || "Link not found")
                    setTimeout(() => {
                        router.push("/")
                    }, 3000)
                }
            } catch (error) {
                console.error("Error redirecting:", error)
                toast.error("Failed to redirect")
                router.push("/")
            }
        }

        redirect()
    }, [slug])

    return (
        <div className="flex items-center justify-center min-h-screen">
            <div className="text-center space-y-4">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
                <p className="text-muted-foreground">Redirecting you to the destination...</p>
            </div>
        </div>
    )
} 