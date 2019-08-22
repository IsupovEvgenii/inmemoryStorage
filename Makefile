# build application
build:
	go build -o ./bin/main cmd/engine/main.go

# run application
run:
	./bin/main 123

deploy:
	./scripts/deploy.sh