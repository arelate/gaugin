FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o gg \
    -ldflags="-s -w -X 'github.com/arelate/gaugin/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

FROM alpine:latest
COPY --from=build /go/src/app/gg /usr/bin/gg
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1848

#images
VOLUME /var/lib/vangogh/images
#items
VOLUME /var/lib/vangogh/items
#downloads
VOLUME /var/lib/vangogh/downloads

ENTRYPOINT ["/usr/bin/gg"]
CMD ["serve","-port", "1848", "-stderr"]