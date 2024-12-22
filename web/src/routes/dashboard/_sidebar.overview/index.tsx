import { createFileRoute, redirect } from "@tanstack/react-router"
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
} from "@/components/ui/card"
import { api } from "@/lib/api"
import { StatusBadge as StatusCard } from "./-components/status-badge"

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
        <Card className="flex flex-row p-2">
          <CardTitle className="p text-lg">Servers</CardTitle>
          <CardContent>
            <CardDescription className="p"></CardDescription>
          </CardContent>
          <CardContent>
            <CardDescription className="p"></CardDescription>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
