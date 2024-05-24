.PHONY: mocks
mocks:
	sleep 1 && rm -rfd mocks && mockery

swag:
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go

run:  swag mocks
	go run cmd/main.go serve

build:  swag mocks
	go build -o tmp/main cmd/main.go

build-win:  swag
	go build -o tmp/main.exe cmd/main.go
		
docker-run:
	docker run -d -p 8080:8080 \
	--name url-shortener \
	url-shortener:latest

air:
	air -c .air.toml serve

air-win:
	air -c .air.win.toml serve