# stage 1: build application 
FROM golang:1.25.1-alpine AS builder 

# install necessary dependencies for the build 
RUN apk add --no-cache git gcc musl-dev 

# defines the working directory 
WORKDIR /app 

# copy the files of the modules 
COPY go.mod go.sum ./ 
RUN go mod download 

# copies source code 
COPY . . 

# compiles the application
RUN CGO_ENABLED=0 GOOS=linux go build -o chirpy . 

# stage 2: final smaller image 
FROM alpine:latest 

WORKDIR /root/ 

# copies the compiled binary of builder stage 
COPY --from=builder /app/chirpy . 

# expose the port 
EXPOSE 8080 

# command to run the application 
CMD ["./chirpy"]
