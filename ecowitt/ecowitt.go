package ecowitt

import (
	"fmt"
	"net/http"
	"strconv"

	"neverending.dev/weather/measurement/Humidity"
	"neverending.dev/weather/measurement/Moisture"
	"neverending.dev/weather/measurement/Pressure"
	"neverending.dev/weather/measurement/Rainfall"
	"neverending.dev/weather/measurement/Temperature"
	"neverending.dev/weather/measurement/Velocity"
)

/*
 * Sample output from ecowitt custom settings. Updated every 60 seconds.
 * Varys depending on sensors connected to the wireless gateway.
 */
//  POST /weather HTTP/1.1
//  HOST: 192.168.1.124
//  Connection: Close
//  Content-Type: application/x-www-form-urlencoded
//  Content-Length:552
//  PASSKEY=0538D7FAACF0A4E894561405A3D7C56F&stationtype=GW1000_V1.6.8&dateutc=2022-01-04+15:08:22&tempinf=77.5&humidityin=59&baromrelin=29.521&baromabsin=29.521&tempf=75.2&humidity=78&winddir=234&windspeedmph=1.79&windgustmph=4.47&maxdailygust=6.93&solarradiation=0.00&uv=0&rainratein=0.000&eventrainin=0.000&hourlyrainin=0.000&dailyrainin=0.000&weeklyrainin=0.071&monthlyrainin=0.571&yearlyrainin=0.571&totalrainin=0.571&temp1f=75.74&humidity1=79&soilmoisture1=34&soilmoisture2=61&wh65batt=0&batt1=0&soilbatt1=1.4&soilbatt2=1.4&freq=433M&model=GW1000_Pro
//
// map[
// 	"PASSKEY":["0538D7FAACF0A4E894561405A3D7C56F"]	// gateway
// 	"baromabsin":["29.855"]							// gateway
// 	"baromrelin":["29.855"]							// gateway
// 	"batt1":["0"]									// temperaturehumidity
// 	"dailyrainin":["0.000"]							// outdoor
// 	"dateutc":["2021-12-29 01:13:21"]				// gateway
// 	"eventrainin":["0.000"]							// outdoor
// 	"freq":["433M"]									// gateway
// 	"hourlyrainin":["0.000"]						// outdoor
// 	"humidity":["59"]								// outdoor
// 	"humidity1":["74"]								// temperaturehumidity
// 	"humidityin":["62"]								// gateway
// 	"maxdailygust":["9.17"]							// outdoor
// 	"model":["GW1000_Pro"]							// gateway
// 	"monthlyrainin":["9.307"]						// outdoor
// 	"rainratein":["0.000"]							// outdoor
// 	"soilbatt1":["1.4"]								// soil
// 	"soilbatt2":["1.4"]								// soil
// 	"soilmoisture1":["34"]							// soil
// 	"soilmoisture2":["67"]							// soil
// 	"solarradiation":["460.59"]						// outdoor
// 	"stationtype":["GW1000_V1.6.8"] 				// gateway
// 	"temp1f":["72.32"]								// temperaturehumidity
// 	"tempf":["77.4"]								// outdoor
// 	"tempinf":["75.7"]              				// gateway
// 	"totalrainin":["69.988"]						// outdoor
// 	"uv":["4"]										// outdoor
// 	"weeklyrainin":["1.299"]						// outdoor
// 	"wh65batt":["0"]								// outdoor
// 	"winddir":["152"]								// outdoor
// 	"windgustmph":["8.05"]							// outdoor
// 	"windspeedmph":["3.80"]							// outdoor
// 	"yearlyrainin":["69.988"]						// outdoor
// ]

type WeatherStationStatus int8

const (
	Uninitialised WeatherStationStatus = iota
	Ready
	NotReady
)

// EcowittGateway holds the data for the ecowitt GW1000 weather station
type EcowittGateway struct {
	PASSKEY          string                  // PASSKEY
	StationType      string                  // stationtype
	Frequency        string                  // freq
	Model            string                  // model
	DateUTC          string                  // dateutc
	Temperature      Temperature.Temperature // tempinf
	Humidity         Humidity.Humidity       // humidityin
	PressureRelative Pressure.Pressure       // baromrelin
	PressureAbsolute Pressure.Pressure       // baromabsin
}

func NewEcowittGateway() EcowittGateway {
	eg := EcowittGateway{
		PASSKEY:          "",
		StationType:      "",
		Frequency:        "",
		Model:            "",
		DateUTC:          "",
		Temperature:      Temperature.New(0.0, Temperature.Undefined),
		Humidity:         Humidity.New(0),
		PressureRelative: Pressure.Pressure{},
		PressureAbsolute: Pressure.Pressure{},
	}

	return eg
}

// Outdoor Sensor Array (WS65 7-in-1: Wind Speed and Direction, UV, Solar Radiation, Temperature, Humidity, Rainfall
type OutdoorSensorArray struct {
	RainRate       Rainfall.Rainfall       // rainratein
	RainEvent      Rainfall.Rainfall       // eventrainin
	RainHourly     Rainfall.Rainfall       // hourlyrainin
	RainDaily      Rainfall.Rainfall       // dailyrainin
	RainWeekly     Rainfall.Rainfall       // weeklyrainin
	RainMonthly    Rainfall.Rainfall       // monthlyrainin
	RainYearly     Rainfall.Rainfall       // yearlyrainin
	RainTotal      Rainfall.Rainfall       // totalrainin
	Temperature    Temperature.Temperature // tempf
	Humidity       Humidity.Humidity       // humidity
	SolarRadiation float64                 // solarradiation
	UV             int64                   // uv
	WindDirection  int64                   // winddir
	WindSpeed      Velocity.Velocity       // windspeedmph
	WindGust       Velocity.Velocity       // windgustmph
	Battery        int64                   // wh65batt
}

// TemperatureSensor holds the data for an ecowitt temperature sensor (WH31, WH32)
type TemperatureHumiditySensor struct {
	ID          int
	Temperature Temperature.Temperature
	Humidity    Humidity.Humidity
	Battery     float64
}

// SoilSensor holds the data for an ecowitt WH51 soil moisture sensor
type SoilSensor struct {
	ID       int
	Moisture Moisture.Moisture
	Battery  float64
}

type WeatherStation struct {
	Status              WeatherStationStatus
	Gateway             EcowittGateway
	Outdoor             OutdoorSensorArray
	TemperatureHumidity []TemperatureHumiditySensor
	SoilMoisture        []SoilSensor
}

var WS = WeatherStation{
	Status: NotReady,
	Gateway: EcowittGateway{
		PASSKEY:          "",
		StationType:      "",
		Frequency:        "",
		Model:            "",
		DateUTC:          "",
		Temperature:      Temperature.Temperature{},
		Humidity:         Humidity.Humidity{},
		PressureRelative: Pressure.Pressure{},
		PressureAbsolute: Pressure.Pressure{},
	},
	Outdoor: OutdoorSensorArray{
		RainRate:       Rainfall.Rainfall{},
		RainEvent:      Rainfall.Rainfall{},
		RainHourly:     Rainfall.Rainfall{},
		RainDaily:      Rainfall.Rainfall{},
		RainWeekly:     Rainfall.Rainfall{},
		RainMonthly:    Rainfall.Rainfall{},
		RainYearly:     Rainfall.Rainfall{},
		RainTotal:      Rainfall.Rainfall{},
		Temperature:    Temperature.Temperature{},
		Humidity:       Humidity.Humidity{},
		SolarRadiation: 0,
		UV:             0,
		WindDirection:  0,
		WindSpeed:      Velocity.Velocity{},
		WindGust:       Velocity.Velocity{},
		Battery:        0,
	},
	TemperatureHumidity: []TemperatureHumiditySensor{},
	SoilMoisture:        []SoilSensor{},
}

func ReportHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Printf("ParseForm() err: %v", err)
		return
	}

	//fmt.Printf("%q\n", req.PostForm)

	if req.PostForm.Get("stationtype") != "" {
		//WS := new(WeatherStation)

		// Indicate the structure is being updated
		WS.Status = NotReady

		WS.Gateway.PASSKEY = req.PostForm.Get("PASSKEY")
		WS.Gateway.StationType = req.PostForm.Get("stationtype")
		WS.Gateway.Model = req.PostForm.Get("model")
		WS.Gateway.Frequency = req.PostForm.Get("freq")
		WS.Gateway.DateUTC = req.PostForm.Get("dateutc")

		if f, err := strconv.ParseFloat(req.PostForm.Get("tempinf"), 32); err == nil {
			WS.Gateway.Temperature = Temperature.New(f, Temperature.Farenheit)
		}
		if h, err := strconv.ParseInt(req.PostForm.Get("humidityin"), 10, 64); err == nil {
			WS.Gateway.Humidity = Humidity.New(h)
		}
		if b, err := strconv.ParseFloat(req.PostForm.Get("baromrelin"), 32); err == nil {
			WS.Gateway.PressureRelative = Pressure.New(b, Pressure.InchOfMercury)
		}
		if b, err := strconv.ParseFloat(req.PostForm.Get("baromabsin"), 32); err == nil {
			WS.Gateway.PressureAbsolute = Pressure.New(b, Pressure.InchOfMercury)
		}

		// Outdoor Sensor Array
		if v, err := strconv.ParseFloat(req.PostForm.Get("tempf"), 32); err == nil {
			WS.Outdoor.Temperature = Temperature.New(v, Temperature.Farenheit)
		}
		if v, err := strconv.ParseInt(req.PostForm.Get("humidity"), 10, 64); err == nil {
			WS.Outdoor.Humidity = Humidity.New(v)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("windspeedmph"), 32); err == nil {
			WS.Outdoor.WindSpeed = Velocity.New(v, Velocity.MilesPerHour)
		}
		if v, err := strconv.ParseInt(req.PostForm.Get("winddir"), 10, 64); err == nil {
			WS.Outdoor.WindDirection = v
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("windgustmph"), 32); err == nil {
			WS.Outdoor.WindGust = Velocity.New(v, Velocity.MilesPerHour)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("solarradiation"), 32); err == nil {
			WS.Outdoor.SolarRadiation = v
		}
		if v, err := strconv.ParseInt(req.PostForm.Get("uv"), 10, 64); err == nil {
			WS.Outdoor.UV = v
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("rainratein"), 32); err == nil {
			WS.Outdoor.RainRate = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("eventrainin"), 32); err == nil {
			WS.Outdoor.RainEvent = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("hourlyrainin"), 32); err == nil {
			WS.Outdoor.RainHourly = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("dailyrainin"), 32); err == nil {
			WS.Outdoor.RainDaily = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("weeklyrainin"), 32); err == nil {
			WS.Outdoor.RainWeekly = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("monthlyrainin"), 32); err == nil {
			WS.Outdoor.RainMonthly = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("yearlyrainin"), 32); err == nil {
			WS.Outdoor.RainYearly = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseFloat(req.PostForm.Get("totalrainin"), 32); err == nil {
			WS.Outdoor.RainTotal = Rainfall.New(v, Rainfall.Inch)
		}
		if v, err := strconv.ParseInt(req.PostForm.Get("wh65batt"), 10, 64); err == nil {
			WS.Outdoor.Battery = v
		}

		// Multi-channel Temperature/Humidity Sensors
		WS.TemperatureHumidity = nil
		for i := 1; i < 8; i++ {
			if req.PostForm.Get(fmt.Sprintf("temp%df", i)) != "" {
				ts := new(TemperatureHumiditySensor)
				ts.ID = i
				if f, err := strconv.ParseFloat(req.PostForm.Get(fmt.Sprintf("temp%df", i)), 32); err == nil {
					ts.Temperature = Temperature.New(f, Temperature.Farenheit)
				}
				if h, err := strconv.ParseInt(req.PostForm.Get(fmt.Sprintf("humidity%d", i)), 10, 64); err == nil {
					ts.Humidity = Humidity.New(h)
				}
				if b, err := strconv.ParseFloat(req.PostForm.Get(fmt.Sprintf("batt%d", i)), 32); err == nil {
					ts.Battery = b
				}
				WS.TemperatureHumidity = append(WS.TemperatureHumidity, *ts)
			}
		}

		// Multi-channel Soil Moisture Sensors
		WS.SoilMoisture = nil
		for i := 1; i < 8; i++ {
			if req.PostForm.Get(fmt.Sprintf("soilmoisture%d", i)) != "" {
				ss := new(SoilSensor)
				ss.ID = i
				if f, err := strconv.ParseInt(req.PostForm.Get(fmt.Sprintf("soilmoisture%d", i)), 10, 64); err == nil {
					ss.Moisture = Moisture.New(f)
				}
				if b, err := strconv.ParseFloat(req.PostForm.Get(fmt.Sprintf("soilbatt%d", i)), 32); err == nil {
					ss.Battery = b
				}
				WS.SoilMoisture = append(WS.SoilMoisture, *ss)
			}
		}

		// Indicate the structure has finished updating
		WS.Status = Ready
	}
}

/*
 * ApparentTemperature (FeelsLike, RealFeel, etc)
 * Ta = Dry bulb temperature in degrees C
 * e  = Water vapour pressure (humidity) in hPa
 * ws = Wind speed in m/s
 * Q  = Net radiation absorbed per unit of body surface (solar radiation) in W/m2
 */
func ApparentTemperature(Ta float64, e float64, ws float64, Q float64) float64 {
	AT := Ta + (0.348 * e) - (0.70 * ws) + (0.7 * (Q / (ws + 10))) - 4.25

	return AT
}
