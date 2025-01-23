FROM golang:1.23 AS build

ENV BIN_FILE=/opt/filter-bot/filter-bot
ENV CODE_DIR=/go/src

WORKDIR ${CODE_DIR}

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ${BIN_FILE} -mod vendor ${CODE_DIR}/cmd/botcleaner/*

FROM alpine:3.9

LABEL SERVICE="filter-bot"
ENV BIN_FILE=/opt/filter-bot/filter-bot
ENV CODE_DIR=/go/src

WORKDIR ${CODE_DIR}
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
