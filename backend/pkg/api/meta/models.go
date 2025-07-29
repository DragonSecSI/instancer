package meta

type MetaConfigResponse struct {
	WebSuffix    string `json:"web_suffix"`
	SocketSuffix string `json:"socket_suffix"`
	SocketPort   int    `json:"socket_port"`
}
