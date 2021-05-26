package weather

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Content struct {
	Title       string `xml:"channel>title"`
	Subtitle    string `xml:"channel>item>title"`
	Description string `xml:"channel>item>description"`
}

func GetWeather(b *gotgbot.Bot, ctx *ext.Context) error {
	content := GetFromRSS()
	response := Transform(content)
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		response,
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}

func Transform(report Content) string {
	Ret := fmt.Sprintf("<b>%s</b>\n%s\n!!%s", report.Title, report.Subtitle, report.Description)
	brRegex, _ := regexp.Compile(`(<br/?>+\s+|<br/>)`)
	spRegex, _ := regexp.Compile(`!![\r\n\s]+`)
	match := brRegex.ReplaceAllString(Ret, "\n")
	match = spRegex.ReplaceAllString(match, "\n")
	return match
}

func GetFromRSS() Content {
	resp, err := http.Get("https://rss.weather.gov.hk/rss/LocalWeatherForecast_uc.xml")
	if err != nil {
		fmt.Println("failed" + err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed" + err.Error())
	}
	var report Content
	if err := xml.Unmarshal([]byte(body), &report); err != nil {
		fmt.Println(err)
	}
	return report
}
