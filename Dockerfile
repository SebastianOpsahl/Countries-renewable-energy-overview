# base image we build our image upon 
FROM golang:1.19

LABEL maintainer="sebastop@ntnu.no" \ alsoMaintainer="sondre.m.eggan@ntnu.no"

#sets up execution enviroment, where it is ran from
WORKDIR /go/src/app/cmd/server

#copies relevant folders into it's placement in the image
COPY ./cmd /go/src/app/cmd
COPY ./functions /go/src/app/functions
COPY ./go.mod /go/src/app/go.mod

# builds the application
RUN GOOS=linux go build -a -o server ./cmd/server/

#indicates port on which server listens
EXPOSE 8080

# executes the workdirectory
CMD ["./server"]
