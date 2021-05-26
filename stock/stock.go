package stock

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/piquette/finance-go/equity"
)

var stocks = []string{"AAPL", "GOOG", "MSFT"}

func SetStock(b *gotgbot.Bot, ctx *ext.Context) error {
	const query = "Please reply with stocks, space-separated, e.g.: \nAAPL GOOG MSFT"
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		query,
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}

func GenerateStockMessage(stock []string) string {
	iter := equity.List(stocks)
	var msg strings.Builder
	var t string
	for iter.Next() {
		q := iter.Equity()
		t = fmt.Sprint(time.Unix(int64(q.RegularMarketTime), 0)) + "\n"
		msg.WriteString(MakeLine(q.Symbol, fmt.Sprint(q.RegularMarketPrice), MakeStockChanges(q.RegularMarketChangePercent)))
		msg.WriteString("\n")
	}
	if iter.Err() != nil {
		panic(iter.Err())
	}
	return t + msg.String()
}

func GetStock(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := GenerateStockMessage(stocks)
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		msg,
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}

func ReceiveStock(b *gotgbot.Bot, ctx *ext.Context) error {
	text := strings.Fields(strings.ToUpper(ctx.EffectiveMessage.Text))
	stocks = text
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		"Done.",
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	return nil
}

func MakeStockChanges(price float64) string {
	var msg strings.Builder
	if price > 0 {
		msg.WriteString("+")
	}
	msg.WriteString(fmt.Sprintf("%.2f%%", price))
	if price > 0 {
		msg.WriteString("📈")
	} else if price < 0 {
		msg.WriteString("📉")
	}
	return msg.String()
}

func MakeLine(items ...string) string {
	var msg strings.Builder
	for _, item := range items {
		msg.WriteString(item)
		msg.WriteString("   ")
	}
	return msg.String()
}
