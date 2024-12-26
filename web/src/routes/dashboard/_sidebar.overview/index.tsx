import { createFileRoute, redirect } from "@tanstack/react-router"
import { api } from "@/lib/api"
import { ModerationCard } from "./-components/moderation-card"
import { QNACard } from "./-components/qna-card"
import { StatusCard } from "./-components/status-card"
import { UptimeCard } from "./-components/uptime-card"
import { VerificationCard } from "./-components/verification-card"

export const Route = createFileRoute("/dashboard/_sidebar/overview/")({
  beforeLoad: async () => {
    const auth = await api.isAuth()
    if (!auth) {
      throw redirect({
        to: "/",
        search: {
          redirect: location.href,
        },
      })
    }

    const admin = await api.isAdmin()
    if (!admin) {
      throw redirect({
        to: "/unauthorized",
        search: {
          redirect: location.href,
        },
      })
    }
  },
  component: Page,
})

function Page() {
  return (
    <div className="container mx-auto space-y-6 p-4">
      <div className="grid grid-cols-1 gap-2 md:grid-cols-2">
        <StatusCard />
        <UptimeCard />
      </div>
      <div className="grid grid-cols-1 gap-2 md:grid-cols-2">
        <VerificationCard />
        <QNACard />
      </div>
      <div className="grid grid-cols-1 gap-2 md:grid-cols-2">
        <ModerationCard />
      </div>
    </div>
  )
}
