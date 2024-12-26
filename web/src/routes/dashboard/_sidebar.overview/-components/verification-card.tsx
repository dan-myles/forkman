"use client"

import { useCallback, useEffect, useState } from "react"
import { toast } from "sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Skeleton } from "@/components/ui/skeleton"
import { Switch } from "@/components/ui/switch"
import { api } from "@/lib/api"

export function VerificationCard() {
  const [isEnabled, setIsEnabled] = useState<boolean | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [channelId, setChannelId] = useState("")
  const [isSending, setIsSending] = useState(false)

  useEffect(() => {
    const getInitialState = async () => {
      try {
        const initialState = await api.getVerificationStatus()
        setIsEnabled(initialState)
      } catch (error) {
        console.error("Failed to fetch initial state:", error)
        toast.error(
          "Failed to load verification status. Please refresh the page."
        )
        setIsEnabled(false) // Default to disabled on error
      }
    }

    getInitialState()
  }, [])

  const handleToggle = useCallback(async () => {
    if (isLoading || isEnabled === null) return

    setIsLoading(true)
    const newState = !isEnabled
    let success = false
    if (isEnabled) {
      success = await api.disableVerification()
    } else {
      success = await api.enableVerification()
    }

    if (success) {
      setIsEnabled(newState)
      toast.success(
        `Verification ${newState ? "enabled" : "disabled"} successfully`
      )
    } else {
      toast.error(
        `Failed to ${newState ? "enable" : "disable"} verification. Please try again.`
      )
    }

    setIsLoading(false)
  }, [isEnabled, isLoading])

  const handleSendVerificationPanel = async () => {
    if (!channelId.trim()) {
      toast.error("Please enter a valid channel ID")
      return
    }

    setIsSending(true)
    try {
      await api.sendVerificationPanel(channelId)
      toast.success("Verification panel sent successfully")
      setChannelId("")
    } catch (error) {
      console.error("Failed to send verification panel:", error)
      toast.error("Failed to send verification panel. Please try again.")
    } finally {
      setIsSending(false)
    }
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader
        className={
          isEnabled
            ? `rounded-md bg-gradient-to-r from-blue-500 to-purple-600
              text-white`
            : "bg-gray-500 rounded-md"
        }
      >
        <CardTitle className="text-2xl font-bold">
          Verification Module
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6 p-6">
        {isEnabled === null ? (
          <div className="flex items-center justify-between">
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-6 w-11" />
          </div>
        ) : (
          <div className="flex items-center justify-between">
            <Label
              htmlFor="verification-toggle"
              className="text-lg font-medium"
            >
              {isEnabled ? "Enabled" : "Disabled"}
            </Label>
            <Switch
              id="verification-toggle"
              checked={isEnabled}
              onCheckedChange={handleToggle}
              disabled={isLoading}
              className={isLoading ? "cursor-not-allowed opacity-50" : ""}
            />
          </div>
        )}
        {isLoading && (
          <p className="mt-2 text-sm text-gray-500">
            Updating verification status...
          </p>
        )}
        <div className="space-y-2">
          <Label htmlFor="channel-id" className="text-sm font-medium">
            Send Verification Panel
          </Label>
          <div className="flex space-x-2">
            <Input
              id="channel-id"
              placeholder="Enter Channel ID"
              value={channelId}
              onChange={(e) => setChannelId(e.target.value)}
              disabled={isSending || !isEnabled}
            />
            <Button
              onClick={handleSendVerificationPanel}
              disabled={isSending || !isEnabled}
            >
              {isSending ? "Sending..." : "Send"}
            </Button>
          </div>
          {!isEnabled && (
            <p className="text-sm text-gray-500">
              Enable verification to send the panel
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
