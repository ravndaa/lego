package domainnameshop

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultBaseURL = "https://api.domeneshop.no/v0"

type dnsDomain struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
}

const (
	errCreatingHTTPRequest = "domainnameshop: error creating http request"
	errDomainNotFound      = "domainnameshop: domain not found"
	errErrorStatusCode     = "domainnameshop: statuscode higher then 299, which is not ok"
)

//DNSRecord ...
type dnsRecord struct {
	ID   int    `json:"id,omitempty"`
	Host string `json:"host,omitempty"`
	TTL  int    `json:"ttl,omitempty"`
	Type string `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}

func (d *DNSProvider) findDomain(domainname string) (*dnsDomain, error) {

	url := fmt.Sprintf("%v%v", defaultBaseURL, "/domains")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(errCreatingHTTPRequest)
	}
	// create a basic auth "token"
	auth := d.config.Token + ":" + d.config.Secret
	token := base64.StdEncoding.EncodeToString([]byte(auth))
	// add token to header.
	req.Header.Add("Authorization", "Basic "+token)

	// run the request and hope for the best.
	resp, err := d.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf(errErrorStatusCode)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}

	domainres := []dnsDomain{}
	err = json.Unmarshal(body, &domainres)
	if err != nil {
		return nil, err
	}

	for _, item := range domainres {
		if item.Domain == domainname {
			return &item, nil
		}
	}

	return nil, fmt.Errorf(errDomainNotFound)
}

func (d *DNSProvider) addTxtRecord(domainname, txtRecord string) error {

	domain, err := d.findDomain(domainname)
	if err != nil {
		return err
	}

	// craft the api url
	path := fmt.Sprintf("/domains/%v/dns", domain.ID)
	url := fmt.Sprintf("%v%v", defaultBaseURL, path)

	data := dnsRecord{
		Host: "_acme-challenge",
		Type: "TXT",
		Data: txtRecord,
	}

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(data)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	// create a basic auth "token"
	auth := d.config.Token + ":" + d.config.Secret
	authtoken := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+authtoken)

	resp, err := d.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return err
	}
	return nil
}

func (d *DNSProvider) findAcmeRecords(domainname string) (*[]dnsRecord, error) {
	domain, err := d.findDomain(domainname)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/domains/%v/dns", domain.ID)

	// search for txt records.
	url := fmt.Sprintf("%v%v?host=%v&type=%v", defaultBaseURL, path, "_acme-challenge", "TXT")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// create a basic auth "token"
	auth := d.config.Token + ":" + d.config.Secret
	token := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", "Basic "+token)

	resp, err := d.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}
	records := []dnsRecord{}
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}

	return &records, nil
}

func (d *DNSProvider) deleteTXTRecords(domainname string) error {

	domain, err := d.findDomain(domainname)
	if err != nil {
		return err
	}

	records, err := d.findAcmeRecords(domainname)
	if err != nil {
		return err
	}

	for _, record := range *records {
		url := fmt.Sprintf("%v/domains/%v/dns/%v", defaultBaseURL, domain.ID, record)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
		auth := d.config.Token + ":" + d.config.Secret
		authtoken := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Basic "+authtoken)
		resp, err := d.config.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 201 {
			return err
		}
	}

	return nil
}
