package forecast

import (
	"testing"
	"time"
)

func TestDayHighestTemp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		hourly HourlyData
		want   float32
	}{
		{"24h cheio", HourlyData{Temperature2M: []float32{1, 5, 3, 7, 2, 9, 4, 6, 0, 8, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 3, 2, 1, 0}}, 10},
		{"menos de 24h", HourlyData{Temperature2M: []float32{3, 2, 1}}, 0}, // sua função devolve 0 nesse caso
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WeatherDataResponse{Hourly: tt.hourly}
			got := w.DayHighestTemp()
			if got != tt.want {
				t.Fatalf("want=%v got=%v", tt.want, got)
			}
		})
	}
}

func TestDayLowestTemp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		hourly HourlyData
		want   float32
	}{
		{"24h cheio", HourlyData{Temperature2M: []float32{1, 5, 3, 7, 2, 9, 4, 6, 0, 8, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 3, 2, 1, 0}}, 0},
		{"menos de 24h", HourlyData{Temperature2M: []float32{3, 2, 1}}, 0}, // sua função devolve 0 nesse caso
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WeatherDataResponse{Hourly: tt.hourly}
			got := w.DayLowestTemp()
			if got != tt.want {
				t.Fatalf("want=%v got=%v", tt.want, got)
			}
		})
	}
}

func TestActualTemp(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2025, 6, 6, 12, 0, 0, 0, time.UTC)
	oldNow := now
	now = func() time.Time { return fixed }
	t.Cleanup(func() { now = oldNow })

	base := time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC)
	times := make([]string, 24)
	temps := make([]float32, 24)
	for i := 0; i < 24; i++ {
		times[i] = base.Add(time.Duration(i) * time.Hour).Format(hourLayout)
		temps[i] = float32(i)
	}

	w := WeatherDataResponse{
		Hourly: HourlyData{
			Time:          times,
			Temperature2M: temps,
		},
	}

	got := w.ActualTemp()
	want := float32(12)
	if got != want {
		t.Fatalf("ActualTemp(): want=%v got=%v", want, got)
	}
}

func TestApparentTempNow(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 6, 6, 12, 0, 0, 0, time.UTC)
	oldNow := now
	now = func() time.Time { return fixedTime }
	t.Cleanup(func() { now = oldNow })

	base := time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC)
	times := make([]string, 24)
	actualTemps := make([]float32, 24)

	for i := 0; i < 24; i++ {
		times[i] = base.Add(time.Duration(i) * time.Hour).Format(hourLayout)
		actualTemps[i] = float32(i)
	}

	w := WeatherDataResponse{
		Hourly: HourlyData{
			Time:          times,
			ApparentTemperature: actualTemps,
		},
	}

	got := w.ApparentTemperatureNow()
	want := float32(12)
	if got != want {
		t.Fatalf("ActualTemp(): want=%v got=%v", want, got)
	}

}

func TestPrecipitationNow(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2025, 6, 6, 12, 0, 0, 0, time.UTC)
	oldNow := now
	now = func() time.Time { return fixedTime }
	t.Cleanup(func() { now = oldNow })

	base := time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC)
	times := make([]string, 24)
	actualPrecipitations := make([]float32, 24)

	for i := 0; i < 24; i++ {
		times[i] = base.Add(time.Duration(i) * time.Hour).Format(hourLayout)
		actualPrecipitations[i] = float32(i)
	}

	w := WeatherDataResponse{
		Hourly: HourlyData{
			Time:          times,
			Precipitation: actualPrecipitations,
		},
	}

	got := w.PrecipitationNow()
	want := float32(12)
	if got != want {
		t.Fatalf("ActualTemp(): want=%v got=%v", want, got)
	}

}
