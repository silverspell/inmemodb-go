FROM golang:1.19-alpine as build-env
# RUN apk --no-cache add build-base git
ADD . ./src
RUN cd src && go build -o martiapp -ldflags "-s -w" .



FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/martiapp /app/
EXPOSE 9001
ENTRYPOINT ["./martiapp"]