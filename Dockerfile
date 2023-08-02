FROM golang:1.20.7-alpine AS build
RUN apk add --no-cache gcc libc-dev

WORKDIR /go/src/app
COPY . .
RUN go test ./...
ARG version=dev
RUN go build -ldflags "-X main.Version=$version" -o /bin/ecr-image-sync

FROM alpine:3.18.2

COPY --from=build /bin/ecr-image-sync /usr/local/bin/ecr-image-sync
ENTRYPOINT ["ecr-image-sync"]
