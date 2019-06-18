package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type STSAkJson struct {
	AccessKeyId     string `json:"AccessKeyId,omitempty"`
	AccessKeySecret string `json:"AccessKeySecret,omitempty"`
	SecurityToken   string `json:"SecurityToken,omitempty"`
	Expiration      string `json:"Expiration,omitempty"`
	LastUpDated     string `json:"LastUpDated,omitempty"`
	Code            string `json:"Code,omitempty"`
}

// for ecs bind ram and get ak by ossutil automaticly
type EcsRoleAK struct {
	lock            sync.Mutex
	HasGet          bool
	url             string //url for get ak,such as http://100.100.100.200/latest/meta-data/Ram/security-credentials/RamRoleName
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
	LastUpDated     string
}

func (ecsRole *EcsRoleAK) String() string {
	return fmt.Sprintf("AccessKeyId:%s,AccessKeySecret:%s,SecurityToken:%s,Expiration:%s,LastUpDated:%s",
		ecsRole.AccessKeyId, ecsRole.AccessKeySecret, ecsRole.SecurityToken, ecsRole.Expiration, ecsRole.LastUpDated)
}

func (ecsRole *EcsRoleAK) GetAccessKeyID() string {
	key, _, _ := ecsRole.GetAk()
	return key
}

func (ecsRole *EcsRoleAK) GetAccessKeySecret() string {
	_, secret, _ := ecsRole.GetAk()
	return secret
}

func (ecsRole *EcsRoleAK) GetSecurityToken() string {
	_, _, token := ecsRole.GetAk()
	return token
}

func (ecsRole *EcsRoleAK) GetAk() (string, string, string) {
	ecsRole.lock.Lock()
	defer ecsRole.lock.Unlock()

	var err error
	bTimeOut := false

	if !ecsRole.HasGet {
		bTimeOut = true
	} else {
		bTimeOut, err = ecsRole.IsTimeOut()
		if err != nil {
			return "", "", ""
		}
	}

	if bTimeOut {
		err = ecsRole.HttpReqAk()
	}

	if err != nil {
		return "", "", ""
	}
	return ecsRole.AccessKeyId, ecsRole.AccessKeySecret, ecsRole.SecurityToken
}

func (ecsRole *EcsRoleAK) IsTimeOut() (bool, error) {
	utcExpirationTime, err := time.Parse("2006-01-02T15:04:05Z", ecsRole.Expiration)
	if err != nil {
		LogError("time.Parse error,Expiration is %s,%s\n", ecsRole.Expiration, err.Error())
		return false, err
	}

	// Now() returns the current local time
	nowLocalTime := time.Now()

	// Unix() returns the number of seconds elapsedsince January 1, 1970 UTC.
	// five minutes in advance
	if utcExpirationTime.Unix()-nowLocalTime.Unix()-5*60 <= 0 {
		return true, nil
	}
	return false, nil
}

func (ecsRole *EcsRoleAK) HttpReqAk() error {
	if ecsRole.url == "" {
		LogError("insight getAK error,url is empty\n")
		return fmt.Errorf("insight getAK error,url is empty")
	}

	//http time out
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	tStart := time.Now().UnixNano() / 1000 / 1000
	resp, err := c.Get(ecsRole.url)
	if err != nil {
		LogError("insight getAK,http client get error,url is %s,%s\n", ecsRole.url, err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	tEnd := time.Now().UnixNano() / 1000 / 1000

	akJson := &STSAkJson{}
	err = json.Unmarshal(body, akJson)
	if err != nil {
		LogError("insight getAK,json.Unmarshal error,body is %s,%s\n", string(body), err.Error())
		return err
	}

	// parsar json,such as
	//{
	//    "AccessKeyId" : "XXXXXXXXX",
	//    "AccessKeySecret" : "XXXXXXXXX",
	//    "Expiration" : "2017-11-01T05:20:01Z",
	//    "SecurityToken" : "XXXXXXXXX",
	//    "LastUpdated" : "2017-10-31T23:20:01Z",
	//    "Code" : "Success"
	// }

	if strings.ToUpper(akJson.Code) != "SUCCESS" {
		LogError("insight getAK,get sts ak error,code:%s\n", akJson.Code)
		return fmt.Errorf("insight getAK,get sts ak error,code:%s", akJson.Code)
	}

	if akJson.AccessKeyId == "" || akJson.AccessKeySecret == "" ||
		akJson.Expiration == "" || akJson.SecurityToken == "" ||
		akJson.LastUpDated == "" {
		LogError("insight getAK,parsar http json body error:\n%s\n", string(body))
		return fmt.Errorf("insight getAK,parsar http json body error:\n%s\n", string(body))
	}

	ecsRole.AccessKeyId = akJson.AccessKeyId
	ecsRole.AccessKeySecret = akJson.AccessKeySecret
	ecsRole.SecurityToken = akJson.SecurityToken
	ecsRole.Expiration = akJson.Expiration
	ecsRole.LastUpDated = akJson.LastUpDated

	LogInfo("get sts ak success,%s,cost:%d(ms)\n", ecsRole.String(), tEnd-tStart)
	ecsRole.HasGet = true
	return nil
}
