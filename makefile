build:
	go build irrverbs/client.go irrverbs/config.go irrverbs/english_verbs.go
test:
	go test irrverbs/db_test.go irrverbs/client.go irrverbs/config.go irrverbs/db.go irrverbs/english_verbs.go
