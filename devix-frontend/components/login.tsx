"use client";

import * as React from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import PhoneInput from "react-phone-input-2";
import { Checkbox } from "@/components/ui/checkbox";
import { Loader2 } from "lucide-react";
import { FcGoogle } from "react-icons/fc";
import { SiGithub } from "react-icons/si";

function MicrosoftIcon() {
  return (
    <span className="inline-grid h-4 w-4 grid-cols-2 gap-0.5">
      <span className="rounded-[1px] bg-[#F25022]" />
      <span className="rounded-[1px] bg-[#7FBA00]" />
      <span className="rounded-[1px] bg-[#00A4EF]" />
      <span className="rounded-[1px] bg-[#FFB900]" />
    </span>
  );
}

// Validation schema for the form
const phoneRegex = /^[+]?[\d\s().-]{7,20}$/;
const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;

type FormSchemaOptions = {
  requireFullName: boolean;
  requireUsername: boolean;
  requireContact: boolean;
  requireDob: boolean;
};

const createFormSchema = (options: FormSchemaOptions) =>
  z
    .object({
      fullName: z.string().trim().optional(),
      username: z.string().trim().optional(),
      dob: z.string().trim().optional(),
      email: z.string().trim().optional(),
      phone: z.string().trim().optional(),
      password: z
        .string()
        .min(8, { message: "Password must be at least 8 characters." }),
      rememberMe: z.boolean().optional(),
    })
    .superRefine((data, ctx) => {
      const trimmedFullName = (data.fullName ?? "").trim();
      const trimmedUsername = (data.username ?? "").trim();
      const trimmedDob = (data.dob ?? "").trim();
      const hasEmail = Boolean(data.email && data.email.length > 0);
      const phoneDigits = (data.phone ?? "").replace(/\D/g, "");
      const hasPhone = phoneDigits.length >= 7;

      if (options.requireFullName && trimmedFullName.length < 2) {
        ctx.addIssue({
          code: "custom",
          message: "Please enter your full name.",
          path: ["fullName"],
        });
      }

      if (options.requireUsername && trimmedUsername.length < 3) {
        ctx.addIssue({
          code: "custom",
          message: "Please choose a username.",
          path: ["username"],
        });
      }

      if (trimmedUsername.length > 0 && !usernameRegex.test(trimmedUsername)) {
        ctx.addIssue({
          code: "custom",
          message: "Username can use letters, numbers, and underscores.",
          path: ["username"],
        });
      }

      if (options.requireDob && trimmedDob.length === 0) {
        ctx.addIssue({
          code: "custom",
          message: "Please select your date of birth.",
          path: ["dob"],
        });
      }

      if (options.requireContact && !hasEmail && !hasPhone) {
        ctx.addIssue({
          code: "custom",
          message: "Enter an email or phone number.",
          path: ["email"],
        });
      }

      if (hasEmail && !z.email().safeParse(data.email).success) {
        ctx.addIssue({
          code: "custom",
          message: "Please enter a valid email.",
          path: ["email"],
        });
      }

      if (!hasEmail && phoneDigits.length > 0 && phoneDigits.length < 7) {
        ctx.addIssue({
          code: "custom",
          message: "Please enter a valid phone number.",
          path: ["phone"],
        });
      }

      if (hasPhone && !phoneRegex.test(data.phone ?? "")) {
        ctx.addIssue({
          code: "custom",
          message: "Please enter a valid phone number.",
          path: ["phone"],
        });
      }
    });

export type FormValues = z.infer<ReturnType<typeof createFormSchema>>;

interface AuthFormSplitScreenProps {
  logo: React.ReactNode;
  title: string;
  description: string;
  rightPanel?: React.ReactNode;
  imageSrc?: string;
  imageAlt?: string;
  showRightPanel?: boolean;
  showSignupFields?: boolean;
  loginWithUsername?: boolean;
  onSubmit: (data: FormValues) => Promise<void>;
  forgotPasswordHref: string;
  createAccountHref: string;
}

/**
 * A responsive, split-screen authentication form component.
 * @param logo - The component to be used as the logo (e.g., an SVG or text).
 * @param title - The main heading for the form.
 * @param description - A short description below the title.
 * @param rightPanel - Optional custom right panel component to render instead of image.
 * @param imageSrc - URL for the image to display on the right panel.
 * @param imageAlt - Alt text for the image for accessibility.
 * @param showRightPanel - Controls whether the right panel is shown.
 * @param showSignupFields - Controls whether full name and username fields are shown.
 * @param loginWithUsername - Shows a username-only login (no email/phone).
 * @param onSubmit - Async function to handle form submission.
 * @param forgotPasswordHref - URL for the "Forgot Password" link.
 * @param createAccountHref - URL for the "Create Account" link.
 */
export function AuthFormSplitScreen({
  logo,
  title,
  description,
  rightPanel,
  imageSrc,
  imageAlt,
  showRightPanel = true,
  showSignupFields = false,
  loginWithUsername = false,
  onSubmit,
  forgotPasswordHref,
  createAccountHref,
}: AuthFormSplitScreenProps) {
  const [isLoading, setIsLoading] = React.useState(false);
  const hasRightPanel = Boolean(showRightPanel && (rightPanel || imageSrc));
  const showUsernameField = showSignupFields || loginWithUsername;
  const showContactFields = !loginWithUsername;
  const showLoginExtras = !showSignupFields;
  const formSchema = React.useMemo(
    () =>
      createFormSchema({
        requireFullName: showSignupFields,
        requireUsername: showUsernameField,
        requireContact: showContactFields,
        requireDob: showSignupFields,
      }),
    [showSignupFields, showUsernameField, showContactFields]
  );

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      fullName: "",
      username: "",
      dob: "",
      email: "",
      phone: "",
      password: "",
      rememberMe: false,
    },
  });

  const handleFormSubmit = async (data: FormValues) => {
    setIsLoading(true);
    try {
      await onSubmit(data);
    } catch (error) {
      console.error("Submission failed:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Animation variants for staggering children
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { y: 0, opacity: 1 },
  };

  return (
    <div
      className={`relative flex min-h-screen w-full flex-col${
        hasRightPanel ? " md:flex-row" : ""
      }`}
    >
      {/* Left Panel: Form */}
      <div
        className={`relative flex w-full flex-col items-center justify-start bg-linear-to-tr from-white via-sky-50/10 to-blue-50/30 p-8 pt-20 pb-12 overflow-hidden${
          hasRightPanel ? " md:w-1/2" : ""
        }`}
      >
        {/* Soft Grid Background */}
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#0284c703_1px,transparent_1px),linear-gradient(to_bottom,#0284c703_1px,transparent_1px)] bg-size-[32px_32px] pointer-events-none" />

        <div className="w-full max-w-lg relative z-10">
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
            className="flex flex-col gap-4"
          >
            <motion.div variants={itemVariants} className="mb-2 flex justify-center">
              {logo}
            </motion.div>
            <motion.div variants={itemVariants} className="text-left">
              <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
              <p className="text-sm text-muted-foreground">{description}</p>
            </motion.div>

            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(handleFormSubmit)}
                className="space-y-3"
              >
                {showSignupFields && (
                  <div className="grid gap-3 md:grid-cols-2">
                    <motion.div variants={itemVariants}>
                      <FormField
                        control={form.control}
                        name="fullName"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Full Name</FormLabel>
                            <FormControl>
                              <Input
                                placeholder="Enter your full name"
                                {...field}
                                disabled={isLoading}
                                className="focus-visible:ring-blue-500"
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </motion.div>

                    {showUsernameField && (
                      <motion.div variants={itemVariants}>
                        <FormField
                          control={form.control}
                          name="username"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Username</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="your_handle"
                                  {...field}
                                  disabled={isLoading}
                                  className="focus-visible:ring-blue-500"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </motion.div>
                    )}
                  </div>
                )}

                {!showSignupFields && showUsernameField && (
                  <motion.div variants={itemVariants}>
                    <FormField
                      control={form.control}
                      name="username"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Username</FormLabel>
                          <FormControl>
                            <Input
                              placeholder="your_handle"
                              {...field}
                              disabled={isLoading}
                              className="focus-visible:ring-blue-500"
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </motion.div>
                )}

                {showContactFields && (
                  <div className="space-y-3">
                    <div className={`grid gap-3 ${showSignupFields ? "md:grid-cols-2" : ""}`}>
                      <motion.div variants={itemVariants}>
                        <FormField
                          control={form.control}
                          name="email"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Email Address</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="email@example.com"
                                  {...field}
                                  disabled={isLoading}
                                  className="focus-visible:ring-blue-500"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </motion.div>

                      {showSignupFields && (
                        <motion.div variants={itemVariants}>
                          <FormField
                            control={form.control}
                            name="dob"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Date of Birth</FormLabel>
                                <FormControl>
                                  <Input
                                    type="date"
                                    {...field}
                                    disabled={isLoading}
                                    className="focus-visible:ring-blue-500"
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </motion.div>
                      )}
                    </div>

                    {showSignupFields ? (
                      <motion.div variants={itemVariants}>
                        <FormField
                          control={form.control}
                          name="phone"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Phone Number (optional)</FormLabel>
                              <FormControl>
                                <PhoneInput
                                  country="in"
                                  value={field.value ?? ""}
                                  onChange={(value) => field.onChange(value)}
                                  enableSearch
                                  placeholder="Enter phone number"
                                  inputProps={{
                                    name: field.name,
                                    disabled: isLoading,
                                  }}
                                  containerClass="!w-full"
                                  inputClass="!w-full !h-10 !text-sm !bg-background !border !border-input !rounded-md !pl-12 !pr-3 !py-2 !text-foreground focus:!ring-2 focus:!ring-ring focus:!ring-offset-2"
                                  buttonClass="!border-input !bg-background !rounded-md !rounded-r-none !h-10 !w-11"
                                  dropdownClass="!text-sm"
                                  searchClass="!text-sm !py-2"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </motion.div>
                    ) : (
                      <motion.div variants={itemVariants}>
                        <FormField
                          control={form.control}
                          name="phone"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Phone Number (optional)</FormLabel>
                              <FormControl>
                                <Input
                                  type="tel"
                                  placeholder="+1 555 000 1234"
                                  {...field}
                                  disabled={isLoading}
                                  className="focus-visible:ring-blue-500"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </motion.div>
                    )}
                  </div>
                )}

                <motion.div variants={itemVariants}>
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Password</FormLabel>
                        <FormControl>
                          <Input
                            type="password"
                            placeholder="••••••••••••"
                            {...field}
                            disabled={isLoading}
                            className="focus-visible:ring-blue-500"
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </motion.div>

                {showLoginExtras && (
                  <motion.div
                    variants={itemVariants}
                    className="flex items-center justify-between"
                  >
                    <FormField
                      control={form.control}
                      name="rememberMe"
                      render={({ field }) => (
                        <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                          <FormControl>
                            <Checkbox
                              checked={field.value}
                              onCheckedChange={field.onChange}
                              disabled={isLoading}
                              className="border-slate-350 data-[state=checked]:bg-blue-600 data-[state=checked]:border-blue-600"
                            />
                          </FormControl>
                          <div className="space-y-1 leading-none">
                            <FormLabel className="font-normal">
                              Remember Me
                            </FormLabel>
                          </div>
                        </FormItem>
                      )}
                    />
                    <a
                      href={forgotPasswordHref}
                      className="text-sm font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
                    >
                      Forgotten Password
                    </a>
                  </motion.div>
                )}

                <motion.div variants={itemVariants}>
                  <Button
                    type="submit"
                    className="w-full bg-linear-to-r from-blue-600 to-sky-600 hover:from-blue-700 hover:to-sky-700 text-white shadow-md shadow-blue-500/10 hover:shadow-blue-500/20 active:scale-[0.98] transition-all cursor-pointer"
                    disabled={isLoading}
                  >
                    {isLoading && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Continue
                  </Button>
                </motion.div>

                {showSignupFields && (
                  <motion.p
                    variants={itemVariants}
                    className="text-center text-[11px] text-muted-foreground leading-relaxed"
                  >
                    By signing up, you agree to our{" "}
                    <Link
                      href="/terms"
                      className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
                    >
                      Terms &amp; Conditions
                    </Link>{" "}
                    and{" "}
                    <Link
                      href="/privacy"
                      className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
                    >
                      Privacy Policy
                    </Link>
                    .
                  </motion.p>
                )}

                <motion.div variants={itemVariants} className="flex items-center gap-3">
                  <div className="h-px flex-1 bg-border" />
                  <span className="text-[11px] text-muted-foreground">or continue with</span>
                  <div className="h-px flex-1 bg-border" />
                </motion.div>

                <motion.div variants={itemVariants} className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <Button
                    type="button"
                    variant="outline"
                    className="gap-2"
                    disabled={isLoading}
                  >
                    <FcGoogle className="h-4 w-4" />
                    Google
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    className="gap-2"
                    disabled={isLoading}
                  >
                    <SiGithub className="h-4 w-4 text-[#181717]" />
                    GitHub
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    className="gap-2"
                    disabled={isLoading}
                  >
                    <MicrosoftIcon />
                    Microsoft
                  </Button>
                </motion.div>
              </form>
            </Form>

            <motion.p
              variants={itemVariants}
              className="px-8 text-center text-sm text-muted-foreground"
            >
              {createAccountHref === "/login" ? (
                <>
                  Already have an account?{" "}
                  <a
                    href={createAccountHref}
                    className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
                  >
                    Sign in here
                  </a>
                </>
              ) : (
                <>
                  Don&apos;t have an account?{" "}
                  <a
                    href={createAccountHref}
                    className="font-medium text-blue-600 hover:text-blue-700 hover:underline transition-colors"
                  >
                    Create one here
                  </a>
                </>
              )}
              .
            </motion.p>
          </motion.div>
        </div>
      </div>

      {hasRightPanel && (
        <div className="relative hidden w-1/2 md:block">
          {rightPanel ? (
            rightPanel
          ) : (
            <>
              {imageSrc && (
                <img
                  src={imageSrc}
                  alt={imageAlt}
                  className="h-full w-full object-cover"
                />
              )}
              <div className="absolute inset-0 bg-linear-to-t from-black/20 to-transparent" />
            </>
          )}
        </div>
      )}
    </div>
  );
}