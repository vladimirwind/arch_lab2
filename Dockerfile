FROM golang:1.21.3
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["./app"]