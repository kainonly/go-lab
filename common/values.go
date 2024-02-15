package common

type Values struct {
	MYSQL string `yaml:"mysql"`
	REDIS string `yaml:"redis"`

	KUBERNETES struct {
		Host     string `yaml:"host"`
		CAData   string `yaml:"ca_data"`
		CertData string `yaml:"cert_data"`
		KeyData  string `yaml:"key_data"`
	} `yaml:"kubernetes"`

	NATS struct {
		Url  string `yaml:"url"`
		NKey string `yaml:"nkey"`
	} `yaml:"nats"`

	Apigw struct {
		Ip struct {
			SecretID  string `yaml:"secret_id"`
			SecretKey string `yaml:"secret_key"`
		} `yaml:"ip"`
	} `yaml:"apigw"`
}
