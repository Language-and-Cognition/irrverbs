package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func pretestInit() *sql.DB {
	var db *sql.DB
	db, _ = sql.Open("sqlite3", ":memory:")
	initIfNotExists(db)
	createUser(db, 1)
	return db
}

func TestDoesUserExsist(t *testing.T) {
	db := pretestInit()
	defer db.Close()
	userExists := doesUserExist(db, 1)
	if !userExists {
		t.Error("User with id 1 does not exist")
	}
	userExists = doesUserExist(db, 2)
	if userExists {
		t.Error("User with id 2 exists")
	}
}

func TestAddAnswers(t *testing.T) {
	db := pretestInit()
	defer db.Close()

	addRightAnswer(db, 1)
	addRightAnswer(db, 1)
	addRightAnswer(db, 1)
	addWrongAnswer(db, 1)
	addWrongAnswer(db, 1)
	addWrongAnswer(db, 1)
	addWrongAnswer(db, 1)

	right, wrong, _ := getLastStatistics(db, 1)
	rightOverall, wrongOverall := getOverallStatistics(db, 1)
	t.Log(right, rightOverall, wrong, wrongOverall)

	if right != rightOverall {
		t.Error("Statistics are not equal")
		t.Log(right, " != ", rightOverall)
	}

	if wrong != wrongOverall {
		t.Error("Statistics are not equal")
		t.Log(wrong, " != ", wrongOverall)
	}

	if right != 3 {
		t.Error("There should be 3 right answers")
		t.Log(right, " != ", 3)
	}

	if wrong != 4 {
		t.Error("There should be 4 wrong answers")
		t.Log(right, " != ", 4)
	}
}
