package cmd_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"os"
	"strings"
)

type stock struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Change    string `json:"change"`
	ChangePct string `json:"changePct"`
	Market    string `json:"market"`
}

var stockApiUrl = os.Getenv("STOCK_API_URL")
var stonkUpGifUrl = os.Getenv("STONK_UP_GIF_URL")
var stonkDownGifUrl = os.Getenv("STONK_DOWN_GIF_URL")

func GetSymbolAndMarket(input string) (string, string) {
	var market string
	var symbol string
	symbolMarketSplitRaw := strings.Split(input, ".")

	symbol = symbolMarketSplitRaw[0]

	if len(symbolMarketSplitRaw) == 1 {
		market = "us"
	} else {
		// only supports US and CAN
		// TODO: make it more robust and support more markets
		switch symbolMarketSplitRaw[1] {
		case "to":
			market = "ca"
		default:
			log.Printf("Invalid market: %s", symbolMarketSplitRaw[1])
		}

	}

	return symbol, market
}

func StockHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options

	symbol, market := GetSymbolAndMarket(strings.ToLower(options[0].StringValue()))

	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", stockApiUrl, symbol, market))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	stock := stock{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&stock)

	if stock.ChangePct[0] != '-' {
		stock.ChangePct = "+" + stock.ChangePct
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("_%s_ (%s) %s **$%s** (**%s**)", stock.Name, stock.Symbol, stock.Market, stock.Price, stock.ChangePct),
		},
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Error responding: %s", err),
		})
	} else {
		if stock.ChangePct[0] == '+' {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf(stonkUpGifUrl),
			})
		} else {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf(stonkDownGifUrl),
			})
		}
	}
}
