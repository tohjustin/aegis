FROM golang:1.10-alpine AS build

# Using the following WORKDIR in order to allow the packages in `/vendor` &
# `/pkg` to be located correctly during the build phase
WORKDIR /go/src/github.com/tohjustin/aegis
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM scratch
EXPOSE 8080
COPY --from=build /go/src/github.com/tohjustin/aegis/main /aegis-microservice
CMD ["./aegis-microservice"]