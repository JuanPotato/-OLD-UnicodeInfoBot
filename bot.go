package main

import (
	"fmt"
	"log"
	"regexp"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	BotToken = ""

	aboutText = `Hi, I'm @UnicodeInfoBot. I provide you with unicode information about the contents of your message.

<a href="https://github.com/JuanPotato/UnicodeInfoBot">UnicodeInfoBot v1.0</a> by @JuanPotato`
)

var bot *tgbotapi.BotAPI

func main() {
	b, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot = b

	username := bot.Self.UserName
	log.Printf("Authorized on account %s", username)

	upConfig := tgbotapi.NewUpdate(0)
	upConfig.Timeout = 60
	startRegex, _ := regexp.Compile(fmt.Sprintf("^\\/(start|about)(?:@%s)?", username))
	updates, err := bot.GetUpdatesChan(upConfig)

	for update := range updates {
		if update.Message != nil {
			switch true {
			case startRegex.MatchString(update.Message.Text):
				go about(update)
			default:
				go unicodeInfo(update)
			}
		}
	}
}

func unicodeInfo(update tgbotapi.Update) {
	text := ""
	for _, e := range update.Message.Text {
		t := fmt.Sprintf("<a href=\"http://www.fileformat.info/info/unicode/char/%X/index.htm\">%v</a>\n",
			e, CodePoints[e])
		if len(text)+len(t) > 3991 { // This will provide just enough space to add the Notice
			text = text + "\nYour message was truncated, please send large texts in pieces if you would like to see each character."
			break
		}
		text = text + t
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
}

func about(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, aboutText)
	msg.ParseMode = "HTML"

	bot.Send(msg)
}
