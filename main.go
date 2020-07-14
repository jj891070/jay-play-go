package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
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

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				log.Printf("message ---> %+v \n", message)
				log.Printf("quota ---> %+v \n", quota)
				var res string
				switch {
				case message.Text == "":
					res = `
						😉 您好，請問您需要什麼服務呢？
						1. 請輸入hello
						2. 請輸入寶哥好
					`
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res)).Do(); err != nil {
						log.Print(err)
					}
				case strings.Contains(message.Text, "hello"):
					res = "nice to meet you！😌 "
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res)).Do(); err != nil {
						log.Print(err)
					}
				case strings.Contains(message.Text, "寶哥"):
					res = "老大好！🙋"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res)).Do(); err != nil {
						log.Print(err)
					}
				default:
					originalContentURL := "https://images.app.goo.gl/M55DNLhKRPGwN2ss7"
					previewImageURL := "https://images.app.goo.gl/M55DNLhKRPGwN2ss7"

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(originalContentURL, previewImageURL)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
