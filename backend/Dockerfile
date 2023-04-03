FROM golang:1.20.2-alpine as builder
RUN apk add --no-cache git make curl build-base
ENV GOOS=linux

WORKDIR /src
COPY backend/go.mod /src/
COPY backend/go.sum /src/
RUN go mod download
COPY backend /src/
RUN make build

FROM alpine:3.15
RUN apk add --no-cache ca-certificates tzdata
RUN export PATH=$PATH:/app
WORKDIR /app
COPY --from=builder /src/bin/tracing-demo-backend /app/backend
CMD ["/app/backend"]
