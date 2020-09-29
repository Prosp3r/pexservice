#FROM golang:1.12.0-alpine3.9
FROM golang:alpine 
RUN apk add git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN ls -aslt
RUN go build -o pexservice
RUN ls -aslt /app
EXPOSE 8081
CMD ["/app/pexservice"]