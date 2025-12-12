package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/drenk83/WeatherTracker/internal/client/http/geocoding"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	geocodingClient := geocoding.NewClinet(httpClient)

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")
		fmt.Println(city)

		res, err := geocodingClient.GetCoords(city)

		raw, err := json.Marshal(res)
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

	jobs, err := initJobs(s)
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

func initJobs(s gocron.Scheduler) (gocron.Job, error) {
	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("test")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return j, nil
}
