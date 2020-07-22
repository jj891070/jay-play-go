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
	log.Println(" ================ ")
	log.Printf("event --> %+v \n", events)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {

		memberID := bot.GetGroupMemberIDs(event.Source.GroupID, os.Getenv("ChannelAccessToken"))
		log.Println("groupID --> ", event.Source.GroupID)
		log.Println("roomID --> ", event.Source.RoomID)
		log.Println("memberIDs --> ", memberID.NewScanner().ID())

		// if _, err := bot.Multicast([]string{event.Source.UserID}, linebot.NewTextMessage("hello my jay")).Do(); err != nil {
		// 	log.Println("Multicast Err -> ", err)
		// }

		log.Println("userID --> ", event.Source.UserID)
		var res *linebot.UserProfileResponse
		res, err = bot.GetProfile(event.Source.UserID).Do()
		log.Println("useraName:", res.DisplayName)
		log.Println("language:", res.Language)
		log.Println("status:", res.StatusMessage)
		log.Println("pic:", res.PictureURL)
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
				case message.Text == "e":
					res = `
						😉 您好，請問您需要什麼服務呢？
						1. 請輸入hello
						2. 請輸入寶哥好
					`

					a := linebot.NewFlexMessage("我愛你", &linebot.BubbleContainer{
						Type: linebot.FlexContainerTypeBubble,
						Body: &linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Color: "#00ff00",
									Size:  "lg",
									Type:  linebot.FlexComponentTypeText,
									Text:  "單號：123456789",
								},
								&linebot.TextComponent{
									Color: "#ff0000",
									Align: "center",
									Size:  "xl",
									Type:  linebot.FlexComponentTypeText,
									Text:  "[Succes]",
								},
								&linebot.TextComponent{
									Type:  linebot.FlexComponentTypeText,
									Align: "center",
									Text:  "內容...💮 🌸 🏵",
								},
							},
						},
					})
					// bot.BroadcastMessage(a)
					if _, err = bot.ReplyMessage(event.ReplyToken, a).Do(); err != nil {
						log.Print(err)
					}
				case strings.Contains(message.Text, "hello"):
					res = "nice to meet you！😌~ "
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res)).Do(); err != nil {
						log.Print(err)
					}
					// bot.Narrowcast(linebot.NewTextMessage(res))

				case strings.Contains(message.Text, "寶哥"):
					res = "老大好！🙋"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res)).Do(); err != nil {
						log.Print(err)
					}
				case message.Text == "123": // carousel
					resp := linebot.NewTemplateMessage(
						"this is a carousel template with imageAspectRatio,  imageSize and imageBackgroundColor",
						linebot.NewCarouselTemplate(
							linebot.NewCarouselColumn(
								"https://farm5.staticflickr.com/4849/45718165635_328355a940_m.jpg",
								"this is menu",
								"description",
								linebot.NewPostbackAction("Buy", "action=buy&itemid=111", "", ""),
								linebot.NewPostbackAction("Add to cart", "action=add&itemid=111", "", ""),
								linebot.NewURIAction("View detail", "http://example.com/page/111"),
							).WithImageOptions("#FFFFFF"),
							linebot.NewCarouselColumn(
								"https://farm5.staticflickr.com/4849/45718165635_328355a940_m.jpg",
								"this is menu",
								"description",
								linebot.NewPostbackAction("Buy", "action=buy&itemid=111", "", ""),
								linebot.NewPostbackAction("Add to cart", "action=add&itemid=111", "", ""),
								linebot.NewURIAction("View detail", "http://example.com/page/111"),
							).WithImageOptions("#FFFFFF"),
						).WithImageOptions("rectangle", "cover"),
					)
					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				case message.Text == "789": // quicklyresponse
					resp := linebot.NewTextMessage(
						"Select your favorite food category or send me your location!",
					).WithQuickReplies(
						linebot.NewQuickReplyItems(
							linebot.NewQuickReplyButton("https://example.com/sushi.png", linebot.NewMessageAction("Sushi", "Sushi")),
							linebot.NewQuickReplyButton("https://example.com/tempura.png", linebot.NewMessageAction("Tempura", "Tempura")),
							linebot.NewQuickReplyButton("", linebot.NewLocationAction("Send location")),
						),
					)

					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				case message.Text == "456": //confirm
					resp := linebot.NewTemplateMessage(
						"this is a confirm template",
						linebot.NewConfirmTemplate(
							"Are you sure?",
							linebot.NewMessageAction("Yes", "yes"),
							linebot.NewMessageAction("No", "no"),
						),
					)

					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				case message.Text == "qa":
					resp := linebot.NewTemplateMessage(
						"this is a buttons template",
						linebot.NewButtonsTemplate(
							"https://farm5.staticflickr.com/4849/45718165635_328355a940_m.jpg",
							"Menu",
							"Please select",
							linebot.NewPostbackAction("Buy", "action=buy&itemid=123", "", "displayText"),
							linebot.NewPostbackAction("Buy", "action=buy&itemid=123", "text", ""),
							linebot.NewURIAction("View detail", "http://example.com/page/123"),
						),
					)

					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				case message.Text == "location":
					resp := linebot.NewLocationMessage("現在地", "宮城県多賀城市", 38.297807, 141.031)

					_, err = bot.ReplyMessage(event.ReplyToken, resp).Do()
					if err != nil {
						log.Print(err)
					}
				default:
					// originalContentURL := "https://i.pinimg.com/736x/65/14/e8/6514e88d1bc17011c076cd525ac8e7df.jpg"
					// previewImageURL := "https://i.pinimg.com/736x/65/14/e8/6514e88d1bc17011c076cd525ac8e7df.jpg"

					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(originalContentURL, previewImageURL)).Do(); err != nil {
					// 	log.Print(err)
					// }
				}
			}
		}
	}
}
