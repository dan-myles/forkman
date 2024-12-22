import { serverListSchema } from "@/schemas/server"
import { uptimeSchema } from "@/schemas/uptime"
import { useUserStore } from "@/stores/userStore"

export const api = {
  isAuth,
  getUptime,
  getServers,
  getStatus,
  restartServer,
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

async function getUptime() {
  const res = await fetch("/uptime", {
    method: "GET",
    credentials: "include",
  })
  const json = await res.json()
  const uptime = uptimeSchema.parse(json)
  return uptime
}

async function getServers() {
  const res = await fetch("/api/v1/user/servers", {
    method: "GET",
    credentials: "include",
  })
  const json = await res.json()
  const servers = serverListSchema.parse(json)
  return servers
}

async function getStatus() {
  try {
    await fetch("/health", {
      method: "GET",
      credentials: "include",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function restartServer() {
  try {
    await fetch("/api/v1/user/restart", {
      method: "POST",
      credentials: "include",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}
