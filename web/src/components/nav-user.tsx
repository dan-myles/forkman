import { ChevronsUpDown, LogOut } from "lucide-react"
import { useNavigate } from "@tanstack/react-router"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useUserStore } from "@/stores/userStore"

export function NavUser() {
  const user = useUserStore((state) => state.user)
  const navigate = useNavigate()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className="w-full rounded-md outline-none ring-ring hover:bg-accent
          focus-visible:ring-2 data-[state=open]:bg-accent"
      >
        <div
          className="flex items-center gap-2 px-2 py-1.5 text-left text-sm
            transition-all"
        >
          <Avatar className="h-7 w-7 rounded-md border">
            <AvatarImage
              src={user?.AvatarURL}
              alt={user?.Name}
              className="animate-in fade-in-50 zoom-in-90"
            />
            <AvatarFallback className="rounded-md">CN</AvatarFallback>
          </Avatar>
          <div className="grid flex-1 leading-none">
            <div className="font-medium">{user?.Name}</div>
            <div className="overflow-hidden text-xs text-muted-foreground">
              <div className="truncate">{user?.Email}</div>
            </div>
          </div>
          <ChevronsUpDown
            className="ml-auto mr-0.5 h-4 w-4 text-muted-foreground/50"
          />
        </div>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-56"
        align="end"
        side="right"
        sideOffset={4}
      >
        <a href="/auth/discord/logout">
          <DropdownMenuItem
            className="gap-2"
          >
            <LogOut className="h-4 w-4 text-muted-foreground" />
            Log out
          </DropdownMenuItem>
        </a>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
