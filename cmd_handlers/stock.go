package cmd_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type stock struct {
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	Price               string `json:"price"`
	Change              string `json:"change"`
	ChangePct           string `json:"changePct"`
	Market              string `json:"market"`
	FiveDayPerf         string `json:"fivedayPerf"`
	OneMonthPerf        string `json:"onemonthPerf"`
	ThreeMonthPerf      string `json:"threemonthPerf"`
	OneYearPerf         string `json:"oneyearPerf"`
	AfterHoursPrice     string `json:"afterHoursPrice"`
	AfterHoursChange    string `json:"afterHoursChange"`
	AfterHoursChangePct string `json:"afterHoursChangePct"`
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

func renderStockInfo(stock stock) string {
	if stock.ChangePct[0] != '-' {
		stock.ChangePct = "+" + stock.ChangePct
	}

	formatMd := "_%s_ (%s) %s **$%s** (**%s**)"

	if stock.AfterHoursPrice != "" {
		formatMd += " After Hours: **$%s** (**%s**) "
		return fmt.Sprintf(formatMd, stock.Name, stock.Symbol, stock.Market, stock.Price, stock.ChangePct, stock.AfterHoursPrice, stock.AfterHoursChangePct)
	}

	return fmt.Sprintf(formatMd, stock.Name, stock.Symbol, stock.Market, stock.Price, stock.ChangePct)
}

func renderPerformance(stock stock) string {
	formatMd := "_%s_ Performance 5D: **%s** 1M: **%s** 3M: **%s** 1Y: **%s**"
	return fmt.Sprintf(formatMd, stock.Symbol, stock.FiveDayPerf, stock.OneMonthPerf, stock.ThreeMonthPerf, stock.OneYearPerf)
}

func StockHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options

	symbol, market := GetSymbolAndMarket(strings.ToLower(options[0].StringValue()))

	log.Printf("Received Interaction - Symbol: %s, Market: %s\n", symbol, market)

	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", stockApiUrl, symbol, market))
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %s\n", err)
			return
		}
	}(resp.Body)

	stock := stock{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stock)
	if err != nil {
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: renderStockInfo(stock),
		},
	})

	if err != nil {
		log.Printf("Error responding: %s\n", err)
		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Error responding: %s", err),
		})
		if err != nil {
			return
		}
		return
	}

	if stock.FiveDayPerf != "" {
		// ETF quotes may not have performance data
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: renderPerformance(stock),
		})

		if err != nil {
			return
		}
	}

	if stock.ChangePct[0] == '-' {
		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf(stonkDownGifUrl),
		})
		if err != nil {
			return
		}
	} else {
		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf(stonkUpGifUrl),
		})
		if err != nil {
			return
		}
	}
}
