package runner

import (
	"log"
	"net"
	"time"

	"github.com/alinowrouzii/service-health-check/models"
	"gorm.io/gorm"
)

type Runner struct {
	DB               *gorm.DB
	RunnerIntervalMs int
}

func connectionStatus(host string, port string) bool {
	timeout := 500 * time.Millisecond
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func (runner *Runner) registerJob(link models.Link) {
	log.Println("register job for link", link.ID, link.UserID)
	ticker := time.NewTicker(time.Duration(runner.RunnerIntervalMs) * time.Millisecond)
	done := make(chan bool)
	errorCounter := 0
	linkAddress := link.URL
	linkPort := "80"
	linkErrorThreshold := link.ErrorThreshold
	linkId := link.ID
	dbConn := runner.DB
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("starting at", t)
			connectionEstablished := connectionStatus(linkAddress, linkPort)
			log.Println("connection", connectionEstablished, linkAddress, linkPort)
			if !connectionEstablished {
				errorCounter++

				if errorCounter > int(linkErrorThreshold) {
					warning := models.Warning{
						LinkID: linkId,
					}
					err := warning.CreateWarning(dbConn)
					if err != nil {
						log.Fatal("killing goroutine cause of db error", err)
					}
					log.Printf("warning created successfully")
					errorCounter = 0
				}
			}
		}
	}
}

func (runner *Runner) Run() {
	var links []models.Link
	result := runner.DB.Find(&links)
	if result.Error != nil {
		log.Fatal("error in job runner", result.Error)
	}

	log.Println("here is fetched links", links)
	go func() {
		for _, link := range links {
			go runner.registerJob(link)
		}
	}()
}
