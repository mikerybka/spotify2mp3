FROM alpine:latest

RUN apk update
RUN apk add go yt-dlp
RUN go install github.com/mikerybka/dl-mp3@latest

WORKDIR /app
COPY . .

RUN go build -o /bin/server main.go

