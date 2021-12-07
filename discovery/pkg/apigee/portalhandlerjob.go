package apigee

import (
	"fmt"

	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agent-sdk/pkg/jobs"
	"github.com/Axway/agent-sdk/pkg/util/log"
)

//portalHandler - job that waits for
type portalHandler struct {
	jobs.Job
	apigeeClient      *GatewayClient
	newPortalChan     chan string
	removedPortalChan chan string
	apiChan           chan *apiDocData
	stopChan          chan interface{}
	isRunning         bool
	runningChan       chan bool
	jobsMap           map[string]*pollPortalAPIsJob // keep map of all jobs created for the portals
}

func newPortalHandlerJob(apigeeClient *GatewayClient, newPortalChan, removedPortalChan chan string, apiChan chan *apiDocData) *portalHandler {
	job := &portalHandler{
		apigeeClient:      apigeeClient,
		newPortalChan:     newPortalChan,
		removedPortalChan: removedPortalChan,
		stopChan:          make(chan interface{}),
		isRunning:         false,
		apiChan:           apiChan,
		runningChan:       make(chan bool),
		jobsMap:           make(map[string]*pollPortalAPIsJob),
	}
	go job.statusUpdate()
	return job
}

func (j *portalHandler) Ready() bool {
	if j.apigeeClient.accessToken == "" {
		return false
	}
	return true
}

func (j *portalHandler) Status() error {
	if !j.isRunning {
		return fmt.Errorf("portal handler not running")
	}
	return nil
}

func (j *portalHandler) statusUpdate() {
	for {
		select {
		case update := <-j.runningChan:
			j.isRunning = update
		}
	}
}

func (j *portalHandler) started() {
	j.runningChan <- true
}

func (j *portalHandler) stopped() {
	j.runningChan <- false
}

func (j *portalHandler) Execute() error {
	j.started()
	defer j.stopped()
	for {
		select {
		case newPortal, ok := <-j.newPortalChan:
			if !ok {
				err := fmt.Errorf("New portal channel was closed")
				return err
			}
			j.handleNewPortal(newPortal)
		case removedPortal, ok := <-j.removedPortalChan:
			if !ok {
				err := fmt.Errorf("New portal channel was closed")
				return err
			}
			j.handleRemovedPortal(removedPortal)
		case <-j.stopChan:
			log.Info("Stopping the portal handler")
			return nil
		}
	}
}

func (j *portalHandler) handleNewPortal(newPortal string) {
	log.Tracef("Handling new portal %s", newPortal)

	portalName, err := j.getPortalNameByID(newPortal)
	if err != nil {
		log.Errorf("could not find portal name with portal ID %s", newPortal)
		return
	}

	// register a new job to poll for apis in this portal
	portalAPIsJob := newPollPortalAPIsJob(j.apigeeClient, newPortal, portalName, j.apiChan)
	err = portalAPIsJob.Register()
	if err != nil {
		log.Errorf("error hit starting job for portal ID %s", newPortal)
		return
	}
	j.jobsMap[newPortal] = portalAPIsJob
}

func (j *portalHandler) handleRemovedPortal(removedPortal string) {
	log.Tracef("Handling removed portal %s", removedPortal)

	// unregister the job to polling for apis in this portal
	if job, ok := j.jobsMap[removedPortal]; ok {
		job.PortalRemoved()
	}
}

func (j *portalHandler) getPortalNameByID(newPortal string) (string, error) {
	portalMapInterface, err := cache.GetCache().Get(portalMapCacheKey)
	if err != nil {
		log.Error("error hit getting the portal map from the cache")
		return "", err
	}
	portalMap := portalMapInterface.(map[string]string)
	if portalName, ok := portalMap[newPortal]; ok {
		return portalName, nil
	}
	return "", fmt.Errorf("portal id %s not in cache", newPortal)
}
