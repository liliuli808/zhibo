package mysql

type Config struct {
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}
