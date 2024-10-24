package queue

type Spec struct {
	Auto      bool   `json:"auto"`
	Bucket    string `json:"bucket"`
	Endpoint  string `json:"endpoint"`
	Port      int    `json:"port"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

func (spec Spec) UpdateEndpoint(endpoint string) Spec {
	spec.Endpoint = endpoint
	return spec
}

func (spec Spec) NoAuto() Spec {
	spec.Auto = false
	return spec
}
