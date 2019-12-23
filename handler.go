package main

import (
    // "fmt"
    "log"
    "os"
    "net/http"
    "strings"
    "math/rand"
    "sync"
    "strconv"
    "io/ioutil"
    "time"
    "encoding/json"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Bot *tgbotapi.BotAPI

type KeysData struct {
    sync.RWMutex
    ChatKeys    map[int64]string
}

func (data *KeysData) Save(path string) error {
    encodedKeys, err := json.Marshal(data.ChatKeys)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(path, encodedKeys, 0644)
    return err
}

func (data *KeysData) Load(path string) error {
    encodedKeys, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    err = json.Unmarshal(encodedKeys, &data.ChatKeys)
    return err
}

const DataPath = "/data/chat_keys.txt"

var Keys KeysData

func httpHandler(w http.ResponseWriter, r *http.Request) {
    tokens := r.URL.Query()["token"]
    if len(tokens) != 1 {
        // error, should be exactly one token specified
        http.Error(w, "Should be exactly one token specified", 500)
        return
    }
    tokenParts := strings.Split(tokens[0], "-")
    var chatIdStr, chatKey string
    if len(tokenParts) == 2 {
        chatIdStr = tokenParts[0]
        chatKey = tokenParts[1]
    }
    chatId, err := strconv.ParseInt(chatIdStr, 36, 64)
    Keys.RLock()
    checkPassed := (len(tokenParts) == 2 && err == nil && Keys.ChatKeys[chatId] == chatKey)
    Keys.RUnlock()
    if !checkPassed {
        http.Error(w, "Invalid token", 500)
        log.Printf("Recieved an invalid token '%s'", tokens[0])
        return
    }
    text, err := ioutil.ReadAll(r.Body)
    if err != nil {
        // error reading body
        http.Error(w, "Error reading body", 500)
        log.Print(err)
        return
    }
    _, err = Bot.Send(tgbotapi.NewMessage(chatId, string(text)))
    if err != nil {
        http.Error(w, "Error sending message", 500)
        log.Print(err)
    }
}

func generateKey() string {
    const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
    res := make([]byte, 10)
    for i := 0; i < 10; i++ {
        res[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(res)
}

func formatToken(chatId int64, chatKey string) string {
    return strconv.FormatInt(chatId, 36) + "-" + chatKey
}

func incommingMessagesHandler() {
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
		    log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
            chatId := update.Message.Chat.ID
            msg := tgbotapi.NewMessage(chatId, "")
			switch update.Message.Command() {
			case "help":
                msg.Text =  "Usages:\n" +
                            "/help - help\n" +
                            "/start - start the bot\n" +
                            "/token - get your chat token\n" +
                            "/stop - stop the bot\n"
            case "start":
                chatKey := generateKey()
                Keys.Lock()
                _, exists := Keys.ChatKeys[chatId]
                if exists {
                    msg.Text = "The bot has been already started!"
                } else {
                    Keys.ChatKeys[chatId] = chatKey
                    if err := Keys.Save(DataPath); err != nil {
                        log.Panic(err)
                    }
                    log.Printf("ChatKey [%s] for ChatId [%v] was successfully generated", chatKey, chatId)
                    msg.Text = "Welcome! Your chat token is " + formatToken(chatId, chatKey) + ". Use it for requests in our API."
                }
                Keys.Unlock()
            case "token":
                Keys.RLock()
                chatKey, exists := Keys.ChatKeys[chatId]
                Keys.RUnlock()
                if exists {
                    msg.Text = "Your chat token is " + formatToken(chatId, chatKey)
                } else {
                    msg.Text = "The bot hasn't been started! Use /start"
                }
            case "stop":
                Keys.Lock()
                _, exists := Keys.ChatKeys[chatId]
                if exists {
                    delete(Keys.ChatKeys, chatId)
                    if err := Keys.Save(DataPath); err != nil {
                        log.Panic(err)
                    }
                    log.Printf("ChatKey [%s] for ChatId [%v] was successfully deleted")
                    msg.Text = "Bot successfully stoped! Use /start to start it again."
                } else {
                    msg.Text = "Ooops, the bot hasn't been started! Use /start"
                }
                Keys.Unlock()
			default:
				msg.Text = "Unexpected command, Use /help to see all commands."
			}
			Bot.Send(msg)
		}
	}
}

func main() {
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/", httpHandler)
    Keys.ChatKeys = make(map[int64]string)
    
    if _, err := os.Stat(DataPath); err == nil {
        err = Keys.Load(DataPath)
        if err != nil {
            log.Panic(err)
        }
        log.Printf("Chat keys were successfully loaded")
    } else if !os.IsNotExist(err) {
        log.Panic(err)
    }

    var err error
    Bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
    }
    Bot.Debug = true
    log.Printf("Authorized on account %s", Bot.Self.UserName)
    
    go incommingMessagesHandler()
	
    log.Panic(http.ListenAndServe(":8080", nil))
}