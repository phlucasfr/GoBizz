"use client"

import Link from "next/link"
import PasswordRecoveryModal from "@/components/password-recovery-modal"

import type React from "react"

import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { motion } from "framer-motion"
import { useAuth } from "@/context/auth-context"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { Eye, EyeOff } from "lucide-react"

export default function Login() {
  const router = useRouter()
  const { login } = useAuth()

  const [email, setEmail] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [password, setPassword] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [showPasswordRecovery, setShowPasswordRecovery] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      await login(email, password)
      router.push("/dashboard")
    } catch (error) {
      console.error("Login failed:", error)
      setError("Invalid email or password. Please try again.")
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

  const logoVariants = {
    hidden: { scale: 0.8, opacity: 0 },
    visible: {
      scale: 1,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 100,
        damping: 15,
      },
    },
  }

  return (
    <main className="flex-1 flex items-center justify-center p-4">
      <motion.div
        className="w-full max-w-md bg-white rounded-xl shadow-sm p-8"
        initial="hidden"
        animate="visible"
        variants={containerVariants}
      >
        <motion.div className="flex justify-center mb-6" variants={logoVariants}>
          <div className="w-12 h-12 rounded-full bg-primary flex items-center justify-center">
            <span className="text-white font-bold">G</span>
          </div>
        </motion.div>

        <motion.h1 className="text-2xl font-bold text-center mb-8 text-primary" variants={itemVariants}>
          Welcome Back
        </motion.h1>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <motion.div variants={itemVariants}>
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="your@email.com"
              className="mt-1"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </motion.div>

          <motion.div variants={itemVariants}>
            <div className="flex items-center justify-between">
              <Label htmlFor="password">Password</Label>
              <button
                type="button"
                className="text-sm text-primary hover:underline"
                onClick={() => setShowPasswordRecovery(true)}
              >
                Forgot password?
              </button>
            </div>
            <div className="relative mt-1">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder="Your password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
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
              {isLoading ? "Signing in..." : "Sign in"}
            </Button>
            {error && (
              <motion.p
                className="mt-2 text-sm text-red-500 text-center"
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
              >
                {error}
              </motion.p>
            )}
          </motion.div>

          <motion.div className="text-center text-sm" variants={itemVariants}>
            Don't have an account?{" "}
            <Link href="/register" className="text-primary font-medium hover:underline">
              Register
            </Link>
          </motion.div>
        </form>
      </motion.div>

      <PasswordRecoveryModal isOpen={showPasswordRecovery} onClose={() => setShowPasswordRecovery(false)} />
    </main>
  )
}

