package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getContainerName(imageName string) string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("docker ps -a | grep %s | head -n 1", imageName))
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	cleanOut := strings.TrimSpace(string(out))

	splitOut := strings.Fields(cleanOut)
	return splitOut[len(splitOut)-1]
}

func main() {
	// get last running container name from that image
	if len(os.Args) < 2 {
		log.Error("Usage: ./lastlog gobuster")
		os.Exit(1)
	}

	services := []string{
		"gobuster",
		"sublist3r",
		"altdns",
	}

	if !stringInSlice(os.Args[1], services) {
		log.Error("Valid services are gobuster, sublist3r, and altdns")
		os.Exit(1)
	}

	containerName := getContainerName(os.Args[1])
	fmt.Printf(containerName)
}
