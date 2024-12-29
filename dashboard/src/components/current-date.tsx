"use client"

import { useState, useEffect } from 'react'

export function CurrentDate() {
  const [date, setDate] = useState<string | null>(null)

  useEffect(() => {
    setDate(new Date().toLocaleDateString())
  }, [])

  if (!date) return null

  return <div>{date}</div>
}


