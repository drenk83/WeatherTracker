package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
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
