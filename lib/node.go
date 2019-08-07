package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Node ...
type Node struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Health  string `mapstructure:"health"` // TODO
	healthy bool
}

func (n *Node) SetHealthy() {
	n.healthy = true
}

func (n *Node) SetUnhealthy() {
	n.healthy = false
}

func (n *Node) IsHealthy() bool {
	return n.healthy
}

func (n *Node) IsUnhealthy() bool {
	return !n.IsHealthy()
}

// Handler forwards request to node. Based on https://stackoverflow.com/a/34725635
func (n *Node) Handler(w http.ResponseWriter, req *http.Request) {
	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%d %s", http.StatusInternalServerError, err.Error())
		return
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("http://%s:%s/%s", n.Address, n.Port, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow
	// copy proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	client := &http.Client{}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		log.Printf("%d %s", http.StatusBadGateway, err.Error())
		n.SetUnhealthy()
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		log.Printf("%d %s", http.StatusBadGateway, err.Error())
		n.SetUnhealthy()
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
	}

	w.Write(b)
	n.SetHealthy()
}
