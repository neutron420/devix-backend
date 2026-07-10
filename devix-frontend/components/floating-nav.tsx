"use client"

import { useState } from "react"
import { Home, Compass, PlusCircle, MessageSquare, LogIn, UserPlus, Menu, X } from "lucide-react"
import { MenuBar } from "@/components/ui/bottom-menu"
import { motion, AnimatePresence } from "motion/react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { UpgradeBanner } from "@/components/ui/upgrade-banner"

export function FloatingNav() {
  const [isOpen, setIsOpen] = useState(false)
  const [showBanner, setShowBanner] = useState(true)
  const pathname = usePathname()

  const hideNavRoutes = [
    "/login",
    "/signup",
    "/forgot-password",
    "/terms",
    "/privacy",
  ]

  const shouldHide = hideNavRoutes.some((route) => pathname?.startsWith(route))

  if (shouldHide) return null

  const menuItems = [
    { icon: Home, label: "Home", href: "/" },
    { icon: Compass, label: "Explore", href: "/explore" },
    { icon: PlusCircle, label: "Create", href: "/post/new" },
    { icon: MessageSquare, label: "Chat", href: "/chat" },
    { icon: LogIn, label: "Sign In", href: "/login" },
    { icon: UserPlus, label: "Sign Up", href: "/signup" },
  ]

  return (
    <>
      {/* Desktop Navigation */}
      <div className="hidden md:block fixed top-4 sm:top-6 left-1/2 -translate-x-1/2 z-[70] select-none">
        <MenuBar items={menuItems} />
      </div>

      {/* Mobile Hamburger Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="md:hidden fixed top-4 sm:top-6 right-4 sm:right-6 z-[70] p-2.5 rounded-full bg-background/85 backdrop-blur-md border border-border/50 shadow-lg text-foreground hover:bg-muted/80 active:scale-95 transition-all select-none"
        aria-label="Toggle menu"
      >
        {isOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
      </button>

      {/* Upgrade Banner (Desktop & Mobile) */}
      <AnimatePresence>
        {showBanner && (
          <motion.div
            initial={{ opacity: 0, y: -10, x: "-50%" }}
            animate={{ opacity: 1, y: 0, x: "-50%" }}
            exit={{ opacity: 0, y: -10, x: "-50%" }}
            transition={{ duration: 0.3 }}
            className="fixed top-24 sm:top-28 md:top-32 left-1/2 z-40 select-none w-[90%] sm:w-auto flex justify-center"
          >
            <UpgradeBanner
              buttonText="Building in Progress"
              description="Devix platform is currently under active construction"
              onClose={() => setShowBanner(false)}
            />
          </motion.div>
        )}
      </AnimatePresence>

      {/* Mobile Menu Overlay */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.25, ease: "easeInOut" }}
            className="fixed inset-0 z-40 bg-background/98 backdrop-blur-xl md:hidden flex flex-col justify-between px-8 pt-28 pb-12"
          >
            {/* Background pattern */}
            <div className="w-full absolute inset-0 z-0 bg-[radial-gradient(circle,_black_1px,_transparent_1px)] dark:bg-[radial-gradient(circle,_white_1px,_transparent_1px)] opacity-[0.03] [background-size:20px_20px] pointer-events-none" />

            <div className="relative z-10 flex flex-col justify-center h-full">
              <nav className="flex flex-col gap-6 pl-4">
                {menuItems.map((item, index) => {
                  const Icon = item.icon
                  return (
                    <motion.div
                      key={item.label}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.05, duration: 0.3 }}
                    >
                      <Link
                        href={item.href}
                        onClick={() => setIsOpen(false)}
                        className="flex items-center gap-4 text-3xl font-light text-foreground/80 hover:text-foreground active:text-primary transition-colors py-2"
                      >
                        <Icon className="w-6 h-6 text-muted-foreground" />
                        <span>{item.label}</span>
                      </Link>
                    </motion.div>
                  )
                })}
              </nav>
            </div>

            <div className="relative z-10 flex flex-col gap-4 border-t border-border/45 pt-6">
              <div className="flex items-center justify-between text-[10px] tracking-wider text-muted-foreground">
                <span>DEVELOPER ECOSYSTEM</span>
                <span className="font-bold italic text-orange-600">BUILD PUBLIC</span>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  )
}
