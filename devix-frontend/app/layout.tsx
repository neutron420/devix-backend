import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { FloatingNav } from "@/components/floating-nav";
import { CookieBanner } from "@/components/cookie-banner";
import "./globals.css";
import "react-phone-input-2/lib/style.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Devix - Engineering Knowledge Sharing Platform",
  description: "A modern, interactive, and real-time social platform built for developers to share technical knowledge, collaborate, and network.",
  icons: {
    icon: "/devix1q.png",
    apple: "/devix1q.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">
        {children}
        <FloatingNav />
        <CookieBanner />
      </body>
    </html>
  );
}
