############################
# STEP 1 build executable binary
############################
# FROM golang:alpine AS builder
FROM public.ecr.aws/bitnami/golang:latest as builder
# Install git.
# Git is required for fetching the dependencies.
RUN update-ca-certificates
# RUN apk update && apk add --no-cache git
# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/github.com/prabhatsharma/open-telemetry1/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Using go mod.
# RUN go mod download
# RUN go mod verify
# Build the binary.
# to tackle error standard_init_linux.go:207: exec user process caused "no such file or directory" set CGO_ENABLED=0   
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o otel1
############################
# STEP 2 build a small image
############################
# FROM public.ecr.aws/lts/ubuntu:latest
# FROM public.ecr.aws/amazonlinux/amazonlinux:latest
# debugging
FROM public.ecr.aws/bitnami/aws-cli:latest 
# FROM scratch
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the ssl certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy our static executable.
COPY --from=builder  /go/src/github.com/prabhatsharma/open-telemetry1/otel1 /go/bin/otel1

# Use an unprivileged user.
USER appuser:appuser
# Port on which the service will be exposed.
EXPOSE 6080
# Run the otel1 binary.
ENTRYPOINT ["/go/bin/otel1"]
