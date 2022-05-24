package gokontrol

//Object is basic entity
type Object struct {
	ID          string
	GlobalID    string
	ExternalID  string
	ServiceID   string
	Status      string
	Attributes  map[string]interface{} // ignore for now, extension
	Token       string
	ExpiryDate  int64
	ApplyPolicy []*Policy
}

//Service is a registered serviced
type Service struct {
	ID            string
	ServiceID     string
	Name          string
	Key           string
	Status        string
	ExpiryDate    int64
	DefaultPolicy []*Policy
	EnforcePolicy []*Policy
}

//ObjectPermission Contains object and it's permission
type ObjectPermission struct {
	ObjectId string `json:"object_id"`
	Token    string `json:"token"`
}

type Policy struct {
	ID         string
	Name       string
	ServiceID  string
	Permission map[string]int
	Status     string
	ApplyFrom  int64
	ApplyTo    int64
}

type CertForSign struct {
	ID         string                     `json:"id"`
	GlobalID   string                     `json:"global_id"`
	ExternalID string                     `json:"external_id"`
	ServiceID  string                     `json:"service_id"`
	ExpiryDate int64                      `json:"expiry_date"`
	Attributes map[string]interface{}     `json:"attributes"`
	Permission map[string]map[string]bool `json:"permission"`
}

type Certificate struct {
	CertForSign
	Token      string `json:"token"`
	ExpiryDate int64  `json:"expiry_date"`
}
