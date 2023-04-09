FROM golang:1.20-buster  as build

WORKDIR /postLogger

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o postLogger ./cmd/postLogger

CMD ./postLogger

FROM alpine:latest as release

RUN apk --no-cache add ca-certificates

COPY --from=build /postLogger ./ 

RUN chmod +x ./postLogger

ENTRYPOINT [ "./postLogger" ]

EXPOSE 3200