package apigee

import (
	"fmt"
	"strings"
	"time"

	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	management "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/management/v1alpha1"
	defs "github.com/Axway/agent-sdk/pkg/apic/definitions"
	"github.com/Axway/agent-sdk/pkg/apic/provisioning"
	prov "github.com/Axway/agent-sdk/pkg/apic/provisioning"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/Axway/agents-apigee/client/pkg/apigee"
	"github.com/Axway/agents-apigee/client/pkg/apigee/models"
)

const (
	credRefKey = "credentialReference"
	appRefName = "appName"
)

type provisioner struct {
	client       client
	credExpDays  int
	cacheManager cacheManager
}

type cacheManager interface {
	GetAccessRequestsByApp(managedAppName string) []*v1.ResourceInstance
	GetAPIServiceInstanceByName(apiName string) (*v1.ResourceInstance, error)
}

type client interface {
	CreateDeveloperApp(newApp models.DeveloperApp) (*models.DeveloperApp, error)
	RemoveDeveloperApp(appName, developerID string) error
	GetDeveloperID() string
	GetDeveloperApp(name string) (*models.DeveloperApp, error)
	GetAppCredential(appName, devID, key string) (*models.DeveloperAppCredentials, error)
	CreateAppCredential(appName, devID string, expDays int) (*models.DeveloperAppCredentials, error)
	RemoveAppCredential(appName, devID, key string) error
	AddProductCredential(appName, devID, key string, cpr apigee.CredentialProvisionRequest) (*models.DeveloperAppCredentials, error)
	RemoveProductCredential(appName, devID, key, productName string) error
	UpdateAppCredential(appName, devID, key string, enable bool) error
}

// NewProvisioner creates a type to implement the SDK Provisioning methods for handling subscriptions
func NewProvisioner(client client, credExpDays int, cacheMan cacheManager) prov.Provisioning {
	return &provisioner{
		client:       client,
		credExpDays:  credExpDays,
		cacheManager: cacheMan,
	}
}

// AccessRequestDeprovision - removes an api from an application
func (p provisioner) AccessRequestDeprovision(req prov.AccessRequest) prov.RequestStatus {
	instDetails := req.GetInstanceDetails()
	apiID := util.ToString(instDetails[defs.AttrExternalAPIID])

	log.Infof("deprovisioning access request for api %s from app %s ", apiID, req.GetApplicationName())
	ps := prov.NewRequestStatusBuilder()
	devID := p.client.GetDeveloperID()

	appName := req.GetApplicationName()
	if appName == "" {
		return failed(ps, fmt.Errorf("application name not found"))
	}

	app, err := p.client.GetDeveloperApp(appName)
	if err != nil {
		if ok := strings.Contains(err.Error(), "404"); ok {
			return ps.Success()
		}

		return failed(ps, fmt.Errorf("failed to retrieve app: %s", err))
	}

	if apiID == "" {
		return failed(ps, fmt.Errorf("%s not found", defs.AttrExternalAPIID))
	}

	var cred *models.DeveloperAppCredentials
	// find the credential that the api is linked to
	for _, c := range app.Credentials {
		for _, p := range c.ApiProducts {
			if p.Apiproduct == apiID {
				cred = &c
			}
		}
	}

	if cred == nil {
		return failed(ps, fmt.Errorf("app %s does not contain credentials for api %s", appName, apiID))
	}

	err = p.client.RemoveProductCredential(appName, devID, cred.ConsumerKey, apiID)
	if err != nil {
		return failed(ps, fmt.Errorf("failed to remove api %s from app: %s", "api-product-name", err))
	}

	log.Infof("removed access for api %s from app %s", apiID, req.GetApplicationName())

	return ps.Success()
}

// AccessRequestProvision - adds an api to an application
func (p provisioner) AccessRequestProvision(req prov.AccessRequest) (prov.RequestStatus, prov.AccessData) {
	instDetails := req.GetInstanceDetails()
	apiID := util.ToString(instDetails[defs.AttrExternalAPIID])

	log.Infof("processing access request for api %s to app %s", apiID, req.GetApplicationName())
	ps := prov.NewRequestStatusBuilder()
	devID := p.client.GetDeveloperID()

	if apiID == "" {
		return failed(ps, fmt.Errorf("%s name not found", defs.AttrExternalAPIID)), nil
	}

	appName := req.GetApplicationName()
	if appName == "" {
		return failed(ps, fmt.Errorf("application name not found")), nil
	}

	app, err := p.client.GetDeveloperApp(appName)
	if err != nil {
		return failed(ps, fmt.Errorf("failed to retrieve app %s: %s", appName, err)), nil
	}

	if len(app.Credentials) == 0 {
		// no credentials to add access too
		return ps.Success(), nil
	}

	// add api to credentials that are not associated with it
	for _, cred := range app.Credentials {
		addProd := true
		for _, p := range cred.ApiProducts {
			if p.Apiproduct == apiID {
				addProd = false
				break // already has the
			}
		}

		// add the product to this credential
		if addProd {
			cpr := apigee.CredentialProvisionRequest{
				ApiProducts: []string{apiID},
			}

			_, err = p.client.AddProductCredential(appName, devID, cred.ConsumerKey, cpr)
			if err != nil {
				return failed(ps, fmt.Errorf("error: %s", err)), nil
			}
		}
	}

	log.Infof("granted access for api %s to app %s", apiID, req.GetApplicationName())

	return ps.Success(), nil
}

// ApplicationRequestDeprovision - removes an app from apigee
func (p provisioner) ApplicationRequestDeprovision(req prov.ApplicationRequest) prov.RequestStatus {
	log.Infof("removing app %s", req.GetManagedApplicationName())
	ps := prov.NewRequestStatusBuilder()

	appName := req.GetManagedApplicationName()
	if appName == "" {
		return failed(ps, fmt.Errorf("managed application %s not found", appName))
	}

	err := p.client.RemoveDeveloperApp(appName, p.client.GetDeveloperID())
	if err != nil {
		return failed(ps, fmt.Errorf("failed to delete app: %s", err))
	}

	log.Infof("removed app %s", req.GetManagedApplicationName())

	return ps.Success()
}

// ApplicationRequestProvision - creates an apigee app
func (p provisioner) ApplicationRequestProvision(req prov.ApplicationRequest) prov.RequestStatus {
	log.Infof("provisioning app %s", req.GetManagedApplicationName())
	ps := prov.NewRequestStatusBuilder()
	app := models.DeveloperApp{
		Attributes: []models.Attribute{
			apigee.ApigeeAgentAttribute,
		},
		DeveloperId: p.client.GetDeveloperID(),
		Name:        req.GetManagedApplicationName(),
	}

	newApp, err := p.client.CreateDeveloperApp(app)
	if err != nil {
		return failed(ps, fmt.Errorf("failed to create app: %s", err))
	}

	// remove the credential created by default for the application, the credential request will create a new one
	p.client.RemoveAppCredential(app.Name, p.client.GetDeveloperID(), newApp.Credentials[0].ConsumerKey)

	log.Infof("provisioned app %s", req.GetManagedApplicationName())

	return ps.Success()
}

// CredentialDeprovision - Return success because there are no credentials to remove until the app is deleted
func (p provisioner) CredentialDeprovision(req prov.CredentialRequest) prov.RequestStatus {
	log.Infof("removing credentials for app %s", req.GetApplicationName())
	ps := prov.NewRequestStatusBuilder()

	appName := req.GetCredentialDetailsValue(appRefName)
	if appName == "" {
		return failed(ps, fmt.Errorf("application name not found"))
	}

	app, err := p.client.GetDeveloperApp(appName)
	if err != nil {
		log.Trace("application had previously been removed")
		ps.Success()
	}

	credKey := ""
	curHash := req.GetCredentialDetailsValue(credRefKey)
	if curHash == "" {
		return failed(ps, fmt.Errorf("credential reference not found"))
	}
	for _, cred := range app.Credentials {
		thisHash, _ := util.ComputeHash(cred.ConsumerKey)
		if curHash == fmt.Sprintf("%v", thisHash) {
			credKey = cred.ConsumerKey
			break
		}
	}

	if credKey == "" {
		return ps.Success()
	}

	// remove the credential created by default for the application, the credential request will create a new one
	err = p.client.RemoveAppCredential(app.Name, p.client.GetDeveloperID(), credKey)
	if err != nil {
		return failed(ps, fmt.Errorf("unexpected error removing the credential"))
	}
	return ps.Success()
}

// CredentialProvision - retrieves the app credentials for oauth or api key authentication
func (p provisioner) CredentialProvision(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	log.Infof("provisioning credentials for app %s", req.GetApplicationName())
	ps := prov.NewRequestStatusBuilder()

	appName := req.GetApplicationName()
	if appName == "" {
		return failed(ps, fmt.Errorf("application name not found")), nil
	}

	app, err := p.client.GetDeveloperApp(appName)
	if err != nil {
		return failed(ps, fmt.Errorf("error retrieving app: %s", err)), nil
	}
	cred, err := p.client.CreateAppCredential(app.Name, p.client.GetDeveloperID(), p.credExpDays)
	if err != nil {
		return failed(ps, fmt.Errorf("error creating app credential: %s", err)), nil
	}

	// get the cred expiry time if it is set
	credBuilder := prov.NewCredentialBuilder()
	if p.credExpDays > 0 {
		credBuilder = credBuilder.SetExpirationTime(time.UnixMilli(int64(cred.ExpiresAt)))
	}

	// associate all products
	accReqs := p.cacheManager.GetAccessRequestsByApp(appName)
	products := []string{}
	for _, arInst := range accReqs {
		accReq := management.NewAccessRequest("", "")
		accReq.FromInstance(arInst)
		inst, _ := p.cacheManager.GetAPIServiceInstanceByName(accReq.Spec.ApiServiceInstance)
		apiID, err := util.GetAgentDetailsValue(inst, defs.AttrExternalAPIID)
		if err == nil && apiID != "" {
			products = append(products, apiID)
		}
	}

	// update the credential
	cpr := apigee.CredentialProvisionRequest{
		ApiProducts: products,
	}

	_, err = p.client.AddProductCredential(appName, p.client.GetDeveloperID(), cred.ConsumerKey, cpr)
	if err != nil {
		return failed(ps, fmt.Errorf("error: %s", err)), nil
	}

	var cr prov.Credential
	t := req.GetCredentialType()
	if t == provisioning.APIKeyCRD {
		cr = credBuilder.SetAPIKey(cred.ConsumerKey)
	} else {
		cr = credBuilder.SetOAuthIDAndSecret(cred.ConsumerKey, cred.ConsumerSecret)
	}

	log.Infof("created credentials for app %s", req.GetApplicationName())

	hash, _ := util.ComputeHash(cred.ConsumerKey)
	return ps.AddProperty(credRefKey, fmt.Sprintf("%v", hash)).AddProperty(appRefName, appName).Success(), cr
}

// CredentialUpdate -
func (p provisioner) CredentialUpdate(req prov.CredentialRequest) (prov.RequestStatus, prov.Credential) {
	log.Infof("updating credentials for app %s", req.GetApplicationName())
	ps := prov.NewRequestStatusBuilder()

	appName := req.GetCredentialDetailsValue(appRefName)
	if appName == "" {
		return failed(ps, fmt.Errorf("application name not found")), nil
	}

	app, err := p.client.GetDeveloperApp(appName)
	if err != nil {
		return failed(ps, fmt.Errorf("error retrieving app: %s", err)), nil
	}

	credKey := ""
	curHash := req.GetCredentialDetailsValue(credRefKey)
	if curHash == "" {
		return failed(ps, fmt.Errorf("credential reference not found")), nil
	}

	for _, cred := range app.Credentials {
		thisHash, _ := util.ComputeHash(cred.ConsumerKey)
		if curHash == fmt.Sprintf("%v", thisHash) {
			credKey = cred.ConsumerKey
			break
		}
	}

	if credKey == "" {
		return failed(ps, fmt.Errorf("error retrieving the requested credential")), nil
	}

	if req.GetCredentialAction() == prov.Suspend {
		err = p.client.UpdateAppCredential(app.Name, p.client.GetDeveloperID(), credKey, false)
	} else if req.GetCredentialAction() == prov.Enable {
		err = p.client.UpdateAppCredential(app.Name, p.client.GetDeveloperID(), credKey, true)
	} else {
		return failed(ps, fmt.Errorf("cound not perform the requested action: %s", req.GetCredentialAction())), nil
	}

	if err != nil {
		return failed(ps, fmt.Errorf("error updating the app credential: %s", err)), nil
	}

	log.Infof("updated credentials for app %s", req.GetApplicationName())

	return ps.Success(), nil
}

func failed(ps prov.RequestStatusBuilder, err error) prov.RequestStatus {
	ps.SetMessage(err.Error())
	log.Error(fmt.Sprintf("subscription provisioning - %s", err))
	return ps.Failed()
}
