package main

import "fmt"

// https://refactoring.guru/design-patterns/proxy/go/example

/*
	server.go: Subject
*/
type server interface {
	handleRequest(url string, method string) (int, string)
}

/*
	application.go: Real subject
*/
type Application struct{}

func (a *Application) handleRequest(url string, method string) (int, string) {
	if url == "/order/v1" && method == "GET" {
		return 200, "Ok"
	}

	if url == "/order/create/v1" && method == "POST" {
		return 201, "Order Created"
	}
	return 404, "Not Ok"
}

/*
	nginx.go: Proxy
*/
type Nginx struct {
	application *Application
	maxRequests int
	rateLimiter map[string]int
}

func (n *Nginx) handleRequest(url, method string) (int, string) {
	allowed := n.checkRateLimiting(url)
	if !allowed {
		return 403, "Not Allowed"
	}
	return n.application.handleRequest(url, method)
}

func (n *Nginx) checkRateLimiting(url string) bool {
	if n.rateLimiter[url] == 0 {
		n.rateLimiter[url] = 1
	}
	if n.rateLimiter[url] > n.maxRequests {
		return false
	}
	n.rateLimiter[url] = n.rateLimiter[url] + 1
	return true
}

func NewNgnixServer() *Nginx {
	return &Nginx{
		application: &Application{},
		maxRequests: 2,
		rateLimiter: make(map[string]int),
	}
}

func main() {
	nginxServer := NewNgnixServer()

	createOrderUrl := "/order/create/v1"
	getOrderListUrl := "/order/v1"

	httpCode, body := nginxServer.handleRequest(getOrderListUrl, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", getOrderListUrl, httpCode, body)

	httpCode, body = nginxServer.handleRequest(getOrderListUrl, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", getOrderListUrl, httpCode, body)

	httpCode, body = nginxServer.handleRequest(getOrderListUrl, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", getOrderListUrl, httpCode, body)

	httpCode, body = nginxServer.handleRequest(createOrderUrl, "POST")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", createOrderUrl, httpCode, body)

	httpCode, body = nginxServer.handleRequest(createOrderUrl, "GET")
	fmt.Printf("\nUrl: %s\nHttpCode: %d\nBody: %s\n", createOrderUrl, httpCode, body)

}
