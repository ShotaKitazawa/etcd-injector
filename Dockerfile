# build stage
FROM golang:1.14 as builder
## init setting
WORKDIR /workdir
ENV GO111MODULE="on"
ARG REPO_NAME
ARG IMAGE_TAG
## download packages
COPY go.mod go.sum ./
RUN go mod download
## build
COPY . ./
RUN GOOS=linux go build -ldflags "-X main.appName=${REPO_NAME} -X main.appVersion=${IMAGE_TAG}"

# run stage
FROM gcr.io/distroless/base
## copy binary
COPY --from=builder /workdir/etcd-injector .
## Run
ENTRYPOINT ["./etcd-injector"]

