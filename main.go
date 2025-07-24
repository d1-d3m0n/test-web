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

	if country == "IN" {
		http.ServeFile(w, r, "india.html")
	} else {
		http.ServeFile(w, r, "global.html")
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
	token := os.Getenv("IPINFO_TOKEN") // Add this in Render dashboard
	url := fmt.Sprintf("https://ipinfo.io/%s/country?token=%s", ip, token)

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
