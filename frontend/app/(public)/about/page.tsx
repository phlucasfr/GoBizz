"use client"

import Image from "next/image"
import { Badge } from "@/components/ui/badge"
import { motion } from "framer-motion"
import { Card, CardContent } from "@/components/ui/card"
import { Briefcase, GraduationCap, Languages, Code, Github, Linkedin, Mail } from "lucide-react"

export default function AboutPage() {
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
            <main className="flex-1 py-12 px-4 max-w-6xl mx-auto">
                <motion.div className="space-y-12" initial="hidden" animate="visible" variants={containerVariants}>
                    {/* Hero Section */}
                    <motion.section className="flex flex-col md:flex-row gap-8 items-center" variants={itemVariants}>
                        <div className="w-48 h-48 relative rounded-full overflow-hidden border-4 border-primary">
                            <Image
                                src="/phelipe.jpg?height=192&width=192"
                                alt="Phelipe Lucas França da Silva"
                                width={192}
                                height={192}
                                className="object-cover"
                            />
                        </div>
                        <div className="text-center md:text-left">
                            <h1 className="text-3xl md:text-4xl font-bold mb-2">Phelipe Lucas França da Silva</h1>
                            <p className="text-xl text-muted-foreground mb-4">Full Stack & RPA Developer</p>
                            <div className="flex flex-wrap gap-2 justify-center md:justify-start">
                                <Badge variant="outline" className="flex items-center gap-1">
                                    <Mail size={14} /> phlucasfr@gmail.com
                                </Badge>
                                <a
                                    href="https://linkedin.com/in/phlucasfr"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center gap-1 border border-primary rounded px-2 py-1 text-primary hover:bg-primary hover:text-white transition"
                                >
                                    <Linkedin size={14} /> linkedin.com/in/phlucasfr
                                </a>
                                <a
                                    href="https://github.com/phlucasfr"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex items-center gap-1 border border-primary rounded px-2 py-1 text-primary hover:bg-primary hover:text-white transition"
                                >
                                    <Github size={14} /> github.com/phlucasfr
                                </a>
                            </div>
                        </div>
                    </motion.section>

                    {/* Professional Summary */}
                    <motion.section variants={itemVariants}>
                        <h2 className="text-2xl font-bold mb-4 border-b pb-2">Professional Summary</h2>
                        <p className="text-lg">
                            Full Stack and RPA Developer with over 3 years of experience, working on process automation (RPA)
                            projects, web development, support and maintenance of legacy systems.
                        </p>
                    </motion.section>



                    {/* Languages */}

                    {/* Technical Skills */}
                    <motion.section variants={itemVariants}>
                        <h2 className="text-2xl font-bold mb-6 border-b pb-2 flex items-center gap-2">
                            <Code className="h-6 w-6" /> Technical Skills
                        </h2>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <h3 className="text-lg font-semibold mb-3">Programming Languages</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>C#</Badge>
                                    <Badge>JavaScript</Badge>
                                    <Badge>PHP</Badge>
                                    <Badge>Golang</Badge>
                                    <Badge>Rust</Badge>
                                    <Badge>Python</Badge>
                                    <Badge>Visual Basic</Badge>
                                    <Badge>Visual FoxPro</Badge>
                                    <Badge>Delphi</Badge>
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">RPA Tools</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>UiPath</Badge>
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">Front-end</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>React</Badge>
                                    <Badge>Next.js</Badge>
                                    <Badge>Solid.js</Badge>
                                 
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">Back-end</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>Node.js</Badge>
                                    <Badge>REST API</Badge>
                                    <Badge>GraphQL</Badge>
                                    <Badge>gRPC</Badge>
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">Databases</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>MySQL</Badge>
                                    <Badge>PostgreSQL</Badge>
                                    <Badge>MongoDB</Badge>
                                    <Badge>SQL Server</Badge>
                                    <Badge>Oracle</Badge>
                                    <Badge>DynamoDB</Badge>
                                    <Badge>Redis</Badge>
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">DevOps & Others</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>Docker</Badge>
                                    <Badge>Kubernetes</Badge>
                                    <Badge>Git</Badge>
                                    <Badge>CI/CD</Badge>
                                    <Badge>Jira</Badge>
                                    <Badge>Scrum</Badge>
                                    <Badge>Kanban</Badge>
                                </div>
                            </div>
                        </div>
                    </motion.section>
                </motion.div>
            </main>

        </div>
    )
}
