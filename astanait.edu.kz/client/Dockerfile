FROM golang
RUN mkdir /app
COPY . /app
ENV CGO_ENABLED=0 \
    GOOS=linux
WORKDIR /app
RUN go build -o client
