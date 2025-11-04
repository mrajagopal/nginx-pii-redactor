.PHONY: all clean build run

build: main.go
	go build -o pii-redactor

clean: 
	rm -f pii-redactor test.json nginx.redacted.conf

run: 
	./pii-redactor nginx.conf