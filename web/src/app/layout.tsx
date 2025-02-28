import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import Navbar from "@/components/Navbar";

const inter = Inter({ subsets: ["latin"] });
const globalStyles = 'min-h-screen font-sans antialiased';

export const metadata: Metadata = {
  title: "Data Explorer",
  description: "Uncover insights in your data",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="light">
      <body className={cn(globalStyles, inter.className)}>
        <Navbar/>
        {children}
      </body>
    </html>
  );
}
