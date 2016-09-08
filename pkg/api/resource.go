package api

type LoginSpec struct {
	AuthName string `json:"authname,omitempty"`
	Auth     string `json:"auth,ommitempty"`
	Token    string `json:"token,omitempty"`
}

type Login struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec LoginSpec `json:"spec,omitempty"`
}

type NodeServer struct {
	Host   string `json:"host,omitempty"`
	Status bool   `json:"status,omitempty"`
}

type NodeAccout struct {
	ID     int64  `json:"id,omitempty"`
	Port   int64  `json:"port,omitempty"`
	Method string `json:"method,omitempty"`
}

type NodeSpec struct {
	Server  NodeServer `json:"server,omitempty"`
	Account NodeAccout `json:"account,omitempty"`
}

type Node struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec NodeSpec `json:"spec,omitempty"`
}

type NodeList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Node `json:"items"`
}

type APIServerInfor struct {
	ID   int64  `json:"id, omitempty"`
	Host string `json:"host, omitempty"`
	Port int64  `json:"port, omitempty"`
}

type APIServerSpec struct {
	Server APIServerInfor `json:"server, omitempty"`
}

type APIServer struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec APIServerSpec `json:"spec,omitempty"`
}

type APIServerList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []APIServer `json:"items"`
}

type HardWareCodeSpec struct {
	Server APIServerInfor `json:"server, omitempty"`
}

type HardWareCode struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec APIServerSpec `json:"spec,omitempty"`
}
