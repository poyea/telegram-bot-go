package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stocks = []string{"AAPL", "GOOG", "MSFT"}

func main() {
	b, err := gotgbot.NewBot(os.Getenv("TOKEN"), &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		panic("Failed to create new bot: " + err.Error())
	}
	log := zap.NewProductionEncoderConfig()
	log.EncodeLevel = zapcore.CapitalLevelEncoder
	log.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(log), os.Stdout, zap.InfoLevel))

	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("start", Start))
	dispatcher.AddHandler(handlers.NewCommand("keyboard", Start))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("close_callback"), CloseKeyboard))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("time_callback"), SendTime))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("weather_callback"), GetWeather))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("set_stock_callback"), SetStock))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("get_stock_callback"), GetStock))
	dispatcher.AddHandler(handlers.NewMessage(filters.Reply, ReceiveStock))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	logger.Sugar().Infof("@%s has been started...\n", b.User.Username)

	updater.Idle()
}

func GenKeyboardLayout() [][]gotgbot.InlineKeyboardButton {
	return [][]gotgbot.InlineKeyboardButton{{
		{Text: "Time", CallbackData: "time_callback"},
		{Text: "Weather", CallbackData: "weather_callback"},
	}, {
		{Text: "Set Stock", CallbackData: "set_stock_callback"},
		{Text: "Get Stock", CallbackData: "get_stock_callback"},
		{Text: "Close", CallbackData: "close_callback"},
	}}
}

func Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Hello, please select a command you want.", &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: GenKeyboardLayout(),
		},
	})
	if err != nil {
		fmt.Println("failed to send: " + err.Error())
	}
	return nil
}

func CloseKeyboard(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	cb.Message.EditText(b, "Use /keyboard anytime to re-open the keyboard.", nil)
	return nil
}

func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	return nil
}

func SendTime(b *gotgbot.Bot, ctx *ext.Context) error {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		fmt.Sprintf("The time now is:\n<b>%s</b>", time.Now().In(loc).Format(time.RFC850)),
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}

func MakeLine(items ...string) string {
	var msg strings.Builder
	for _, item := range items {
		msg.WriteString(item)
		msg.WriteString("   ")
	}
	return msg.String()
}
