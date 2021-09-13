package zentaoService

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
)

func Login(baseUrl string, account string, password string) bool {

	ok := GetConfig(baseUrl)

	if !ok {
		return false
	}

	// $referer = '', $from = ''
	uri := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url := baseUrl + uri

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	var log string
	_, ok = client.PostStr(url, params)
	if ok {
		if vari.Verbose {
			log = i118Utils.I118Prt.Sprintf("success_to_login")
			logUtils.PrintToCmd("success_to_login", color.FgRed)
			logUtils.Screen(log)
		}
	} else {
		log = i118Utils.I118Prt.Sprintf("fail_to_login")
		logUtils.PrintToCmd(log, color.FgRed)
	}
	return ok
}

// func GetCookie(baseUrl string, account string, password string) bool {

// 	ok := GetConfig(baseUrl)

// 	if !ok {
// 		return false
// 	}

// 	// $referer = '', $from = ''
// 	uri := ""
// 	if vari.RequestType == constant.RequestTypePathInfo {
// 		uri = "user-login.json"
// 	} else {
// 		uri = "index.php?m=user&f=login&t=json"
// 	}
// 	url := baseUrl + uri

// 	params := make(map[string]string)
// 	params["account"] = account
// 	params["password"] = password
// 	// params := map[string] interface {"C":5, "Go":4.5, "Python":4.5, "C++":2 }

// 	var log string
// 	respdata, ok := client.PostStr2(url, params)
// 	logUtils.PrintToCmd(respdata, color.FgRed)
// 	if ok {
// 		if vari.Verbose {
// 			log = i118Utils.I118Prt.Sprintf("success_to_login")
// 			logUtils.PrintToCmd("success_to_login", color.FgRed)
// 			logUtils.Screen(log)
// 		}
// 	} else {
// 		log = i118Utils.I118Prt.Sprintf("fail_to_login")
// 		logUtils.PrintToCmd(log, color.FgRed)
// 	}
// 	return ok
// }

func GetConfig(baseUrl string) bool {
	if vari.RequestType != "" {
		return true
	}

	url := baseUrl + "?mode=getconfig"

	body, ok := client.Get(url)

	if ok {
		json, _ := simplejson.NewJson([]byte(body))

		vari.SessionId, _ = json.Get("sessionID").String()
		vari.SessionVar, _ = json.Get("sessionVar").String()
		vari.RequestType, _ = json.Get("requestType").String()
		vari.RequestFix, _ = json.Get("requestFix").String()

		return true
	} else {
		return false
	}
}
