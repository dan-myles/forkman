"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { api } from "@/lib/api"

export function UptimeCard() {
  const [uptime, setUptime] = useState<number>(0)
  const days = Math.floor(uptime / 86400)
  const hours = Math.floor((uptime % 86400) / 3600)
  const minutes = Math.floor((uptime % 3600) / 60)
  const seconds = Math.floor(uptime % 60)

  useEffect(() => {
    const interval = setInterval(() => {
      if (uptime) {
        setUptime(uptime + 1)
      }
    }, 1000)

    return () => clearInterval(interval)
  }, [uptime])

  useEffect(() => {
    const init = async () => {
      const res = await api.getUptime()
      setUptime(res.uptime)
    }
    init()
  }, [])

  // const progressValue = (uptime / (30 * 24 * 60 * 60)) * 100 // Assuming 30 days as max

  return (
    <Card className="max-w-md overflow-hidden">
      <CardHeader
        className="bg-gradient-to-r from-blue-500 to-purple-600 text-white"
      >
        <CardTitle className="text-center text-2xl font-bold">
          Uptime Monitor
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          {[
            { label: "Days", value: days },
            { label: "Hours", value: hours },
            { label: "Minutes", value: minutes },
            { label: "Seconds", value: seconds },
          ].map(({ label, value }) => (
            <div key={label} className="text-center">
              <div className="mb-1 text-3xl font-bold tabular-nums">
                {value < 10 ? `0${value}` : value}
              </div>
              <div className="text-sm text-gray-500">{label}</div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
