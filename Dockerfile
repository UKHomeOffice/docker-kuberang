FROM quay.io/ukhomeofficedigital/alpine-glibc:3.6

RUN apk upgrade --no-cache && apk add --no-cache bash curl coreutils
RUN adduser -h /kuberang -D -u 1000 kuberang

ENV KUBECTL_VERSION 1.18.17
ENV KUBERANG_VERSION 1.4.0-rc1

RUN curl -s https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    -o /usr/bin/kubectl && chmod +x /usr/bin/kubectl

RUN curl -L -s https://github.com/UKHomeOffice/kuberang/releases/download/v${KUBERANG_VERSION}/kuberang-linux-amd64 \
    -o /usr/local/bin/kuberang && chmod +x /usr/local/bin/kuberang

COPY ./bin/smoketest /usr/local/bin/smoketest

USER 1000
ENTRYPOINT ["/usr/local/bin/smoketest"]
