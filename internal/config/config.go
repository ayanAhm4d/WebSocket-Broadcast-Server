package config

//Configuration for the websocket server
type ServerConfig struct {
	Host string
	Port string
}

//Configuration for the websocket client
type ClientConfig struct {
	ServerAddr string
	Username   string
}
