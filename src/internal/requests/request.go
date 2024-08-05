package requests

import (
	"errors"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	timeout             = 60 // timeout in seconds for connections to urls.
	rateLimited         = 0
	successfulRequests  = 0
	serverLoadReached   = 0
	proxyBackendFailure = 0
	requestFailures     = 0
	timedOut            = 0
)

func Request(url string, parameters map[string]string, proxyIterator *entities.ProxyIterator, retryCount int) (string, error) {
	if retryCount > 100 { // Limiting the number of retries
		return "", errors.New("retry limit exceeded")
	}

	retryCount++

	req, err := buildRequest(url, parameters)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	client, clientErr := build_client(proxyIterator)
	if clientErr != nil {
		log.Println(clientErr.Error())
		return "", clientErr
	}

	response, responseError := client.Do(req)
	if responseError != nil {
		if strings.Contains(responseError.Error(), "Client.Timeout") {
			timedOut++
			if timedOut%100 == 0 {
				fmt.Printf("%d timeouts\n", timedOut)
			}
		} else {
			fmt.Println(responseError.Error())
		}

		// Proxy failure
		// timeouts / offline proxies... just retry with a new proxy
		time.Sleep(time.Duration(100) * time.Millisecond)
		requestFailures++
		if requestFailures%1000 == 0 {
			fmt.Printf("[%d] Client failed to send request.\n", requestFailures)
		}
		return Request(url, parameters, proxyIterator, retryCount)
	}

	switch response.StatusCode {
	case 404:
		// page not found
		// user banned or username changed
		return "", errors.New("page not found")

	case 429:
		// Proxy rate limited
		time.Sleep(time.Duration(rand.Intn(500)+1000) * time.Millisecond)
		rateLimited++
		if rateLimited%100 == 0 {
			fmt.Printf("%d Rate limit reached.\n", rateLimited)
		}
		return Request(url, parameters, proxyIterator, retryCount)
	case 502:
		// Proxy backend issue
		time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)
		proxyBackendFailure++
		if proxyBackendFailure%100 == 0 {
			fmt.Printf("%d Proxy backend issue.\n", proxyBackendFailure)
		}
		return Request(url, parameters, proxyIterator, retryCount)
	case 503:
		// Server under high load
		time.Sleep(time.Duration(rand.Intn(1000)+1500) * time.Millisecond)
		serverLoadReached++
		if serverLoadReached%100 == 0 {
			fmt.Printf("%d Server load reached.\n", serverLoadReached)
		}
		return Request(url, parameters, proxyIterator, retryCount)
	case 504:
		// Proxy backend issue
		time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)
		proxyBackendFailure++
		if proxyBackendFailure%100 == 0 {
			fmt.Printf("%d Proxy backend issue.\n", proxyBackendFailure)
		}
		return Request(url, parameters, proxyIterator, retryCount)

	case 200:
		bodyBytes, readErr := io.ReadAll(response.Body)
		successfulRequests++

		if successfulRequests%1000 == 0 {
			fmt.Printf("%d Valid requests\n", successfulRequests)
		}

		if readErr != nil {
			// local Error reading body
			log.Println(readErr.Error())
			return Request(url, parameters, proxyIterator, retryCount)
		}

		bodyString := string(bodyBytes)

		// 429 but HTML Page
		if strings.Contains(bodyString, "Due to excessive use of the Hiscore system, your IP has been temporarily blocked") {
			rateLimited++
			if rateLimited%100 == 0 {
				fmt.Printf("%d Rate limit reached.\n", rateLimited)
			}
			// rate limits reached on this Script should Sleep for long duration.
			// This will ensure the Skill Updater has priority over these Requests.
			time.Sleep(10 * time.Second)
			return Request(url, parameters, proxyIterator, retryCount)
		} else {
			return bodyString, nil
		}
	default:
		fmt.Println(fmt.Sprintf("Unknown status codesss: %d", response.StatusCode))
		break
	}

	return "", errors.New(fmt.Sprintf("Unknown status code: %d", response.StatusCode))

}

func buildRequest(api_url string, parameters map[string]string) (*http.Request, error) {
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
