import { Power } from "lucide-react"
import { useEffect, useState } from "react"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { api } from "@/lib/api"

export function StatusCard() {
  const [isOnline, setIsOnline] = useState(true)

  useEffect(() => {
    const interval = setInterval(async () => {
      const online = await api.getStatus()
      setIsOnline(online)
    }, 1000)

    return () => clearInterval(interval)
  }, [])

  return (
    <Card className="max-w-52 p-4 flex-grow">
      <CardTitle className="flex justify-between pb-2 text-lg">
        <p>Status</p>
        <div>
          <Dialog>
            <DialogTrigger>
              <Power className="h-4 w-4" />
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Are you sure?</DialogTitle>
                <DialogDescription>
                  <p>
                    This will completely restart the server. This can take
                    anywhere from a few seconds to a few minutes. Please ensure
                    this is what you want to do.
                  </p>
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <DialogClose>
                  <Button
                    variant="destructive"
                    onClick={() => {
                      api.restartServer()
                      toast.success("Restarting server...")
                    }}
                  >
                    Restart
                  </Button>
                </DialogClose>
                <DialogClose>
                  <Button variant="outline">Cancel</Button>
                </DialogClose>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </CardTitle>
      <CardHeader className="p-2">
        {isOnline ? (
          <Badge
            className="flex h-12 items-center justify-center bg-green-400
              hover:bg-green-400"
          >
            <p className="animate-pulse">Online</p>
          </Badge>
        ) : (
          <Badge
            className="flex h-12 items-center justify-center bg-red-400
              hover:bg-red-400"
          >
            <p>Offline</p>
          </Badge>
        )}
      </CardHeader>
    </Card>
  )
}
