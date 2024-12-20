import { BadgeCheck, Brain, SearchCheck } from "lucide-react"
import { useLocation } from "@tanstack/react-router"

export function NavMain() {
  return (
    <ul className="grid gap-0.5">
      <li>
        <Overview />
      </li>
      <li>
        <Verification />
      </li>
      <li>
        <QNA />
      </li>
    </ul>
  )
}

function Overview() {
  const location = useLocation().pathname

  return (
    <div className="relative flex items-center">
      <a
        className={
          `flex h-8 min-w-8 flex-1 cursor-pointer items-center gap-2
          overflow-hidden rounded-md px-1.5 text-sm font-medium outline-none
          ring-ring transition-all hover:bg-accent hover:text-accent-foreground
          focus-visible:ring-2` +
          (location === "/dashboard" ? " bg-accent" : "")
        }
      >
        <div className="flex flex-1 overflow-hidden">
          <SearchCheck className="mr-1 mt-[2px] h-4 w-4 shrink-0 self-center" />
          <div className="line-clamp-1 pr-6">Overview</div>
        </div>
      </a>
    </div>
  )
}

function Verification() {
  const location = useLocation().pathname
  return (
    <div className="relative flex items-center">
      <a
        className={
          `flex h-8 min-w-8 flex-1 cursor-pointer items-center gap-2
          overflow-hidden rounded-md px-1.5 text-sm font-medium outline-none
          ring-ring transition-all hover:bg-accent hover:text-accent-foreground
          focus-visible:ring-2` +
          (location === "/dashboard/verification" ? " bg-accent" : "")
        }
      >
        <div className="flex flex-1 overflow-hidden">
          <BadgeCheck className="mr-1 mt-[2px] h-4 w-4 shrink-0 self-center" />
          <div className="line-clamp-1 pr-6">Verification</div>
        </div>
      </a>
    </div>
  )
}

function QNA() {
  const location = useLocation().pathname
  return (
    <div className="relative flex items-center">
      <a
        className={
          `flex h-8 min-w-8 flex-1 cursor-pointer items-center gap-2
          overflow-hidden rounded-md px-1.5 text-sm font-medium outline-none
          ring-ring transition-all hover:bg-accent hover:text-accent-foreground
          focus-visible:ring-2` +
          (location === "/dashboard/qna" ? " bg-accent" : "")
        }
      >
        <div className="flex flex-1 overflow-hidden">
          <Brain className="mr-1 mt-[2px] h-4 w-4 shrink-0 self-center" />
          <div className="line-clamp-1 pr-6">Q&A</div>
        </div>
      </a>
    </div>
  )
}
