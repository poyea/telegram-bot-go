package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var letters = []rune("_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func gen_random_string(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetRandom(b *gotgbot.Bot, ctx *ext.Context) error {
	rand.Seed(time.Now().UTC().UnixNano())
	msg := "0"
	for msg[0] >= '0' && msg[0] <= '9' {
		msg = gen_random_string(32)
	}
	if _, err := b.SendMessage(ctx.EffectiveChat.Id,
		msg,
		&gotgbot.SendMessageOpts{ParseMode: "html"}); err != nil {
		fmt.Println("failed: " + err.Error())
	}
	cb := ctx.Update.CallbackQuery
	cb.Answer(b, nil)
	return nil
}
