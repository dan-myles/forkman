"use client"

import { useCallback, useEffect, useState } from "react"
import { toast } from "sonner"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Skeleton } from "@/components/ui/skeleton"
import { Switch } from "@/components/ui/switch"
import { api } from "@/lib/api"

export function ModerationCard() {
  const [isEnabled, setIsEnabled] = useState<boolean | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    const getInitialState = async () => {
      try {
        const initialState = await api.getModerationStatus()
        setIsEnabled(initialState)
      } catch (error) {
        console.error("Failed to fetch initial state:", error)
        toast.error(
          "Failed to load Moderation status. Please refresh the page."
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
      success = await api.disableModeration()
    } else {
      success = await api.enableModeration()
    }

    if (success) {
      setIsEnabled(newState)
      toast.success(
        `Moderation module ${newState ? "enabled" : "disabled"} successfully`
      )
    } else {
      toast.error(
        `Failed to ${newState ? "enable" : "disable"} Moderation module. Please try again.`
      )
    }

    setIsLoading(false)
  }, [isEnabled, isLoading])

  return (
    <Card className="max-h-fit w-full max-w-md">
      <CardHeader
        className={
          isEnabled
            ? `rounded-md bg-gradient-to-r from-blue-500 to-purple-600
              text-white`
            : "rounded-md bg-gray-500"
        }
      >
        <CardTitle className="text-2xl font-bold">Moderation Module</CardTitle>
      </CardHeader>
      <CardContent className="p-6">
        {isEnabled === null ? (
          <div className="flex items-center justify-between">
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-6 w-11" />
          </div>
        ) : (
          <div className="flex items-center justify-between">
            <Label htmlFor="moderation-toggle" className="text-lg font-medium">
              {isEnabled ? "Enabled" : "Disabled"}
            </Label>
            <Switch
              id="moderation-toggle"
              checked={isEnabled}
              onCheckedChange={handleToggle}
              disabled={isLoading}
              className={isLoading ? "cursor-not-allowed opacity-50" : ""}
            />
          </div>
        )}
        {isLoading && (
          <p className="mt-2 text-sm text-gray-500">
            Updating Moderation module status...
          </p>
        )}
      </CardContent>
    </Card>
  )
}
