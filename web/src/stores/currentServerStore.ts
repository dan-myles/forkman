import { create } from "zustand"
import { Server } from "@/schemas/server"

// TODO:
// Move away from an embedded object and into a server "switcher"
// you can get this info from the URL: /api/v1/user/servers
const D2D: Server = {
  id: "1187144343400751234",
  name: "Devil2Devil",
  icon: "a_0be46d9072c97a20702a6e567ecc1717",
  owner: false,
  permissions: "2251799813685247",
  features: [
    "ANIMATED_ICON",
    "MEMBER_PROFILES",
    "AUTO_MODERATION",
    "GUILD_WEB_PAGE_VANITY_URL",
    "SEVEN_DAY_THREAD_ARCHIVE",
    "CHANNEL_ICON_EMOJIS_GENERATED",
    "LINKED_TO_HUB",
    "THREE_DAY_THREAD_ARCHIVE",
    "DISCOVERABLE",
    "ANIMATED_BANNER",
    "GUILD_ONBOARDING_EVER_ENABLED",
    "PREVIEW_ENABLED",
    "VANITY_URL",
    "ACTIVITY_FEED_DISABLED_BY_USER",
    "BANNER",
    "HAS_DIRECTORY_ENTRY",
    "MEMBER_VERIFICATION_GATE_ENABLED",
    "ENABLED_DISCOVERABLE_BEFORE",
    "GUILD_ONBOARDING",
    "NEWS",
    "ROLE_ICONS",
    "INVITE_SPLASH",
    "PRIVATE_THREADS",
    "COMMUNITY",
    "GUILD_ONBOARDING_HAS_PROMPTS",
    "SOUNDBOARD",
  ],
  approximate_member_count: 0,
  approximate_presence_count: 0,
}

interface CurrentServerStore {
  currentServer: Server | null
}

interface CurrentServerActions {
  setCurrentServer: (currentServer: Server) => void
  clearCurrentServer: () => void
}

export const useCurrentServerStore = create<
  CurrentServerStore & CurrentServerActions
>()((set) => ({
  currentServer: D2D,
  setCurrentServer: (currentServer) => set({ currentServer }),
  clearCurrentServer: () => set({ currentServer: null }),
}))
