FROM golang:1.19.5-alpine

WORKDIR /app

COPY . ./
RUN go mod download


RUN go build -o /god-of-right-go

EXPOSE 8080

CMD [ "/god-of-right-go" ]