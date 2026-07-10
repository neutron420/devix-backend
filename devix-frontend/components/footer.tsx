"use client";

import Image from "next/image";
import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { FaGithub, FaTwitter, FaLinkedin, FaYoutube, FaInstagram } from "react-icons/fa";

const footerConfig = {
  description:
    "Devix is a modern, real-time social platform built for software developers to share technical questions, concepts, and build logs in public. Connect with the global community of builders and shape the future of technical knowledge.",
  contact: {
    email: "support@devix.in",
    phone: "+1 (555) 123-4567",
  },
  socials: [
    { icon: FaGithub, href: "#" },
    { icon: FaTwitter, href: "#" },
    { icon: FaLinkedin, href: "#" },
    { icon: FaYoutube, href: "#" },
    { icon: FaInstagram, href: "#" },
  ],
  columns: [
    {
      title: "Company",
      links: [
        { label: "About Us", href: "#" },
        { label: "Careers", href: "#" },
        { label: "Press & Media", href: "#" },
        { label: "Blog", href: "#" },
        { label: "Our Team", href: "#" },
        { label: "Events", href: "#" },
      ],
    },
    {
      title: "Platform",
      links: [
        { label: "Features", href: "#" },
        { label: "Pricing", href: "#" },
        { label: "Docs", href: "#" },
        { label: "API Reference", href: "#" },
        { label: "CLI", href: "#" },
        { label: "Changelog", href: "#" },
      ],
    },
    {
      title: "Resources",
      links: [
        { label: "Community Forum", href: "#" },
        { label: "Guides", href: "#" },
        { label: "Help Center", href: "#" },
        { label: "System Status", href: "#" },
        { label: "Case Studies", href: "#" },
      ],
    },
    {
      title: "Partners",
      links: [
        { label: "Become a Partner", href: "#" },
        { label: "Affiliate Program", href: "#" },
        { label: "Reseller Program", href: "#" },
        { label: "Technology Partners", href: "#" },
      ],
    },
    {
      title: "Legal",
      links: [
        { label: "Privacy Policy", href: "/privacy" },
        { label: "Terms of Service", href: "/terms" },
        { label: "Security", href: "#" },
        { label: "Cookie Policy", href: "#" },
      ],
    },
  ],
};

export default function Footer() {
  return (
    <footer className="bg-white dark:bg-black text-black dark:text-white px-6 py-14 border-t border-gray-200 dark:border-gray-800">
      <div className="max-w-7xl mx-auto">
        <div className="mb-6">
          {/* Logo */}
          <div className="relative mb-4">
            <Link href="/" className="flex items-center hover:opacity-80 transition-opacity">
              <Image
                src="/devix1.png"
                alt="Devix Logo"
                width={220}
                height={220}
                className="h-20 md:h-24 w-auto object-contain dark:invert"
              />
            </Link>
          </div>

          <p className="text-sm text-gray-600 dark:text-gray-300 leading-relaxed max-w-2xl">
            {footerConfig.description}
          </p>
        </div>

        <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start">
          {/* Left Side: Links */}
          <div className="grid grid-cols-2 md:grid-cols-5 gap-8 flex-1">
            {footerConfig.columns.map((col, idx) => (
              <div key={idx}>
                <h3 className="text-sm font-medium mb-3">{col.title}</h3>
                <ul className="space-y-2">
                  {col.links.map((link, i) => (
                    <li key={i}>
                      <Link
                        href={link.href}
                        className="text-[0.85rem] text-gray-600 dark:text-gray-300 hover:text-blue-500 transition"
                      >
                        {link.label}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          {/* Right Side: Newsletter and Contact */}
          <div className="lg:w-1/4">
            {/* Contact */}
            <Card className="shadow-none border-none mb-4">
              <CardContent className="p-0 space-y-3">
                <p className="text-sm font-medium">For Corporates & Universities</p>
                <form className="flex flex-col gap-3">
                  <Button variant="default" type="submit" className="w-full bg-gray-200 border border-gray-400 text-gray-600 hover:text-white">
                    Get In Touch
                  </Button>
                </form>
              </CardContent>
            </Card>
            {/* Quick Links & Resources */}
            <Card className="shadow-none border-none mb-4">
              <CardContent className="p-0">
                <p className="text-sm font-medium mb-3">
                  Quick Links
                </p>
                <div className="space-y-2">
                  <Link href="#" className="block text-sm text-gray-600 dark:text-gray-300 hover:text-blue-500 transition">
                    Documentation
                  </Link>
                  <Link href="#" className="block text-sm text-gray-600 dark:text-gray-300 hover:text-blue-500 transition">
                    Getting Started
                  </Link>
                  <Link href="#" className="block text-sm text-gray-600 dark:text-gray-300 hover:text-blue-500 transition">
                    Community Forum
                  </Link>
                  <Link href="#" className="block text-sm text-gray-600 dark:text-gray-300 hover:text-blue-500 transition">
                    Tutorials
                  </Link>
                  <Link href="#" className="block text-sm text-gray-600 dark:text-gray-300 hover:text-blue-500 transition">
                    Status Page
                  </Link>
                </div>
                
                {/* Social Links */}
                <div className="mt-4 pt-3 border-t border-gray-200 dark:border-gray-700">
                  <p className="text-sm font-medium mb-2">Follow Us</p>
                  <div className="flex gap-3">
                    {footerConfig.socials.map(({ icon: Icon, href }, idx) => (
                      <Link key={idx} href={href} className="text-gray-500 hover:text-blue-500 transition">
                        <Icon className="w-4 h-4" />
                      </Link>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Bottom Section */}
        <div className="mt-12 pt-6 flex flex-col md:flex-row justify-between items-center text-xs text-gray-500 dark:text-gray-400 gap-4">
          <p>© {new Date().getFullYear()} devix.in. All rights reserved.</p>
          <div className="flex gap-6">
            <Link href="#">Privacy</Link>
            <Link href="#">Terms</Link>
            <Link href="#">Sitemap</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
