package main

type registerRequest struct {
	LocalIP   string `json:"localIp"`
	LocalPort int    `json:"localPort"`
}

type registerResponse struct {
	Status   string `json:"status"`
	ClientID string `json:"clientId"`
	Error    string `json:"error"`
}

type discoverResponse struct {
	Status  string           `json:"status"`
	Clients []clientSnapshot `json:"clients"`
	Error   string           `json:"error"`
}

type clientSnapshot struct {
	ClientID   string `json:"clientId"`
	PublicIP   string `json:"publicIp"`
	PublicPort int    `json:"publicPort"`
	LocalIP    string `json:"localIp"`
	LocalPort  int    `json:"localPort"`
}
