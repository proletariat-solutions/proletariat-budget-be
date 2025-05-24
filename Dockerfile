FROM golang:1.24-alpine as build

COPY build-context /root/

WORKDIR /app
RUN apk update --no-cache && apk add --no-cache git upx ca-certificates openssh-client
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

ENV GOPRIVATE=github.com/ghorkov32

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY go.mod go.sum ./

# Allow the build container to access SSH keys via SSH agents: https://docs.docker.com/engine/reference/builder/#run---mounttypessh
RUN --mount=type=ssh go mod download
COPY . ./

# Install oapi-codegen and generate go-chi boilerplate
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.0 \
    && mkdir -p openapi && oapi-codegen -package openapi -config oapi-gen.cfg.yaml api-v1.yaml

# Build and run the app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /api
RUN upx -9 /api

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata && apk upgrade --no-cache
COPY --from=build /api /api

ENTRYPOINT ["/api"]
