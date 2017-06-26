FROM golang:1.8
ENV GOPATH /gopath
ADD makerelease.sh /makerelease.sh
ADD config /config
ADD secrets /secrets
ENTRYPOINT ["/makerelease.sh"]
