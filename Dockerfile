FROM golang as go
WORKDIR /go/src/app
RUN curl https://glide.sh/get | sh
COPY main.go ./
COPY glide* ./
RUN glide install
RUN  CGO_ENABLED=0 go build -a main.go

FROM scratch
COPY --from=go /go/src/app/main /
CMD ["/main"]
