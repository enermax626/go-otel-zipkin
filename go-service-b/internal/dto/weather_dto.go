package dto

type WeatherLocationResponse struct {
	Current CurrentTemperature `json:"current"`
}

type CurrentTemperature struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type WeatherTemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeatherTemperatureResponse(tempC float64, tempF float64, city string) *WeatherTemperatureResponse {
	return &WeatherTemperatureResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempC + 273,
	}
}
