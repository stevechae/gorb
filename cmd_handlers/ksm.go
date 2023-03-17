package cmd_handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// should be a formattable string
var ksmImageSource = os.Getenv("KSM_IMG_SOURCE")
var ksmImageCount = os.Getenv("KSM_IMG_COUNT")

func GetKsmImageUrl() string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	ksmImageCountNum, err := strconv.Atoi(ksmImageCount)
	if err != nil {
		log.Fatalf("Invalid KSM_IMG_COUNT: %s", err)
	}
	n := r.Intn(ksmImageCountNum + 1)

	return fmt.Sprintf(ksmImageSource, n)
}

func KsmHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Received interaction: %s\n", i.ApplicationCommandData().Name)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: GetKsmImageUrl(),
		},
	})
	if err != nil {
		log.Fatalf("Error responding to interaction: %s", err)
	}
}
