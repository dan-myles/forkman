
export function NavTeamSwitcher() {
  const guildAVI =
    "https://media.discordapp.net/attachments/1251224958496145448/1295522921762131979/Devil2Devil_Bot_PFP_Forkman.png?ex=6765f83e&is=6764a6be&hm=7a415eb86dfd4434ed70ece01a6dd75377f114757764117233ea285b5fff8c2c&=&format=webp&quality=lossless&width=1112&height=1112"

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
          <img src={guildAVI} className="h-10 w-10" />
        </div>
        <div className="line-clamp-1 flex-1 pr-2 text-xl font-medium">
          Forkman
        </div>
      </div>
    </div>
  )
}
