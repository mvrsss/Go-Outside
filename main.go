package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
	"github.com/mvrsss/Go-Outside/api"
)

var OWM_API_KEY string = "apikey"
var POSITIONSTACK_KEY string = "apikey"
var LOCATION string
var WeatherDetails api.OpenWeatherModel

type PageVariables struct {
	Weekday       string
	Day           int
	Month         time.Month
	Year          int
	Location      string
	Temperature   float64
	Tuesday       string
	Wednesday     string
	Thursday      string
	Friday        string
	Saturday      string
	Sunday        string
	WeatherDesc   string
	Precipitation int
	Humidity      int
	Visibility    float64
	WindSpeed     float64
	SunriseHour   int
	SunriseMinute int
	SunsetHour    int
	SunsetMinute  int
}

var City string
var Country string

func PostLocation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "index.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		City = r.FormValue("city")
		Country = r.FormValue("country")
		LOCATION = City + ", " + Country
		if LOCATION != ", " {
			WeatherDetails = getWeatherDetails()
			now := time.Now()
			day := now.Day()
			month := now.Month()
			year := now.Year()

			timeSunrise := time.Unix(int64(WeatherDetails.Current.Sunrise), 0)
			hourSunrise := timeSunrise.Hour()
			minuteSunrise := timeSunrise.Minute()
			timeSunset := time.Unix(int64(WeatherDetails.Current.Sunset), 0)
			hourSunset := timeSunset.Hour()
			minuteSunset := timeSunset.Minute()
			temperature := WeatherDetails.Current.Temp
			PageVars := PageVariables{now.Weekday().String(), day, month, year, LOCATION, temperature, fmt.Sprintf("%.2f", (temperature + rand.Float64())), fmt.Sprintf("%.2f", (temperature + rand.Float64())), fmt.Sprintf("%.2f", (temperature - rand.Float64())), fmt.Sprintf("%.2f", (temperature - rand.Float64())), fmt.Sprintf("%.2f", (temperature + rand.Float64())), fmt.Sprintf("%.2f", (temperature - rand.Float64())), WeatherDetails.Current.Weather[len(WeatherDetails.Current.Weather)-1].Description, WeatherDetails.Minutely[len(WeatherDetails.Minutely)-1].Precipitation, WeatherDetails.Current.Humidity, float64(WeatherDetails.Current.Visibility) / 1000, WeatherDetails.Current.WindSpeed, hourSunrise, minuteSunrise, hourSunset, minuteSunset}
			t, err := template.ParseFiles("index.html")
			if err != nil {
				log.Print("template parsing error: ", err)
			}
			err = t.Execute(w, PageVars)
			if err != nil {
				log.Print("template executing error: ", err)
			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", PostLocation)
	i := City + ", " + Country
	LOCATION = i

	print(LOCATION)
	http.ListenAndServe(":8080", nil)
}

func getWeatherDetails() api.OpenWeatherModel {
	locationReturn := make(chan api.PositionStackModel)
	weatherReturn := make(chan api.OpenWeatherModel)

	go api.GeoLocation(POSITIONSTACK_KEY, LOCATION, locationReturn)
	LocationModel := <-locationReturn
	go api.GetWeather(OWM_API_KEY, fmt.Sprintf("%f", LocationModel.Data[0].Latitude), fmt.Sprintf("%f", LocationModel.Data[0].Longitude), "7", weatherReturn)
	WeatherModel := <-weatherReturn

	return WeatherModel
}
