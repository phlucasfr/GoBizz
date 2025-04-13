"use client"

import Link from "next/link"

import type React from "react"

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { motion } from "framer-motion"
import { useAuth } from "@/context/auth-context"
import { useState } from "react"
import { useToast } from "@/components/ui/use-toast"
import { useRouter } from "next/navigation"
import { Eye, EyeOff } from "lucide-react"
import { maskPhone, maskCpfCnpj } from "@/utils/Index"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog"


export default function Register() {
  const router = useRouter()
  const { register } = useAuth()
  const { toast } = useToast()

  const [formData, setFormData] = useState({
    email: "",
    phone: "",
    document: "",
    password: "",
    companyName: "",
  })
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [showSuccessModal, setShowSuccessModal] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value } = e.target
    let maskedValue = value

    if (id === "phone") {
      maskedValue = maskPhone(value)
    } else if (id === "document") {
      maskedValue = maskCpfCnpj(value)
    }

    setFormData((prev) => ({
      ...prev,
      [id]: maskedValue,
    }))
    setError(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      await register(formData)
      setShowSuccessModal(true)
    } catch (error) {
      const errorMessage = error instanceof Error
        ? error.message
        : "Registration error. Please try again later."
      setError(errorMessage)
      toast({
        title: "Registration Error",
        description: errorMessage,
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05,
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
      <main className="flex-1 flex items-center justify-center p-4">
        <motion.div
          className="w-full max-w-md bg-white rounded-xl shadow-sm p-8"
          initial="hidden"
          animate="visible"
          variants={containerVariants}
        >
          <motion.h1 className="text-2xl font-bold text-center mb-8 text-primary" variants={itemVariants}>
            Register Your Company
          </motion.h1>

          <form className="space-y-6" onSubmit={handleSubmit}>
            {error && (
              <motion.div
                className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-200"
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                variants={itemVariants}
              >
                {error}
              </motion.div>
            )}

            <motion.div variants={itemVariants}>
              <Label htmlFor="companyName">Company Name</Label>
              <Input
                id="companyName"
                placeholder="Your company"
                className="mt-1"
                value={formData.companyName}
                onChange={handleChange}
                required
              />
            </motion.div>

            <motion.div variants={itemVariants}>
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="your@email.com"
                className="mt-1"
                value={formData.email}
                onChange={handleChange}
                required
              />
            </motion.div>

            <motion.div variants={itemVariants}>
              <Label htmlFor="phone">Phone</Label>
              <Input
                id="phone"
                placeholder="(00) 00000-0000"
                className="mt-1"
                value={formData.phone}
                onChange={handleChange}
                required
              />
            </motion.div>

            <motion.div variants={itemVariants}>
              <Label htmlFor="document">Tax ID</Label>
              <Input
                id="document"
                placeholder="000.000.000-00 or 00.000.000/0000-00"
                className="mt-1"
                value={formData.document}
                onChange={handleChange}
                required
              />
            </motion.div>

            <motion.div variants={itemVariants}>
              <Label htmlFor="password">Password</Label>
              <div className="relative mt-1">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="Your password"
                  value={formData.password}
                  onChange={handleChange}
                  required
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </motion.div>

            <motion.div variants={itemVariants}>
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? "Registering..." : "Register"}
              </Button>
            </motion.div>

            <motion.div className="text-center text-sm" variants={itemVariants}>
              Already have an account?{" "}
              <Link href="/login" className="text-primary font-medium hover:underline">
                Sign in
              </Link>
            </motion.div>
          </form>
        </motion.div>
      </main>

      <Dialog open={showSuccessModal} onOpenChange={setShowSuccessModal}>
        <DialogContent className="sm:max-w-md" hideCloseButton>
          <DialogHeader>
            <DialogTitle>Registration Successful!</DialogTitle>
            <DialogDescription>
              A verification email has been sent to your email address. Please check your inbox and follow the instructions to verify your account.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button onClick={() => router.push('/')} className="w-full">
              Return to Home
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}

