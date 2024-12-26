import { z } from "zod"

export type Status = z.infer<typeof statusSchema>

export const statusSchema = z.object({
  message: z.string(),
  status: z.boolean(),
})
