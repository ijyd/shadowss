package access

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/glog"

	utildbus "k8s.io/kubernetes/pkg/util/dbus"
	utilexec "k8s.io/kubernetes/pkg/util/exec"
	utiliptables "k8s.io/kubernetes/pkg/util/iptables"
)

var IptablesHandler utiliptables.Interface

const (
	jumpAccept = "ACCEPT"
	jumpDrop   = "DROP"
)

func init() {
	IptablesHandler = utiliptables.New(utilexec.New(), utildbus.New(), utiliptables.ProtocolIpv4)
}

// Join all words with spaces, terminate with newline and write to buf.
func writeLine(buf *bytes.Buffer, words ...string) {
	buf.WriteString(strings.Join(words, " ") + "\n")
}

//OpenLocalPort  open port on local host
func OpenLocalPort(port int, protocol string) error {

	portRules := bytes.NewBuffer(nil)
	args := []string{
		"-m", "comment", "--comment", fmt.Sprintf(`"%s hostport %d"`, "shadowss", port),
		"-m", protocol, "-p", protocol,
		"--dport", fmt.Sprintf("%d", port),
		"-j", string(jumpAccept),
	}
	writeLine(portRules, args...)

	natLines := portRules.Bytes()
	glog.V(3).Infof("Restoring iptables rules: %s", natLines)
	ok, err := IptablesHandler.EnsureRule(utiliptables.Append, utiliptables.TableFilter, utiliptables.ChainInput, args...)
	if err != nil {
		return fmt.Errorf("Failed to execute iptables: %v", err)
	}
	glog.V(3).Infof("ensure role result %v\r\n", ok)

	return nil
}

//TurnoffLocalPort  turnoff port on local host
func TurnoffLocalPort(port int, protocol string) error {
	portRules := bytes.NewBuffer(nil)
	args := []string{
		"-m", "comment", "--comment", fmt.Sprintf(`"%s hostport %d"`, "shadowss", port),
		"-m", protocol, "-p", protocol,
		"--dport", fmt.Sprintf("%d", port),
		"-j", string(jumpAccept)}

	writeLine(portRules, args...)

	natLines := portRules.Bytes()
	glog.V(3).Infof("Ensure iptables rules: %s", natLines)
	err := IptablesHandler.DeleteRule(utiliptables.TableFilter, utiliptables.ChainInput, args...)
	if err != nil {
		return fmt.Errorf("Failed to execute iptables : %v", err)
	}

	return nil
}
