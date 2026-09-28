package rdap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RDAPResult struct {
	IP              string
	ObjectClassName string
	Name            string
	Handle          string
	ParentHandle    string
	StartIP         string
	EndIP           string
	IPVersion       string
	Country         string
	Type            string
	Status          []string
	Port43          string
	Conformance     []string

	CIDRs      []RDAPCIDR
	OriginASNs []uint32

	Entities []RDAPEntity
	Events   []RDAPEvent
	Remarks  []RDAPRemark
	Links    []RDAPLink

	Source string
}

type RDAPCIDR struct {
	IPv4Prefix string `json:"v4prefix"`
	IPv6Prefix string `json:"v6prefix"`
	Length     int    `json:"length"`
}

type RDAPEntity struct {
	Handle string
	Roles  []string
	Status []string
	Port43 string

	Name         string
	Kind         string
	Organization string
	Title        string
	Role         string

	Emails    []string
	Phones    []string
	Addresses []string
	URLs      []string

	PublicIDs []RDAPPublicID
	Events    []RDAPEvent
	Remarks   []RDAPRemark
	Links     []RDAPLink

	Entities []RDAPEntity
}

type RDAPPublicID struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type RDAPEvent struct {
	Action string `json:"eventAction"`
	Actor  string `json:"eventActor"`
	Date   string `json:"eventDate"`
}

type RDAPLink struct {
	Value string `json:"value"`
	Rel   string `json:"rel"`
	Href  string `json:"href"`
	Type  string `json:"type"`
}

type RDAPRemark struct {
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	Description []string   `json:"description"`
	Links       []RDAPLink `json:"links"`
}

type RDAPNotice struct {
	Title       string     `json:"title"`
	Type        string     `json:"type"`
	Description []string   `json:"description"`
	Links       []RDAPLink `json:"links"`
}

type rdapResponse struct {
	RDAPConformance []string `json:"rdapConformance"`

	ObjectClassName string   `json:"objectClassName"`
	Name            string   `json:"name"`
	Handle          string   `json:"handle"`
	ParentHandle    string   `json:"parentHandle"`
	StartAddress    string   `json:"startAddress"`
	EndAddress      string   `json:"endAddress"`
	IPVersion       string   `json:"ipVersion"`
	Country         string   `json:"country"`
	Type            string   `json:"type"`
	Status          []string `json:"status"`
	Port43          string   `json:"port43"`

	CIDRs      []RDAPCIDR `json:"cidr0_cidrs"`
	OriginASNs []uint32   `json:"arin_originas0_originautnums"`

	Entities []rdapEntityResponse `json:"entities"`
	Events   []RDAPEvent          `json:"events"`
	Remarks  []RDAPRemark         `json:"remarks"`
	Links    []RDAPLink           `json:"links"`
}

type rdapEntityResponse struct {
	Handle     string          `json:"handle"`
	Roles      []string        `json:"roles"`
	Status     []string        `json:"status"`
	Port43     string          `json:"port43"`
	VCardArray json.RawMessage `json:"vcardArray"`

	PublicIDs []RDAPPublicID       `json:"publicIds"`
	Events    []RDAPEvent          `json:"events"`
	Remarks   []RDAPRemark         `json:"remarks"`
	Links     []RDAPLink           `json:"links"`
	Entities  []rdapEntityResponse `json:"entities"`
}

func Rdap(ip string) (RDAPResult, error) {
	ip = strings.TrimSpace(ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"https://rdap.org/ip/"+ip,
		nil,
	)
	if err != nil {
		return RDAPResult{}, err
	}

	req.Header.Set("Accept", "application/rdap+json")
	req.Header.Set("User-Agent", "OSINT-Master/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return RDAPResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RDAPResult{}, fmt.Errorf(
			"RDAP lookup failed: %s",
			resp.Status,
		)
	}

	var data rdapResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return RDAPResult{}, err
	}

	result := RDAPResult{
		IP:              ip,
		ObjectClassName: data.ObjectClassName,
		Name:            data.Name,
		Handle:          data.Handle,
		ParentHandle:    data.ParentHandle,
		StartIP:         data.StartAddress,
		EndIP:           data.EndAddress,
		IPVersion:       data.IPVersion,
		Country:         data.Country,
		Type:            data.Type,
		Status:          data.Status,
		Port43:          data.Port43,
		Conformance:     data.RDAPConformance,

		CIDRs:      data.CIDRs,
		OriginASNs: data.OriginASNs,

		Entities: parseRDAPEntities(data.Entities),
		Events:   data.Events,
		Remarks:  data.Remarks,
		Links:    data.Links,

		Source: resp.Request.URL.String(),
	}

	result.Entities = enrichEntities(
		client,
		result.Entities,
		result.Source,
	)
	return result, nil
}
