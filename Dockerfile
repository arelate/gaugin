FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -o gaugin \
    main.go

FROM alpine
COPY --from=build /go/src/app/gaugin /usr/bin/gaugin

EXPOSE 1848
#vangogh artifacts: checksums, images, metadata, recycle_bin, videos
VOLUME /var/lib/vangogh

ENTRYPOINT ["/usr/bin/gaugin"]
CMD ["serve","-port", "1848", "-stderr"]