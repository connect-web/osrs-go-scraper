package requests

import (
	"errors"
	"fmt"
	entities "github.com/connect-web/Low-Latency/src/utils/entities"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	timeout               = 60 // timeout in seconds for connections to urls.
	rate_limited          = 0
	server_load_reached   = 0
	proxy_backend_failure = 0
)

func request(url string, parameters map[string]string, proxy_iterator *entities.ProxyIterator, retryCount int) (string, error) {
	if retryCount > 100 { // Limiting the number of retries
		return "", errors.New("retry limit exceeded")
	}
	retryCount++

	req, err := build_request(url, parameters)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	client, clientErr := build_client(proxy_iterator)
	if clientErr != nil {
		log.Println(clientErr.Error())
		return "", clientErr
	}

	response, responseError := client.Do(req)
	if responseError != nil {
		fmt.Println(responseError.Error())
		// Proxy failure
		// timeouts / offline proxies... just retry with a new proxy
		time.Sleep(time.Duration(100) * time.Millisecond)
		fmt.Println("Client failed to send request.")
		return request(url, parameters, proxy_iterator, retryCount)
	}

	switch response.StatusCode {
	case 429:
		// Proxy rate limited
		time.Sleep(time.Duration(rand.Intn(500)+1000) * time.Millisecond)
		rate_limited++
		if rate_limited%100 == 0 {
			fmt.Printf("%d Rate limit reached.\n", rate_limited)
		}
		return request(url, parameters, proxy_iterator, retryCount)
	case 502:
		// Proxy backend issue
		time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)
		proxy_backend_failure++
		if proxy_backend_failure%100 == 0 {
			fmt.Printf("%d Proxy backend issue.\n", proxy_backend_failure)
		}
		return request(url, parameters, proxy_iterator, retryCount)
	case 503:
		// Server under high load
		time.Sleep(time.Duration(rand.Intn(1000)+1500) * time.Millisecond)
		server_load_reached++
		if server_load_reached%100 == 0 {
			fmt.Printf("%d Server load reached.\n", server_load_reached)
		}
		return request(url, parameters, proxy_iterator, retryCount)
	case 504:
		// Proxy backend issue
		time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)
		proxy_backend_failure++
		if proxy_backend_failure%100 == 0 {
			fmt.Printf("%d Proxy backend issue.\n", proxy_backend_failure)
		}
		return request(url, parameters, proxy_iterator, retryCount)

	case 200:
		bodyBytes, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			// local Error reading body
			log.Println(readErr.Error())
			return request(url, parameters, proxy_iterator, retryCount)
		}

		bodyString := string(bodyBytes)

		// 429 but HTML Page
		if strings.Contains(bodyString, "Due to excessive use of the Hiscore system, your IP has been temporarily blocked") {
			rate_limited++
			if rate_limited%100 == 0 {
				fmt.Printf("%d Rate limit reached.\n", rate_limited)
			}
			// rate limits reached on this Script should Sleep for long duration.
			// This will ensure the Skill Updater has priority over these requests.
			time.Sleep(10 * time.Second)
			return request(url, parameters, proxy_iterator, retryCount)
		} else {
			return bodyString, nil
		}
	default:
		fmt.Println(fmt.Sprintf("Unknown status codesss: %d", response.StatusCode))
		break
	}

	return "", errors.New(fmt.Sprintf("Unknown status code: %d", response.StatusCode))

}

func build_request(api_url string, parameters map[string]string) (*http.Request, error) {
	if len(parameters) != 0 {
		// Add Query Params if required
		params := url.Values{}
		for key, value := range parameters {
			params.Add(key, value)
		}
		api_url = fmt.Sprintf("%s?%s", api_url, params.Encode())
	}

	req, err := http.NewRequest("GET", api_url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	return req, err
}

func build_client(proxy_iterator *entities.ProxyIterator) (*http.Client, error) {
	proxy_url, err := proxy_iterator.Next()
	if err != nil {
		fmt.Printf("Error building proxy url : %s\n", err.Error())
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy_url),
			DisableKeepAlives: true},
		Timeout: time.Duration(timeout) * time.Second,
	}
	return client, nil
}
