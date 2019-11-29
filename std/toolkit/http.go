package toolkit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func PostRequest(url string, data interface{}) ([]byte, int, error) {
	jsonData, _ := json.Marshal(data)
	reader := bytes.NewReader(jsonData)
	resp,err := http.DefaultClient.Post(url, "application/json;charset=UTF-8", reader)
	if resp.StatusCode != 200 || err != nil{
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}
