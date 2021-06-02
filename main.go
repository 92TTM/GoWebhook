package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New("4bd0bc91b1cca5fb1c0a621d856b31fa", "M6fBGbOHaSl0MBDC52eLH1kZrkB54aJxd2vQsk89Lpo5YhR+tzM4cOqtYrow6vQFpq1G4/kxA6iv++CehbPchvdLh4k2DPx2Ozmhpl8zi4+RYE8xanKnplRi7js1DrqfiBuyJm3IzznIXIsDbkGwtAdB04t89/1O/w1cDnyilFU=")
	fmt.Println("Hello World")
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := "80"
	addr := fmt.Sprintf(":%s", port)
	fmt.Println(addr)
	http.ListenAndServe(addr, nil)

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	fmt.Println(events)
	quota, err := bot.GetMessageQuota().Do()
	if err != nil {
		log.Println("Quota err:", err)
	}
	fmt.Println(quota)
	for _, event := range events {
		userid := event.Source.UserID
		users, err := bot.GetProfile(userid).Do()
		fmt.Println(users)
		name := users.DisplayName
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:
				replytext := ""
				if message.Text == "123" {
					replytext = "456"
				}
				if message.Text == "恬恬喵" {
					replytext = "好萌！！"
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(name+":"+replytext)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.ImageMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("image")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
