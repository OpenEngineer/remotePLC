package blocks

import (
	"../logger/"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PhilipsHueBridgeOutput struct {
	OutputBlockData
	lightNo string
	uriGet  string
	uriPut  string
	prev    []float64
}

func getHttpJson(uri string) (data map[string]interface{}, err error) {
	err = nil

	response, responseErr := http.Get(uri)
	logger.WriteError("getHttpJson()", responseErr)

	err = responseErr

	defer response.Body.Close()

	raw, readError := ioutil.ReadAll(response.Body)
	logger.WriteError("getHttpJson()", readError)
	if err == nil {
		err = readError
	}

	var rawJson interface{}
	jsonError := json.Unmarshal(raw, &rawJson)
	logger.WriteError("getHttpJson()", jsonError)

	if err == nil {
		err = jsonError
	}

	// depending on the type return a different json
	switch rawJson.(type) {
	case map[string]interface{}:
		data = rawJson.(map[string]interface{})
	case []interface{}:
		data = rawJson.([]interface{})[0].(map[string]interface{})
	}

	return
}

// method can be "PUT" or "POST" (or possibly others: to be checked)
func requestHttpJson(method string, uri string, data map[string]interface{}) error {
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

func putHttpJson(uri string, data map[string]interface{}) error {
	return requestHttpJson("PUT", uri, data)
}

func getPhilipsHueBri(x float64) (on bool, bri int) {
	if x <= 0.0 {
		bri = 1 // minimum value
	} else if x <= 1.0 {
		bri = int(x*253) + 1 // a value between 1 and 254
	} else {
		bri = 254 // maximum value
	}

	if x < 0.0 {
		on = false
	} else {
		on = true
	}

	return
}

func getPhilipsHueHue(x float64) int {
	// no wrapping needed now, can just loop around
	xrel := x - float64(int(x)) // modulo 1.0
	if xrel < 0.0 {             // TODO: really needed?
		xrel += 1.0
	}

	hue := int(65535 * xrel)
	return hue
}

func getPhilipsHueSat(x float64) int {
	var sat int
	if x < 0.0 {
		sat = 0
	} else if x >= 1.0 {
		sat = 254
	} else {
		sat = int(254 * x)
	}

	return sat
}

func updatePhilipsHueBridgeState(oldState map[string]interface{}, x []float64) map[string]interface{} {
	newState := make(map[string]interface{})

	// use only the first value of the input for the brightness
	_, briOk := oldState["bri"]
	_, onOk := oldState["on"]
	var on bool
	if len(x) >= 1 && briOk && onOk {
		var bri int
		on, bri = getPhilipsHueBri(x[0])
		if on && bri != int(oldState["bri"].(float64)) {
			newState["bri"] = bri
		}

		if on != oldState["on"].(bool) {
			newState["on"] = on
		}
	}

	if on {
		// hue-sat couples make more sense from a control point of view than xy values (which in turn better for users selecting values on a map)
		// the second value is the hue value
		_, hueOk := oldState["hue"]
		if len(x) >= 2 && hueOk {
			hue := getPhilipsHueHue(x[1])
			if hue != int(oldState["hue"].(float64)) {
				newState["hue"] = hue
			}
		}

		// the 3rd value is the sat value
		_, satOk := oldState["sat"]
		if len(x) >= 3 && satOk {
			sat := getPhilipsHueSat(x[2])
			if sat != int(oldState["sat"].(float64)) {
				newState["sat"] = sat
			}
		}
	}

	// floats after 3 are simply ignored

	return newState
}

// TODO: general function into blocks
func (b *PhilipsHueBridgeOutput) InputIsUndefined() bool {
	isUndefined := false

	for _, v := range b.in {
		if v == UNDEFINED {
			isUndefined = true
		}
	}

	return isUndefined
}

// herein any http errors are ignored:
func (b *PhilipsHueBridgeOutput) Update() {
	// stops update immediately
	if b.InputIsUndefined() {
		return
	} else {
		b.prev = SafeCopy(len(b.in), b.in, len(b.in)) // TODO: use this safe copy everywhere
	}

	// get the old state
	oldStates, err := getHttpJson(b.uriGet)

	if err == nil {
		oldState := oldStates[b.lightNo].(map[string]interface{})["state"].(map[string]interface{})
		newState := updatePhilipsHueBridgeState(oldState, b.in) // minimal state message that modifies the PhilipsHue

		// now put the state
		if len(newState) > 0 {
			putHttpJson(b.uriPut, newState)
		}

		b.out = SafeCopy(len(b.in), b.in, len(b.in))
	} else {
		b.out = []float64{}
	}
}

func PhilipsHueBridgeOutputConstructor(name string, words []string) Block {
	ipaddr := words[0]
	username := words[1]
	lightNo := words[2]

	// compose the uri string
	uriGet := fmt.Sprintf("http://%s/api/%s/lights", ipaddr, username)
	uriPut := fmt.Sprintf("%s/%s/state", uriGet, lightNo)

	// get the list of lights
	states, getErr := getHttpJson(uriGet)
	if getErr != nil {
		log.Fatal("in PhilipsHueBridgeOutputConstructior(), failed to get states. Could be bad url. ", getErr)
	}

	if _, isError := states["error"]; isError {
		log.Fatal("api error: ", states["error"].(map[string]interface{})["description"])
	}

	// check that the lightNo exists
	isFound := false
	for k, _ := range states {
		if k == lightNo {
			isFound = true
			break
		}
	}

	if !isFound {
		log.Fatal("in PhilipsHueBridgeOutputConstructor, error: couldnt find light ", lightNo)
	}

	// get the general state
	b := &PhilipsHueBridgeOutput{lightNo: lightNo, uriGet: uriGet, uriPut: uriPut}
	return b
}

var PhilipsHueBridgeOutputOk = AddConstructor("PhilipsHueBridgeOutput", PhilipsHueBridgeOutputConstructor)
