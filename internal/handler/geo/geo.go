package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GeoInfo struct {
	City    string `json:"city"`
	Country string `json:"country"`
	ISP     string `json:"isp"`
	Org     string `json:"org"`
}

func GetCityByIP(ip string) (string, error) {
	if ip == "::1" || ip == "127.0.0.1" {
		return "localhost", nil
	}

	client := http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,city,country,isp,org", ip)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result GeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.City == "" || result.Country == "" {
		return "неизвестно", nil
	}
	return fmt.Sprintf("%s, %s", result.City, result.Country), nil
}
