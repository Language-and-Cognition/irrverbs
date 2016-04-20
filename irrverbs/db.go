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
	db.Exec("INSERT INTO overall (id, right, wrong) VALUES (?, 0, 0)", id)
	db.Exec("INSERT INTO last (id, right, wrong, since) VALUES (?, 0, 0, date('now'))", id)
}

func removeLastStatistics(db *sql.DB, id int) {
	db.Exec("UPDATE last SET right = 0, wrong = 0, since = datetime('now') WHERE id = ?", id)
}

func nukeAllStatistics(db *sql.DB, id int) {
	removeLastStatistics(db, id)
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

func getLastStatistics(db *sql.DB, id int) (int, int, time.Time) {
	var right, wrong int
	var since time.Time
	db.QueryRow("SELECT right, wrong, since FROM last WHERE id = ?", id).Scan(&right, &wrong, &since)
	return right, wrong, since
}

func initIfNotExists(db *sql.DB) {
	db.Exec("CREATE TABLE overall (id INTEGER CONSTRAINT id PRIMARY KEY, right INTEGER, wrong INTEGER)")
	db.Exec("CREATE TABLE last (id INTEGER CONSTRAINT id PRIMARY KEY, right INTEGER, wrong INTEGER, since TEXT)")
}
