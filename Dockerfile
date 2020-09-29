#FROM golang:1.12.0-alpine3.9
FROM golang:alpine 
RUN apk add git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main pexservice
EXPOSE 8080
CMD ["./app/main 8080"]