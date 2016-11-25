FROM alpine:3.4

RUN apk upgrade --no-cache && apk add --no-cache bash curl coreutils
RUN adduser -h /kuberang -D kuberang

ENV KUBECTL_VERSION 1.4.3

RUN curl -s https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    -o /usr/bin/kubectl && chmod +x /usr/bin/kubectl

# TODO: Pull Kuberang release from github once it's published. 
#       For the time being use precompiled binary.
COPY ./bin/kuberang /usr/local/bin/kuberang
COPY ./bin/smoketest.sh /usr/local/bin/smoketest.sh

USER kuberang
ENTRYPOINT ["/usr/local/bin/smoketest.sh"]
