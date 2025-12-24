import SignInForm from "@/components/auth/SignInForm";
import { Metadata } from "next";
import { Suspense } from "react";

export const metadata: Metadata = {
  title: "Sign In - Permit App",
  description: "Sign in to your Permit App account to manage your permissions and access control.",
};

export default function SignIn() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center flex-1"><div className="text-center"><div className="inline-block w-8 h-8 border-4 border-gray-300 border-t-brand-500 rounded-full animate-spin"></div></div></div>}>
      <SignInForm />
    </Suspense>
  );
}
