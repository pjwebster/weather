package airgradient

import (
	"fmt"
	"net/http"
	"strconv"

	"neverending.dev/weather/measurement/Humidity"
	"neverending.dev/weather/measurement/Temperature"
)

type AirGradientStationStatus int8

const (
	Uninitialised AirGradientStationStatus = iota
	Ready
	NotReady
)

/*
 * Sample JSON data posted to the normal AirGradient API. Note the URL includes the sensors unique ID. The
 * AirGradient API server also reports back an extended data set
 */
// POST URL: http://hw.airgradient.com/sensors/airgradient:dcf074/measures
// Payload: {"wifi":-50,"pm02":0,"rco2":582,"atmp":27.10,"rhum":49}
// Response Code: 200
// Response Data: {"timestamp":1641301497000,"date":"23:04:57","pm02":0,"pm02_clr":"green","pm02_lbl":"Good","pm02_idx":1,"pm02_raw":1,"pi02":0,"pi02_min":0,"pi02_max":0,"pi02_clr":"green","pi02_lbl":"Good","pi02_idx":1,"atmp":27.1,"rhum":49,"rco2":582,"rco2_clr":"green","rco2_lbl":"Excellent","rco2_idx":1,"wifi":-50,"heatindex":27.4,"heatindex_clr":"green","heatindex_lbl":"Good","heatindex_idx":1,"heat_index_fahrenheit":81,"heat_index_celsius":27.4,"atmp_fahrenheit":81,"c19_score":0,"c19_score_lbl":"very low"}

type AirGradientStation struct {
	Status         AirGradientStationStatus
	ID             string
	SignalStrength int64
	PM2dot5        uint64
	CO2            uint64
	Temperature    Temperature.Temperature
	Humidity       Humidity.Humidity
}

var AG = AirGradientStation{
	Status:         Uninitialised,
	ID:             "",
	SignalStrength: 0,
	PM2dot5:        0,
	CO2:            0,
	Temperature:    Temperature.Temperature{},
	Humidity:       Humidity.Humidity{},
}

func ReportHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Printf("ParseForm() err: %v", err)
		return
	}

	fmt.Printf("%q\n", req.PostForm)

	if req.PostForm.Get("station_id") != "" {
		// Indicate the structure is being updated
		AG.Status = NotReady

		AG.ID = req.PostForm.Get("station_id")

		if rssi, err := strconv.ParseInt(req.PostForm.Get("wifi"), 10, 64); err == nil {
			AG.SignalStrength = rssi
		}

		if value, err := strconv.ParseUint(req.PostForm.Get("pm02"), 10, 64); err == nil {
			AG.PM2dot5 = value
		}

		if value, err := strconv.ParseUint(req.PostForm.Get("rco2"), 10, 64); err == nil {
			AG.CO2 = value
		}

		if value, err := strconv.ParseFloat(req.PostForm.Get("atmp"), 64); err == nil {
			AG.Temperature = Temperature.New(value, Temperature.Celsius)
		}

		if value, err := strconv.ParseInt(req.PostForm.Get("rhum"), 10, 64); err == nil {
			AG.Humidity = Humidity.New(value)
		}

		// Indicate the structure has finished updating
		AG.Status = Ready
	}
}
