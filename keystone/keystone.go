package keystone

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/luoyancn/dubhe/logging"
	"github.com/luoyancn/merak"

	"github.com/luoyancn/ocean/common"
	"github.com/luoyancn/ocean/config"
)

var (
	re          = regexp.MustCompile(`("password":\s*)"[^"]+"`)
	secure_pass = `"password":******`
)

func Authorization() (*common.RespToken, string, error) {
	auth := common.NewAuth(config.USERNAME, config.PASSWORD,
		config.PROJECT_NAME, config.USER_DOMAIN_NAME,
		config.PROJECT_DOMAIN_NAME)
	auth_json, _ := merak.EasyJson.MarshalToString(auth)

	logging.LOG.Debugf("The auth context is %s\n",
		re.ReplaceAllString(auth_json, secure_pass))

	url := config.AUTH_URL + "/auth/tokens"
	logging.LOG.Infof("The keystone auth url is %s\n", url)

	resp, err := common.Post(url, common.HEADERS, &auth_json)
	if nil != err {
		logging.LOG.Errorf("Unexpected Error occured:%v\n", err)
		return nil, "", err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if http.StatusCreated != resp.StatusCode {
		logging.LOG.Errorf("%s\n", string(body))
		return nil, "", errors.New(string(body))
	}

	logging.LOG.Debugf(string(body))
	token := common.RespToken{}
	merak.EasyJson.Unmarshal(body, &token)
	return &token, resp.Header.Get("X-Subject-Token"), nil
}
