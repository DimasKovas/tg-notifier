package main

import (
    // "fmt"
    "log"
    "os"
    "net/http"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Bot *tgbotapi.BotAPI

var IdTable map[string]int64

func handler(w http.ResponseWriter, r *http.Request) {
    _, ok := r.Header["UserID"]
    if ok {
        id, ok := IdTable[r.Header["UserID"][0]]
        if ok {
            msg := tgbotapi.NewMessage(id, "Done")
            Bot.Send(msg)
        }
    } 
    // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func httpHandler() {
    log.Panic(http.ListenAndServe(":8080", nil))
}

func main() {
    http.HandleFunc("/", handler)
    go httpHandler()
    var err error
    Bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
    }
    
    Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

    updates, err := Bot.GetUpdatesChan(u)
    if err != nil {
        log.Panic(err)
    }

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
        }
        if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "Cant help you"
            case "start":
                msg.Text = "OK"
			case "stop":
				msg.Text = "Hi :)"
			default:
				msg.Text = "I don't know that command"
			}
			Bot.Send(msg)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
        log.Printf("%d", update.Message.Chat.ID)
		msg.ReplyToMessageID = update.Message.MessageID

		Bot.Send(msg)
	}

}