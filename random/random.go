package random

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func GetRandom(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := "random_string"
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		msg,
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}
