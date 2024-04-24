FROM golang:1.22 
WORKDIR /usr/src/app
COPY . .
RUN go mod download 
RUN go build -v -o /usr/local/bin/app
CMD ["app"]