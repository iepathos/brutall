package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{"docker": "host"}).Info("One brave binary!")
}
