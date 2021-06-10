package stock

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/piquette/finance-go/equity"
)

func SetStock(b *gotgbot.Bot, ctx *ext.Context) error {
	const query = "Please reply with stocks, space-separated, e.g.: \nAAPL GOOG MSFT"
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	cb.Message.EditText(b, query, nil)
	return nil
}

func GenerateStockMessage() string {
	data, err := ioutil.ReadFile("stock/stocks.txt")
	if err != nil {
		panic(err)
	}
	iter := equity.List(strings.Fields(string(data)))
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
	ret := t + msg.String()
	f, err := os.OpenFile("stock/stocks.log", os.O_APPEND|os.O_WRONLY, 0644)
	_, err = f.WriteString(ret)
	f.Close()
	return ret
}

func GetStock(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := GenerateStockMessage()
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
	text := []byte(strings.ToUpper(ctx.EffectiveMessage.Text))
	if err := ioutil.WriteFile("stock/stocks.txt", text, 0644); err != nil {
		panic(err)
	}
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
	msg.WriteString("| ")
	for _, item := range items {
		msg.WriteString(item)
		msg.WriteString(" | ")
	}
	return msg.String()
}
