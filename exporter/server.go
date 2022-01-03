package exporter

import (
	"fmt"
	"net/http"

	"neverending.dev/weather/ecowitt"
	"neverending.dev/weather/measurement/Pressure"
	"neverending.dev/weather/measurement/Rainfall"
	"neverending.dev/weather/measurement/Temperature"
	"neverending.dev/weather/measurement/Velocity"
)

func generateWeatherReport() map[string]string {
	report := make(map[string]string)

	if ecowitt.WS.Status == ecowitt.Ready {
		// These values don't go into prometheus
		// report["ecowitt_gw_timestamp"] = ecowitt.WS.Gateway.DateUTC
		// report["ecowitt_gw_model"] = ecowitt.WS.Gateway.Model
		// report["ecowitt_gw_station_type"] = ecowitt.WS.Gateway.StationType
		report["ecowitt_gw_temperature"] = ecowitt.WS.Gateway.Temperature.ToStringAs(Temperature.Celsius)
		report["ecowitt_gw_humidity"] = ecowitt.WS.Gateway.Humidity.ToString()
		report["ecowitt_gw_pressure_rel"] = ecowitt.WS.Gateway.PressureRelative.ToStringAs(Pressure.Hectopascal)
		report["ecowitt_gw_pressure_abs"] = ecowitt.WS.Gateway.PressureAbsolute.ToStringAs(Pressure.Hectopascal)

		report["ecowitt_outdoor_temperature"] = ecowitt.WS.Outdoor.Temperature.ToStringAs(Temperature.Celsius)
		report["ecowitt_outdoor_humidity"] = ecowitt.WS.Outdoor.Humidity.ToString()
		report["ecowitt_outdoor_wind_speed"] = ecowitt.WS.Outdoor.WindSpeed.ToStringAs(Velocity.KilometresPerHour)
		report["ecowitt_outdoor_wind_direction"] = fmt.Sprintf("%d", ecowitt.WS.Outdoor.WindDirection)
		report["ecowitt_outdoor_wind_gust"] = ecowitt.WS.Outdoor.WindGust.ToStringAs(Velocity.KilometresPerHour)
		report["ecowitt_outdoor_solar_radiation"] = fmt.Sprintf("%.2f", ecowitt.WS.Outdoor.SolarRadiation)
		report["ecowitt_outdoor_uv"] = fmt.Sprintf("%d", ecowitt.WS.Outdoor.UV)
		report["ecowitt_outdoor_rain_rate"] = ecowitt.WS.Outdoor.RainRate.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_event"] = ecowitt.WS.Outdoor.RainEvent.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_hourly"] = ecowitt.WS.Outdoor.RainHourly.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_daily"] = ecowitt.WS.Outdoor.RainDaily.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_weekly"] = ecowitt.WS.Outdoor.RainWeekly.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_monthly"] = ecowitt.WS.Outdoor.RainMonthly.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_yearly"] = ecowitt.WS.Outdoor.RainYearly.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_rain_total"] = ecowitt.WS.Outdoor.RainTotal.ToStringAs(Rainfall.Millimetre)
		report["ecowitt_outdoor_battery"] = fmt.Sprintf("%d", ecowitt.WS.Outdoor.Battery)

		for i, sensor := range ecowitt.WS.TemperatureHumidity {
			keystr := fmt.Sprintf("ecowitt_th_sensor_%d_temperature", i)
			report[keystr] = sensor.Temperature.ToStringAs(Temperature.Celsius)
			keystr = fmt.Sprintf("ecowitt_th_sensor_%d_humidity", i)
			report[keystr] = sensor.Humidity.ToString()
			keystr = fmt.Sprintf("ecowitt_th_sensor_%d_battery", i)
			report[keystr] = fmt.Sprintf("%.2v", sensor.Battery)
		}

		for i, sensor := range ecowitt.WS.SoilMoisture {
			keystr := fmt.Sprintf("ecowitt_soil_sensor_%d_moisture", i)
			report[keystr] = sensor.Moisture.ToString()
			keystr = fmt.Sprintf("ecowitt_soil_sensor_%d_battery", i)
			report[keystr] = fmt.Sprintf("%.2v", sensor.Battery)
		}
	}

	return report
}

func Serve(w http.ResponseWriter, r *http.Request) {
	weatherReport := generateWeatherReport()
	output := ""

	for metric, value := range weatherReport {
		output += fmt.Sprintf("%s %s\n", metric, value)
	}

	fmt.Fprint(w, output)
}
