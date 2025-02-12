package sendrequest

import (
	"bytes"
	"ddns-updater/readenv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Cloudflare_ID struct {
	Result []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type Cloudflare_DNSRECORDS struct {
	Result []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	Success bool `json:"success"`
}

type CloudFlare_UPDATERECORDS struct {
	Result struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content string `json:"content"`
	}
	Success bool `json:"success"`
}

type MyIpInfo struct {
	Ip string `json:"ip"`
}

func SendPutRequest(url string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return nil, err
	}

	key, email := readenv.ReadEnvKeys()

	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", key)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, nil
}

func SendGetRequest(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return nil, err
	}

	key, email := readenv.ReadEnvKeys()

	req.Header.Add("X-Auth-Email", email)
	req.Header.Add("X-Auth-Key", key)
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, nil
}

func GetZoneId() (string, string, error) {

	url := "https://api.cloudflare.com/client/v4/zones"

	resp, err := SendGetRequest(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", "", err
	}

	var data Cloudflare_ID
	if err := json.Unmarshal(body, &data); err != nil {
		return "", "", err
	}

	if len(data.Result) == 0 {
		return "", "", fmt.Errorf("no zones found")
	}

	zone := data.Result[0]

	name := zone.Name
	id := zone.ID

	return name, id, nil
}

func ListDnsRecords(zone_id string) (Cloudflare_DNSRECORDS, error) {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zone_id)

	resp, err := SendGetRequest(url)
	if err != nil {
		return Cloudflare_DNSRECORDS{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return Cloudflare_DNSRECORDS{}, err
	}

	var data Cloudflare_DNSRECORDS
	if err := json.Unmarshal(body, &data); err != nil {
		return Cloudflare_DNSRECORDS{}, err
	}

	return data, nil
}

func GetPublicIp() (string, error) {
	url := "https://ipinfo.io"

	resp, err := SendGetRequest(url)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	var data MyIpInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return "", nil
	}

	ip := data.Ip
	return ip, nil
}

func OverwritteDnsrecords(zone_id, dns_record_id, dns_record_name, new_ip string) (CloudFlare_UPDATERECORDS, error) {

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zone_id, dns_record_id)

	var json_body = []byte(`{
		"comment": "Domain verification record",
		"content": "` + new_ip + `",
		"name": "` + dns_record_name + `",
		"proxied": false,
		"type": "A"
	  }`)

	resp, err := SendPutRequest(url, json_body)
	if err != nil {
		fmt.Println("Request failed", err)
		return CloudFlare_UPDATERECORDS{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return CloudFlare_UPDATERECORDS{}, err
	}

	var data CloudFlare_UPDATERECORDS
	if err := json.Unmarshal(body, &data); err != nil {
		return CloudFlare_UPDATERECORDS{}, err
	}
	success := data.Success

	if !success {
		fmt.Println("Records did not update, full response: ", data)
		return CloudFlare_UPDATERECORDS{}, err
	}

	return data, nil

}
