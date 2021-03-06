package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tj/go-spin"
	"github.com/urfave/cli"
	"gopkg.in/resty.v1"
)

func IfZero(i, j int) int {
	if i == 0 {
		return j
	}
	return i
}

func showSpinnerWhile(c *cli.Context) func() {
	// no spinner while verbose logging
	if c.GlobalBool("v") {
		return func() {}
	}
	spinner := spin.New()
	spinner.Set(spin.Box1)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
			default:
				// reprint new spinner state
				fmt.Fprintf(os.Stderr, "\r%s", spinner.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return func() {
		done <- true
		// remove spinner
		fmt.Fprintf(os.Stderr, "\033[%dD", 1)
	}
}

func optionJSON(c *cli.Context, v interface{}) bool {
	wantsJSON := c.Bool("json")
	if wantsJSON {
		data, _ := json.MarshalIndent(v, "", "\t")
		fmt.Println(string(data))
	}
	return wantsJSON
}

func newRequest(c Config) *resty.Request {
	client := resty.New()
	return client.R().
		SetHeader("accept", "application/json").
		SetHeader("content-type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
}
