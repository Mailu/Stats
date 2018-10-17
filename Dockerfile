FROM golang:alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /go/src/app
COPY . .

RUN go build -v ./...
RUN go install -v ./...


FROM scratch

COPY --from=builder /go/bin/app /app

CMD ["/app"]

# ENV LOGPATH=/output.log # where to store the log
# ENV HOST=0.0.0.0 # which IP or Hostname to listen on (supports IPv4 and IPv6)
# ENV PORT=53 # which port to listen on
# ENV DOMAIN=stats.mailu.io. # which domain we serve. MUST end with a dot!
# ENV VALUECOUNT=2 # how many subdomains (=values) we want
