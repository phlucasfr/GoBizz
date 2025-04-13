import { Metadata } from "next"

export const metadata: Metadata = {
    title: 'Login - GoBizz',
    description: 'Sign in to your GoBizz account',
}

export default function LoginLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return children
} 