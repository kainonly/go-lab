package common

type Values struct {
	REDIS    string `yaml:"redis"`
	MONGO    string `yaml:"mongo"`
	MYSQL    string `yaml:"mysql"`
	POSTGRES string `yaml:"postgres"`

	KUBERNETES struct {
		Host     string `yaml:"host"`
		CAData   string `yaml:"ca_data"`
		CertData string `yaml:"cert_data"`
		KeyData  string `yaml:"key_data"`
	} `yaml:"kubernetes"`

	ELASTICSEARCH struct {
		Hosts    string `yaml:"hosts"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"elasticsearch"`

	STMP struct {
		Addr     string `yaml:"addr"`
		Host     string `yaml:"host"`
		Identity string `yaml:"identity"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"stmp"`

	CLS struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"access_key_id"`
		AccessKeySecret string `yaml:"access_key_secret"`
		TopicId         string `yaml:"topic_id"`
	} `yaml:"cls"`

	INFLUX struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"influx"`

	NATS struct {
		Url  string `yaml:"url"`
		NKey string `yaml:"nkey"`
	} `yaml:"nats"`

	PULSAR struct {
		Url   string `yaml:"url"`
		Token string `yaml:"token"`
		Topic string `yaml:"topic"`
	} `yaml:"pulsar"`

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
