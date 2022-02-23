package relayState

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RelayState struct {
	Ison            bool `json:"ison"`
	Overpower       bool `json:"overpower"`
	Overtemperature bool `json:"overtemperature"`
	// {"ison":true,"has_timer":false,"timer_started":0,"timer_duration":0,"timer_remaining":0,"overpower":false,"overtemperature":false,"is_valid":true,"source":"input"}
}

func (relay *RelayState) IsonInt() int {
	var result int
	if relay.Ison {
		result = 1
	}
	return result
}

func (relay *RelayState) OverpowerInt() int {
	var result int
	if relay.Overpower {
		result = 1
	}
	return result
}

func (relay *RelayState) OvertemperatureInt() int {
	var result int
	if relay.Overtemperature {
		result = 1
	}
	return result
}

func (relay *RelayState) Collect(shellyUrl string, relayId int) {
	requestUrl := shellyUrl + "/relay/" + fmt.Sprint(relayId)

	shellyClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "shelly2prometheus")
	res, getErr := shellyClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, relay)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Printf("Parsed {ison:%t, ovepower:%t, overtemperature:%t} from %s", relay.Ison, relay.Overpower, relay.Overtemperature, requestUrl)
}
