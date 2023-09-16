FROM golang:1.20 as build

ENV BIN_FILE /opt/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .. ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} ./cmd/service/main.go


FROM alpine:3.9

ENV BIN_FILE "/opt/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
