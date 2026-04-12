# build stage
FROM golang:1.26.2-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN apk update && \
    apk upgrade && \
    apk add --no-cache git && \
    apk add --no-cache make
RUN make build

# final stage
FROM scratch
COPY --from=builder /app/aegis .
EXPOSE 8080
CMD ["sh", "-c", "./aegis --github-access-token ${GITHUB_ACCESS_TOKEN}"]
