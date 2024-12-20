import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getGuildIcon(guildId: string, iconHash: string) {
  if (iconHash.length > 0 && iconHash.startsWith("a_")) {
    return `https://cdn.discordapp.com/icons/${guildId}/${iconHash}.gif`
  }

  return `https://cdn.discordapp.com/icons/${guildId}/${iconHash}.png`
}
