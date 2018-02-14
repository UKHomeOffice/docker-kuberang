FROM alpine:3.4

RUN apk upgrade --no-cache && apk add --no-cache bash curl coreutils
RUN adduser -h /kuberang -D kuberang

ENV KUBECTL_VERSION 1.8.2
ENV KUBERANG_VERSION v1.2.2

RUN curl -s https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    -o /usr/bin/kubectl && chmod +x /usr/bin/kubectl

RUN curl -s https://github.com/apprenda/kuberang/releases/download/v${KUBERANG_VERSION}/kuberang-linux-amd64 \
    -o /usr/local/bin/kuberang && chmod +x /usr/local/bin/kuberang

COPY ./bin/smoketest /usr/local/bin/smoketest

USER kuberang
ENTRYPOINT ["/usr/local/bin/smoketest"]
