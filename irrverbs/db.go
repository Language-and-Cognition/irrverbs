package main

import (
	"database/sql"
	"time"
)

func doesUserExist(db *sql.DB, id int) bool {
	var temp int
	err := db.QueryRow("SELECT id FROM overall WHERE id = ?", id).Scan(&temp)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func createUser(db *sql.DB, id int) {
	now := time.Now().UTC()
	db.Exec("INSERT INTO overall (id, right, wrong) VALUES (?, 0, 0)", id)
	db.Exec("INSERT INTO last (id, right, wrong, since) VALUES (?, 0, 0, ?)", id, now.Format(time.RFC822))
}

func clearLastStatistics(db *sql.DB, id int) {
	now := time.Now().UTC()
	db.Exec("UPDATE last SET right = 0, wrong = 0, since = ? WHERE id = ?", now.Format(time.RFC822), id)
}

func nukeAllStatistics(db *sql.DB, id int) {
	clearLastStatistics(db, id)
	db.Exec("UPDATE overall SET right = 0, wrong = 0 WHERE id = ?", id)
}

func addRightAnswer(db *sql.DB, id int) {
	db.Exec("UPDATE overall SET right = right + 1 WHERE id = ?", id)
	db.Exec("UPDATE last SET right = right + 1 WHERE id = ?", id)
}

func addWrongAnswer(db *sql.DB, id int) {
	db.Exec("UPDATE overall SET wrong = wrong + 1 WHERE id = ?", id)
	db.Exec("UPDATE last SET wrong = wrong + 1 WHERE id = ?", id)
}

func getOverallStatistics(db *sql.DB, id int) (int, int) {
	var right, wrong int
	db.QueryRow("SELECT right, wrong FROM overall WHERE id = ?", id).Scan(&right, &wrong)
	return right, wrong
}

func getLastStatistics(db *sql.DB, id int) (int, int, string) {
	var right, wrong int
	var since string
	db.QueryRow("SELECT right, wrong, since FROM last WHERE id = ?", id).Scan(&right, &wrong, &since)
	return right, wrong, since
}

func initIfNotExists(db *sql.DB) {
	db.Exec("CREATE TABLE IF NOT EXISTS overall (id INTEGER CONSTRAINT id PRIMARY KEY, right INTEGER, wrong INTEGER)")
	db.Exec("CREATE TABLE IF NOT EXISTS last (id INTEGER CONSTRAINT id PRIMARY KEY, right INTEGER, wrong INTEGER, since TEXT)")
}
