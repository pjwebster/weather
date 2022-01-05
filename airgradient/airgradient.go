package airgradient

import (
	"encoding/json"
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

// POST /airgradient HTTP/1.1
// Host: 192.168.1.124:8090
// User-Agent: ESP8266HTTPClient
// Accept-Encoding: identity;q=1,chunked;q=0.1,*;q=0
// Connection: keep-alive
// content-type: application/json
// Content-Length: 87
// {"station_id":"dcf074","wifi":"-45","pm02":"0","rco2":"566","atmp":"26.50","rhum":"53"}

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

type AirGradientJSON struct {
	ID             string `json:"station_id"`
	SignalStrength string `json:"wifi"`
	PM2dot5        string `json:"pm02"`
	CO2            string `json:"rco2"`
	Temperature    string `json:"atmp"`
	Humidity       string `json:"rhum"`
}

func ReportHandler(w http.ResponseWriter, req *http.Request) {
	var m AirGradientJSON

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	// fmt.Printf("%q\n", req.Body)

	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	j, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(j))
	// if req.PostForm.Get("station_id") != "" {
	// Indicate the structure is being updated
	AG.Status = NotReady

	AG.ID = m.ID

	if rssi, err := strconv.ParseInt(m.SignalStrength, 10, 64); err == nil {
		AG.SignalStrength = rssi
	}

	if value, err := strconv.ParseUint(m.PM2dot5, 10, 64); err == nil {
		AG.PM2dot5 = value
	}

	if value, err := strconv.ParseUint(m.CO2, 10, 64); err == nil {
		AG.CO2 = value
	}

	if value, err := strconv.ParseFloat(m.Temperature, 64); err == nil {
		AG.Temperature = Temperature.New(value, Temperature.Celsius)
	}

	if value, err := strconv.ParseInt(m.Humidity, 10, 64); err == nil {
		AG.Humidity = Humidity.New(value)
	}

	// Indicate the structure has finished updating
	AG.Status = Ready
	// }
	w.Write([]byte("OK"))
}
