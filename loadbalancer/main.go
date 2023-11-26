package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Server represents a backend server in the load balancer.
type Server struct {
	URL *url.URL
}

// NewServer creates a new backend server.
func NewServer(urlStr string) (*Server, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return &Server{URL: u}, nil
}

// LoadBalancer represents a round-robin load balancer.
type LoadBalancer struct {
	servers []*Server
	index   int
}

// NewLoadBalancer creates a new round-robin load balancer.
func NewLoadBalancer(serverURLs []string) (*LoadBalancer, error) {
	var servers []*Server
	for _, urlStr := range serverURLs {
		server, err := NewServer(urlStr)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return &LoadBalancer{servers: servers, index: 0}, nil
}

// ServeHTTP implements the http.Handler interface for the LoadBalancer.
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.NextServer()
	proxy := httputil.NewSingleHostReverseProxy(server.URL)
	proxy.ServeHTTP(w, r)
}

// NextServer returns the next backend server in a round-robin fashion.
func (lb *LoadBalancer) NextServer() *Server {
	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request recieved on server ", s.URL)
}

func main() {
	// Define backend server URLs
	serverURLs := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		// Add more backend servers as needed
	}

	// Create a new load balancer
	lb, err := NewLoadBalancer(serverURLs)
	if err != nil {
		fmt.Println("Error creating load balancer:", err)
		return
	}

	// Start the load balancer server
	port := 8080
	fmt.Printf("Load Balancer listening on :%d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), lb)
	if err != nil {
		fmt.Println("Error starting load balancer:", err)
	}
}
