package blocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PhilipsHueOutput struct {
	OutputBlockData
	uri string
}

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

	client := &http.Client{}
	_, errResp := client.Do(req)
	if errResp != nil {
		return errResp
	}

	return nil
}

// herein any http errors are ignored
func (b *PhilipsHueOutput) Update() {
	// use only the first value of the input
	if len(b.in) > 0 {
		x := b.in[0]

		// this is converted into a int between 0 and 254
		var i int
		if x <= 0.0 {
			i = 0
		} else if x <= 1.0 {
			i = int(x * 254)
		} else {
			i = 254
		}

		// create the state map
		state := make(map[string]interface{})
		if i == 0 {
			state["on"] = false
			state["bri"] = 1
		} else {
			state["on"] = true
			state["bri"] = i
		}

		// now put the state
		sendHttpJson("PUT", b.uri, state)
	} else {
		b.out = b.in
	}
}

func PhilipsHueOutputConstructor(words []string) Block {
	// compile the uri string
	uri := fmt.Sprintf("http://%s/api/%s/lights", words[0], words[1])

	// get the list of lights
	lights, getErr := getHttpJson(uri)

	// get the light id, based on its unique id
	var lightId string
	lightIdFound := false
	for k, m := range lights {
		if m.(map[string]interface{})["uniqueid"].(string) == words[2] {
			lightId = k
			lightIdFound = true
			break
		}
	}

	if !lightIdFound || getErr != nil {
		log.Fatal("in PhilipsHueOutputConstructor, error: couldnt find device")
	}

	// the uri used for setting the state
	uri = uri + "/" + lightId + "/state"

	// get the general state
	b := &PhilipsHueOutput{uri: uri}
	return b
}

var PhilipsHueOutputOk = AddConstructor("PhilipsHueOutput", PhilipsHueOutputConstructor)
