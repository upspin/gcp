FROM buildpack-deps:jessie-curl
COPY bin/upspinserver-gcp /upspinserver
ENTRYPOINT ["/upspinserver", "-serverconfig=/upspin/server", "-letscache=/upspin/letscache", "-https=:8443"]
EXPOSE 8443
