package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/drenk83/WeatherTracker/internal/client/http/geocoding"
	openmeteo "github.com/drenk83/WeatherTracker/internal/client/http/open_meteo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5"
)

const (
	httpPort = ":3001"
	city     = "moscow"
)

type Reading struct {
	Name        string    `db:"name"`
	Timestamp   time.Time `db:"timestamp"`
	Temperature float64   `db:"temperature"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://drenk83:password@localhost:54321/weather")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		cityName := chi.URLParam(r, "city")
		fmt.Println(cityName)

		var readings Reading
		err = conn.QueryRow(
			ctx,
			"SELECT name, timestamp, temperature FROM reading WHERE name = $1 ORDER BY timestamp desc limit 1",
			city,
		).Scan(&readings.Name, &readings.Timestamp, &readings.Temperature)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}

		raw, err := json.Marshal(readings)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write(raw)
		if err != nil {
			log.Println(err)
		}
	})

	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(ctx, s, conn)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Горутины
	go func() {
		defer wg.Done()
		fmt.Println("starting server on port:", httpPort)
		err := http.ListenAndServe(httpPort, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Println("starting cron with job:", jobs.ID())
		s.Start()
	}()

	wg.Wait()

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		panic(err)
	}
}

func initJobs(ctx context.Context, s gocron.Scheduler, conn *pgx.Conn) (gocron.Job, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	geocodingClient := geocoding.NewClinet(httpClient)
	openMeteoClient := openmeteo.NewClinet(httpClient)

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				geoRes, err := geocodingClient.GetCoords(city)
				if err != nil {
					log.Println(err)
				}
				openRes, err := openMeteoClient.GetTemperature(geoRes.Latitude, geoRes.Longitude)
				if err != nil {
					log.Println(err)
				}

				timestamp, err := time.Parse("2006-01-02T15:04", openRes.Current.Time)
				if err != nil {
					log.Println(err)
					return
				}

				_, err = conn.Exec(
					ctx,
					"INSERT INTO reading (name, timestamp, temperature) VALUES ($1, $2, $3)",
					city,                          // $1 - name
					timestamp,                     // $2 - timestamp
					openRes.Current.Temperature2m, // $3 - temperature
				)
				if err != nil {
					log.Println("Failed to insert reading for", city)
				}

				log.Println("Update data for city:", city, "time:", timestamp, "temperature", openRes.Current.Temperature2m)
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return j, nil
}
