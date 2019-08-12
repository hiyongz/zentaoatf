package zentao

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetProductInfo(baseUrl string, productId string) *simplejson.Json {
	params := map[string]string{"productID": productId}

	myurl := baseUrl + utils.GenSuperApiUri("product", "getById", params)
	body, err := http.Get(myurl, nil)

	if err == nil {
		json, _ := simplejson.NewJson([]byte(body))

		status, _ := json.Get("status").String()
		if status == "success" {
			dataStr, _ := json.Get("data").String()
			data, _ := simplejson.NewJson([]byte(dataStr))

			return data
		}
	}

	return nil
}
