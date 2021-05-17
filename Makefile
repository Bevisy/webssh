.PHONY: build clean

APP_NAME = webssh

build:
	GOOS=linux GOARCH=x86_64 go build -o ${APP_NAME}

clean:
	rm ${APP_NAME}