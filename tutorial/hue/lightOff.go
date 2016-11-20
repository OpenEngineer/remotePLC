package main

import (
	//"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	//"reflect"
)

var client = &http.Client{}

func getHttpJson(uri string) (data map[string]interface{}, err error) {
	err = nil

	response, responseErr := http.Get(uri)
	err = responseErr

	defer response.Body.Close()

	raw, readError := ioutil.ReadAll(response.Body)
	if err == nil {
		err = readError
	}

	jsonError := json.Unmarshal(raw, &data)
	if err == nil {
		err = jsonError
	}

	return
}

// method is "PUT" or "POST"
func sendHttpJson(method string, uri string, data map[string]interface{}) error {
	raw, errJson := json.Marshal(data)
	if errJson != nil {
		return errJson
	}

	req, errReq := http.NewRequest(method, uri, bytes.NewReader(raw))
	if errReq != nil {
		return errReq
	}

	_, errResp := client.Do(req)
	if errResp != nil {
		return errResp
	}

	return nil
}

func getLightState(ipaddr string, username string, lightid string) (state map[string]interface{}) {
	uri := fmt.Sprintf("http://%s/api/%s/lights", ipaddr, username)

	data, err := getHttpJson(uri)
	if err != nil {
		return
	}

	state = data[lightid].(map[string]interface{})["state"].(map[string]interface{})
	return
}

func setLightState(ipaddr string, username string, lightid string, state map[string]interface{}) error {
	uri := fmt.Sprintf("http://%s/api/%s/lights/%s/state", ipaddr, username, lightid)

	err := sendHttpJson("PUT", uri, state)
	return err
}

func scanLights(ipaddr string, username string) error {
	uri := fmt.Sprintf("http://%s/api/%s/lights", ipaddr, username)

	m := make(map[string]interface{}) // dummy empty map
	err := sendHttpJson("POST", uri, m)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func main() {
	//locators, _ := hue.DiscoverBridges(true)
	//locator := locators[0]
	//deviceType := "christian#hp"

	//bridge, err := locator.CreateUser(deviceType)
	//if err != nil {
	//log.Fatal(err)
	//}
	ipaddr := "192.168.1.100"
	uname := "gawdlP-23CzKxbGc6IkNJdwNNSCTCI40y2RbBc-G"
	//lightid := "2"

	//params := make(map[string]interface{})
	//params["on"] = true
	//params["bri"] = 254

	//data, err := json.Marshal(params)
	//if err != nil {
	//log.Fatal(err)
	//}

	//client := &http.Client{}
	//uri := fmt.Sprintf("http://%s/api/%s/lights/%s/state", ipaddr, uname, lightid)
	//data, err := getHttpJson(uri)
	//if err != nil {
	//log.Fatal(err)
	//}

	state := getLightState(ipaddr, uname, "1")
	fmt.Println(state)
	if !state["reachable"].(bool) {
		//scanLights(ipaddr, uname)
		time.Sleep(1000 * time.Millisecond)
		state = getLightState(ipaddr, uname, "1")
		fmt.Println(state)
	}

	stateIn := make(map[string]interface{})
	stateIn["on"] = true
	stateIn["bri"] = 1
	setLightState(ipaddr, uname, "1", stateIn)
	state = getLightState(ipaddr, uname, "1")
	fmt.Println(state)
}
