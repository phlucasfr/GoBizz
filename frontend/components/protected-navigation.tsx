'use client'

import Navbar from '@/components/navbar'
import Footer from '@/components/footer'

import { usePathname } from 'next/navigation'
import { publicRoutes } from '@/middleware'

export default function ProtectedNavigation({ children }: { children: React.ReactNode }) {
    const pathname = usePathname()
    const isPublicRoute = publicRoutes.some(route => route.path === pathname)

    return (
        <div className="min-h-screen flex flex-col">
            {isPublicRoute && <Navbar />}
            <main className="flex-grow">
                {children}
            </main>
            {isPublicRoute && <Footer />}
        </div>
    )
} 