# build application
build:
	go build -o ./bin/main cmd/storage/main.go
	go build -o ./bin/mainCli cmd/storage-cli/main.go

# run storage
run-storage:
	./bin/main
# run cli
run-storage-cli:
	./bin/mainCli