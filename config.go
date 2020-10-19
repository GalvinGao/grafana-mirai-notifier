package main

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	MiraiHTTP struct {
		Address string `yaml:"address"`
		AuthKey string `yaml:"auth_key"`
		QQNumber uint `yaml:"qq_number"`
	} `yaml:"mirai_http"`

	Log struct {
		Name string `yaml:"name"`
	} `yaml:"log"`

	QQ struct{
		Group uint `yaml:"group"`
	} `yaml:"qq"`
}
