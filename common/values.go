package common

type Values struct {
	REDIS string `yaml:"redis"`
	MONGO string `yaml:"mongo"`
	MYSQL string `yaml:"mysql"`

	KUBERNETES struct {
		Host     string `yaml:"host"`
		CAData   string `yaml:"ca_data"`
		CertData string `yaml:"cert_data"`
		KeyData  string `yaml:"key_data"`
	} `yaml:"kubernetes"`

	ELASTIC struct {
		Hosts    []string `yaml:"hosts"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"elastic"`

	STMP struct {
		Addr     string `yaml:"addr"`
		Host     string `yaml:"host"`
		Identity string `yaml:"identity"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"stmp"`

	INFLUX struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"influx"`

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

	Emqx struct {
		Host   string `yaml:"host"`
		ApiKey string `yaml:"api_key"`
	} `yaml:"emqx"`
}
