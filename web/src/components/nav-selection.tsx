import { useCurrentServerStore } from "@/stores/currentServerStore"

export function NavSelection() {
  const server = useCurrentServerStore().currentServer

  return (
    <div className="flex flex-col gap-2">
      <p className="text-xs font-medium text-muted-foreground">
        {" "}
        Currently Selected Server
      </p>
      <div className="flex items-center gap-2">
        <img
          className="h-8 w-8 rounded-full bg-accent"
          src="https://cdn.discordapp.com/icons/1187144343400751234/a_0be46d9072c97a20702a6e567ecc1717.gif?size=1024"
        />
        <p className="text-sm font-medium text-muted-foreground">
          {" " + server?.name}
        </p>
      </div>
    </div>
  )
}
