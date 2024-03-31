# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22 AS build
WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/status_checker ./main.go

# Final stage
FROM scratch
COPY --from=build /bin/status_checker /bin/status_checker
CMD ["/bin/status_checker"]
