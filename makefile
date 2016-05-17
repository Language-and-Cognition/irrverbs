build:
	go build irrverbs/client.go irrverbs/config.go irrverbs/english_verbs.go irrverbs/db.go
test:
	go test irrverbs/db_test.go irrverbs/db.go
