FROM golang:1.18.1-alpine3.15 as builder

RUN apk add --no-cache git jq

ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build . && \
    mv atlassian-automator /usr/local/bin/

FROM debian:stable-20220622
COPY --from=builder /usr/local/bin/atlassian-automator /usr/local/bin/atlassian-automator

ENTRYPOINT ["atlassian-automator"]
