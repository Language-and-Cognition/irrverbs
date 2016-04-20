package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rockneurotiko/go-tgbot"
)

var globalDb *sql.DB

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

	globalDb, _ = sql.Open("sqlite3", "./statistics.db")
	defer globalDb.Close()

	bot := tgbot.NewTgBot(cfg.Telegram.Token)
	bot.CommandFn(`echo (.+)`, echoHandler)
	bot.SimpleCommandFn(`learning`, startLearningHandler)
	bot.NotCalledFn(answerHandler)
	bot.SimpleStart()
}

func getRandomVerb() string {
	for key := range getAllVerbs() {
		return key
	}
	return "cut"
}

func startLearningHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	verb := getRandomVerb()
	usersAnswering.set(msg.Chat.ID, verb)
	bot.Answer(msg).Text(verb).End()
	return nil
}

func answerHandler(bot tgbot.TgBot, msg tgbot.Message) {
	if msg.Text == nil {
		return
	}
	s, ok := usersAnswering.get(msg.Chat.ID)
	if *msg.Text == "/stop" {
		usersAnswering.del(msg.Chat.ID)
		return
	}
	if !ok {
		bot.Answer(msg).Text("You need to start /learning first").End()
		return
	}
	verbs := getAllVerbs()
	userVerbs := strings.Split(*msg.Text, " ")
	if len(userVerbs) != 2 {
		bot.Answer(msg).Text("Answer should be two verbs separated by space").End()
		return
	}
	v2, v3 := verbs[s][0], verbs[s][1]
	if strings.ToLower(userVerbs[0]) == v2 && strings.ToLower(userVerbs[1]) == v3 {
		bot.Answer(msg).Text("Correct!").End()
	} else {
		bot.Answer(msg).Text(fmt.Sprintf("Incorrect. The right answer is %s %s", v2, v3)).End()
	}
	verb := getRandomVerb()
	usersAnswering.set(msg.Chat.ID, verb)
	bot.Answer(msg).Text(verb).End()
}

func echoHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	fmt.Println(vals, kvals)
	newmsg := fmt.Sprintf("[Echoed]: %s", vals[1])
	return &newmsg
}

func getAllVerbs() map[string][]string {
	return GetEnglishVerbs()
}
