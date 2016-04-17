package main

import (
	"fmt"
	"sync"

	"github.com/rockneurotiko/go-tgbot"
)

type usersAnsweringStruct struct {
	*sync.RWMutex
	Users map[int]string
}

var usersAnswering = usersAnsweringStruct{&sync.RWMutex{}, make(map[int]string)}

func (users *usersAnsweringStruct) get(user int) (string, bool) {
	users.RLock()
	s, ok := users.Users[user]
	users.RUnlock()
	return s, ok
}

func (users *usersAnsweringStruct) set(user int, value string) {
	users.Lock()
	users.Users[user] = value
	users.Unlock()
}

func (users *usersAnsweringStruct) del(user int) {
	users.Lock()
	delete(users.Users, user)
	users.Unlock()
}

func main() {
	cfg, _ := getConfig()
	bot := tgbot.NewTgBot(cfg.Telegram.Token)
	bot.CommandFn(`echo (.+)`, echoHandler)
	bot.SimpleCommandFn(`learninig`, learnHandler)
	bot.NotCalledFn(answerHandler)
	bot.SimpleStart()
}

func learnHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	key := "understand" // TODO: must be random
	usersAnswering.set(msg.Chat.ID, key)
	bot.Answer(msg).Text(key).ReplyToMessage(msg.ID).End()
	return nil
}

func answerHandler(bot tgbot.TgBot, msg tgbot.Message) {
	if msg.Text == nil {
		return
	}
	s, ok := usersAnswering.get(msg.Chat.ID)
	usersAnswering.del(msg.Chat.ID)
	if !ok {
		bot.Answer(msg).Text("You need to start /learninig first").End()
		return
	}
	verbs := getAllVerbs()
	// TODO: Check answer
	bot.Answer(msg).Text(fmt.Sprintf("%s %s", verbs[s][0], verbs[s][1])).End()
}

func echoHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	fmt.Println(vals, kvals)
	newmsg := fmt.Sprintf("[Echoed]: %s", vals[1])
	return &newmsg
}

func getAllVerbs() map[string][]string {
	return GetEnglishVerbs()
}
