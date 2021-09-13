package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"

	"github.com/ajg/form"
	"github.com/easysoft/zentaoatf/src/model"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/yosssi/gohtml"
)

var gCurCookieJar *cookiejar.Jar
var client *http.Client
var zentao_cookies string

func init() {
	gCurCookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		Jar: gCurCookieJar,
	}
}

func Get(url string) (string, bool) {
	// client := &http.Client{}
	if vari.RequestType == constant.RequestTypePathInfo {
		url = url + "?" + vari.SessionVar + "=" + vari.SessionId
	} else {
		url = url + "&" + vari.SessionVar + "=" + vari.SessionId
	}

	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_address") + url)
	}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("server_return")+reqErr.Error(), color.FgRed)
		}
		return "", false
	}

	req = set_cookies(req) // 设置cookies

	resp, respErr := client.Do(req)
	if respErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("server_return")+respErr.Error(), color.FgRed)
		}
		return "", false
	}

	// cookies := resp.Cookies() //遍历cookies
	// for _, cookie := range cookies {
	// 	fmt.Println("cookie:", cookie)
	// }

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + logUtils.ConvertUnicode(bodyStr))
	}

	var bodyJson model.ZentaoResponse
	jsonErr := json.Unmarshal(bodyStr, &bodyJson)
	if jsonErr != nil {
		if strings.Index(string(bodyStr), "<html>") > -1 {
			if vari.Verbose {
				logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + " HTML - " +
					gohtml.FormatWithLineNo(string(bodyStr)))
			}
			return "", false
		} else {
			if vari.Verbose {
				logUtils.PrintToCmd(jsonErr.Error(), color.FgRed)
			}
			return "", false
		}
	}

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		return string(bodyStr), true
	} else { // 嵌套结构
		dataStr := bodyJson.Data
		return dataStr, status == "success"
	}
}

func PostObject(url string, params interface{}, useFormFormat bool) (string, bool) {
	if vari.RequestType == constant.RequestTypePathInfo {
		url = url + "?" + vari.SessionVar + "=" + vari.SessionId
	} else {
		url = url + "&" + vari.SessionVar + "=" + vari.SessionId
	}
	url = url + "&XDEBUG_SESSION_START=PHPSTORM"

	jsonStr, _ := json.Marshal(params)
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_address") + url)
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_params") + string(jsonStr))
	}

	client := &http.Client{}

	val := string(jsonStr)
	if useFormFormat {
		val, _ = form.EncodeToString(params)
		// convert data to post fomat
		re3, _ := regexp.Compile(`([^&]*?)=`)
		val = re3.ReplaceAllStringFunc(string(val), replacePostData)
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(val))
	if reqErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("server_return")+reqErr.Error(), color.FgRed)
		}
		return "", false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = set_cookies(req) // 设置cookies

	resp, respErr := client.Do(req)
	if respErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("server_return")+respErr.Error(), color.FgRed)
		}
		return "", false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + logUtils.ConvertUnicode(bodyStr))
	}

	var bodyJson model.ZentaoResponse
	jsonErr := json.Unmarshal(bodyStr, &bodyJson)
	if jsonErr != nil {
		if strings.Index(string(bodyStr), "<html>") > -1 { // some api return a html
			if vari.Verbose {
				logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + " HTML - " +
					gohtml.FormatWithLineNo(string(bodyStr)))
			}
			return string(bodyStr), true
		} else {
			if vari.Verbose {
				logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("server_return")+jsonErr.Error(), color.FgRed)
			}
			return "", false
		}
	}

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		return string(bodyStr), true
	} else { // 嵌套结构
		dataStr := bodyJson.Data
		return dataStr, status == "success"
	}
}

func PostStr(url string, params map[string]string) (msg string, ok bool) {
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_address") + url)
	}
	client := &http.Client{}

	paramStr := ""
	idx := 0
	for pkey, pval := range params {
		if idx > 0 {
			paramStr += "&"
		}
		paramStr = paramStr + pkey + "=" + pval
		idx++
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if reqErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(reqErr.Error(), color.FgRed)
		}
		ok = false
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", vari.SessionVar+"="+vari.SessionId)
	req = set_cookies(req) // 设置cookies
	resp, respErr := client.Do(req)

	cookies := resp.Cookies() //遍历cookies
	for _, cookie := range cookies {
		fmt.Println("cookie:", cookie)
	}

	if respErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(respErr.Error(), color.FgRed)
		}
		ok = false
		return
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + logUtils.ConvertUnicode(bodyStr))
	}

	var bodyJson model.ZentaoResponse
	jsonErr := json.Unmarshal(bodyStr, &bodyJson)
	if jsonErr != nil {
		if vari.Verbose {
			if strings.Index(url, "login") == -1 { // jsonErr caused by login request return a html
				logUtils.PrintToCmd(jsonErr.Error(), color.FgRed)
			}
		}
		ok = false
		return
	}

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		msg = string(bodyStr)
		return
	} else { // 嵌套结构
		msg = bodyJson.Data
		ok = status == "success"
		return
	}
}

func PostStr2(url string, params map[string]string) (msg string, ok bool) {

	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_address") + url)
	}

	paramStr := ""
	idx := 0
	for pkey, pval := range params {
		if idx > 0 {
			paramStr += "&"
		}
		paramStr = paramStr + pkey + "=" + pval
		idx++
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if reqErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(reqErr.Error(), color.FgRed)
		}
		ok = false
		return
	}

	// zentao_cookies = "zentaosid=eae03bd53490944f4d798b195d19fcc0; lang=zh-cn; device=desktop; theme=default; keepLogin=on; za=admin; zp=d9df39f1a372fc1d92e2e4ba484017cfb22a0071; openApp=admin; windowWidth=1824; windowHeight=227"

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", vari.SessionVar+"="+vari.SessionId)
	// req.Header.Set("Cookie", zentao_cookies)
	req = set_cookies(req)

	rc := req.Cookies()
	client.Jar.SetCookies(req.URL, rc)

	resp, respErr := client.Do(req)

	if respErr != nil {
		if vari.Verbose {
			logUtils.PrintToCmd(respErr.Error(), color.FgRed)
		}
		ok = false
		return
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if vari.Verbose {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("server_return") + logUtils.ConvertUnicode(bodyStr))
	}

	var bodyJson model.ZentaoResponse
	jsonErr := json.Unmarshal(bodyStr, &bodyJson)
	if jsonErr != nil {
		if vari.Verbose {
			if strings.Index(url, "login") == -1 { // jsonErr caused by login request return a html
				logUtils.PrintToCmd(jsonErr.Error(), color.FgRed)
			}
		}
		ok = false
		return
	}

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		msg = string(bodyStr)
		return
	} else { // 嵌套结构
		msg = bodyJson.Data
		ok = status == "success"
		return
	}
}

func set_cookies(resp *http.Request) *http.Request{
	zentao_cookies = "zentaosid=eae03bd53490944f4d798b195d19fcc0; lang=zh-cn; device=desktop; theme=default; keepLogin=on; za=admin; zp=d9df39f1a372fc1d92e2e4ba484017cfb22a0071; openApp=admin; windowWidth=1824; windowHeight=227"
	resp.Header.Set("Cookie", zentao_cookies)
	return resp
}

func replacePostData(str string) string {
	str = strings.ToLower(str[:1]) + str[1:]

	if strings.Index(str, ".") > -1 {
		str = strings.Replace(str, ".", "[", -1)
		str = strings.Replace(str, "=", "]=", -1)
	}
	return str
}
