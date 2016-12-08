# Kuberang - smoke test for Kubernetes

Kuberang in a docker image with all necessary scripts to run in
kubernetes cluster.

Entrypoint to the container is a bash script which executes Kuberang in a loop with given interval and pushes metrics to Sysdig.

## Getting Started

You need to make sure that your kubernetes cluster supports service accounts and 
that `kuberang` service account has sufficient permissions to create objects in `smoke-test` namespace as well as listing cluster nodes.

### Build

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=<version>" -o bin/smoketest
```

### Configuration

* `KUBE_NAMESPACE` - defaults to `smoke-test`.
* `INTERVAL` - defaults to `300` (seconds).
* `DEBUG` - when set to `true` it'll log details of individual checks.
* `PUSH_METRICS` - when set to `true` it'll push gauge metrics to Sysdig statsd.

### Deployment

#### Create smoke-test namespace

```
kubectl create -f kube/kuberang-namespace.yaml
``` 

#### Create service account

```
kubectl create -f kube/kuberang-serviceaccount.yaml
```

Ensure that `kuberang` service account has necessary permissions to `smoke-test` namespace.

```
{"apiVersion":"abac.authorization.kubernetes.io/v1beta1","kind":"Policy","spec":{"user":"system:serviceaccount:smoke-test:kuberang","namespace":"smoke-test","resource":"*","apiGroup":"*"}}

{"apiVersion":"abac.authorization.kubernetes.io/v1beta1","kind":"Policy","spec":{"user":"system:serviceaccount:smoke-test:kuberang","readonly":true,"resource":"nodes"}}
```

#### Create Kuberang deployment

Ensure (docker-registry) `registrykey` secret is set, otherwise kuberang won't be able to pull an image from Artifactory!

```
kubectl create -f kube/kuberang-deployment.yaml
```

## Contributing

Pull requests welcome. Please check issues and existing PRs before submitting a patch.

## Author

Marcin Ciszak [marcinc](https://github.com/marcinc)

## License

[MIT](LICENSE)
