"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Mail, Smartphone } from "lucide-react";
import { motion } from "framer-motion";
import { AuthShell } from "@/components/auth-shell";

export default function ForgotPasswordPage() {
  const router = useRouter();

  const handleMethod = (channel: "email" | "sms") => {
    router.push(`/forgot-password/method?channel=${channel}`);
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { opacity: 1, transition: { staggerChildren: 0.1 } },
  };

  const itemVariants = {
    hidden: { y: 15, opacity: 0 },
    visible: { y: 0, opacity: 1 },
  };

  return (
    <AuthShell
      logo={
        <Image
          src="/devix1.png"
          alt="Devix Logo"
          width={220}
          height={70}
          className="h-14 w-auto object-contain"
          priority
        />
      }
      title="Forgot Password"
      description="How would you like to receive your OTP?"
    >
      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="grid gap-3"
      >
        <motion.button
          variants={itemVariants}
          type="button"
          onClick={() => handleMethod("email")}
          className="group flex items-center gap-4 w-full rounded-xl border border-border bg-white p-4 text-left shadow-sm hover:border-blue-300 hover:bg-blue-50/40 hover:shadow-md active:scale-[0.98] transition-all cursor-pointer"
        >
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-blue-50 text-blue-600 group-hover:bg-blue-100 transition-colors">
            <Mail className="h-5 w-5" />
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">
              Send OTP via Email
            </p>
            <p className="text-xs text-muted-foreground mt-0.5">
              We&apos;ll send a code to your registered email address
            </p>
          </div>
        </motion.button>

        <motion.button
          variants={itemVariants}
          type="button"
          onClick={() => handleMethod("sms")}
          className="group flex items-center gap-4 w-full rounded-xl border border-border bg-white p-4 text-left shadow-sm hover:border-emerald-300 hover:bg-emerald-50/40 hover:shadow-md active:scale-[0.98] transition-all cursor-pointer"
        >
          <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600 group-hover:bg-emerald-100 transition-colors">
            <Smartphone className="h-5 w-5" />
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">
              Send OTP via SMS
            </p>
            <p className="text-xs text-muted-foreground mt-0.5">
              We&apos;ll send a code to your registered phone number
            </p>
          </div>
        </motion.button>
      </motion.div>

      <p className="text-center text-sm text-muted-foreground">
        Remembered your password?{" "}
        <Link
          href="/login"
          className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
        >
          Sign in here
        </Link>
        .
      </p>
    </AuthShell>
  );
}
