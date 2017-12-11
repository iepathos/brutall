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
	log.WithFields(log.Fields{}).Info("Gathering tools from github")
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
			"error": err.Error(),
			"path":  servicePath,
		}).Error("An error occurred trying to build")
		os.Exit(1)
	}
}

func buildVolumes() {
	// go into volumes and build.sh
	log.WithFields(log.Fields{}).Info("Building volumes")
	baseDir := getBaseDir()
	volumesDir := filepath.Join(baseDir, "volumes")
	buildService(volumesDir)
}

func buildServices() {
	// go into services repos and build.sh
	log.WithFields(log.Fields{}).Info("Building services")
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

func runGobuster(domain string) {
	// $HOME/work/bin/gobuster -m dns -u $TARGETS -w $finalLOC -t $gobusterthreads -fw > /tmp/gobuster.txt
	// ./gobuster.sh --m dns -u $domain -w /words/allwords.txt -t 100 -fw > /words/gobuster.txt
	cmd := "./services/gobuster/gobuster.sh"
	args := []string{"--m", "dns", "-u", domain, "-t", "100", "-w", "/words/allwords.txt", "-fw"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"cmd":   cmd,
			"args":  args,
		}).Error("An error occurred trying to execute")
		os.Exit(1)
	}
}

func runSublist3r(domain string) {
	// ./sublist3r.sh -d $domain -t $sublist3rthreads -v -o $sublist3rfile
	cmd := "./services/sublist3r/sublist3r.sh"
	args := []string{"-d", domain, "-t", "100", "-v", "-o", "/words/sublist3r.txt"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"cmd":   cmd,
			"args":  args,
		}).Error("An error occurred trying to execute")
		os.Exit(1)
	}
}

func runServices(domain string) {
	baseDir := getBaseDir()
	os.Chdir(baseDir)
	runGobuster(domain)
	runSublist3r(domain)
}

func main() {
	if len(os.Args) < 2 {
		log.Error("Usage: ./brutall domain.com")
		os.Exit(1)
	} else {
		log.Info("Starting one brave binary!")
	}
	gatherTools()
	build()
	runServices(os.Args[1])
}
