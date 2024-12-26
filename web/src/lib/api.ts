import { serverListSchema } from "@/schemas/server"
import { statusSchema } from "@/schemas/status"
import { uptimeSchema } from "@/schemas/uptime"
import { useUserStore } from "@/stores/userStore"

const D2D_SERVER_ID = "1187144343400751234"

export const api = {
  isAuth,
  isAdmin,
  getUptime,
  getServers,
  getStatus,
  restartServer,
  logout,
  disableVerification,
  enableVerification,
  getVerificationStatus,
  sendVerificationPanel,
  disableQna,
  enableQna,
  getQnaStatus,
  disableModeration,
  enableModeration,
  getModerationStatus,
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

async function isAdmin() {
  const res = await fetch("/api/v1/user/servers")
  if (res.status !== 200) {
    return false
  }

  const json = await res.json()
  const servers = serverListSchema.parse(json)

  for (const server of servers) {
    if (server.id === D2D_SERVER_ID) {
      return true
    }
  }

  return false
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
    await fetch("/api/v1/sentinel/restart", {
      method: "POST",
      credentials: "include",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function logout() {
  await fetch("/auth/discord/logout", {
    method: "GET",
    credentials: "include",
    cache: "no-cache",
  })
}

async function disableVerification() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/verification/disable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function enableVerification() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/verification/enable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function getVerificationStatus() {
  try {
    const res = await fetch(
      `/api/v1/${D2D_SERVER_ID}/module/verification/status`,
      {
        method: "GET",
        credentials: "include",
        cache: "no-cache",
      }
    )
    const json = await res.json()
    const status = statusSchema.parse(json)
    return status.status
  } catch (e) {
    console.error(e)
    return false
  }
}

async function sendVerificationPanel(channelId: string) {
  try {
    await fetch(
      `/api/v1/${D2D_SERVER_ID}/module/verification/panel/send/${channelId}`,
      {
        method: "POST",
        credentials: "include",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
        },
      }
    )
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function disableQna() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/qna/disable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function enableQna() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/qna/enable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function getQnaStatus() {
  try {
    const res = await fetch(`/api/v1/${D2D_SERVER_ID}/module/qna/status`, {
      method: "GET",
      credentials: "include",
      cache: "no-cache",
    })
    const json = await res.json()
    const status = statusSchema.parse(json)
    return status.status
  } catch (e) {
    console.error(e)
    return false
  }
}

async function enableModeration() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/moderation/enable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}
async function disableModeration() {
  try {
    await fetch(`/api/v1/${D2D_SERVER_ID}/module/moderation/disable`, {
      method: "POST",
      credentials: "include",
      cache: "no-cache",
    })
    return true
  } catch (e) {
    console.error(e)
    return false
  }
}

async function getModerationStatus() {
  try {
    const res = await fetch(
      `/api/v1/${D2D_SERVER_ID}/module/moderation/status`,
      {
        method: "GET",
        credentials: "include",
        cache: "no-cache",
      }
    )
    const json = await res.json()
    const status = statusSchema.parse(json)
    return status.status
  } catch (e) {
    console.error(e)
    return false
  }
}
