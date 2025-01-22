package qna

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/bwmarrin/discordgo"
)

var (
	CIDAdditionalAssistanceBtn = "additional_assistance_button"
	HelperRoleID               = "1213199035968528434"
)

func (m *QNA) handleQNARequest(s *discordgo.Session, msg *discordgo.MessageCreate) {
	channel, err := m.session.Channel(msg.ChannelID)
	if err != nil {
		m.log.Error().Err(err).Msg("critical error getting channel")
		return
	}

	if channel.Type != discordgo.ChannelTypeGuildPublicThread {
		return
	}

	if channel.MessageCount != 0 {
		return
	}

	if channel.ParentID != m.forumChannelId {
		return
	}

	userId := msg.Author.ID
	channelID := msg.ChannelID
	content := "Hi <@" + userId + ">, I'm Forkman, you're friendly support bot. I'm looking through our knowledge base to see if I can answer your question. :wave:"

	message, _ := s.ChannelMessageSend(msg.ChannelID, content)
	query := channel.Name + " " + msg.Content

	input := &bedrockagentruntime.RetrieveAndGenerateInput{
		Input: &types.RetrieveAndGenerateInput{
			Text: aws.String(query),
		},
		RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
			Type: types.RetrieveAndGenerateTypeKnowledgeBase,
			KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
				ModelArn:        aws.String("us.anthropic.claude-3-5-sonnet-20241022-v2:0"),
				KnowledgeBaseId: aws.String(m.knowledgeBaseId),
			},
		},
	}

	response, err := m.bedrock.RetrieveAndGenerate(context.Background(), input)
	if err != nil || response.Output == nil {
		m.log.Error().Err(err).Msg("failed to retrieve and generate")
		s.ChannelMessageEdit(channelID, message.ID, "Uh oh, I couldn't find an answer to your question. Please try again later.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Description: "<@" + userId + ">, we're still improving our answers! Let us know if you still need assistance.",
		Color:       0x00FF00, // green color
	}

	buttonRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "I still need help!",
				Style: discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: "👆",
				},
				CustomID: CIDAdditionalAssistanceBtn, // Custom ID for button to trigger modal
			},
		},
	}

	content = content + "\n----------------------\n" + *response.Output.Text

	_, err = s.ChannelMessageEditComplex(
		&discordgo.MessageEdit{
			Content:    &content,
			Channel:    channelID,
			ID:         message.ID,
			Embed:      embed,
			Components: &[]discordgo.MessageComponent{buttonRow},
		},
	)
	if err != nil {
		m.log.Error().Err(err).Msg("error editing message")
		return
	}
}

func (m *QNA) handleCIDAdditionalAssistanceBtn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ping := fmt.Sprintf("<@%s> Assistance requested.", HelperRoleID)
	content := i.Message.Content

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ping,
		},
	})

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &content,
		Channel:    i.ChannelID,
		ID:         i.Message.ID,
		Embeds:     &[]*discordgo.MessageEmbed{},
		Components: &[]discordgo.MessageComponent{},
	})
	if err != nil {
		m.log.Error().Err(err).Msg("error editing message")
		return
	}
}
