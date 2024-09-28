FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
RUN ls -l
COPY internal ./internal
COPY ui ./ui
COPY ddl.sql ./
COPY Makefile ./
COPY engine.db ./
RUN ls -l

RUN CGO_ENABLED=0 GOOS=linux make build

CMD ["/pse"]