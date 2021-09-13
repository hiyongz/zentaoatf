package zentaoService

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/emirpasic/gods/maps"
)

func LoadTestCases(productIdStr, moduleIdStr, suiteIdStr, taskIdStr string) (testcases []model.TestCase, loginFail bool) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	// ok := GetCookie(config.Url, config.Account, config.Password)
	if !ok {
		loginFail = true
		return
	}

	if moduleIdStr != "" {
		testcases = ListCaseByModule(config.Url, productIdStr, moduleIdStr)
	} else if suiteIdStr != "" {

		testcases = ListCaseBySuite(config.Url, suiteIdStr)
	} else if taskIdStr != "" {
		testcases = ListCaseByTask(config.Url, taskIdStr)
	} else if productIdStr != "" {
		testcases = ListCaseByProduct(config.Url, productIdStr)
	} else {
		logUtils.PrintUsage()
	}

	return
}

func ListCaseByProduct(baseUrl string, productId string) []model.TestCase {
	// $productID=productId, $branch = '', $browseType = 'byModule', $param=moduleId,
	// $orderBy='id_desc', $recTotal=0, $recPerPage=10000, $pageID=1)

	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s--byModule-all-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%s&branch=&browseType=byModule&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", productId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "browse", params)
	dataStr, ok := client.Get(url)

	if ok {
		var product model.Product
		json.Unmarshal([]byte(dataStr), &product)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range product.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByModule(baseUrl string, productId string, moduleId string) []model.TestCase {
	// $productID=productId, $branch = '', $browseType = 'byModule', $param=moduleId,
	// $orderBy='id_desc', $recTotal=0, $recPerPage=10000, $pageID=1)

	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s--byModule-%s-id_asc-0-10000-1", productId, moduleId)
	} else {
		params = fmt.Sprintf("productID=%s&branch=&browseType=byModule&param=%s&orderBy=id_desc&recTotal=0&recPerPage=10000", productId, moduleId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "browse", params)
	dataStr, ok := client.Get(url)

	if ok {
		var module model.Module
		json.Unmarshal([]byte(dataStr), &module)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range module.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseBySuite(baseUrl string, suiteId string) []model.TestCase {
	// $suiteID, $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s-id_asc-0-10000-1", suiteId)
	} else {
		params = fmt.Sprintf("suiteID=%s&orderBy=id_desc&recTotal=0&recPerPage=10000", suiteId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testsuite", "view", params)
	dataStr, ok := client.Get(url)

	if ok {
		var suite model.TestSuite
		json.Unmarshal([]byte(dataStr), &suite)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range suite.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByTask(baseUrl string, taskId string) []model.TestCase {
	// $taskID, $browseType = 'all', $param = 0,
	// $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s-all-0-id_asc-0-10000-1", taskId)
	} else {
		params = fmt.Sprintf("taskID=%s&browseType=all&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", taskId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testtask", "cases", params)
	dataStr, ok := client.Get(url)

	if ok {
		var task model.TestTask
		json.Unmarshal([]byte(dataStr), &task)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range task.Runs {
			caseId := cs.Case

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func genCaseSteps(csWithSteps model.TestCase) (ret []model.TestStep) {
	// get order keys
	keys := make([]int, 0, len(csWithSteps.Steps))
	for k := range csWithSteps.Steps {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		step := csWithSteps.Steps[key]
		ret = append(ret, step)
	}

	return
}

func GetCaseById(baseUrl string, caseId string) model.TestCase {
	// $caseID, $version = 0, $from = 'testcase', $taskID = 0

	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s-0-testcase-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%s&version=0&$from=testcase&taskID=0", caseId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "view", params)
	dataStr, ok := client.Get(url)

	if ok {
		var csw model.TestCaseWrapper
		json.Unmarshal([]byte(dataStr), &csw)

		cs := csw.Case
		return cs
	}

	return model.TestCase{}
}

func GetCaseIdsBySuite(suiteId string, idMap *map[int]string) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	testcases := ListCaseBySuite(config.Url, suiteId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}

func GetCaseIdsByTask(taskId string, idMap *map[int]string) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	testcases := ListCaseByTask(config.Url, taskId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}

func CommitCase(caseId int, title string, stepMap maps.Map, stepTypeMap maps.Map, expectMap maps.Map) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	// $caseID, $comment = false
	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%d-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%d&comment=0", caseId)
	}

	url := config.Url + zentaoUtils.GenApiUri("testcase", "edit", params)

	requestObj := map[string]interface{}{"title": title,
		"steps":    commonUtils.LinkedMapToMap(stepMap),
		"stepType": commonUtils.LinkedMapToMap(stepTypeMap),
		"expects":  commonUtils.LinkedMapToMap(expectMap)}

	json, _ := json.Marshal(requestObj)
	logUtils.PrintToCmd(string(json), -1)

	var yes bool
	logUtils.PrintToWithColor("\n"+i118Utils.I118Prt.Sprintf("case_update_confirm", caseId, title), -1)
	stdinUtils.InputForBool(&yes, true, "want_to_continue")

	if yes {
		_, ok = client.PostObject(url, requestObj, true)

		if ok {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_commit_case", caseId) + "\n")
		}
	}
}

func IsMultiLine(step model.TestStep) bool {
	if strings.Index(step.Desc, "\n") > -1 || strings.Index(step.Expect, "\n") > -1 {
		return true
	}

	return false
}

func GetCaseContent(stepObj model.TestStep, seq string, independentFile bool, isChild bool) []string {
	lines := make([]string, 0)

	step := strings.TrimSpace(stepObj.Desc)
	expect := strings.TrimSpace(stepObj.Expect)

	stepStr := getStepContent(step, isChild)
	expectStr := getExpectContent(expect, isChild, independentFile)

	lines = append(lines, stepStr+expectStr)

	return lines
}

func getStepContent(str string, isChild bool) (ret string) {
	str = strings.TrimSpace(str)

	rpl := "\n"
	if isChild {
		rpl = "\n" + "  "
	}
	ret = strings.ReplaceAll(str, "\r\n", rpl)
	if isChild {
		ret = "  " + ret
	}

	return
}
func getExpectContent(str string, isChild bool, independentFile bool) (ret string) {
	str = strings.TrimSpace(str)
	if str == "" {
		return
	}

	isMultiLine := strings.Count(str, "\r\n") > 0
	if !isMultiLine {
		if independentFile {
			ret = str
		} else {
			ret = " >> " + str
		}
	} else {
		rpl := "\r\n" + "  "

		if independentFile {
			ret = ">>\n" + strings.ReplaceAll(str, "\r\n", rpl) + "\n<<"
		} else {
			ret = " >> " + strings.ReplaceAll(str, "\r\n", rpl) + "\n<<"
		}
	}

	return
}

func addPrefixSpace(str string, numb int) string {
	arr := strings.Split(str, "\r\n")

	ret := make([]string, 0)
	for _, line := range arr {
		line = fmt.Sprintf("%*s", numb, " ") + line

		ret = append(ret, line)
	}

	return strings.Join(ret, "\n")
}
