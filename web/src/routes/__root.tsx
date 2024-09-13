import React from "react"
import { Outlet, createRootRoute } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@/components/tanstack-router-devtools"

export const Route = createRootRoute({
  component: () => (
    <>
      <Outlet />
      <React.Suspense>
        <TanStackRouterDevtools />
      </React.Suspense>
    </>
  ),
})
