import { serverListSchema } from "@/schemas/server"
import { useUserStore } from "@/stores/userStore"

export const api = {
  isAuth,
  getServers,
}

async function isAuth() {
  const res = await fetch("/api/v1/user/session", {
    method: "GET",
    credentials: "include",
  })

  if (res.status === 401) {
    return false
  }

  const user = await res.json()
  useUserStore.setState({ user })
  return true
}

async function getServers() {
  const res = await fetch("/api/v1/user/servers", {
    method: "GET",
    credentials: "include",
  })
  const json = await res.json()
  console.log(json)
  const servers = serverListSchema.parse(json)
  return servers
}
