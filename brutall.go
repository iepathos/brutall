package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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

func getLastLog(containerName string) string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("docker logs %s", containerName))
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	cleanOut := strings.TrimSpace(string(out))
	return cleanOut
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

func runGobuster(domain string) bool {
	// $HOME/work/bin/gobuster -m dns -u $TARGETS -w $finalLOC -t $gobusterthreads -fw > /tmp/gobuster.txt
	// ./gobuster.sh --m dns -u $domain -w /words/allwords.txt -t 100 -fw > /words/gobuster.txt
	baseDir := getBaseDir()
	gobuster := filepath.Join(baseDir, "services", "gobuster", "gobuster.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s --m dns -t 100 -w /words/allwords.txt -u %s -q -fw", gobuster, domain))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return true
}

func handleGobusterOutput() {
	// get last gobuster image log
	lastLog := getLastLog("gobuster")
	// cleanup log
	// cat /tmp/gobuster.txt | grep Found | sed 's/Found: //'
	domains := []string{}
	line := ""
	scanner := bufio.NewScanner(strings.NewReader(lastLog))
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "Found") {
			lsplit := strings.Fields(line)
			domain := lsplit[len(lsplit)-1]
			domains = append(domains, domain)
		}
	}
	cleanedLogs := strings.Join(domains, "\n")

	// write to /tmp/gobuster.txt
	tmpGobuster := "/tmp/gobuster.txt"
	f, err := os.Create(tmpGobuster)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = f.WriteString(cleanedLogs)
	if err != nil {
		log.Fatal(err.Error())
	}

	// add gobuster.txt to docker words volume
	baseDir := getBaseDir()
	addFileToWordsVolume := filepath.Join(baseDir, "volumes", "add_file_to_words_volume.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", addFileToWordsVolume, tmpGobuster))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func runSublist3r(domain string) bool {
	// ./sublist3r.sh -d $domain -t $sublist3rthreads -v -o $sublist3rfile
	baseDir := getBaseDir()
	sublist3r := filepath.Join(baseDir, "services", "sublist3r", "sublist3r.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -d %s -t 50 -v -o /words/sublist3r.txt", sublist3r, domain))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return true
}

func runEnumall(domain string) bool {
	baseDir := getBaseDir()
	enumall := filepath.Join(baseDir, "services", "enumall", "enumall.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s -w /words/gobuster.txt", enumall, domain))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return true
}

func runAltdns(domain string) bool {
	// /usr/bin/python /opt/subscan/altdns/altdns.py -i $finaloutputbeforealtdns -o data_output -w words.txt -r -e -d $altdnsserver -s $altdnsoutput -t $altdnsthreads

	baseDir := getBaseDir()
	altdns := filepath.Join(baseDir, "services", "altdns", "altdns.sh")
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -d %s -i /words/finaloutputbeforealtdns.txt -o data_output -w words.txt -r -e -d 8.8.8.8 -s /words/altdnsoutput.txt -t 100", altdns, domain))
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return true
}

func runServices(domain string) {
	baseDir := getBaseDir()
	os.Chdir(baseDir)
	serviceStatuses := make(chan bool)
	var wg sync.WaitGroup
	// run these services in parallel
	wg.Add(2)
	go func() {
		defer wg.Done()
		serviceStatuses <- runGobuster(domain)
		// grab log from gobuster container that just ran
		// cleanup the log, output to /tmp/gobuster.txt and add to words volume
		handleGobusterOutput()

		// run enumall, pass the gobuster output as the wordslist
		runEnumall(domain)
	}()
	go func() {
		defer wg.Done()
		serviceStatuses <- runSublist3r(domain)
	}()
	wg.Wait()

	// run altdns after the other services have completed
	runAltdns(domain)
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
