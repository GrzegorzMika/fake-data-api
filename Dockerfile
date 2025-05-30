FROM golang:1.24 AS development

WORKDIR /fake-data-api

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/fake-data-api ./main.go
RUN chmod a+x /fake-data-api

FROM golang:1.24-alpine AS app

EXPOSE 8080

COPY --from=development /fake-data-api/build/fake-data-api /fake-data-api

CMD [ "/fake-data-api" ]