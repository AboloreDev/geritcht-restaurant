import type { Metadata } from "next";
import "./globals.css";
import Providers from "./providers";
import { Figtree } from "next/font/google";
import { cn } from "@/lib/utils";
import { Toaster } from "sonner";

const figtree = Figtree({ subsets: ["latin"], variable: "--font-sans" });

export const metadata: Metadata = {
  title: "Geritcht Restaurant",
  description:
    "Book a table, order takeout, and track your order in real time at Geritcht Restaurant.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html className={cn("font-sans", figtree.variable)}>
      {" "}
      <link rel="preconnect" href="https://fonts.googleapis.com" />
      <link
        rel="preconnect"
        href="https://fonts.gstatic.com"
        crossOrigin="anonymous"
      />
      <link
        href="https://fonts.googleapis.com/css2?family=Lato:ital,wght@0,100;0,300;0,400;0,700;0,900;1,100;1,300;1,400;1,700;1,900&display=swap"
        rel="stylesheet"
      />
      <body className="min-h-full flex flex-col">
        <Providers>
          <Toaster
            position="top-right"
            duration={3000}
            closeButton
            richColors
          />
          {children}
        </Providers>
      </body>
    </html>
  );
}
