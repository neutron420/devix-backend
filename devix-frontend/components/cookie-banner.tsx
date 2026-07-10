"use client";

import { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

const STORAGE_KEY = "devix_cookie_consent";

export function CookieBanner() {
  const [isVisible, setIsVisible] = useState(() => {
    if (typeof window === "undefined") {
      return false;
    }
    return !window.localStorage.getItem(STORAGE_KEY);
  });

  const handleChoice = (value: "accepted" | "rejected") => {
    window.localStorage.setItem(STORAGE_KEY, value);
    setIsVisible(false);
  };

  if (!isVisible) {
    return null;
  }

  return (
    <div className="fixed bottom-6 left-6 right-6 z-60 sm:right-auto sm:max-w-md">
      <div className="rounded-2xl border border-border/70 bg-white/95 p-5 shadow-[0_12px_30px_rgba(15,23,42,0.18)] backdrop-blur">
        <p className="text-sm text-slate-700">
          This website uses cookies to improve your experience and remember your
          preferences. You can accept or reject optional cookies anytime. See our{" "}
          <Link
            href="/cookie-policy"
            className="font-medium text-blue-600 hover:text-blue-700 hover:underline"
          >
            cookie policy
          </Link>
          .
        </p>
        <div className="mt-4 flex items-center gap-3">
          <Button
            type="button"
            className="bg-slate-900 text-white hover:bg-slate-800"
            onClick={() => handleChoice("accepted")}
          >
            Accept all
          </Button>
          <Button
            type="button"
            variant="ghost"
            className="text-slate-700 hover:text-slate-900"
            onClick={() => handleChoice("rejected")}
          >
            Reject all
          </Button>
        </div>
      </div>
    </div>
  );
}
