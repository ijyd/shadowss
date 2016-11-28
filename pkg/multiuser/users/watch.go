package users

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"cloud-keeper/pkg/api"
	"shadowss/pkg/multiuser/apiserverproxy"
	"shadowss/pkg/multiuser/common"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const (
	watchURL = "/api/v1/nodes"
)

type nodeEvent struct {
	Type   string       `json:"type,omitempty"`
	Object api.NodeUser `json:"object,omitempty"`
}

type rawMsg struct {
	msgType int
	data    []byte
}

func (u *Users) syncUsers(nodeev *nodeEvent, nodeName string) error {

	nodeUser := &nodeev.Object
	glog.V(5).Infof("event(%v) node user %v\r\n", nodeev.Type, nodeUser)
	phase := nodeUser.Spec.Phase
	//userRefer := &nodeev.Object.Spec.User

	// nodeUser.Name = userRefer.User.Name
	// nodeUser.Spec.NodeName = userRefer.NodeName
	// nodeUser.Spec.Phase = userRefer.Phase
	// nodeUser.Spec.User = userRefer.User

	lables, ok := nodeUser.Labels[nodeName]
	if !ok {
		return fmt.Errorf("user(%s) not for node(%s)\r\n", nodeUser.Spec.User.Name, nodeName)
	}

	switch nodeev.Type {
	case "MODIFIED":
		fallthrough
	case "ADDED":
		switch lables {
		case api.NodeUserPhaseAdd:
			glog.V(5).Infof("add new node user %v\r\n", *nodeUser)
			u.AddUsers(nodeUser)
		case api.NodeUserPhaseDelete:
			glog.V(5).Infof("delete node user %v\r\n", *nodeUser)
			u.DelUsers(nodeUser)
		default:
			glog.Warningf("ignore phase %v for user %v \r\n", phase, *nodeUser)
		}
	default:
		glog.V(5).Infof("ignore event %v", nodeev.Type)
	}

	return nil
}

func (u *Users) WatchUserLoop(nodeName string) error {
	url := url.URL{
		Scheme: "wss",
		Host:   fmt.Sprintf("%s:%d", apiserverproxy.ApiServerList[0].Host, apiserverproxy.ApiServerList[0].Port),
		Path:   fmt.Sprintf("/api/v1/watch/nodeusers"),
	}

	url.Query().Add("lableSelector", fmt.Sprintf("%s", nodeName))

	glog.V(5).Infof("start watch on %+v\r\n", url.String())

	websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	websocket.DefaultDialer.ReadBufferSize = 1024 * 5

	wsHeaders := http.Header{
		"Origin":        {"http://localhost"},
		"Authorization": {common.Token},
	}
	wsc, _, err := websocket.DefaultDialer.Dial(url.String(), wsHeaders)
	if err != nil {
		glog.Errorf("cant watch node users error(%v) at(%v)\r\n", err, url)
		return err
	}

	defer wsc.Close()

	done := make(chan struct{})
	recv := make(chan *rawMsg)

	go func() {
		defer wsc.Close()
		defer close(done)
		for {
			recvmessage := &rawMsg{}
			recvmessage.msgType, recvmessage.data, err = wsc.ReadMessage()
			if err != nil {
				glog.Errorf("read error:%v\r\n", err)
				return
			}
			recv <- recvmessage
		}
	}()

	for {
		select {
		case msg := <-recv:
			switch msg.msgType {
			case websocket.TextMessage:
				nodeev := &nodeEvent{}
				err = json.Unmarshal(msg.data, nodeev)
				if err != nil {
					return err
				}
				u.syncUsers(nodeev, nodeName)
			case websocket.CloseMessage:
				return fmt.Errorf("recevice close message")
			case websocket.BinaryMessage:
				fallthrough
			case websocket.PingMessage:
				fallthrough
			case websocket.PongMessage:
				glog.Infof("got message(%d) data:%v\r\n", msg.msgType, string(msg.data))
			}
		case <-done:
			return fmt.Errorf("receive error shutdown connect")
		}
	}
}
