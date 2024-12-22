import { DiscordLogoIcon } from "@radix-ui/react-icons"
import { createFileRoute, redirect } from "@tanstack/react-router"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { api } from "@/lib/api"

export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    const auth = await api.isAuth()
    if (auth) {
      throw redirect({
        to: "/dashboard/overview",
        search: {
          redirect: location.href,
        },
      })
    }
  },
  component: Page,
})

export default function Page() {
  return (
    <div className="flex h-screen items-center justify-center">
      <Card className="mx-auto max-w-sm">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Forkman</CardTitle>
          <CardDescription className="text-sm">
            Login with your Discord account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <a href="/auth/discord/login">
            <Button className="w-full" variant="outline">
              <DiscordLogoIcon className="mr-2 h-4 w-4" />
              Continue with Discord
            </Button>
          </a>
          <p className="mt-4 text-center text-sm text-muted-foreground">
            Don&apos;t have an account?{" "}
            <a href="https://discord.com/register" className="underline">
              Sign up
            </a>
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
