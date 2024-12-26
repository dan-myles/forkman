
export function NavTeamSwitcher() {
  return (
    <div>
      <div
        className="flex items-center gap-1.5 overflow-hidden px-2 py-1.5
          text-left text-sm transition-all"
      >
        <div
          className="flex items-center justify-center rounded-sm bg-transparent
            text-primary-foreground"
        >
          <img src="/forkman.png" className="h-10 w-10" />
        </div>
        <div className="line-clamp-1 flex-1 pr-2 text-xl font-medium">
          Forkman
        </div>
      </div>
    </div>
  )
}
