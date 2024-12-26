import { Outlet, createFileRoute, useLocation } from "@tanstack/react-router"
import { AppSidebar } from "@/components/app-sidebar"
import { SidebarLayout, SidebarTrigger } from "@/components/ui/sidebar"
import { useSidebar } from "@/hooks/use-sidebar"

export const Route = createFileRoute("/dashboard/_sidebar")({
  component: Layout,
})

function Layout() {
  const { open } = useSidebar()
  const location = useLocation().pathname.split("/").pop()
  if (!location) return null
  const title = location.charAt(0).toUpperCase() + location.slice(1)

  return (
    <SidebarLayout defaultOpen={open}>
      <AppSidebar />
      <main
        className="flex flex-1 flex-col p-2 transition-all duration-300
          ease-in-out"
      >
        <div className="flex flex-row items-center space-y-2 rounded-md p-2">
          <SidebarTrigger />
          <div className="pb-2 pl-2 text-2xl">{title}</div>
        </div>
        <Outlet />
      </main>
    </SidebarLayout>
  )
}
