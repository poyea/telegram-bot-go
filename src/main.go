package main

import (
	"fmt"
	"time"
	"net/http"
	"os"
    "encoding/xml"
	"io"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

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
	logger.Sugar().Infof("Server started")

	dispatcher.AddHandler(handlers.NewCommand("start", Start))
	dispatcher.AddHandler(handlers.NewCommand("keyboard", Start))
	dispatcher.AddHandler(handlers.NewCommand("else", other))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("close_callback"), CloseKeyboard))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("time_callback"), SendTime))
	dispatcher.AddHandler(handlers.NewCallback(filters.Equal("weather_callback"), GetWeather))
	dispatcher.AddHandler(handlers.NewMessage(filters.All, echo))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	logger.Sugar().Infof("%s has been started...\n", b.User.Username)

	updater.Idle()
}

func GenKeyboardLayout() [][]gotgbot.InlineKeyboardButton {
	return [][]gotgbot.InlineKeyboardButton{{
		{Text: "Time", CallbackData: "time_callback"},
		{Text: "Weather", CallbackData: "weather_callback"},
		{Text: "Close", CallbackData: "close_callback"},
	},{
		{Text: "Weather", CallbackData: "other_callback"},
		{Text: "Weather", CallbackData: "other_callback"},
		{Text: "Weather", CallbackData: "other_callback"},
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
	ctx.EffectiveMessage.Reply(b, time.Now().In(loc).Format(time.RFC850), nil)
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}


func GetWeather(b *gotgbot.Bot, ctx *ext.Context) error {
	fmt.Println("In Weather")
	resp, err := http.Get("https://rss.weather.gov.hk/rss/LocalWeatherForecast_uc.xml")
	if err != nil {
		fmt.Println("failed" + err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed" + err.Error())
	}
	type Things struct {
		String  string    `xml:"rss>channel>item>author"`
		// String  string    `xml:"rss"`
	}
	var m Things
	if err := xml.Unmarshal(body, &m); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", m.String)
	ctx.EffectiveMessage.Reply(b, m.String, nil)
	return nil
}
