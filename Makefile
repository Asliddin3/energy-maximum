swag:
	swag init -g cmd/crm/main.go

migrate:
	go run ./migrate/migrate.go

run:
	 go run cmd/crm/main.go

getswag:
	export PATH=$(go env GOPATH)/bin:$PATH
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

golanglint:
	golangci-lint run

air:
	kill -9 $(lsof -ti :8080)
	air