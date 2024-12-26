import { LayoutDashboard } from "lucide-react"
import { useEffect, useState } from "react"
import { DiscordLogoIcon } from "@radix-ui/react-icons"
import { createFileRoute, Link } from "@tanstack/react-router"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { api } from "@/lib/api"
import { useUserStore } from "@/stores/userStore"

export const Route = createFileRoute("/")({
  component: Page,
})

export default function Page() {
  const [isAuth, setIsAuth] = useState(false)
  const user = useUserStore((state) => state.user)

  useEffect(() => {
    const init = async () => {
      const auth = await api.isAuth()
      setIsAuth(auth)
    }
    init()
  }, [])

  return (
    <div className="flex h-screen items-center justify-center">
      <Card className="mx-auto max-w-sm">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">Forkman</CardTitle>
          <CardDescription className="text-sm">
            {isAuth ? "Welcome back " + user?.Name : "Login to continue"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isAuth ? (
            <div className="mb-4">
              <Link
                to="/dashboard/overview"
              >

                <Button
                  className="w-full bg-green-300 text-black"
                  variant="outline"
                >
                  <LayoutDashboard className="mr-2 h-4 w-4" />
                  Go to Dashboard
                </Button>
              </Link>
            </div>
          ) : (
            <>
              <a href="/auth/discord/login">
                <Button className="w-full" variant="outline">
                  <DiscordLogoIcon className="mr-2 h-4 w-4" />
                  Login with Discord
                </Button>
              </a>
              <p className="mt-4 text-center text-sm text-muted-foreground">
                Don&apos;t have an account?{" "}
                <a href="https://discord.com/register" className="underline">
                  Sign up
                </a>
              </p>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
