# RediSearch Shopping

> A sample project using RediSearch

## Running the project

### Fetching Dependencies

```sh
go mod download
```

### Building the project

```sh
go build -o ./shopping cmd/redisearch-shopping/main.go
```

### Generating mock products data

```sh
./shopping mock
```

### Running the project

Make sure that RediSearch is up and running on your local machine.

```sh
./shopping serve -m "debug" -p "8080"
```