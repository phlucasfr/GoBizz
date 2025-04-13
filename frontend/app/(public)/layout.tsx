import { ReactNode } from "react"
import ClientLayout from "@/components/client-layout"

export default function PublicLayout({ children }: { children: ReactNode }) {
    return <ClientLayout>{children}</ClientLayout>
} 