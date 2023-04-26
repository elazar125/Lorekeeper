package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dg, err := discordgo.New("Bot " + os.Getenv("LOREKEEPER_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	re := regexp.MustCompile(`\[\[\!?([^|\]]+)\|?([^\]]*)\]\]`)
	cardQuerys := re.FindAllStringSubmatch(m.Content, -1)

	for _, cardQuery := range cardQuerys {
		cardName := strings.Replace(strings.TrimSpace(cardQuery[1]), "’", "'", -1)
		imageFormat := strings.TrimSpace(cardQuery[2])

		if strings.ToLower(cardName) == "help" {
			s.ChannelMessageSend(m.ChannelID, getHelpText())
		} else {
			parsedImageFormat, errMessage := validateImageFormat(imageFormat)
			if errMessage != "" {
				s.ChannelMessageSend(m.ChannelID, errMessage)
			}
			cardUrl, err := getCardUrl(cardName, parsedImageFormat)
			if err != nil {
				log.Println(err)
			}
			s.ChannelMessageSend(m.ChannelID, cardUrl)
		}
	}
}

func validateImageFormat(imageFormat string) (string, string) {
	switch imageFormat {
	case "foil_full_art", "ffa":
		return "foil_full_art", ""
	case "foil_text", "ft":
		return "foil_text", ""
	case "full_art", "fa":
		return "full_art", ""
	case "text", "t", "":
		return "text", ""
	default:
		return "text", "Unknown image format, using 'text' instead"
	}
}

func getCardUrl(cardName, imageFormat string) (string, error) {
	result, err := searchSpellslingerer(fmt.Sprintf(`(name~"%s")`, cardName))
	if err != nil {
		return "", err
	}

	if result.TotalItems < 1 {
		return fmt.Sprintf(`No cards found for "%s", please check your spelling`, cardName), nil
	}

	if result.TotalItems == 1 {
		return fmt.Sprintf("https://spellslingerer.com/cards/%s?image=%s", result.Items[0]["id"], imageFormat), nil
	}

	for _, item := range result.Items {
		if strings.ToLower(item["name"].(string)) == strings.ToLower(cardName) {
			return fmt.Sprintf("https://spellslingerer.com/cards/%s?image=%s", item["id"], imageFormat), nil
		}
	}

	return fmt.Sprintf(`More than one card found for "%s", please be more specific`, cardName), nil
}

// Copied from https://github.com/pocketbase/pocketbase/blob/v0.13.2/tools/search/provider.go#L27-L33
type Result struct {
	Page       int                      `json:"page"`
	PerPage    int                      `json:"perPage"`
	TotalItems int                      `json:"totalItems"`
	TotalPages int                      `json:"totalPages"`
	Items      []map[string]interface{} `json:"items"`
}

func searchSpellslingerer(filter string) (Result, error) {
	url := fmt.Sprintf(`https://spellslingerer.com/api/collections/cards/records?filter=%s`, filter)
	res, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()

	var result Result
	err = json.NewDecoder(res.Body).Decode(&result)
	return result, err
}

func getHelpText() string {
	return `Hey spellslinger! Here's how to ask me for a card:
- [[card name]] or [[!card name]] will work
- Only part of the name will work, as long as it uniquely identifies the card
- If you want fancy art, specify it like [[card name|format]]
- Format can be one of "text", "full_art", "foil_text", "foil_full_art"
- You can also use the abbreviations "t", "fa", "ft", "ffa" instead
- If you specify anything else I'll use "text"
- You can add spaces around the name and format if you like ([[ card name | format ]]), but don't put extra spaces in the middle of the name
`
}
