package entities

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

var ProxyList = NewProxyIterator("proxies.txt")

var proxyCount int

type Proxy struct {
	IP       string
	Port     string
	User     string
	Password string
}

type ProxyIterator struct {
	proxies []Proxy
	index   int
	mutex   sync.Mutex
}

func NewProxyIterator(filename string) *ProxyIterator {
	proxies := readProxiesFromFile(filename)
	proxyCount = len(proxies)
	fmt.Printf("Loaded %d proxies.\n", proxyCount)
	if proxyCount == 0 {
		log.Fatal("You require proxies to run this program.")
	}
	return &ProxyIterator{
		proxies: proxies,
	}
}

func (p *ProxyIterator) Next() (*url.URL, error) {
	if len(p.proxies) == 0 {
		return nil, fmt.Errorf("no more proxies")
	}
	proxyPointer, err := p.next()
	proxy := *proxyPointer

	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", proxy.User, proxy.Password, proxy.IP, proxy.Port)
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Printf("Proxy parsing error: %s", err.Error())
		return nil, err
	}

	return proxyURL, nil
}

func (p *ProxyIterator) next() (*Proxy, error) {
	/*
		Locks the mutex on the lowest level method of proxy iterator.
	*/
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.proxies) == 0 {
		return nil, fmt.Errorf("no more proxies")
	}

	proxy := p.proxies[p.index]
	p.index = (p.index + 1) % len(p.proxies)

	return &proxy, nil
}

func readProxiesFromFile(filename string) []Proxy {
	file, fileNotFoundError := os.Open(filename)
	if fileNotFoundError != nil {
		log.Fatal(fileNotFoundError.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var proxies []Proxy
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		proxies = append(proxies, Proxy{
			IP:       parts[0],
			Port:     parts[1],
			User:     parts[2],
			Password: parts[3],
		})
	}

	return proxies
}
