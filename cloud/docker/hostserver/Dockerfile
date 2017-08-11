FROM buildpack-deps:jessie-curl
COPY config /upspin/config
COPY public.upspinkey /upspin/public.upspinkey
COPY secret.upspinkey /upspin/secret.upspinkey
COPY bin/hostserver-gcp /hostserver-gcp
ENTRYPOINT ["/hostserver-gcp", "-letscache=", "-https=:443", "-config=/upspin/config"]
EXPOSE 443
