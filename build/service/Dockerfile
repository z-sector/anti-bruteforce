FROM golang:1.20 as build

ENV CODE_DIR /go/src/
WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .. ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o /opt/app ./cmd/service/main.go


FROM scratch

COPY --from=build /opt/app /opt/app

ENTRYPOINT ["/opt/app"]
