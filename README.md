## Getting Started <a name = "getting_started"></a>

### API Specification

-   Swagger Docs: /docs

### Prerequisites

-   Setup env variable
-   Install dependencies

### Run Migration

```
./server migrate up
./server migrate down
./server migrate fresh
```

## Build <a name="build"></a>

```
go build -o server ./cmd/main.go
```

## Usage <a name="usage"></a>

### Run Server

```
./server serve
```

## Development <a name="development"></a>

```bash
# swag command
# refer https://github.com/swaggo/swag
go install github.com/swaggo/swag/cmd/swag@latest

# mockery command
# refer https://github.com/vektra/mockery
go install go install github.com/vektra/mockery/v2@v2.43.1

# air (golang hot reload)
# refer https://github.com/cosmtrek/air
go install github.com/cosmtrek/air@latest
```

```
air serve
```

for windows

```
make air-win
```

### Generate Swagger Docs

install [swaggo](https://github.com/swaggo/swag)

```
make swag
```
