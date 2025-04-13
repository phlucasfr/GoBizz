"use client"

import Link from "next/link"

import { motion } from "framer-motion"
import { Button } from "@/components/ui/button"
import { ArrowRight, FileText, Link2 } from "lucide-react"

export default function Home() {
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

  const cardVariants = {
    hidden: { scale: 0.95, opacity: 0 },
    visible: {
      scale: 1,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 100,
        damping: 15,
        delay: 0.2,
      },
    },
    hover: {
      y: -5,
      boxShadow: "0 10px 25px rgba(0, 0, 0, 0.1)",
      transition: {
        type: "spring",
        stiffness: 400,
        damping: 10,
      },
    },
  }

  return (
    <div className="flex min-h-screen flex-col">

      <main className="flex-1">
        <motion.section
          className="py-20 md:py-28 text-center px-4"
          initial="hidden"
          animate="visible"
          variants={containerVariants}
        >
          <motion.h1
            className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 max-w-4xl mx-auto"
            variants={itemVariants}
          >
            Transform Your Business Management
          </motion.h1>

          <motion.p className="text-lg text-muted-foreground max-w-2xl mx-auto mb-10" variants={itemVariants}>
            Discover intelligent tools and complete resources to drive your business success.
          </motion.p>

          <motion.div variants={itemVariants}>
            <Button asChild size="lg" className="rounded-full px-8 py-6 text-lg font-medium">
              <Link href="/register">
                Get Started <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
            </Button>
          </motion.div>
        </motion.section>

        <section className="py-16 px-4">
          <div className="max-w-6xl mx-auto">
            <h2 className="text-3xl font-bold text-center mb-16">Our Solutions</h2>

            <div className="grid md:grid-cols-2 gap-8">
              <motion.div
                className="bg-card rounded-xl p-8 shadow-sm"
                variants={cardVariants}
                initial="hidden"
                whileInView="visible"
                whileHover="hover"
                viewport={{ once: true, amount: 0.3 }}
              >
                <div className="flex items-start justify-between mb-4">
                  <h3 className="text-xl font-semibold">Report Generator</h3>
                  <FileText className="text-primary h-6 w-6" />
                </div>
                <p className="text-muted-foreground mb-6">
                  Create detailed reports and valuable insights for your business.
                </p>

                <h4 className="font-medium mb-3">Highlights:</h4>
                <ul className="space-y-3">
                  <li className="flex items-start">
                    <div className="bg-amber-100 p-1 rounded mr-3 mt-0.5 dark:bg-amber-900">
                      <svg
                        className="h-4 w-4 text-amber-600 dark:text-amber-300"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M13 10V3L4 14h7v7l9-11h-7z"
                        />
                      </svg>
                    </div>
                    <span>Fast processing of large data sets for real-time reporting.</span>
                  </li>
                  <li className="flex items-start">
                    <div className="bg-green-100 p-1 rounded mr-3 mt-0.5 dark:bg-green-900">
                      <svg
                        className="h-4 w-4 text-green-600 dark:text-green-300"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                        />
                      </svg>
                    </div>
                    <span>Granular access control and log auditing to protect sensitive information.</span>
                  </li>
                </ul>
              </motion.div>

              <motion.div
                className="bg-card rounded-xl p-8 shadow-sm"
                variants={cardVariants}
                initial="hidden"
                whileInView="visible"
                whileHover="hover"
                viewport={{ once: true, amount: 0.3 }}
              >
                <div className="flex items-start justify-between mb-4">
                  <h3 className="text-xl font-semibold">Link Shortener</h3>
                  <Link2 className="text-primary h-6 w-6" />
                </div>
                <p className="text-muted-foreground mb-6">Simplify and track your links with our intuitive tool.</p>

                <h4 className="font-medium mb-3">Highlights:</h4>
                <ul className="space-y-3">
                  <li className="flex items-start">
                    <div className="bg-amber-100 p-1 rounded mr-3 mt-0.5 dark:bg-amber-900">
                      <svg
                        className="h-4 w-4 text-amber-600 dark:text-amber-300"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M13 10V3L4 14h7v7l9-11h-7z"
                        />
                      </svg>
                    </div>
                    <span>Ultra-fast redirection for an uninterrupted user experience.</span>
                  </li>
                  <li className="flex items-start">
                    <div className="bg-green-100 p-1 rounded mr-3 mt-0.5 dark:bg-green-900">
                      <svg
                        className="h-4 w-4 text-green-600 dark:text-green-300"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                        />
                      </svg>
                    </div>
                    <span>Protection against spam and malware for secure and reliable links.</span>
                  </li>
                </ul>
              </motion.div>
            </div>
          </div>
        </section>
      </main>

    </div>
  )
}

