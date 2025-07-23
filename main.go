package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", geoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func geoHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)

	country, err := getCountryFromIP(ip)
	if err != nil {
		http.Error(w, "Failed to detect location", http.StatusInternalServerError)
		return
	}

	switch country {
	case "IN":
		fmt.Fprint(w, "Namaste from India 🇮🇳!")
	case "US":
		fmt.Fprint(w, "Hello from the USA 🇺🇸!")
	default:
		fmt.Fprintf(w, "Hello from %s!", country)
	}
}

func getIP(r *http.Request) string {
	hdr := r.Header.Get("X-Forwarded-For")
	if hdr != "" {
		ips := strings.Split(hdr, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func getCountryFromIP(ip string) (string, error) {
	apikey := os.Getenv("IPAPI_KEY")
	url := fmt.Sprintf("https://ipapi.co/%s/country/?key=%s", ip, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bodyBytes)), nil
}
