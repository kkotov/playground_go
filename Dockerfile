######### build stage #########
FROM golang:1.16-alpine AS GO_BUILD

# Copy code over assuming build is ran from the top module directory (adjust COPY if ran manually)
COPY tokenizer/tokenizer.go /go/src/
RUN echo -e "module tokenizer\\ngo 1.16" > /go/src/go.mod
WORKDIR /go/src/

# Build the app
RUN go get
RUN go build -o /go/bin/tokenizer

######### execution stage #########
FROM alpine:3.10
WORKDIR /app

# Copy built app from the build stage
COPY --from=GO_BUILD /go/bin/ ./

# Container will listen to that port. For documentation purposes. This does not do any real binding of the port.
EXPOSE 8180
CMD ["/app/tokenizer"]
