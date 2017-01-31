package main

import (
	"strconv"
	"os/exec"
	"bufio"
	"bytes"
)

// runCmd executes passed command and pipes output through handler.
func runCmd(command string, fn ...handler) error {
	cmd := exec.Command("bash", "-c", command)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logError.Print("Error creating StdoutPipe for Cmd: ", stderr.String())
		return err
	}

	for _, f := range fn {
		go f(bufio.NewScanner(cmdReader))
	}

	err = cmd.Start()
	if err != nil {
		logError.Print("Error starting Cmd: ", stderr.String())
		return err
	}

	err = cmd.Wait()
	if err != nil {
		logError.Print("Error waiting for Cmd: ", stderr.String())
		return err
	}

	return nil
}

// cleanupServices removes hanging services in working namespace.
// It's triggered at the end of each failing Kuberang iteration.
func cleanupServices() error {
	logInfo.Println("Cleanup services")

	err := runCmd("kubectl delete svc --selector run=kuberang-nginx --ignore-not-found --namespace " + c.String("namespace"))
	if err != nil {
		return err
	}
	return nil
}

// cleanupDeployments removes hanging deployment in working namespace.
// It's triggered at the end of each failing Kuberang iteration.
func cleanupDeployments() error {
	logInfo.Println("Cleanup deployments")

	err := runCmd("kubectl delete deployment kuberang-busybox kuberang-nginx --ignore-not-found --namespace " + c.String("namespace"))
	if err != nil {
		return err
	}
	return nil
}

// pushMetric sends metric with a given value to Sysdig via Statsd (udp)
func pushMetric(metric string, value int) error {
	if c.Bool("push-metrics") {
		logInfo.Println("Sending metric to sysdig: " + metric + ":" + strconv.Itoa(value))

		// Sending metrics to Sysdig via Statsd
		err := runCmd("echo '" + metric + ":" + strconv.Itoa(value) + "|g' > " + sStatsd)
		if err != nil {
			return err
		}
	}
	return nil
}

// metricValue returns a metric value depending of test result
func metricValue() int {
	value := 0
	if rErrorTest.MatchString(logLine) {
		value = 1
	}
	return value
}

// resetCounters clears out counters.
// Triggered at the end of each Kuberang iteration.
func resetCounters() {
	totalFromWorker = 0
	totalFromNode = 0
	failingFromWorker = 0
	failingFromNode = 0
}
