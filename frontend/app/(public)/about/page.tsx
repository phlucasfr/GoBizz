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
                                <Badge variant="outline" className="flex items-center gap-1">
                                    <Linkedin size={14} /> linkedin.com/in/phlucasfr
                                </Badge>
                                <Badge variant="outline" className="flex items-center gap-1">
                                    <Github size={14} /> github.com/phlucasfr
                                </Badge>
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

                    {/* Work Experience */}
                    <motion.section variants={itemVariants}>
                        <h2 className="text-2xl font-bold mb-6 border-b pb-2 flex items-center gap-2">
                            <Briefcase className="h-6 w-6" /> Work Experience
                        </h2>

                        <div className="space-y-8">
                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                        <h3 className="text-xl font-semibold">RPA Developer</h3>
                                        <div className="text-muted-foreground">January 2025 - Present</div>
                                    </div>
                                    <div className="text-primary font-medium mb-4">Capgemini, Brazil</div>
                                    <ul className="list-disc list-inside space-y-2">
                                        <li>
                                            Development and maintenance of automations using C#, UiPath and Blue Prism for process automation.
                                        </li>
                                    </ul>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                        <h3 className="text-xl font-semibold">Full Stack Developer</h3>
                                        <div className="text-muted-foreground">September 2024 - December 2024</div>
                                    </div>
                                    <div className="text-primary font-medium mb-4">FW7 Soluções, Blumenau, SC, Brazil</div>
                                    <ul className="list-disc list-inside space-y-2">
                                        <li>
                                            Development of features with Node.js, TypeScript, React, Apollo GraphQL and Sequelize for sales
                                            automation and lead management.
                                        </li>
                                    </ul>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                        <h3 className="text-xl font-semibold">Full Stack Developer</h3>
                                        <div className="text-muted-foreground">September 2023 - July 2024</div>
                                    </div>
                                    <div className="text-primary font-medium mb-4">Benner Sistemas, Blumenau, SC, Brazil</div>
                                    <ul className="list-disc list-inside space-y-2">
                                        <li>
                                            Development of custom reports for major corporate clients, such as BRB, using C# and Visual Basic
                                            for process integration and automation.
                                        </li>
                                        <li>
                                            Support to the maintenance team, providing specific fixes and ensuring application stability in
                                            Delphi, focusing on legacy system maintenance and continuous improvement.
                                        </li>
                                    </ul>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                        <h3 className="text-xl font-semibold">Full Stack Developer</h3>
                                        <div className="text-muted-foreground">August 2022 - August 2023</div>
                                    </div>
                                    <div className="text-primary font-medium mb-4">Sances Sistemas, Blumenau, SC, Brazil</div>
                                    <ul className="list-disc list-inside space-y-2">
                                        <li>
                                            Maintenance and integration of APIs in Golang, handling large volumes of data for vehicle
                                            manufacturers.
                                        </li>
                                        <li>Development of a project for used vehicle commercialization, using React, C# and MongoDB.</li>
                                    </ul>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                        <h3 className="text-xl font-semibold">IT Technical Support</h3>
                                        <div className="text-muted-foreground">March 2022 - August 2022</div>
                                    </div>
                                    <div className="text-primary font-medium mb-4">Sances Sistemas, Blumenau, SC, Brazil</div>
                                    <ul className="list-disc list-inside space-y-2">
                                        <li>Fixes in HTML, CSS and Visual FoxPro code.</li>
                                        <li>Resolution of problems related to business, fiscal and accounting rules.</li>
                                    </ul>
                                </CardContent>
                            </Card>
                        </div>
                    </motion.section>

                    {/* Education */}
                    <motion.section variants={itemVariants}>
                        <h2 className="text-2xl font-bold mb-6 border-b pb-2 flex items-center gap-2">
                            <GraduationCap className="h-6 w-6" /> Education
                        </h2>

                        <Card className="hover:shadow-md transition-shadow">
                            <CardContent className="pt-6">
                                <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
                                    <h3 className="text-xl font-semibold">Systems Analysis and Development</h3>
                                    <div className="text-muted-foreground">January 2022 - August 2025 (Expected)</div>
                                </div>
                                <div className="text-primary font-medium">Estácio, Blumenau, SC, Brazil</div>
                            </CardContent>
                        </Card>
                    </motion.section>

                    {/* Languages */}
                    <motion.section variants={itemVariants}>
                        <h2 className="text-2xl font-bold mb-6 border-b pb-2 flex items-center gap-2">
                            <Languages className="h-6 w-6" /> Languages
                        </h2>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <h3 className="text-xl font-semibold mb-2">English</h3>
                                    <p>Proficiency B2 (Reading, Writing, Speaking)</p>
                                </CardContent>
                            </Card>

                            <Card className="hover:shadow-md transition-shadow">
                                <CardContent className="pt-6">
                                    <h3 className="text-xl font-semibold mb-2">Spanish</h3>
                                    <p>Intermediate</p>
                                </CardContent>
                            </Card>
                        </div>
                    </motion.section>

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
                                    <Badge>TypeScript</Badge>
                                    <Badge>Golang</Badge>
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
                                    <Badge>Blue Prism</Badge>
                                </div>
                            </div>

                            <div>
                                <h3 className="text-lg font-semibold mb-3">Front-end</h3>
                                <div className="flex flex-wrap gap-2">
                                    <Badge>React</Badge>
                                    <Badge>HTML</Badge>
                                    <Badge>CSS</Badge>
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
