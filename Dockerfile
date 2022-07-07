FROM golang:1.18.3 as dev

RUN apt-get install -y exiftool git jq

ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build . && \
    mv atlassian-automator /usr/local/bin/

FROM debian:stable-20220622
COPY --from=dev /usr/local/bin/atlassian-automator /usr/local/bin/atlassian-automator

ENTRYPOINT ["atlassian-automator"]
