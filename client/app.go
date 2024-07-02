package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	daprHost := os.Getenv("DAPR_HOST")
	if daprHost == "" {
		daprHost = "http://localhost"
	}
	daprHttpPort := os.Getenv("DAPR_HTTP_PORT")
	if daprHttpPort == "" {
		daprHttpPort = "3500"
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	for i := 1; i <= 20; i++ {
		if err := func() error {
			order := `{"orderId":` + strconv.Itoa(i) + "}"
			req, err := http.NewRequest("POST", daprHost+":"+daprHttpPort+"/orders", strings.NewReader(order))
			if err != nil {
				return err
			}

			// Adding app id as part of the header
			req.Header.Add("dapr-app-id", "server")

			// Invoking a service
			response, err := client.Do(req)
			if err != nil {
				time.Sleep(1 * time.Second)
				return err
			}

			// Read the response
			result, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			response.Body.Close()

			fmt.Println("Order passed:", string(result))
			return nil
		}(); err != nil {
			log.Println(err)
		}
	}
}
