FROM golang:1.18.3-bullseye as dev

RUN apt-get update && apt-get install -y \
    exiftool \
    git \
    jq \
    && rm -rf /var/lib/apt/lists/*

ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM debian:stable-20220622-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=dev /app/atlassian-automator /app/atlassian-automator

ENTRYPOINT ["/app/atlassian-automator"]
