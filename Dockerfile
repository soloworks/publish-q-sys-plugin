# STEP 1 build executable binary
FROM golang:alpine as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -d -v
# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main .
RUN apk --update add ca-certificates
# STEP 2 build action image
FROM scratch

# Copy our static executable
COPY --from=builder /main /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


# Code file to execute when the docker container starts up
ENTRYPOINT ["/main"]