package stock

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/piquette/finance-go/equity"
)

var settingStocks = false

func SetStock(b *gotgbot.Bot, ctx *ext.Context) error {
	data, err := ioutil.ReadFile("stock/stocks.txt")
	if err != nil {
		panic(err)
	}
	settingStocks = true
	query := "Please reply with stocks, space-separated, e.g.: \n" + string(data)
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	cb.Message.EditText(b, query, nil)
	return nil
}

func Deduplicate(strSlice []string) []string {
    keys := make(map[string]bool)
    outputList := []string{}
    for _, item := range strSlice {
        if _, value := keys[item]; !value {
            keys[item] = true
            outputList = append(outputList, item)
        }
    }
    return outputList
}

func GenerateStockMessage() string {
	data, err := ioutil.ReadFile("stock/stocks.txt")
	if err != nil {
		panic(err)
	}
	iter := equity.List(Deduplicate(strings.Fields(string(data))))
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
	ret := t + msg.String() + "\n"
	LogStockMessage(ret)
	return ret
}

func LogStockMessage(msg string) {
	const MAX_SIZE = 10000
	f, err := os.OpenFile("stock/stocks.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()
	if int(fi.Size()) > MAX_SIZE {
		if err := os.Truncate("stock/stocks.log", 0); err != nil {
			log.Printf("Failed to truncate: %v", err)
		}
	}
	_, err = f.WriteString(msg)
	if err != nil {
		panic(err)
	}
	f.Close()
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
	if !settingStocks {
		if _, err := b.SendMessage(ctx.EffectiveChat.Id,
			"You're not setting stocks!",
			&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
			fmt.Println("failed: " + err.Error())
		}
		return nil
	}
	text := []byte(strings.ToUpper(ctx.EffectiveMessage.Text))
	if err := ioutil.WriteFile("stock/stocks.txt", text, 0644); err != nil {
		panic(err)
	}
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		"Done.",
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	settingStocks = false
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
