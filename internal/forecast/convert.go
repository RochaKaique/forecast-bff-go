package forecast

import (
	"time"
)

const hourLayout = "2006-01-02T15:00"

func (w WeatherDataResponse) DayHighestTemp() float32 {
	if len(w.Hourly.Temperature2M) < 24 {
		return 0
	}
	max := w.Hourly.Temperature2M[0]
	for _, v := range w.Hourly.Temperature2M[:24] {
		if v > max {
			max = v
		}
	}
	return max
}

func (w WeatherDataResponse) DayLowestTemp() float32 {
	if len(w.Hourly.Temperature2M) < 24 {
		return 0
	}
	min := w.Hourly.Temperature2M[0]
	for _, v := range w.Hourly.Temperature2M[:24] {
		if v < min {
			min = v
		}
	}
	return min
}

func (w WeatherDataResponse) ActualTemp() float32 {
	idx := w.nowIndexUTC()
	if idx < 0 || idx >= len(w.Hourly.Temperature2M) {
		return 0
	}
	return w.Hourly.Temperature2M[idx]
}

func (w WeatherDataResponse) ApparentTemperatureNow() float32 {
	idx := w.nowIndexUTC()
	if idx < 0 || idx >= len(w.Hourly.ApparentTemperature) {
		return 0
	}
	return w.Hourly.ApparentTemperature[idx]
}

func (w WeatherDataResponse) PrecipitationNow() float32 {
	idx := w.nowIndexUTC()
	if idx < 0 || idx >= len(w.Hourly.Precipitation) {
		return 0
	}
	return w.Hourly.Precipitation[idx]
}

func (w WeatherDataResponse) NextDaysForecast() []NextDayForecast {
	n := len(w.Hourly.Temperature2M)
	if n <= 24 || len(w.Hourly.Time) != n || len(w.Hourly.Precipitation) != n {
		return nil
	}
	days := n / 24
	out := make([]NextDayForecast, 0, max(0, days-1))
	for d := 1; d < days; d++ {
		start, end := d*24, d*24+24
		if end > n {
			break
		}
		maxT, minT := w.Hourly.Temperature2M[start], w.Hourly.Temperature2M[start]
		var precSum float32
		for j := start; j < end; j++ {
			t := w.Hourly.Temperature2M[j]
			if t > maxT {
				maxT = t
			}
			if t < minT {
				minT = t
			}
			precSum += w.Hourly.Precipitation[j]
		}
		out = append(out, NextDayForecast{
			Date: w.Hourly.Time[start], MaxTemperature: maxT,
			MinTemperature: minT, PrecipitationProb: precSum,
		})
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var now = time.Now

func (w WeatherDataResponse) nowIndexUTC() int {
	nowUTC := now().UTC().Truncate(time.Hour).Format(hourLayout)
	for i, ts := range w.Hourly.Time {
		if ts == nowUTC {
			return i
		}
	}
	return -1
}
