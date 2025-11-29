FROM golang:1.25.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY ui ./ui
COPY sql ./sql
COPY Makefile ./
COPY engine.db ./
RUN ls -l

RUN CGO_ENABLED=0 GOOS=linux make build

CMD ["/pse"]