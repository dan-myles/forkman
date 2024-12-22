import { z } from "zod"

export type Uptime = z.infer<typeof uptimeSchema>

export const uptimeSchema = z.object({
  message: z.string(),
  uptime: z.number(),
})
