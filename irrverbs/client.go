package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rockneurotiko/go-tgbot"

	_ "github.com/mattn/go-sqlite3"
)

var globalDb *sql.DB

func main() {
	if !checkArgs() {
		usage()
		os.Exit(1)
	}
	globalDb, _ = sql.Open("sqlite3", "./verbs.db")
	defer globalDb.Close()

	if os.Args[1] == "learn" {
		learn()
	}

	cfg, _ := getConfig()
	bot := tgbot.NewTgBot(cfg.Telegram.Token)
	bot.CommandFn(`echo (.+)`, echoHandler)
	bot.SimpleStart()

}

func echoHandler(bot tgbot.TgBot, msg tgbot.Message, vals []string, kvals map[string]string) *string {
	fmt.Println(vals, kvals)
	newmsg := fmt.Sprintf("[Echoed]: %s", vals[1])
	return &newmsg
}

func checkArgs() bool {
	if len(os.Args) < 2 {
		return false
	} else if os.Args[1] != "add" && os.Args[1] != "learn" {
		return false
	} else if os.Args[1] == "add" && len(os.Args) != 5 {
		return false
	} else if os.Args[1] == "learn" && len(os.Args) != 2 {
		return false
	}
	return true
}

func getAllVerbs() map[string][]string {
	var verbs map[string][]string
	verbs = make(map[string][]string)
	rows, _ :=
		globalDb.Query("SELECT v1, v2, v3 FROM english_irregular_verbs")
	for rows.Next() {
		var v1, v2, v3 string
		rows.Scan(&v1, &v2, &v3)
		verbs[v1] = []string{v2, v3}
	}
	return verbs
}

func learn() {
	verbs := getAllVerbs()
	for key, value := range verbs {
		fmt.Println(key, value)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s add v1 v2 v3\n", os.Args[0])
	fmt.Printf("\t%s learn\n", os.Args[0])
}
