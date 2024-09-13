import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@/components/tanstack-router-devtools";
import React from "react";

export const Route = createRootRoute({
  component: () => (
    <>
      <Outlet />
      <React.Suspense>
        <TanStackRouterDevtools />
      </React.Suspense>
    </>
  ),
});
