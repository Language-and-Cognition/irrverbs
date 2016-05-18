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

	initIfNotExists(globalDb)

	bot := tgbot.NewTgBot(cfg.Telegram.Token)
	bot.SimpleCommandFn(`learning`, startLearningHandler)
	bot.SimpleCommandFn(`statistics`, statisticsHandler)
	bot.SimpleCommandFn(`help`, helpHandler)
	bot.SimpleCommandFn(`clear_last_statistics`, clearStatisticsHandler)
	bot.SimpleCommandFn(`nuke_all_statistics`, nukeStatisticsHandler)
	bot.NotCalledFn(answerHandler)
	bot.SimpleStart()
}

func getRandomVerb() string {
	for key := range getAllVerbs() {
		return key
	}
	return "cut"
}

func helpHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	message := `You can start /learning
You can see your /statistics
You can /clear_last_statistics
You can /nuke_all_statistics
Have fun!`
	bot.Answer(msg).Text(message).End()
	return nil
}

func startLearningHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	if !doesUserExist(globalDb, msg.Chat.ID) {
		createUser(globalDb, msg.Chat.ID)
	}
	verb := getRandomVerb()
	usersAnswering.set(msg.Chat.ID, verb)
	bot.Answer(msg).Text(verb).End()
	return nil
}

func statisticsHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	if !doesUserExist(globalDb, msg.Chat.ID) {
		createUser(globalDb, msg.Chat.ID)
	}
	right, wrong, since := getLastStatistics(globalDb, msg.Chat.ID)
	rightOverall, wrongOverall := getOverallStatistics(globalDb, msg.Chat.ID)
	format := "Last statistics since %s:\nRight: %d\nWrong: %d\nRatio: %f\n\nOverall:\nRight: %d\nWrong: %d\nRatio: %f"
	ratio := float64(right) / float64(right+wrong)
	ratioOverall := float64(rightOverall) / float64(rightOverall+wrongOverall)
	answer := fmt.Sprintf(format, since, right, wrong, ratio, rightOverall, wrongOverall, ratioOverall)
	bot.Answer(msg).Text(answer).End()
	return nil
}

func clearStatisticsHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	clearLastStatistics(globalDb, msg.Chat.ID)
	bot.Answer(msg).Text("Last statistics has been cleared").End()
	return nil
}

func nukeStatisticsHandler(bot tgbot.TgBot, msg tgbot.Message, text string) *string {
	nukeAllStatistics(globalDb, msg.Chat.ID)
	bot.Answer(msg).Text("All statistics has been cleared").End()
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
		addRightAnswer(globalDb, msg.Chat.ID)
		bot.Answer(msg).Text("Correct!").End()
	} else {
		addWrongAnswer(globalDb, msg.Chat.ID)
		bot.Answer(msg).Text(fmt.Sprintf("Incorrect. The right answer is %s %s", v2, v3)).End()
	}
	verb := getRandomVerb()
	usersAnswering.set(msg.Chat.ID, verb)
	bot.Answer(msg).Text(verb).End()
}

func getAllVerbs() map[string][]string {
	return GetEnglishVerbs()
}
