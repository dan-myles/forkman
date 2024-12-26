"use client"

import { useEffect, useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
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
    }, 30000)

    return () => clearInterval(interval)
  }, [])

  return (
    <Card className="max-w-md overflow-hidden">
      <CardHeader
        className="bg-gradient-to-r from-blue-500 to-purple-600 text-white"
      >
        <CardTitle className="text-center text-2xl font-bold">
          Server Status
        </CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        <div className="flex flex-col items-center justify-center space-y-4">
          <p
            className={`text-lg font-semibold
              ${isOnline ? "text-green-600" : "text-red-600"}`}
          >
            {isOnline ? "Online" : "Offline"}
          </p>
          <Dialog>
            <DialogTrigger asChild>
              <Button variant="outline" className="mt-4">
                Restart Server
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogTitle>Are you sure?</DialogTitle>
              <DialogHeader>
                <DialogDescription>
                  This will completely restart the server. This can take
                  anywhere from a few seconds to a few minutes. Please ensure
                  this is what you want to do.
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <DialogClose asChild>
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
                <DialogClose asChild>
                  <Button variant="outline">Cancel</Button>
                </DialogClose>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </CardContent>
    </Card>
  )
}
