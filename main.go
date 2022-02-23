package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/szymonrychu/shelly2prometheus/config"
	"github.com/szymonrychu/shelly2prometheus/powerMeter"
	"github.com/szymonrychu/shelly2prometheus/relayState"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type shelly25Collector struct {
	relay0      *prometheus.Desc
	relay1      *prometheus.Desc
	power0      *prometheus.Desc
	power1      *prometheus.Desc
	totalPower0 *prometheus.Desc
	totalPower1 *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newShelly25Collector(namePrefix string) *shelly25Collector {
	return &shelly25Collector{
		relay0: prometheus.NewDesc(namePrefix+"_relay0",
			"Shows relay0 state [1/0] (on/off)",
			nil, nil),
		power0: prometheus.NewDesc(namePrefix+"_power0",
			"Shows relay0 power [W]",
			nil, nil),
		totalPower0: prometheus.NewDesc(namePrefix+"_total_power0",
			"Shows relay0 total power (summary of all readings) [Wm]",
			nil, nil),

		relay1: prometheus.NewDesc(namePrefix+"_relay1",
			"Shows relay1 state [1/0] (on/off)",
			nil, nil),
		power1: prometheus.NewDesc(namePrefix+"_power1",
			"Shows relay1 power [W]",
			nil, nil),
		totalPower1: prometheus.NewDesc(namePrefix+"_total_power1",
			"Shows relay1 total power (summary of all readings) [Wm]",
			nil, nil),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *shelly25Collector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.relay0
	ch <- collector.relay1
	ch <- collector.power0
	ch <- collector.power1
}

//Collect implements required collect function for all promehteus collectors
func (collector *shelly25Collector) Collect(ch chan<- prometheus.Metric) {

	shelly25url := "http://192.168.1.201:80"

	relay0 := relayState.RelayState{}
	meter0 := powerMeter.PowerMeter{}

	relay1 := relayState.RelayState{}
	meter1 := powerMeter.PowerMeter{}

	relay0.Collect(shelly25url, 0)
	meter0.Collect(shelly25url, 0)

	relay1.Collect(shelly25url, 1)
	meter1.Collect(shelly25url, 1)

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	relay0State := relay0.IsonInt()
	relay0Power := meter0.Power
	relay0TotalPower := meter0.TotalPower

	relay1State := relay1.IsonInt()
	relay1Power := meter1.Power
	relay1TotalPower := meter1.TotalPower

	relay0Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.relay0, prometheus.GaugeValue, float64(relay0State)))
	power0Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.power0, prometheus.GaugeValue, relay0Power))
	totalPower0Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.totalPower0, prometheus.GaugeValue, relay0TotalPower))

	relay1Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.relay1, prometheus.GaugeValue, float64(relay1State)))
	power1Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.power1, prometheus.GaugeValue, relay1Power))
	totalPower1Gauge := prometheus.NewMetricWithTimestamp(time.Now(), prometheus.MustNewConstMetric(collector.totalPower1, prometheus.GaugeValue, relay1TotalPower))

	ch <- relay0Gauge
	ch <- power0Gauge
	ch <- totalPower0Gauge

	ch <- relay1Gauge
	ch <- power1Gauge
	ch <- totalPower1Gauge
}

func main() {
	conf := config.Config{}
	err := conf.Load()

	if err != nil {
		log.Fatal(err)
	}

	shelly25collector := newShelly25Collector(conf.MetricsPrefix)
	prometheus.MustRegister(shelly25collector)

	http.Handle(conf.MetricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(conf.MetricsPort), nil))
}
