package main

import (
	// "bytes"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	// "strings"
)

func getBaseDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func gatherTools() {
	// ./services/clone_tools.sh
	cmd := "./services/clone_tools.sh"
	if err := exec.Command(cmd).Run(); err != nil {
		log.WithFields(log.Fields{
			"docker": "localhost",
			"error":  err.Error(),
		}).Error("An error occurred trying to gather the necessary tools")
		os.Exit(1)
	}
}

func buildService(servicePath string) {
	// ./services/servicePath/build.sh
	os.Chdir(servicePath)
	cmd := "./build.sh"
	if err := exec.Command(cmd).Run(); err != nil {
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
	log.Info("Building services")
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
	baseDir := getBaseDir()
	gobuster := filepath.Join(baseDir, "services", "gobuster", "gobuster.sh")

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s --m dns -t 100 -w /words/allwords.txt -u %s -fw", gobuster, domain))

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func runSublist3r(domain string) {
	// ./sublist3r.sh -d $domain -t $sublist3rthreads -v -o $sublist3rfile
	baseDir := getBaseDir()
	sublist3r := filepath.Join(baseDir, "services", "sublist3r", "sublist3r.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -d %s -t 100 -v -o /words/sublist3r.txt", sublist3r, domain))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func runServices(domain string) {
	baseDir := getBaseDir()
	os.Chdir(baseDir)
	runGobuster(domain)
	runSublist3r(domain)
}

func main() {
	if len(os.Args) < 2 {
		log.Error("Usage: ./brutall domain.com --build")
		os.Exit(1)
	} else {
		log.Info("Starting one brave binary!")
	}

	gatherTools()
	if stringInSlice("--build", os.Args) {
		build()
	}
	runServices(os.Args[1])
}
