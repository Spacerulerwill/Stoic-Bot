package bot

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	quoteCmdChoices []*discordgo.ApplicationCommandOptionChoice
	commands        []*discordgo.ApplicationCommand
	commandHandlers = make(map[string]func(bot *discordgo.Session, interact *discordgo.InteractionCreate))
)

func init() {
	// Add quote command option choices using yaml data
	for stoicName := range stoicData {
		var slashCommandOptionChoice = new(discordgo.ApplicationCommandOptionChoice)
		slashCommandOptionChoice.Name = cases.Title(language.English, cases.NoLower).String(stoicName)
		slashCommandOptionChoice.Value = stoicName
		quoteCmdChoices = append(quoteCmdChoices, slashCommandOptionChoice)
	}

	// Add commands
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "quote",
			Description: "Recieve wisdom from a stoic philosopher",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "stoic",
					Description: "Stoic philosopher to recieve wisdom from",
					Required:    false,
					Choices:     quoteCmdChoices,
				},
			},
		},
	}

	commandHandlers = map[string]func(bot *discordgo.Session, interact *discordgo.InteractionCreate){
		"quote": quoteCommand,
	}
}

func quoteCommand(bot *discordgo.Session, interact *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := interact.ApplicationCommandData().Options
	var chosenStoicData stoic
	var chosenStoicName string

	if len(options) > 0 {
		// type assertion for string
		chosenStoicName = options[0].Value.(string)
		chosenStoicData = stoicData[chosenStoicName]

	} else {
		// other wise show a random stoic
		k := rand.Intn(len(stoicData))
		i := 0

		for name, iterStoic := range stoicData {
			if i == k {
				chosenStoicData = iterStoic
				chosenStoicName = name
			}
			i++
		}
	}

	// show quote in embed
	randomQuote := chosenStoicData.Quotes[rand.Intn(len(chosenStoicData.Quotes))]
	embed := new(discordgo.MessageEmbed)
	embed.Title = fmt.Sprintf("%s - %s", cases.Title(language.English, cases.NoLower).String(chosenStoicName), randomQuote.Source)
	embed.Description = randomQuote.Quote
	embed.Color = bot.State.UserColor(interact.Member.User.ID, interact.ChannelID)
	thumbnail := new(discordgo.MessageEmbedThumbnail)
	thumbnail.URL = chosenStoicData.ImageURL
	embed.Thumbnail = thumbnail
	bot.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
