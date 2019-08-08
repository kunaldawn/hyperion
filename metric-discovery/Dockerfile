# BUILD STAGE
FROM golang:alpine AS build_stage
RUN apk --no-cache add ca-certificates git
WORKDIR /hyperion/
ADD discovery.go .
RUN go get "github.com/docker/docker/client"
RUN go build -o discovery discovery.go

# FINAL STAGE
FROM alpine as final_state
WORKDIR /hyperion/
COPY --from=build_stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build_stage /hyperion/discovery .
ENTRYPOINT ["/hyperion/discovery"]
