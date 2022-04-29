package check

import (
	"bufio"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip"
)

// BlockList searches the ipaddr in lists.blocklist.de/lists/dnsbl/all.list.
func BlockList(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "blocklist.de",
		Type: checkip.TypeSec,
	}

	file, err := getDbFilesPath("blocklist.de_all.list")
	if err != nil {
		return result, err
	}

	u := "https://lists.blocklist.de/lists/dnsbl/all.list"
	if err := updateFile(file, u, ""); err != nil {
		return result, newCheckError(err)
	}

	f, err := os.Open(file)
	if err != nil {
		return result, err
	}
	defer f.Close()

	input := bufio.NewScanner(f)
	for input.Scan() {
		fields := strings.Split(input.Text(), ":")
		if net.ParseIP(fields[0]).Equal(ipaddr) {
			result.Malicious = true
			break
		}
	}
	if err := input.Err(); err != nil {
		return result, err
	}

	return result, nil
}
