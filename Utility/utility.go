package utility

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
)

func HttpGet(partUrl string, retData chan []byte) {

	var movieUrl string = "http://data.moviebuff.com/"
	movieUrl = movieUrl + partUrl

	rawUrl, _ := url.ParseRequestURI(movieUrl)
	reqUrl := fmt.Sprintf("%+v", rawUrl)

	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		fmt.Printf("\n Error in Creating HTTp Request Data. \n Error is %s\n", err)
		//return nil
		retData <- nil
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Access-Control-Allow-Origin", "*")
	request.Header.Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("\n Error in HTTP Get. Error is %s\n", err)
		//return nil
		retData <- nil
		return
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\n Error while reading response data. Error is %s\n", err)
		//return nil
		retData <- nil
		return
	}
	//return respData
	retData <- respData
	return
}
