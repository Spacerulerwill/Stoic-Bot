package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

type quote struct {
	Quote  string
	Source string
}

type stoic struct {
	ImageURL string
	Quotes   []quote
}

var stoicData = make(map[string]stoic)

func init() {
	// Load stoic data from stoic.yml
	content, error := os.ReadFile("res/stoic.yml")
	if error != nil {
		log.Fatal("Failed to find stoic.yml")
		return
	}

	yaml_err := yaml.Unmarshal(content, &stoicData)
	if yaml_err != nil {
		log.Fatalf("error: %v", yaml_err)
	}
}

func Run(token string, guildID string) {
	// Create discord bot, add handlers and add slash commands
	discord, bot_create_err := discordgo.New("Bot " + token)
	discord.Identify.Intents |= discordgo.IntentsAll

	if bot_create_err != nil {
		log.Fatal(bot_create_err)
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	discord.AddHandler(ready)

	discord_open_err := discord.Open()
	if discord_open_err != nil {
		log.Fatalf("Cannot open the session: %v", discord_open_err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer discord.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

func ready(bot *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", bot.State.User.Username, bot.State.User.Discriminator)
	bot.UpdateGameStatus(0, "Use /quote!")
}
