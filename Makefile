swag:
	swag init -g cmd/main.go

migrate:
	go run ./migrate/migrate.go

run:
	 go run cmd/main.go

getswag:
	export PATH=$(go env GOPATH)/bin:$PATH
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

server:
	nohup go run cmd/main.go &

stop:
	pgrep -f "go run cmd/main.go"

golanglint:
	golangci-lint run
kill:
	kill -9 $(lsof -ti :8000)

build:
	go build ./cmd/main.go
	pm2 restart ./main 
	pm2 log
air:
	kill -9 $(lsof -ti :8080)
	air