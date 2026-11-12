.PHONY: ktlo

ktlo:
	mkdir -p ./bin
	go build -o ./bin/ktlo .
