FROM alpine
RUN apk add --no-cache ca-certificates
ADD beeru /usr/bin/beeru
ENTRYPOINT ["/usr/bin/beeru"]
EXPOSE 8000
