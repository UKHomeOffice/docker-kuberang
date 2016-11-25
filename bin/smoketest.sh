#!/bin/bash

set -o errexit
set -o pipefail

: ${KUBE_NAMESPACE:-smoke-test}
: ${INTERVAL:-300}
: ${DEBUG:-false}
: ${PUSH_METRICS:-false}

SYSDIG_METRIC_NAME=k8s_smoke_test
SYSDIG_STATSD=/dev/udp/127.0.0.1/8125

push_metric() {
  if [[ $1 == 0 ]]; then
    info 'OK'
  else
    error "$2"
  fi
  if ${PUSH_METRICS}; then
    info "Pushing metric to Sysdig: $1"
    echo "$SYSDIG_METRIC_NAME:$1|g" > $SYSDIG_STATSD
  fi      
}

log() {
  echo "[${1}] $(date): ${2}"
}

info() {
  log 'INFO' "$1"
}

debug() {
  log 'DEBUG' "$1"
}

error() {
  log 'ERROR' "$1"
}

cleanup_services() {
  info "clean up services in ${KUBE_NAMESPACE} namespace..."
  for service in $(kubectl get svc --no-headers --namespace ${KUBE_NAMESPACE} | awk '{print $1}'); do
    kubectl delete svc $service --namespace ${KUBE_NAMESPACE}
  done
}

cleanup_deployments() {
  info "clean up deployments in ${KUBE_NAMESPACE} namespace..."
  kubectl delete deployment kuberang-busybox kuberang-nginx --ignore-not-found --namespace ${KUBE_NAMESPACE} 
}

run() {
  while true; do
    ERROR=0
    while read line; do
      if ${DEBUG}; then debug "$line"; fi
      
      if [[ $line == *'[ERROR]'* ]]; then
        push_metric 1 "$line"
        ERROR=1
        break
      fi  
    done < <(/usr/local/bin/kuberang --namespace ${KUBE_NAMESPACE})

    if [[ $ERROR == 0 ]]; then
      push_metric 0 'OK'
    else
      cleanup_services
      cleanup_deployments
    fi

    sleep $INTERVAL
  done
}

run
