package usecase

type AuthStatusResponse struct {
	Authenticated   bool     `json:"authenticated"`
	BackendState    string   `json:"backend_state"`
	AuthURL         string   `json:"auth_url,omitempty"`
	NodeKey         string   `json:"node_key,omitempty"`
	Hostname        string   `json:"hostname,omitempty"`
	DNSName         string   `json:"dns_name,omitempty"`
	IPs             []string `json:"ips,omitempty"`
	UserDisplayName string   `json:"user_display_name,omitempty"`
	UserEmail       string   `json:"user_email,omitempty"`
	UserProfilePic  string   `json:"user_profile_pic,omitempty"`
}

type ServeStatusResponse struct {
	Services []ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Port      uint16 `json:"port"`
	Type      string `json:"type"`
	LocalURL  string `json:"local_url"`
	PublicURL string `json:"public_url,omitempty"`
}

type ServeStartResponse struct {
	Message string `json:"message"`
	Port    uint16 `json:"port"`
}

type FunnelStartResponse struct {
	Message string `json:"message"`
	Port    uint16 `json:"port"`
}

type StatusResponse struct {
	Config   map[string]interface{} `json:"config"`
	DNSName  string                 `json:"dns_name"`
	Hostname string                 `json:"hostname"`
}
