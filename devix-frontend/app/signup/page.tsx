"use client";

import { AuthFormSplitScreen, type FormValues } from "../../components/login";
import Image from "next/image";
import { useRouter } from "next/navigation";

export default function SignupPage() {
  const router = useRouter();

  const handleSubmit = async (data: FormValues) => {
    // Simulate API call delay for registration
    await new Promise((resolve) => setTimeout(resolve, 1500));
    console.log("Signed up successfully:", data);
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
      title="Create an Account"
      description="Join the Devix community of developers to share, discuss, and build"
      showRightPanel={false}
      showSignupFields
      onSubmit={handleSubmit}
      forgotPasswordHref="#"
      createAccountHref="/login"
    />
  );
}
