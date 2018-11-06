package nova

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/luoyancn/dubhe/logging"
	"github.com/luoyancn/merak"

	"github.com/luoyancn/ocean/common"
	"github.com/luoyancn/ocean/config"
)

var once sync.Once
var baseurl string = ""

func baseUrl(context *common.RespToken) string {
	once.Do(func() {
		for _, catalog := range context.Token.Catalog {
			if catalog.Type == "compute" {
				for _, endpoint := range catalog.Endpoints {
					if endpoint.Interface == "public" &&
						endpoint.Region == config.OS_REGION_NAME {
						baseurl = endpoint.Url
						return
					}
				}
			}
		}
		return
	})
	return baseurl
}

func novaheaders(authtoken string) map[string]string {
	headers := common.HEADERS
	headers[common.CONTENT_TYPE] = common.APPLICATION_JSON
	headers[common.X_OPENSTACK_NOVA_API_VERSION] = config.API_VERSION
	headers[common.X_AUTH_TOKEN] = authtoken
	return headers
}

func ForceDownAndComputeSrv(context *common.RespToken, authtoken string,
	nodename string, forceDown bool) error {
	headers := novaheaders(authtoken)
	base := baseUrl(context)
	srvid, err := getServiceId(base+"/os-services", headers, nodename)
	if nil != err {
		logging.LOG.Errorf("Unexpected Error occured:%v\n", err)
		return err
	}
	var downctx string
	if forceDown {
		downctx = `{"forced_down": true}`
	} else {
		downctx = `{"forced_down": false}`
	}

	_, err = merak.Put(base+"/os-services/"+srvid, headers, &downctx)
	return err
}

func EvacuateVmOnHost(context *common.RespToken, authtoken string,
	nodename string) error {
	headers := novaheaders(authtoken)
	base := baseUrl(context)
	vms := listVmOnHosts(base, headers, nodename)

	evaBody := `{"evacuate": {"onSharedStorage": "True"}}`
	evaReq := make(chan struct{}, len(vms))

	for _, id := range vms {
		go func(vmid string) {
			resp, err := merak.Post(
				base+"/servers/"+vmid+"/action", headers, &evaBody)
			if nil != err || resp.StatusCode > 300 {
				logging.LOG.Errorf("Failed to evacuate vm with uuid %s: %v\n", vmid, err)
			} else {
				logging.LOG.Infof("Accept the evacuate request for vm with uuid %s\n", vmid)
			}
			evaReq <- struct{}{}
		}(id)
	}

	for range vms {
	}

	return nil
}

func listVmOnHosts(url string, headers map[string]string, nodename string) []string {
	srvurl := url + "/servers?all_tenants=1&status=active&host=" + nodename
	logging.LOG.Infof("Request url is %s\n", srvurl)
	resp, err := merak.Get(srvurl, headers)
	if nil != err {
		logging.LOG.Errorf("Unexpected Error occured:%v\n", err)
		return []string{}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if http.StatusOK != resp.StatusCode {
		logging.LOG.Errorf("%s\n", string(body))
		return []string{}
	}
	logging.LOG.Debugf(string(body))
	servers := merak.EasyJson.Get(body, "servers")
	size := servers.Size()
	uuids := make([]string, 0)
	for i := 0; i < size; i++ {
		uuids = append(uuids, servers.Get(i, "id").ToString())
	}
	logging.LOG.Infof("The servers on host are %v\n", uuids)
	return uuids
}

func getServiceId(url string, headers map[string]string, nodename string) (string, error) {
	srv := url + "?binary=nova-compute&host=" + nodename
	logging.LOG.Infof("Try to get the id of compute service on host %s using request url %s \n",
		nodename, srv)
	resp, err := merak.Get(srv, headers)
	if nil != err {
		logging.LOG.Errorf("Unexpected Error occured:%v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if http.StatusOK != resp.StatusCode {
		logging.LOG.Errorf("%s\n", string(body))
		return "", errors.New(string(body))
	}

	logging.LOG.Debugf(string(body))
	return merak.EasyJson.Get(body, "services", 0, "id").ToString(), nil
}
