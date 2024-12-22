import { createFileRoute, redirect } from "@tanstack/react-router"
import { api } from "@/lib/api"
import { StatusCard } from "./-components/status-card"
import { UptimeCard } from "./-components/uptime-card"

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
  },
  component: Page,
})

function Page() {
  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-row gap-2">
        <StatusCard />
        <UptimeCard />
      </div>
    </div>
  )
}
