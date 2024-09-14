import React from "react"
import { Outlet, createRootRoute } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@/components/tanstack-router-devtools"
import { ThemeProvider } from "@/components/theme-provider"

export const Route = createRootRoute({
  component: () => (
    <>
      <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
        <Outlet />
      </ThemeProvider>
      <React.Suspense>
        <TanStackRouterDevtools initialIsOpen={false} position="bottom-right" />
      </React.Suspense>
    </>
  ),
})
