package apigee

import (
	"fmt"
	"strconv"

	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agent-sdk/pkg/jobs"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/Axway/agents-apigee/client/pkg/apigee"
)

// job that will poll for any new apis in the portal on APIGEE Edge
type pollPortalAPIsJob struct {
	jobs.Job
	apigeeClient   *apigee.ApigeeClient
	portalID       string
	portalName     string
	portalAPIsMap  map[string]string
	processAPIChan chan *apigee.APIDocData
	removedAPIChan chan string
	jobID          string
}

func newPollPortalAPIsJob(apigeeClient *apigee.ApigeeClient, portalID, portalName string, processAPIChan chan *apigee.APIDocData, removedAPIChan chan string) *pollPortalAPIsJob {
	return &pollPortalAPIsJob{
		apigeeClient:   apigeeClient,
		portalID:       portalID,
		portalName:     portalName,
		portalAPIsMap:  make(map[string]string),
		processAPIChan: processAPIChan,
		removedAPIChan: removedAPIChan,
	}
}

func (j *pollPortalAPIsJob) Register() error {
	jobID, err := jobs.RegisterIntervalJobWithName(j, j.apigeeClient.GetConfig().GetPollInterval(), fmt.Sprintf("%s Portal Poller", j.portalName))
	if err != nil {
		return err
	}
	j.jobID = jobID
	return nil
}

func (j *pollPortalAPIsJob) Ready() bool {
	return j.apigeeClient.IsReady()
}

func (j *pollPortalAPIsJob) Status() error {
	return nil
}

func (j *pollPortalAPIsJob) Execute() error {
	log.Tracef("Executing %s Portal poller", j.portalName)
	allPortalAPIs := j.apigeeClient.GetPortalAPIs(j.portalID)
	log.Tracef("%s Portal APIs: %+v", j.portalName, allPortalAPIs)
	apisFound := make(map[string]string)
	for _, api := range allPortalAPIs {
		id := strconv.Itoa(api.ID)
		apisFound[id] = api.ProductName
		if _, ok := j.portalAPIsMap[id]; !ok {
			log.Debugf("Found new api product %s", api.ProductName)
			j.portalAPIsMap[id] = api.ProductName
		}
		changed, err := cache.GetCache().HasItemChanged(id, *api)
		if err != nil || changed {
			api.SetPortalTitle(j.portalName)
			// send to new api handler
			j.processAPIChan <- api
		}
	}

	// check if any api has been removed
	for id := range j.portalAPIsMap {
		if _, ok := apisFound[id]; !ok {
			j.removedAPIChan <- id
			defer func(id string) {
				delete(j.portalAPIsMap, id)
			}(id) // remove apis from the map
		}
	}
	log.Tracef("Finished %s Portal poller", j.portalName)
	return nil
}

func (j *pollPortalAPIsJob) PortalRemoved() {
	// Loop all apis in this portal and remove them
	for _, id := range j.portalAPIsMap {
		j.removedAPIChan <- id
	}

	// Unregister this portals job
	jobs.UnregisterJob(j.jobID)
}