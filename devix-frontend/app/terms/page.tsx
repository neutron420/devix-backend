"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { ArrowLeft } from "lucide-react";

export default function TermsPage() {
  const router = useRouter();

  const handleBack = () => {
    if (typeof window !== "undefined" && window.history.length > 1) {
      router.back();
    } else {
      router.push("/");
    }
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { opacity: 1, transition: { staggerChildren: 0.05 } },
  };

  const itemVariants = {
    hidden: { y: 10, opacity: 0 },
    visible: { y: 0, opacity: 1 },
  };

  return (
    <div className="min-h-screen bg-white dark:bg-neutral-950">
      {/* Content */}
      <main className="max-w-3xl mx-auto px-5 sm:px-6 md:px-8 pt-8 sm:pt-12 pb-12">
        {/* Back Button */}
        <div className="mb-6 flex justify-start">
          <button
            type="button"
            onClick={handleBack}
            className="flex items-center justify-center w-9 h-9 rounded-full border border-border/60 bg-white/90 dark:bg-neutral-900/90 backdrop-blur-md hover:bg-slate-50 dark:hover:bg-neutral-800 text-muted-foreground hover:text-foreground transition-all active:scale-95 cursor-pointer shadow-md"
            aria-label="Go back"
          >
            <ArrowLeft className="w-4 h-4" />
          </button>
        </div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="space-y-8"
        >
          {/* Title */}
          <motion.div variants={itemVariants} className="space-y-3">
            <div className="flex items-center gap-2">
              <span className="px-2.5 py-1 rounded-full bg-blue-50 dark:bg-blue-950/50 text-blue-600 dark:text-blue-400 text-xs font-semibold border border-blue-100 dark:border-blue-900/50">
                Legal
              </span>
              <span className="text-xs text-muted-foreground">
                Last updated: June 12, 2026
              </span>
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
              Terms &amp; Conditions
            </h1>
            <p className="text-muted-foreground text-base leading-relaxed max-w-2xl">
              Please read these terms and conditions carefully before using the Devix platform.
            </p>
          </motion.div>

          <motion.hr variants={itemVariants} className="border-border/60" />

          {/* Sections */}
          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">1. Acceptance of Terms</h2>
            <p className="text-muted-foreground leading-relaxed">
              By accessing or using the Devix platform (&quot;Service&quot;), you agree to be bound by these Terms and Conditions. If you do not agree to all the terms, you may not access the Service. These terms apply to all users, contributors, and visitors of the platform.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">2. Description of Service</h2>
            <p className="text-muted-foreground leading-relaxed">
              Devix is a real-time social platform designed for software developers to share technical questions, concepts, build logs, and collaborate with other engineers. The platform includes features such as:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li>Posting technical questions and code snippets</li>
              <li>Sharing build logs and project progress</li>
              <li>Real-time direct messaging and chat</li>
              <li>Community discussions and upvoting</li>
              <li>Developer profiles and portfolio showcasing</li>
            </ul>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">3. User Accounts</h2>
            <p className="text-muted-foreground leading-relaxed">
              To use certain features, you must register for an account. You are responsible for maintaining the confidentiality of your account credentials and for all activities that occur under your account. You must provide accurate and complete information during registration and keep it updated.
            </p>
            <p className="text-muted-foreground leading-relaxed">
              You must be at least 13 years of age to create an account. By creating an account, you represent that you meet this age requirement.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">4. User Content</h2>
            <p className="text-muted-foreground leading-relaxed">
              You retain ownership of the content you post on Devix. By posting content, you grant Devix a non-exclusive, worldwide, royalty-free license to use, display, and distribute your content within the platform. You are solely responsible for the content you publish and must ensure it does not violate any laws or third-party rights.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">5. Prohibited Conduct</h2>
            <p className="text-muted-foreground leading-relaxed">
              You agree not to:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li>Post spam, malware, or malicious content</li>
              <li>Harass, abuse, or threaten other users</li>
              <li>Impersonate any person or entity</li>
              <li>Attempt to gain unauthorized access to the platform</li>
              <li>Use the Service for any illegal or unauthorized purpose</li>
              <li>Scrape, crawl, or use automated means to collect data without permission</li>
            </ul>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">6. Intellectual Property</h2>
            <p className="text-muted-foreground leading-relaxed">
              The Devix platform, including its design, logo, source code, and all associated intellectual property, is owned by Devix and protected by applicable copyright, trademark, and intellectual property laws. You may not reproduce, distribute, or create derivative works from any part of the platform without explicit written permission.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">7. Termination</h2>
            <p className="text-muted-foreground leading-relaxed">
              We reserve the right to suspend or terminate your account at any time, with or without cause, and with or without notice. Upon termination, your right to use the Service will immediately cease. Provisions that by their nature should survive termination shall survive, including ownership provisions, warranty disclaimers, and limitations of liability.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">8. Disclaimer of Warranties</h2>
            <p className="text-muted-foreground leading-relaxed">
              The Service is provided on an &quot;as is&quot; and &quot;as available&quot; basis without warranties of any kind, either express or implied. Devix does not warrant that the Service will be uninterrupted, error-free, or free of viruses or other harmful components.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">9. Limitation of Liability</h2>
            <p className="text-muted-foreground leading-relaxed">
              To the maximum extent permitted by law, Devix shall not be liable for any indirect, incidental, special, consequential, or punitive damages, including but not limited to loss of data, profits, or goodwill, arising out of your use of or inability to use the Service.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">10. Changes to Terms</h2>
            <p className="text-muted-foreground leading-relaxed">
              We reserve the right to modify these Terms at any time. We will notify users of significant changes via the platform or email. Your continued use of the Service after such modifications constitutes acceptance of the updated Terms.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">11. Contact Us</h2>
            <p className="text-muted-foreground leading-relaxed">
              If you have any questions about these Terms, please contact us at{" "}
              <a href="mailto:support@devix.in" className="text-blue-600 hover:text-blue-700 hover:underline font-medium transition-colors">
                support@devix.in
              </a>.
            </p>
          </motion.section>

          <motion.hr variants={itemVariants} className="border-border/60" />

          {/* Footer links */}
          <motion.div variants={itemVariants} className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <Link href="/privacy" className="text-blue-600 hover:text-blue-700 hover:underline font-medium transition-colors">
              Privacy Policy
            </Link>
          </motion.div>
        </motion.div>
      </main>
    </div>
  );
}
