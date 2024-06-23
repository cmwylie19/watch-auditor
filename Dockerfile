FROM golang:1.22.4-alpine3.20 AS builder

LABEL description="watch-auditor" \
      maintainer="Casey Wylie casewylie@gmail.com"

WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN go build 

FROM scratch
WORKDIR /app
COPY --from=builder /app/watch-auditor ./

# # Set the entrypoint to run the application
ENTRYPOINT ["./watch-auditor"]
