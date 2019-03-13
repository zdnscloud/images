FROM alpine

ENV KUBE_LATEST_VERSION="v1.13.4"
COPY ./run-sh /bin/
RUN apk add --update -t deps curl \
 && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && chmod +x /usr/local/bin/kubectl \
 && apk del --purge deps curl \
 && apk add bash \
 && rm /var/cache/apk/*

USER 1001
ENTRYPOINT ["/bin/run-sh"]
