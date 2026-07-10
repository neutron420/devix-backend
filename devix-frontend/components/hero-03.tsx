"use client";

import { Separator } from "@/components/ui/separator";
import { BadgeQuestionMark } from "@aliimam/icons";
import { Instagram, Threads, X } from "@aliimam/logos";
import React, { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { FloatingUIElements } from "./floating-ui-elements";

const badgeColors = [
  { bg: "#fde68a", text: "#0f172a" },
  { bg: "#bfdbfe", text: "#0f172a" },
  { bg: "#fecdd3", text: "#0f172a" },
  { bg: "#bbf7d0", text: "#0f172a" },
  { bg: "#e9d5ff", text: "#0f172a" },
  { bg: "#a7f3d0", text: "#0f172a" },
  { bg: "#bae6fd", text: "#0f172a" },
  { bg: "#fed7aa", text: "#0f172a" },
  { bg: "#fbcfe8", text: "#0f172a" },
  { bg: "#cffafe", text: "#0f172a" },
];

export function HeroSection03() {
  const [badgeColorIndex, setBadgeColorIndex] = useState(0);

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setBadgeColorIndex((current) => (current + 1) % badgeColors.length);
    }, 2000);

    return () => window.clearInterval(intervalId);
  }, []);

  const badgeColor = badgeColors[badgeColorIndex];

  return (
    <div className="min-h-screen relative">
      <div className="w-full absolute h-full z-0 bg-[radial-gradient(circle,black_1px,transparent_1px)] dark:bg-[radial-gradient(circle,white_1px,transparent_1px)] opacity-15 bg-size-[20px_20px]" />
      <header className="fixed top-4 sm:top-6 left-4 sm:left-6 z-70 flex h-10 items-center pointer-events-auto">
        <Link href="/" className="flex items-center hover:opacity-80 transition-opacity">
          <Image
            src="/devix1.png"
            alt="Devix Logo"
            width={180}
            height={180}
            className="h-10 w-auto object-contain"
            priority
          />
        </Link>
      </header>
      <main className="relative pt-24 sm:pt-20 pb-20 overflow-hidden">
        <FloatingUIElements />
        <div className="flex relative z-10 gap-6 sm:gap-4 px-6 sm:px-8 md:items-center w-full flex-col justify-center">
          <div className="sm:flex gap-4 sm:gap-6 items-start sm:items-center">
            <p className="text-[11px] sm:text-sm text-muted-foreground text-center sm:text-right leading-5 max-w-[320px] md:max-w-65 mx-auto sm:mx-0">
              Devix is a modern, interactive, and real-time social platform built for developers to share technical knowledge, collaborate, and network.
            </p>
            <h1 className="text-[clamp(2.6rem,12vw,6rem)] sm:text-6xl md:text-7xl xl:text-[10rem] font-light leading-[0.9] sm:leading-none tracking-[0.04em] sm:tracking-wider">
              DIGITAL
            </h1>
          </div>

          <div className="sm:flex gap-4 sm:gap-6 items-start sm:items-center">
            <h1 className="text-[clamp(2.6rem,12vw,6rem)] sm:text-6xl md:text-7xl xl:text-[10rem] flex font-light leading-[0.9] sm:leading-none tracking-[0.04em] sm:tracking-wider">
              <span>PR</span>
              <BadgeQuestionMark
                type="solid"
                className="lg:size-40 md:size-18 sm:size-14 size-10 text-primary"
              />
              <span>DUCTS</span>
            </h1>
            <p className="text-[11px] sm:text-sm text-muted-foreground sm:pt-8 leading-5 max-w-[320px] md:max-w-65 text-center sm:text-left mx-auto sm:mx-0">
              Devix bridges the gap between technical blogging, real-time collaboration, and community discussion.
            </p>
          </div>

          <div className="sm:flex gap-4 sm:gap-6 items-start sm:items-center">
            <h1 className="text-[clamp(2.6rem,12vw,6rem)] sm:text-6xl md:text-7xl xl:text-[10rem] sm:flex font-light leading-[0.9] sm:leading-none tracking-[0.04em] sm:tracking-wider">
              <span>DESIGN</span>
              <div className="hidden lg:block">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="160"
                  height="160"
                  viewBox="0 0 24 24"
                  fill="#f43f5e"
                >
                  <path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5" />
                </svg>
              </div>
              <div className="block lg:hidden">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="70"
                  height="70"
                  viewBox="0 0 24 24"
                  fill="#f43f5e"
                >
                  <path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5" />
                </svg>
              </div>
              <span>CODE</span>
            </h1>
          </div>
        </div>
        <div className="mx-auto max-w-7xl w-full px-6 gap-3">
          <div className="md:flex md:mx-8 grid md:justify-end items-center gap-3 text-center md:text-left">
            <Separator className="w-full my-6 mx-auto max-w-3xl" />
            <div className="text-xs md:text-sm">
              ENGINEERING KNOWLEDGE SHARING PLATFORM
            </div>
            <div className="flex w-full flex-col sm:flex-row sm:items-end items-center gap-1 sm:gap-3">
              <span className="text-2xl md:text-4xl font-thin">DEVELOPER</span>
              <span className="text-3xl md:text-5xl font-bold italic text-orange-600">
                DEVIX
              </span>
            </div>
          </div>
        </div>

        <div className="md:px-20 px-6 gap-6 items-end md:flex pt-12">
          <div className="w-full sm:w-84 h-48 sm:h-52 shadow-lg border rounded-md overflow-hidden mb-6 md:mb-0">
            <img
              src="/image.png"
              alt="Devix preview"
              className="w-full h-full object-cover"
            />
          </div>
          <p className="text-xs text-muted-foreground md:text-sm pt-8 leading-5">
            Designed to host developer-centric activities, giving engineers a dedicated home to write code journals, ask questions, and chat in real time.
          </p>
        </div>

        <div className="absolute bottom-8 right-8 md:right-12 hidden sm:flex gap-6">
          <Instagram />
          <X />
          <Threads />
        </div>

        <div className="fixed right-2 sm:right-0 top-1/2 h-28 sm:h-36 items-center hidden sm:flex transform -translate-y-1/2 z-55">
          <div
            className="py-6 px-3 text-sm font-bold transition-colors duration-500"
            style={{ backgroundColor: badgeColor.bg, color: badgeColor.text }}
          >
            <span className="rotate-180 [writing-mode:vertical-rl]">
              Devix Hub
            </span>
          </div>
        </div>
      </main>
    </div>
  );
}
