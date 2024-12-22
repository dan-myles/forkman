import { NavMain } from "@/components/nav-main"
import { NavSelection } from "@/components/nav-selection"
import { NavTeamSwitcher } from "@/components/nav-team-switcher"
import { NavUser } from "@/components/nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarItem,
  SidebarLabel,
} from "@/components/ui/sidebar"

export const iframeHeight = "870px"

export const containerClassName = "w-full h-full"

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarHeader>
        <NavTeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <SidebarItem>
          <SidebarLabel>Platform</SidebarLabel>
          <NavMain />
        </SidebarItem>
        <SidebarItem className="mt-auto">
          <NavSelection />
        </SidebarItem>
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
    </Sidebar>
  )
}
