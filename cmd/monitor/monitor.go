package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/iputil"

	"github.com/skycoin/skycoin/src/util/file"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("monitor")

	whitespaceFilter = regexp.MustCompile(`\s`)

	dialTimeout = 5 * time.Second
	urlPeers    = getPathGo() + "/src/github.com/skycoin/skycoin/peers.txt"
)

func getPathGo() string {
	gopath := os.Getenv("GOPATH")
	// by default go uses GOPATH=$HOME/go if it is not set
	if gopath == "" {
		home := filepath.Clean(file.UserHome())
		gopath = filepath.Join(home, "go")
	}
	return gopath
}

// validateAddress returns a sanitized address if valid, otherwise an error
func validateAddress(ipPort string) (string, error) {
	ipPort = whitespaceFilter.ReplaceAllString(ipPort, "")
	pts := strings.Split(ipPort, ":")
	if len(pts) != 2 {
		return "", pex.ErrInvalidAddress
	}
	ip := net.ParseIP(pts[0])
	if ip == nil {
		return "", pex.ErrInvalidAddress
	} else if ip.IsLoopback() {
		if iputil.IsLocalhost(pts[0]) {
			return "", pex.ErrNoLocalhost
		}
	} else if !ip.IsGlobalUnicast() {
		return "", pex.ErrNotExternalIP
	}

	port, err := strconv.ParseUint(pts[1], 10, 16)
	if err != nil {
		return "", pex.ErrInvalidAddress
	}

	if port < 1024 {
		return "", pex.ErrPortTooLow
	}

	return ipPort, nil
}

// parseRemotePeerList parses a remote peers.txt file
// The peers list format is newline separated ip:port
// Any lines that don't parse to an ip:port are skipped
// Localhost ip:port addresses are ignored
func parseRemotePeerList(body string) []string {
	var peers []string
	for _, addr := range strings.Split(string(body), "\n") {
		addr = whitespaceFilter.ReplaceAllString(addr, "")
		if addr == "" {
			continue
		}

		// Never allow localhost addresses from the remote peers list
		a, err := validateAddress(addr)
		if err != nil {
			logger.Errorf("Remote peers list has invalid address %s: %v", addr, err)
			continue
		}

		peers = append(peers, a)
	}

	return peers
}

type peer struct {
	Addr            string `json:"address"`
	Trusted         bool   `json:"trusted"`
	StartConnection string `json:"timeStartConnection"`
	GetConnection   string `json:"timeGetConnection"`
}

func getPeers() ([]string, map[string]peer) {
	file, err := os.Open(urlPeers)
	if err != nil {
		logger.WithError(err).Error(file)
		os.Exit(1)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error(b)
		os.Exit(1)
	}

	pexList := parseRemotePeerList(string(b))

	var output []string
	peers := make(map[string]peer, len(pexList))
	for _, pex := range pexList {
		conn, err := net.DialTimeout("tcp", pex, dialTimeout)
		peers[pex] = peer{
			Addr:            pex,
			Trusted:         false,
			StartConnection: "nil",
			GetConnection:   "nil",
		}
		if err != nil {
			logger.WithError(err).Errorf("net.DialTimeout")
			continue
		}
		defer conn.Close()
		conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

		start := time.Now()
		oneByte := make([]byte, 1)
		_, err = conn.Read(oneByte)
		if err != nil {
			logger.WithError(err).Errorf("Read")
			continue
		}
		timeStart := time.Since(start)

		_, err = ioutil.ReadAll(conn)
		if err != nil {
			logger.WithError(err).Errorf("ReadAll")
			continue
		}
		timeEnd := time.Since(start)

		output = append(output, pex)
		bufferpex := peer{Addr: pex, Trusted: true, StartConnection: timeStart.String(), GetConnection: timeEnd.String()}
		peers[pex] = bufferpex

	}

	return output, peers

}
func main() {

	urlPeer := flag.String("f", urlPeers, "Url the file peer.txt")
	timeDuration := flag.Int64("t", 5, "Time dialout in host in second(s)")
	isFile := flag.Bool("o", false, "If export list the trusted")
	flag.Parse()

	urlPeers = *urlPeer

	timeParse, err := time.ParseDuration(fmt.Sprintf("%ss", fmt.Sprint(*timeDuration)))

	if err != nil {
		logger.WithError(err).Errorf("Convert")
	}

	dialTimeout = timeParse

	output, outputPeers := getPeers()

	outputJSON, err := json.MarshalIndent(outputPeers, "", "     ")
	if err != nil {
		fmt.Println("Error formating wallet to JSON. Error:", err)
		os.Exit(1)
	}

	log.Println(string(outputJSON))

	if *isFile {
		outputFILE, err := json.MarshalIndent(output, "", "     ")
		if err != nil {
			fmt.Println("Error formating wallet to JSON. Error:", err)
			os.Exit(1)
		}

		log.Println(outputFILE)
	}

}