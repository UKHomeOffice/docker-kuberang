# Kuberang - smoke test for Kubernetes

This repo contains [Kuberang](https://github.com/apprenda/kuberang) wrapper program executing kubernetes smoke tests in a loop with given internal, processes test results and pushes metrics to Sysdig. 

## Getting Started

You need to make sure that your kubernetes cluster supports service accounts and 
that `kuberang` service account has sufficient permissions to create objects in `smoke-test` namespace as well as listing cluster nodes. For details see [README](https://gitlab.digital.homeoffice.gov.uk/Devops/kube-kuberang).

### Build

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=<version>" -o bin/smoketest
```

### Configuration

* `KUBE_NAMESPACE` - defaults to `smoke-test`.
* `INTERVAL` - defaults to `5m` (5 minutes).
* `DEBUG` - when set to `true` it'll log details of individual checks.
* `PUSH_METRICS` - when set to `true` it'll push gauge metrics to Sysdig statsd.

### Deployment

To deploy follow instructions in [kube-kuberang](https://gitlab.digital.homeoffice.gov.uk/Devops/kube-kuberang) repo.

## Contributing

Pull requests welcome. Please check issues and existing PRs before submitting a patch.

## Author

Marcin Ciszak [marcinc](https://github.com/marcinc)

## License

[MIT](LICENSE)
