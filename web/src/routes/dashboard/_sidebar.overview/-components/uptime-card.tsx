import { useEffect, useState } from "react"
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

  return (
    <div>
      Uptime Badge
      <div>
        <p>Days: {days < 10 ? "0" + days : days}</p>
        <p>Hours: {hours < 10 ? "0" + hours : hours}</p>
        <p>Minutes: {minutes < 10 ? "0" + minutes : minutes}</p>
        <p>Seconds: {seconds < 10 ? "0" + seconds : seconds}</p>
      </div>
    </div>
  )
}
