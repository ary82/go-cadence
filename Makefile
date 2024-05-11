.PHONY: run
run: build
	@echo "running..."
	@./bin/temp

.PHONY: build
build:
	@echo "building..."
	@go build -o ./bin/temp ./main.go

.PHONY: watch
watch:
	@${HOME}/go/bin/air

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	@echo "removing bin/ files"
	@rm ./bin/*
	@echo "removing tmp/ files"
	@rm ./tmp/*

.PHONY: cadence-init
cadence-init:
	@docker run -it --rm ubercadence/cli:master --address host.docker.internal:7933 --do test-domain d re
