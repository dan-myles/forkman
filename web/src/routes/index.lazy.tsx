import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/")({
  component: Page,
})

export default function Page() {
  return (
    <div className="">
      <div>test</div>
    </div>
  )
}
