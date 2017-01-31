package main

import (
	"bufio"
	"strings"
	"regexp"
)

var (
	// Kuberang output matchers
	rErrorTest = regexp.MustCompile(`\[ERROR\]`)
	rDeployment = regexp.MustCompile(`Both deployments completed successfully within timeout`)
	rWorkerNginxNetwork = regexp.MustCompile(`Accessed Nginx pod at .* from BusyBox`)
	rWorkerNginxServiceIp = regexp.MustCompile(`Accessed Nginx service at .* from BusyBox`)
	rWorkerNginxDNS = regexp.MustCompile(`Accessed Nginx service via DNS .* from BusyBox`)
	rWorkerExternalNetwork = regexp.MustCompile(`Accessed Google.com from BusyBox`)
	rNodeNginxNetwork = regexp.MustCompile(`Accessed Nginx pod at .* from this node`)
	rNodeNginxDNS = regexp.MustCompile(`Accessed Nginx service via DNS .* from this node`)
	rNodeExternalNetwork = regexp.MustCompile(`Accessed Google.com from this node`)

	// Processed log line
	logLine string 
)

type handler func(*bufio.Scanner)

// kuberangOutputHandler processes results of Kuberang command
func kuberangOutputHandler(s *bufio.Scanner) {
	for s.Scan() {
		logLine = s.Text()

		if c.Bool("debug") {
			logDebug.Println(logLine)
		}

		switch {
		case rDeployment.MatchString(logLine): deploymentCheck()
		case rWorkerNginxNetwork.MatchString(logLine): countWorkerNetworkCheck()
		case rWorkerNginxServiceIp.MatchString(logLine): workerNginxViaServiceIpCheck()
		case rWorkerNginxDNS.MatchString(logLine): workerNginxViaDNSCheck()
		case rWorkerExternalNetwork.MatchString(logLine): workerExternalNetworkCheck()
		case rNodeNginxNetwork.MatchString(logLine): countNodeNetworkCheck()
		case rNodeNginxDNS.MatchString(logLine): nodeNginxViaDNSCheck()
		case rNodeExternalNetwork.MatchString(logLine): nodeExternalNetworkCheck()
		}
	}
}

func suspectedNodesHandler(s *bufio.Scanner) {
	arr := []string{}
	for s.Scan() {
		arr = append([]string{s.Text()}, arr...)
	}
	if len(arr) > 0 {
		logError.Print("Suspected nodes: ")
		logError.Println(strings.Join(arr, ", "))
	}
}
