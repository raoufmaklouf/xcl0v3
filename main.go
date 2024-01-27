package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			//payload := "test"
			Port, host, path, err := parseURL(u)
			if err == nil {
				port, _ := strconv.Atoi(Port)
				_, r2 := attackRequest(host, port, path, u)
				if len(r2) > 1 {
					//fmt.Println(r2, "\n=====================================================================================\n")
					h2, b2, err := splitHTTPResponse(r2)
					if err == nil {
						scode2, err := extractStatusCode(h2)
						if err == nil {
							r3 := normalRequest("GET", "/", host, port)
							if len(r3) > 1 {
								//fmt.Println(r3, "\n=====================================================================================\n")
								h3, b3, err := splitHTTPResponse(r3)
								if err == nil {
									scode3, err := extractStatusCode(h3)
									if err == nil {
										if b2 == b3 && scode2 == scode3 && strings.Contains(b2, "z66sm65sdl0vogrdzvr0ogzjjgjgz") {

											//fmt.Println(r4, "\n=====================================================================================\n")

											//fmt.Println("\033[31m", u, "is vuln ", payload, "\033[0m")
											//fmt.Println(r2, "\n=====================================================================================\n")
											//fmt.Println(r3, "\n=====================================================================================\n")
											//fmt.Println(r4, "\n=====================================================================================\n")
											fmt.Println("\033[31m"+u, "is vuln \033[0m")

										}

									}
								}
							}

						}

					}
				}
			}

		}(url)
		wg.Wait()
	}

}

func parseURL(inputURL string) (port, rootURL, path string, err error) {
	// Parse the input URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", "", "", err
	}

	// Extract port from the Host
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) > 1 {
		port = hostParts[1]
	} else {
		// Port not specified in URL, set default based on the scheme
		if parsedURL.Scheme == "https" {
			port = "443"
		} else if parsedURL.Scheme == "http" {
			port = "80"
		}
	}

	// Construct root URL without the port and protocol
	rootURL = hostParts[0]

	// Extract path
	path = parsedURL.Path

	// If path is empty, set it to "/"
	if path == "" {
		path = "/"
	}

	return port, rootURL, path, nil
}

func splitHTTPResponse(response string) (string, string, error) {
	// Find the position of the first double newline
	index := strings.Index(response, "\r\n\r\n")

	// Ensure that a double newline is found
	if index == -1 {
		return "", "", fmt.Errorf("malformed HTTP response")
	}

	// Extract headers and body
	headers := strings.TrimSpace(response[:index])
	body := strings.TrimSpace(response[index+2:])

	return headers, body, nil
}

func extractStatusCode(rawResponse string) (int, error) {
	// Create a scanner to read from the raw response string
	scanner := bufio.NewScanner(strings.NewReader(rawResponse))

	// Read the first line
	if scanner.Scan() {
		// Extract the status code from the status line
		statusLine := scanner.Text()
		var statusCode int
		_, err := fmt.Sscanf(statusLine, "HTTP/1.1 %d", &statusCode)
		if err != nil {
			return 0, err
		}
		return statusCode, nil
	}

	// If the first line cannot be read, return an error
	return 0, fmt.Errorf("malformed HTTP response")
}
