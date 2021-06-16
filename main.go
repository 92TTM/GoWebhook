package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
)

var bot *linebot.Client

func main() {
	viper.SetConfigName("setting")
	//viper.SetConfigName("webhook")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()
	viper.WatchConfig()
	var port string
	var connstring string
	if viper.GetString("env") == "test" {
		port = "80"
		connstring = viper.GetString("databasecon")
	} else {
		port = os.Getenv("PORT")
		connstring = os.Getenv("DATABASE_URL")
	}
	db, err := sql.Open("postgres", connstring)
	checkerror(err)
	rows, err := db.Query("select remark from sys_webhook_config")
	fmt.Println(err)
	checkerror(err)
	for rows.Next() {
		var remark string
		rows.Scan(&remark)
		fmt.Println(remark)
	}
	defer db.Close()
	http.HandleFunc("/callback", callbackHandler)
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	client := r.URL.Query()
	CNO := client.Get("client")
	fmt.Println(CNO)
	secret, TK := checksource(CNO)
	fmt.Println("secret:" + secret + ",TK=" + TK)
	bot, err = linebot.New(secret, TK)
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
				} else if message.Text == "恬恬喵" {
					replytext = "好萌！！"
				} else {
					replytext = "沒答案！"
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

func checksource(client string) (string, string) {
	/*
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	*/
	return viper.GetString(client + ".Sercet"), viper.GetString(client + ".AccessTK")
}

func checkerror(err error) {
	if err != nil {
		panic(err)
	}

}
