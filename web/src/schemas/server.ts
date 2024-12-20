import { z } from "zod"

export type Server = z.infer<typeof serverSchema>
export type ServerList = z.infer<typeof serverListSchema>

export const serverSchema =
    z.object({
    id: z.string(),
    name: z.string(),
    icon: z.string(),
    owner: z.boolean(),
    permissions: z.string(),
    features: z.array(z.string()),
    approximate_member_count: z.number(),
    approximate_presence_count: z.number()
  })
export const serverListSchema = z.array(serverSchema)
