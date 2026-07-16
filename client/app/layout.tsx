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
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin=""
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Gorditas:wght@400;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="gorditas">
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
