FROM golang:1.20-alpine AS build

WORKDIR /app
COPY . .
RUN go build -o tasker

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=build /app/tasker .
RUN chown -R appuser:appgroup /app

EXPOSE 8080

USER appuser
CMD ["./tasker"]