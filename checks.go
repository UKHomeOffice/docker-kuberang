package main

import (
	"strconv"
)

var (
	// Network check counters
	totalFromWorker = 0
	totalFromNode = 0
	failingFromWorker = 0
	failingFromNode = 0
)

const (
	// Sysdig metric names
	metricPrefix = "smoke_test.gauge"

	sDeployment = metricPrefix + ".deployment.status"

	sWorkerNetworkErrorRate = metricPrefix + ".worker.network.error.percent"
	sWorkerNetworkServiceIp = metricPrefix + ".worker.network.service_ip.status"
	sWorkerNetworkDNS = metricPrefix + ".worker.network.dns.status"
	sWorkerNetworkExternal = metricPrefix + ".worker.network.external.status"

	sNodeNetworkErrorRate = metricPrefix + ".node.network.error.percent"
	sNodeNetworkDNS = metricPrefix + ".node.network.dns.status"
	sNodeNetworkExternal = metricPrefix + ".node.network.external.status"

	// Sysdig Statsd
	sStatsd = "/dev/udp/127.0.0.1/8125"

	// Allowed Network error threshold %
	errorThreshold = 20
)

// smoke test deployment status
func deploymentCheck() error {
	if rErrorTest.MatchString(logLine) {
		logError.Println("Smoke test deployments failed. One (or more) nodes might be having some issues or quota exceeded!")

		// Display list of suspeced nodes
		runCmd("kubectl get po --all-namespaces --show-all --no-headers -owide | grep -E '(Terminating|Pending)' | awk '{print $8}' | sort | uniq", suspectedNodesHandler)
	}
	return pushMetric(sDeployment, metricValue())
}

// worker -> nginx via service IP
func workerNginxViaServiceIpCheck() error {
	return pushMetric(sWorkerNetworkServiceIp, metricValue())
}

// worker -> nginx via DNS
func workerNginxViaDNSCheck() error {
	return pushMetric(sWorkerNetworkDNS, metricValue())
}

// worker -> external network
func workerExternalNetworkCheck() error {
	return pushMetric(sWorkerNetworkExternal, metricValue())
}

// kuberang node -> nginx via DNS 
func nodeNginxViaDNSCheck() error {
	return pushMetric(sNodeNetworkDNS, metricValue())
}

// kuberang node -> external network
func nodeExternalNetworkCheck() error {
	return pushMetric(sNodeNetworkExternal, metricValue())
}

// counts worker -> nginx (via pod IP) network tests 
func countWorkerNetworkCheck() {
	totalFromWorker++
	if rErrorTest.MatchString(logLine) {
		failingFromWorker++
	}
}

// counts node -> nginx (via pod IP) network tests
func countNodeNetworkCheck() {
	totalFromNode++
	if rErrorTest.MatchString(logLine) {
		failingFromNode++
	}
}

// calculates error rates for network checks
func networkErrorRateCheck() {
	workerErrorRate := 0
	nodeErrorRate := 0

	if totalFromWorker > 0 {
		workerErrorRate = int(float64(failingFromWorker) * 100 / float64(totalFromWorker))
	}
	if totalFromNode > 0 {
		nodeErrorRate = int(float64(failingFromNode) * 100 / float64(totalFromNode))
	}

	if c.Bool("debug") {
		logDebug.Println("Worker to Nginx network error rate: " + strconv.Itoa(workerErrorRate) + "%")
		logDebug.Println("Node to Nginx network error rate: " + strconv.Itoa(nodeErrorRate) + "%")
	}
	
	pushMetric(sWorkerNetworkErrorRate, workerErrorRate)
	pushMetric(sNodeNetworkErrorRate, nodeErrorRate)
	resetCounters()
}
