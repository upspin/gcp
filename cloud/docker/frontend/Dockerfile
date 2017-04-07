FROM buildpack-deps:jessie-curl
COPY bin/frontend-gcp /frontend-gcp
COPY src/upspin.io/doc /doc
ENTRYPOINT ["/frontend-gcp", "-docpath=/doc", "-http=:80", "-https=:443"]
EXPOSE 80 443
