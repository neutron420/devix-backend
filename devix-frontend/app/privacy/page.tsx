"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { ArrowLeft } from "lucide-react";

export default function PrivacyPolicyPage() {
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
              <span className="px-2.5 py-1 rounded-full bg-emerald-50 dark:bg-emerald-950/50 text-emerald-600 dark:text-emerald-400 text-xs font-semibold border border-emerald-100 dark:border-emerald-900/50">
                Privacy
              </span>
              <span className="text-xs text-muted-foreground">
                Last updated: June 12, 2026
              </span>
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
              Privacy Policy
            </h1>
            <p className="text-muted-foreground text-base leading-relaxed max-w-2xl">
              Your privacy is important to us. This policy explains how Devix collects, uses, and protects your personal data.
            </p>
          </motion.div>

          <motion.hr variants={itemVariants} className="border-border/60" />

          {/* Sections */}
          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">1. Information We Collect</h2>
            <p className="text-muted-foreground leading-relaxed">
              We collect information you provide directly when you create an account, post content, or communicate with other users. This includes:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li><strong className="text-foreground">Account Information:</strong> Full name, username, email address, phone number, date of birth, and password</li>
              <li><strong className="text-foreground">Profile Information:</strong> Bio, profile picture, and social links</li>
              <li><strong className="text-foreground">Content:</strong> Posts, comments, code snippets, messages, and build logs you create</li>
              <li><strong className="text-foreground">Usage Data:</strong> IP address, browser type, device information, and interaction data</li>
            </ul>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">2. How We Use Your Information</h2>
            <p className="text-muted-foreground leading-relaxed">
              We use the information we collect to:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li>Provide, maintain, and improve the Devix platform</li>
              <li>Create and manage your account</li>
              <li>Enable real-time messaging and collaboration features</li>
              <li>Send you important notifications about your account and activity</li>
              <li>Detect and prevent fraud, abuse, and security issues</li>
              <li>Analyze usage patterns to improve user experience</li>
            </ul>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">3. Information Sharing</h2>
            <p className="text-muted-foreground leading-relaxed">
              We do not sell your personal data. We may share your information in the following circumstances:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li><strong className="text-foreground">Public Content:</strong> Posts, comments, and profile information you make public are visible to other users</li>
              <li><strong className="text-foreground">Service Providers:</strong> With trusted third-party services that help us operate the platform (hosting, analytics, email delivery)</li>
              <li><strong className="text-foreground">Legal Requirements:</strong> When required by law, court order, or governmental request</li>
              <li><strong className="text-foreground">Safety:</strong> To protect the rights, safety, and property of Devix and its users</li>
            </ul>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">4. Data Security</h2>
            <p className="text-muted-foreground leading-relaxed">
              We implement industry-standard security measures to protect your personal information, including encryption in transit (TLS/SSL), secure password hashing, and regular security audits. However, no method of transmission or storage is 100% secure, and we cannot guarantee absolute security.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">5. Cookies &amp; Tracking</h2>
            <p className="text-muted-foreground leading-relaxed">
              We use cookies and similar technologies to maintain your login session, remember your preferences, and analyze platform usage. You can control cookie settings through your browser, though disabling cookies may limit some functionality.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">6. Data Retention</h2>
            <p className="text-muted-foreground leading-relaxed">
              We retain your personal data for as long as your account is active or as needed to provide services. If you delete your account, we will remove your personal data within 30 days, except where retention is required by law or for legitimate business purposes (such as fraud prevention).
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">7. Your Rights</h2>
            <p className="text-muted-foreground leading-relaxed">
              Depending on your jurisdiction, you may have the following rights regarding your personal data:
            </p>
            <ul className="list-disc pl-6 space-y-1.5 text-muted-foreground">
              <li><strong className="text-foreground">Access:</strong> Request a copy of the data we hold about you</li>
              <li><strong className="text-foreground">Correction:</strong> Request correction of inaccurate or incomplete data</li>
              <li><strong className="text-foreground">Deletion:</strong> Request deletion of your personal data</li>
              <li><strong className="text-foreground">Portability:</strong> Request a machine-readable copy of your data</li>
              <li><strong className="text-foreground">Objection:</strong> Object to certain types of data processing</li>
            </ul>
            <p className="text-muted-foreground leading-relaxed">
              To exercise any of these rights, contact us at{" "}
              <a href="mailto:privacy@devix.in" className="text-blue-600 hover:text-blue-700 hover:underline font-medium transition-colors">
                privacy@devix.in
              </a>.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">8. Children&apos;s Privacy</h2>
            <p className="text-muted-foreground leading-relaxed">
              Devix is not intended for children under 13 years of age. We do not knowingly collect personal data from children under 13. If we become aware that we have collected data from a child under 13, we will take steps to delete it promptly.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">9. Changes to This Policy</h2>
            <p className="text-muted-foreground leading-relaxed">
              We may update this Privacy Policy from time to time. We will notify you of significant changes via the platform or email. The &quot;Last updated&quot; date at the top reflects the most recent revision.
            </p>
          </motion.section>

          <motion.section variants={itemVariants} className="space-y-3">
            <h2 className="text-xl font-semibold text-foreground">10. Contact Us</h2>
            <p className="text-muted-foreground leading-relaxed">
              If you have any questions about this Privacy Policy or our data practices, please contact us at{" "}
              <a href="mailto:privacy@devix.in" className="text-blue-600 hover:text-blue-700 hover:underline font-medium transition-colors">
                privacy@devix.in
              </a>.
            </p>
          </motion.section>

          <motion.hr variants={itemVariants} className="border-border/60" />

          {/* Footer links */}
          <motion.div variants={itemVariants} className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <Link href="/terms" className="text-blue-600 hover:text-blue-700 hover:underline font-medium transition-colors">
              Terms &amp; Conditions
            </Link>
          </motion.div>
        </motion.div>
      </main>
    </div>
  );
}
