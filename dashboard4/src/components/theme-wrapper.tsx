'use client'

import { useState, useLayoutEffect } from 'react'
import { ThemeProvider } from "@/components/theme-provider"

export function ThemeWrapper({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false)

  useLayoutEffect(() => {
    setMounted(true)
  }, [])

  // Render the children without the ThemeProvider on the server
  if (!mounted) {
    return <>{children}</>
  }

  return (
    <div className={`contents ${mounted ? '' : 'invisible'}`}>
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
      >
        {children}
      </ThemeProvider>
    </div>
  )
}


