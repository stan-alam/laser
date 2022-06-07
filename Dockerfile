FROM golang:1.18.2 as builder
WORKDIR /srv
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . .
RUN cd cmd/beam &&
	go clean &&
	go get &&
	go test -v -race &&
	go test -c &&
	go build -a -installsuffix cgo
RUN cd cmd/integration &&
	go test -c -o integration

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /srv/cmd/beam/beam /go/bin/beam
COPY --from=builder /srv/cmd/beam/beam.test /go/bin/beam.test
COPY --from=builder /srv/cmd/beam/integration /go/bin/integration
CMD ["/go/bin/beam"]
