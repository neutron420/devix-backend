"use client";

import type { FormEvent } from "react";
import { Suspense } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Mail, Smartphone, ShieldCheck } from "lucide-react";
import { motion } from "framer-motion";
import { AuthShell } from "@/components/auth-shell";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

function ForgotPasswordVerifyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const channel = searchParams.get("channel") ?? "email";
  const isEmail = channel === "email";
  const channelLabel = isEmail ? "Email" : "SMS";

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    router.push("/login");
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
      title="Verify OTP"
      description={`Enter the 6-digit code sent to your ${channelLabel}`}
    >
      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        {/* Channel indicator badge */}
        <motion.div
          variants={itemVariants}
          className="flex items-center gap-2 mb-4"
        >
          <div
            className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold ${
              isEmail
                ? "bg-blue-50 text-blue-600 border border-blue-100"
                : "bg-emerald-50 text-emerald-600 border border-emerald-100"
            }`}
          >
            {isEmail ? (
              <Mail className="w-3.5 h-3.5" />
            ) : (
              <Smartphone className="w-3.5 h-3.5" />
            )}
            {channelLabel} Verification
          </div>
        </motion.div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <motion.div variants={itemVariants} className="space-y-1">
            <label className="text-sm font-medium flex items-center gap-1.5">
              <ShieldCheck className="w-4 h-4 text-blue-500" />
              OTP Code
            </label>
            <Input
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              name="otp"
              placeholder="Enter 6-digit code"
              maxLength={6}
              required
              className="focus-visible:ring-blue-500 text-center text-lg tracking-[0.5em] font-mono"
            />
          </motion.div>
          <motion.div variants={itemVariants}>
            <Button
              type="submit"
              className="w-full bg-linear-to-r from-blue-600 to-sky-600 hover:from-blue-700 hover:to-sky-700 text-white shadow-md shadow-blue-500/10 hover:shadow-blue-500/20 active:scale-[0.98] transition-all cursor-pointer"
            >
              Verify OTP
            </Button>
          </motion.div>
        </form>
      </motion.div>

      <div className="flex flex-col items-center gap-2">
        <p className="text-center text-sm text-muted-foreground">
          Didn&apos;t receive the code?{" "}
          <button
            type="button"
            className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors cursor-pointer"
          >
            Resend OTP
          </button>
        </p>
        <p className="text-center text-sm text-muted-foreground">
          <Link
            href="/forgot-password"
            className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
          >
            Choose another method
          </Link>
          .
        </p>
      </div>
    </AuthShell>
  );
}

export default function ForgotPasswordVerifyPage() {
  return (
    <Suspense fallback={null}>
      <ForgotPasswordVerifyContent />
    </Suspense>
  );
}
