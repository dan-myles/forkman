import { create } from "zustand"
import { Server } from "@/schemas/server"

interface ServerStore {
  servers: Server[] | null
}

interface ServerActions {
  setServers: (servers: Server[]) => void
  clearServers: () => void
}

export const useServerStore = create<ServerStore & ServerActions>()((set) => ({
  servers: null,
  setServers: (servers) => set({ servers }),
  clearServers: () => set({ servers: [] }),
}))
