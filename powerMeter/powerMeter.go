package powerMeter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PowerMeter struct {
	Power      float64   `json:"power"`
	Counters   []float64 `json:"counters"`
	TotalPower float64   `json:"total"`
	// {"power":0.00,"overpower":0.00,"is_valid":true,"timestamp":1645533295,"counters":[0.000, 0.000, 0.000],"total":0}
}

func (meter *PowerMeter) Collect(shellyUrl string, relayId int) {
	requestUrl := shellyUrl + "/meter/" + fmt.Sprint(relayId)

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
	jsonErr := json.Unmarshal(body, meter)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Printf("Parsed {power:%f, counters:[%f, %f, %f], total:%f} from %s", meter.Power, meter.Counters[0], meter.Counters[1], meter.Counters[2], meter.TotalPower, requestUrl)
}
