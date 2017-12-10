package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
)

func getBaseDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func gatherTools() {
	// ./services/clone_tools.sh
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Gathering tools from github")
	cmd := "./services/clone_tools.sh"
	args := []string{}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker": "localhost",
			"error":  err.Error(),
		}).Error("An error occurred trying to gather the necessary tools")
		os.Exit(1)
	}
}

func buildService(servicePath string) {
	os.Chdir(servicePath)
	cmd := "./build.sh"
	args := []string{}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker":  "localhost",
			"error":   err.Error(),
			"service": servicePath,
		}).Error("An error occurred trying to build")
		os.Exit(1)
	}
}

func buildVolumes() {
	// go into volumes and build.sh
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Building volumes")
	baseDir := getBaseDir()
	volumesDir := filepath.Join(baseDir, "volumes")
	buildService(volumesDir)
}

func buildServices() {
	// go into services repos and build.sh
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Building services")
	baseDir := getBaseDir()
	servicesDir := filepath.Join(baseDir, "services")
	services := []string{
		"gobuster",
		"sublist3r",
		"altdns",
	}
	for _, service := range services {
		servicePath := filepath.Join(servicesDir, service)
		buildService(servicePath)
	}
}

func build() {
	buildVolumes()
	buildServices()
}

func runService(service string, domain string) {
	// results <- runService("sublist3r")
	threads := 100
	cmd := "./services/" + service + "/" + service + ".sh"
	args := []string{"-d", domain, "-t", string(threads), "-v"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker":  "localhost",
			"error":   err.Error(),
			"service": service,
		}).Error("An error occurred trying to execute")
		os.Exit(1)
	}
}

func main() {
	log.WithFields(log.Fields{"docker": "localhost"}).Info("Starting one brave binary!")

	gatherTools()
	build()
}
