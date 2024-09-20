package define

// GatewayOverrides 获取到的配置
type GatewayOverrides struct {
	Telemetry struct {
		Name     string  `json:"Name"`
		Endpoint string  `json:"Endpoint"`
		Batcher  string  `json:"Batcher"`
		Sampler  float64 `json:"Sampler"`
	} `json:"Telemetry"`

	ShortLink struct {
		Url      string `json:"Url"`
		Key      string `json:"Key"`
		Domain   string `json:"Domain"`
		Protocol string `json:"Protocol"`
	} `json:"ShortLink"`
}

// WebOverrides 获取到的配置
type WebOverrides struct {
	Telemetry struct {
		Name     string  `json:"Name"`
		Endpoint string  `json:"Endpoint"`
		Batcher  string  `json:"Batcher"`
		Sampler  float64 `json:"Sampler"`
	} `json:"Telemetry"`
}

// UserOverrides 获取到的配置
type UserOverrides struct {
	Telemetry struct {
		Name     string  `json:"Name"`
		Endpoint string  `json:"Endpoint"`
		Batcher  string  `json:"Batcher"`
		Sampler  float64 `json:"Sampler"`
	} `json:"Telemetry"`
}
