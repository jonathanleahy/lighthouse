import DashboardGrid from '@/components/dashboard-grid'
import Link from "next/link";
import {ArrowLeft, HelpCircle} from "lucide-react";
import {Button} from "@/components/ui/button";
import {ThemeToggle} from "@/components/theme-toggle";
import React from "react";

export default function Home() {
  return (
      <div>
          {/*<div className="min-h-screen bg-background mb-6">*/}
          <header
              className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
              <div className="container mx-auto px-4 sm:px-4 lg:px-4 flex h-12 items-center justify-between">
                  <div className="flex items-center space-x-4">
                      <Link href="/" className="flex items-center space-x-2">
                          <span className="font-medium">Home Dashboard</span>
                      </Link>
                  </div>
                  <div className="flex items-center space-x-4">
                      <ThemeToggle/>
                  </div>
              </div>
          </header>
          <main className="min-h-screen pb-8">
              <div className="container mx-auto px-4">
                  <h1 className="text-3xl font-bold tracking-tight mb-4 py-4">Microservices Dashboard</h1>
                  <DashboardGrid/>
              </div>
          </main></div>
          )
          }

