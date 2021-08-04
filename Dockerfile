FROM golang:alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

# Copy config file and set ENV CONFIG_PATH to point to it
COPY service.json ./
ENV CONFIG_PATH=/app/service.json

RUN go mod download
COPY *.go ./
RUN go build -o /list-service-build

EXPOSE 4000


CMD [ "/list-service-build" ]




