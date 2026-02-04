package api

// RegisterRequest is the payload sent by a client when registering with the signalling server.
type RegisterRequest struct {
	LocalIP   string `json:"localIp"`
	LocalPort int    `json:"localPort"`
}

// RegisterResponse is the JSON response returned by the register endpoint.
type RegisterResponse struct {
	Status   string `json:"status"`
	ClientID string `json:"clientId"`
	Error    string `json:"error,omitempty"`
}

// ClientSnapshot describes a peer's contact information as returned by discovery.
type ClientSnapshot struct {
	ClientID   string `json:"clientId"`
	PublicIP   string `json:"publicIp"`
	PublicPort int    `json:"publicPort"`
	LocalIP    string `json:"localIp,omitempty"`
	LocalPort  int    `json:"localPort,omitempty"`
}

// DiscoverResponse is the JSON response returned by the discover endpoint.
type DiscoverResponse struct {
	Status  string           `json:"status"`
	Clients []ClientSnapshot `json:"clients"`
	Error   string           `json:"error,omitempty"`
}
