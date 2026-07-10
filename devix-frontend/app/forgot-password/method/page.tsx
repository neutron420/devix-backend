"use client";

import type { FormEvent } from "react";
import { Suspense } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Mail, Smartphone } from "lucide-react";
import { motion } from "framer-motion";
import { AuthShell } from "@/components/auth-shell";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import PhoneInput from "react-phone-input-2";
import "react-phone-input-2/lib/style.css";

function MethodContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const channel = searchParams.get("channel") ?? "email";
  const isEmail = channel === "email";

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    router.push(`/forgot-password/verify?channel=${channel}`);
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
      title={isEmail ? "Enter Your Email" : "Enter Your Phone Number"}
      description={
        isEmail
          ? "We'll send a one-time password to your email address"
          : "We'll send a one-time password to your phone via SMS"
      }
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
            {isEmail ? "Email OTP" : "SMS OTP"}
          </div>
        </motion.div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <motion.div variants={itemVariants} className="space-y-1">
            <label className="text-sm font-medium">
              {isEmail ? "Email Address" : "Phone Number"}
            </label>
            {isEmail ? (
              <Input
                type="email"
                name="email"
                placeholder="email@example.com"
                required
                className="focus-visible:ring-blue-500"
              />
            ) : (
              <PhoneInput
                country="in"
                enableSearch
                placeholder="Enter phone number"
                inputProps={{
                  name: "phone",
                  required: true,
                }}
                containerClass="!w-full"
                inputClass="!w-full !h-10 !text-sm !bg-background !border !border-input !rounded-md !pl-12 !pr-3 !py-2 !text-foreground focus:!ring-2 focus:!ring-blue-500 focus:!ring-offset-2"
                buttonClass="!border-input !bg-background !rounded-md !rounded-r-none !h-10 !w-11"
                dropdownClass="!text-sm"
                searchClass="!text-sm !py-2"
              />
            )}
          </motion.div>
          <motion.div variants={itemVariants}>
            <Button
              type="submit"
              className="w-full bg-linear-to-r from-blue-600 to-sky-600 hover:from-blue-700 hover:to-sky-700 text-white shadow-md shadow-blue-500/10 hover:shadow-blue-500/20 active:scale-[0.98] transition-all cursor-pointer"
            >
              Send OTP
            </Button>
          </motion.div>
        </form>
      </motion.div>

      <p className="text-center text-sm text-muted-foreground">
        Want to use a different method?{" "}
        <Link
          href="/forgot-password"
          className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
        >
          Go back
        </Link>
        .
      </p>
    </AuthShell>
  );
}

export default function ForgotPasswordMethodPage() {
  return (
    <Suspense fallback={null}>
      <MethodContent />
    </Suspense>
  );
}
