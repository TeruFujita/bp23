package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	linebot "github.com/line/line-bot-sdk-go/v7/linebot"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// .envファイルを読み込む（エラーは無視：本番では環境変数を使う想定）
	_ = godotenv.Load("../.env.local")

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	if channelSecret == "" || channelToken == "" {
		log.Fatal("LINE_CHANNEL_SECRET and LINE_CHANNEL_TOKEN must be set in environment")
	}

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		for _, ev := range events {
			switch ev.Type {
			case linebot.EventTypeMessage:
				if msg, ok := ev.Message.(*linebot.TextMessage); ok {
					// 見た目確認用の簡易Flexメッセージ（タイトル+説明+ボタン）
					flex := &linebot.BubbleContainer{
						Type: linebot.FlexContainerTypeBubble,
						Body: &linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{Type: linebot.FlexComponentTypeText, Text: "Detour Search", Weight: linebot.FlexTextWeightTypeBold, Size: linebot.FlexTextSizeTypeLg},
								&linebot.TextComponent{Type: linebot.FlexComponentTypeText, Text: "入力: " + msg.Text, Wrap: true, Size: linebot.FlexTextSizeTypeSm, Color: "#666666"},
							},
						},
						Footer: &linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type:  linebot.FlexComponentTypeButton,
									Style: linebot.FlexButtonStyleTypePrimary,
									Color: "#0EA5E9",
									Action: &linebot.MessageAction{
										Label: "もう一度",
										Text:  "もう一度",
									},
								},
							},
						},
					}
					if _, err := bot.ReplyMessage(ev.ReplyToken, linebot.NewFlexMessage("Detour", flex)).Do(); err != nil {
						log.Println("reply error:", err)
					}
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	port := getEnv("PORT", "3001")
	log.Println("LINE Bot listening on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
