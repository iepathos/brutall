package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)


func gatherTools() {
	// ./services/clone_tools.sh
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Gathering tools from github")
	cmd := "./services/clone_tools.sh"
	args := []string{}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker": "localhost",
			"error": err.Error(),
		}).Error("An error occurred trying to gather the necessary tools")
		os.Exit(1)
	}
}

func runService(service string, domain string) {
	// results <- runService("sublist3r")
	threads := 100
	cmd := "./services/" + service + "/" + service + ".sh"
	args := []string{"-d", domain, "-t", string(threads), "-v"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker": "localhost",
			"error": err.Error(),
			"service": service,
		}).Error("An error occurred trying to execute")
		os.Exit(1)
	}
}


func main() {
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Starting one brave binary!")

	gatherTools()

	log.WithFields(log.Fields{"docker": "gobuster"}).Info("Starting gobuster")

	log.WithFields(log.Fields{"docker": "sublist3r"}).Info("Starting sublist3r")

	log.WithFields(log.Fields{"docker": "altdns"}).Info("Starting altdns")
}
