# ------------------------------------------------------------------------------
# Test image
# ------------------------------------------------------------------------------
    FROM golang:alpine AS test_img

    RUN apk update && apk upgrade && apk add --no-cache git openssh build-base
    
    # setup source directory, put outside GOPATH to enable go modules
    ENV APP_DIR=/src
    RUN mkdir -p $APP_DIR
    WORKDIR $APP_DIR
    
    # allow caching of go modules, only refetch when dependencies change
    COPY go.mod .
    COPY go.sum .
    RUN go mod download
    
    COPY . .
    
    CMD ["go", "test", "./..."]
    
    
    # ------------------------------------------------------------------------------
    # Development image
    # ------------------------------------------------------------------------------
    FROM test_img AS dev_img
    
    WORKDIR $APP_DIR/cmd/loan-engine
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
       go build -gcflags "all=-N -l" -o /loan-engine
    ADD cmd/loan-engine/.env /.env
    
    ENTRYPOINT /loan-engine
    
    
    # ------------------------------------------------------------------------------
    # Production image
    # ------------------------------------------------------------------------------
    FROM alpine:3.7 as prod_img
    COPY --from=dev_img /loan-engine /
    EXPOSE 8080
    ENTRYPOINT /loan-engine