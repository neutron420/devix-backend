"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";

interface AuthShellProps {
  logo: React.ReactNode;
  title: string;
  description?: string;
  children: React.ReactNode;
  showBackButton?: boolean;
}

export function AuthShell({
  logo,
  title,
  description,
  children,
  showBackButton = true,
}: AuthShellProps) {
  const router = useRouter();

  const handleBack = () => {
    if (typeof window !== "undefined" && window.history.length > 1) {
      router.back();
    } else {
      router.push("/login");
    }
  };

  return (
    <div className="relative flex min-h-screen w-full flex-col items-center justify-start bg-linear-to-tr from-white via-sky-50/10 to-blue-50/30 p-8 pt-20 pb-12 overflow-hidden">
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#0284c703_1px,transparent_1px),linear-gradient(to_bottom,#0284c703_1px,transparent_1px)] bg-size-[32px_32px] pointer-events-none" />

      {showBackButton && (
        <button
          type="button"
          onClick={handleBack}
          className="absolute top-5 left-5 z-[50] flex items-center justify-center w-9 h-9 rounded-full border border-border/60 bg-white/90 dark:bg-neutral-900/90 backdrop-blur-md hover:bg-slate-50 dark:hover:bg-neutral-800 text-muted-foreground hover:text-foreground transition-all active:scale-95 cursor-pointer shadow-md"
          aria-label="Go back"
        >
          <ArrowLeft className="w-4 h-4" />
        </button>
      )}

      <div className="w-full max-w-md relative z-10 flex flex-col gap-4">
        <div className="flex justify-center">{logo}</div>
        <div className="text-left">
          <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
          {description && (
            <p className="text-sm text-muted-foreground">{description}</p>
          )}
        </div>
        {children}
      </div>
    </div>
  );
}
