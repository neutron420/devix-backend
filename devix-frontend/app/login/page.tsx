"use client";

import { AuthFormSplitScreen, type FormValues } from "../../components/login";
import Image from "next/image";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const router = useRouter();

  const handleSubmit = async (data: FormValues) => {
    // Simulate API call delay for authentication
    await new Promise((resolve) => setTimeout(resolve, 1500));
    console.log("Logged in successfully:", data);
    router.push("/");
  };

  return (
    <AuthFormSplitScreen
      logo={
        <Image
          src="/devix1.png"
          alt="Devix Logo"
          width={220}
          height={70}
          className="h-16 w-auto object-contain"
          priority
        />
      }
      title="Welcome Back"
      description="Enter your username below to sign in to your Devix account"
      showRightPanel={false}
      loginWithUsername
      onSubmit={handleSubmit}
      forgotPasswordHref="/forgot-password"
      createAccountHref="/signup"
    />
  );
}
