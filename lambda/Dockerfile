FROM golang:1.18-alpine3.16 as builder

# Builder dir
RUN mkdir build
WORKDIR build

# Add dependencies
ADD go.mod go.sum ./
RUN go mod download

# Build
ADD . .
RUN go build -o /main

# Runner image
FROM alpine:3.16.2

ENV NOTION_API_KEY ""
ENV NOTION_DATABASE_ID ""
ENV GITHUB_USER = ""

COPY --from=builder /main /main
ENTRYPOINT [ "/main" ]
