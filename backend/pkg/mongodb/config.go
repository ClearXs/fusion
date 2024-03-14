package mongodb

type Mongodb struct {
	Url      string `yaml:"url"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}
