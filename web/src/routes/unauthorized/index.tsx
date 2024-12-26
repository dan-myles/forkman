import { AlertTriangle } from "lucide-react"
import { createFileRoute } from "@tanstack/react-router"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { api } from "@/lib/api"

export const Route = createFileRoute("/unauthorized/")({
  component: Page,
})

function Page() {
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div
            className="mx-auto mb-4 flex h-16 w-16 items-center justify-center
              rounded-full bg-red-400"
          >
            <AlertTriangle className="h-8 w-8 text-black" />
          </div>
          <CardTitle className="text-2xl font-bold">
            Unauthorized Access
          </CardTitle>
        </CardHeader>
        <CardContent className="text-center">
          <p className="mb-4 text-gray-400">
            Oops! It seems you don't have permission to access this page. If you
            believe this is an error, please contact the administrator.
          </p>
          <Button
            className="mt-2"
            onClick={async () => {
              await api.logout()
              location.href = "/"
            }}
          >
            Logout
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}

export default Page
