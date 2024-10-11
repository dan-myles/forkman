package verification

import (
	"github.com/bwmarrin/discordgo"
)

func (m *Verification) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

}

func (m *Verification) listen(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		return
	}

	go m.verifyEmail(s, i)
}

func (m *Verification) verifyEmail(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
}
