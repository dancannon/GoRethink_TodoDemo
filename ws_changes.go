package main

import (
	r "github.com/dancannon/gorethink"
	"log"
)

func allChanges(ch chan interface{}) {
	// Use goroutine to wait for changes. Prints the first 10 results
	go func() {
		for {
			res, err := r.Db("todo").Table("items").Changes().Run(session)
			if err != nil {
				log.Fatalln(err)
			}

			var response r.WriteChanges
			for res.Next(&response) {
				ch <- response
			}

			if res.Err() != nil {
				log.Println(res.Err())
			}
		}
	}()
}
func activeChanges(ch chan interface{}) {
	// Use goroutine to wait for changes. Prints the first 10 results
	go func() {
		for {
			res, err := r.Db("todo").Table("items").Filter(r.Row.Field("Status").Eq("active")).Changes().Run(session)
			if err != nil {
				log.Fatalln(err)
			}

			var response r.WriteChanges
			for res.Next(&response) {
				ch <- response
			}

			if res.Err() != nil {
				log.Println(res.Err())
			}
		}
	}()
}
func completedChanges(ch chan interface{}) {
	// Use goroutine to wait for changes. Prints the first 10 results
	go func() {
		for {
			res, err := r.Db("todo").Table("items").Filter(r.Row.Field("Status").Eq("complete")).Changes().Run(session)
			if err != nil {
				log.Fatalln(err)
			}

			var response r.WriteChanges
			for res.Next(&response) {
				ch <- response
			}

			if res.Err() != nil {
				log.Println(res.Err())
			}
		}
	}()
}
