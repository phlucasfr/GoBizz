"use client"

import type React from "react"

import { motion } from "framer-motion"
import { useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { MapPin, Phone, Mail } from "lucide-react"


export default function ContactPage() {
    const [formState, setFormState] = useState({
        name: "",
        email: "",
        subject: "",
        message: "",
    })
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [isSubmitted, setIsSubmitted] = useState(false)

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target
        setFormState((prev) => ({ ...prev, [name]: value }))
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        setIsSubmitting(true)

        // Simulate form submission
        setTimeout(() => {
            setIsSubmitting(false)
            setIsSubmitted(true)
            setFormState({
                name: "",
                email: "",
                subject: "",
                message: "",
            })

            // Reset success message after 5 seconds
            setTimeout(() => {
                setIsSubmitted(false)
            }, 5000)
        }, 1500)
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

    return (
        <div className="flex min-h-screen flex-col">
            <main className="flex-1 py-12 px-4">
                <motion.div className="max-w-6xl mx-auto" initial="hidden" animate="visible" variants={containerVariants}>
                    <motion.div className="text-center mb-12" variants={itemVariants}>
                        <h1 className="text-4xl font-bold mb-4">Get in Touch</h1>
                        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                            Have a question or want to work together? Feel free to reach out using the form below or through my
                            contact information.
                        </p>
                    </motion.div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        <motion.div className="md:col-span-1 space-y-6" variants={itemVariants}>
                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex items-start space-x-4">
                                        <MapPin className="h-6 w-6 text-primary mt-1" />
                                        <div>
                                            <h3 className="font-semibold text-lg mb-1">Location</h3>
                                            <p className="text-muted-foreground">Blumenau, SC, Brazil, 89037-035</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex items-start space-x-4">
                                        <Phone className="h-6 w-6 text-primary mt-1" />
                                        <div>
                                            <h3 className="font-semibold text-lg mb-1">Phone</h3>
                                            <p className="text-muted-foreground">+55 47 99631-3303</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex items-start space-x-4">
                                        <Mail className="h-6 w-6 text-primary mt-1" />
                                        <div>
                                            <h3 className="font-semibold text-lg mb-1">Email</h3>
                                            <p className="text-muted-foreground">phlucasfr@gmail.com</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <div className="p-6 bg-primary/10 rounded-lg">
                                <h3 className="font-semibold text-lg mb-3">Connect with me</h3>
                                <p className="text-muted-foreground mb-4">
                                    Find me on LinkedIn and GitHub to see more of my work and professional background.
                                </p>
                                <div className="flex space-x-4">
                                    <a
                                        href="https://www.linkedin.com/in/phlucasfr/"
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="p-2 bg-[#0077B5] text-white rounded-full hover:opacity-90 transition-opacity"
                                    >
                                        <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                                            <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
                                        </svg>
                                    </a>
                                    <a
                                        href="https://github.com/phlucasfr"
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="p-2 bg-[#333] text-white rounded-full hover:opacity-90 transition-opacity"
                                    >
                                        <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                                            <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
                                        </svg>
                                    </a>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </motion.div>
            </main>
        </div>
    )
}
