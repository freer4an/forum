# __Stage 1 (build)__

# Get minimal image of linux for Go (to run Go commands)
FROM golang:alpine AS build

# install gcc to run go-sqlite3
RUN apk --no-cache add gcc musl-dev

# Choose root dir in the image (if path doesn't exist create it automatically)
WORKDIR /app

#install 3rd party packages
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# COPY [from local] [to container]
# First dot "." means that we will copy everything from this dockerfile's root dir
# Second dot "." means copy from first to root dir in image which is our WORKDIR

COPY . .

# Create executable file of our main.go in the image's "." root dir
ENV CGO_ENABLED=1
RUN go build ./cmd/main.go

# __Final stage (2)__

# Get clear version of alpine linux (read abot it here: https://hub.docker.com/_/alpine)
FROM alpine:latest

# Choose root dir in the final image
WORKDIR /app

# COPY --from=<stage> 
# Copy from our first stage (build) only our root dir /app where located our project and it's executable
# Copy /app to root dir of final image "./"
COPY --from=build /app ./

# Command to run when launching the final image 
CMD [ "./main" ]