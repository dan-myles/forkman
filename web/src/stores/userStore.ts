import { create } from "zustand"
import { User } from "@/schemas/user"

interface UserStore {
  user: User | null
}

interface UserActions {
  setUser: (user: User) => void
  clearUser: () => void
}

export const useUserStore = create<UserStore & UserActions>()((set) => ({
  user: null,
  setUser: (user: User) => set({ user }),
  clearUser: () => set({ user: null }),
}))
