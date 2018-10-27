package skycoin

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"log"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/useragent"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	help = false
)

// Config records skycoin node and build config
type Config struct {
	Node  NodeConfig
	Build readable.BuildInfo
}

// NodeConfig records the node's configuration
type NodeConfig struct {
	// Name of the coin
	CoinName string

	// Disable peer exchange
	DisablePEX bool
	// Download peer list
	DownloadPeerList bool
	// Download the peers list from this URL
	PeerListURL string
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Enable GUI
	EnableGUI bool
	// Disable CSRF check in the wallet API
	DisableCSRF bool
	// Enable unversioned API endpoints (without the /api/v1 prefix)
	EnableUnversionedAPI bool
	// Disable CSP disable content-security-policy in http response
	DisableCSP bool
	// Comma separated list of API sets enabled on the remote web interface
	EnabledAPISets string
	// Comma separated list of API sets disabled on the remote web interface
	DisabledAPISets string
	// Enable all of API sets. Applies before disabling individual sets
	EnableAllAPISets bool

	enabledAPISets map[string]struct{}
	// Comma separate list of hostnames to accept in the Host header, used to bypass the Host header check which only applies to localhost addresses
	HostWhitelist string
	hostWhitelist []string

	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	// gnet uses this for TCP incoming and outgoing
	Port int
	// Maximum outgoing connections to maintain
	MaxOutgoingConnections int
	// Maximum default outgoing connections
	MaxDefaultPeerOutgoingConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// PeerlistSize represents the maximum number of peers that the pex would maintain
	PeerlistSize int
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface bool
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string
	// Remote web interface certificate
	WebInterfaceCert string
	// Remote web interface key
	WebInterfaceKey string
	// Remote web interface HTTPS support
	WebInterfaceHTTPS bool
	// Remote web interface username and password
	WebInterfaceUsername string
	WebInterfacePassword string
	// Allow web interface auth without HTTPS
	WebInterfacePlaintextAuth bool

	// Enable the deprecated JSON 2.0 RPC interface
	RPCInterface bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string
	// GUI directory contains assets for the HTML interface
	GUIDirectory string

	// Timeouts for the HTTP listener
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	// Remark to include in user agent sent in the wire protocol introduction
	UserAgentRemark string
	userAgent       useragent.Data

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Disable "Reply to ping", "Received pong" log messages
	DisablePingPong bool

	// Verify the database integrity after loading
	VerifyDB bool
	// Reset the database if integrity checks fail, and continue running
	ResetCorruptDB bool

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	// Wallet crypto type
	WalletCryptoType string

	// Disable the hardcoded default peers
	DisableDefaultPeers bool
	// Load custom peers from disk
	CustomPeersFile string

	RunMaster bool

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	DBPath      string
	DBReadOnly  bool
	Arbitrating bool
	LogToFile   bool
	Version     bool // show node version

	GenesisSignatureStr string
	GenesisAddressStr   string
	BlockchainPubkeyStr string
	BlockchainSeckeyStr string
	GenesisTimestamp    uint64
	GenesisCoinVolume   uint64
	DefaultConnections  []string

	genesisSignature cipher.Sig
	genesisAddress   cipher.Address

	blockchainPubkey cipher.PubKey
	blockchainSeckey cipher.SecKey
}

// NewNodeConfig returns a new node config instance
func NewNodeConfig(mode string, node NodeParameters) NodeConfig {
	nodeConfig := NodeConfig{
		CoinName:            node.CoinName,
		GenesisSignatureStr: node.GenesisSignatureStr,
		GenesisAddressStr:   node.GenesisAddressStr,
		GenesisCoinVolume:   node.GenesisCoinVolume,
		GenesisTimestamp:    node.GenesisTimestamp,
		BlockchainPubkeyStr: node.BlockchainPubkeyStr,
		BlockchainSeckeyStr: node.BlockchainSeckeyStr,
		DefaultConnections:  node.DefaultConnections,
		// Disable peer exchange
		DisablePEX: false,
		// Don't make any outgoing connections
		DisableOutgoingConnections: false,
		// Don't allowing incoming connections
		DisableIncomingConnections: false,
		// Disables networking altogether
		DisableNetworking: false,
		// Enable GUI
		EnableGUI: false,
		// Enable unversioned API
		EnableUnversionedAPI: false,
		// Disable CSRF check in the wallet API
		DisableCSRF: false,
		// DisableCSP disable content-security-policy in http response
		DisableCSP: false,
		// Only run on localhost and only connect to others on localhost
		LocalhostOnly: false,
		// Which address to serve on. Leave blank to automatically assign to a
		// public interface
		Address: "",
		//gnet uses this for TCP incoming and outgoing
		Port: node.Port,
		// MaxOutgoingConnections is the maximum outgoing connections allowed.
		MaxOutgoingConnections: 8,
		// MaxDefaultOutgoingConnections is the maximum default outgoing connections allowed.
		MaxDefaultPeerOutgoingConnections: 1,
		DownloadPeerList:                  true,
		PeerListURL:                       node.PeerListURL,
		// How often to make outgoing connections, in seconds
		OutgoingConnectionsRate: time.Second * 5,
		PeerlistSize:            65535,
		// Wallet Address Version
		//AddressVersion: "test",
		// Remote web interface
		WebInterface:      true,
		WebInterfacePort:  node.WebInterfacePort,
		WebInterfaceAddr:  "127.0.0.1",
		WebInterfaceCert:  "",
		WebInterfaceKey:   "",
		WebInterfaceHTTPS: false,
		EnabledAPISets:    api.EndpointsRead + "," + api.EndpointsTransaction,
		DisabledAPISets:   "",
		EnableAllAPISets:  false,

		RPCInterface: false,

		LaunchBrowser: false,
		// Data directory holds app data
		DataDirectory: node.DataDirectory,
		// Web GUI static resources
		GUIDirectory: "./src/gui/static/",
		// Logging
		ColorLog:        true,
		LogLevel:        "INFO",
		LogToFile:       false,
		DisablePingPong: false,

		VerifyDB:       false,
		ResetCorruptDB: false,

		// Wallets
		WalletDirectory:  "",
		WalletCryptoType: string(wallet.CryptoTypeScryptChacha20poly1305),

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		HTTPReadTimeout:  time.Second * 10,
		HTTPWriteTimeout: time.Second * 60,
		HTTPIdleTimeout:  time.Second * 120,

		// Centralized network configuration
		RunMaster: false,
		/* Developer options */

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",
	}

	nodeConfig.applyConfigMode(mode)

	return nodeConfig
}

func (c *Config) postProcess() error {
	if help {
		flag.Usage()
		fmt.Println("Additional environment variables:")
		fmt.Println("* COINHOUR_BURN_FACTOR - Set the coin hour burn factor required for transactions. Must be > 1.")
		os.Exit(0)
	}

	var err error
	if c.Node.GenesisSignatureStr != "" {
		c.Node.genesisSignature, err = cipher.SigFromHex(c.Node.GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}

	if c.Node.GenesisAddressStr != "" {
		c.Node.genesisAddress, err = cipher.DecodeBase58Address(c.Node.GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}
	if c.Node.BlockchainPubkeyStr != "" {
		c.Node.blockchainPubkey, err = cipher.PubKeyFromHex(c.Node.BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey, err = cipher.SecKeyFromHex(c.Node.BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		c.Node.BlockchainSeckeyStr = ""
	}
	if c.Node.BlockchainSeckeyStr != "" {
		c.Node.blockchainSeckey = cipher.SecKey{}
	}

	home := file.UserHome()
	c.Node.DataDirectory, err = file.InitDataDir(replaceHome(c.Node.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.Node.WebInterfaceCert == "" {
		c.Node.WebInterfaceCert = filepath.Join(c.Node.DataDirectory, "skycoind.cert")
	} else {
		c.Node.WebInterfaceCert = replaceHome(c.Node.WebInterfaceCert, home)
	}

	if c.Node.WebInterfaceKey == "" {
		c.Node.WebInterfaceKey = filepath.Join(c.Node.DataDirectory, "skycoind.key")
	} else {
		c.Node.WebInterfaceKey = replaceHome(c.Node.WebInterfaceKey, home)
	}

	if c.Node.WalletDirectory == "" {
		c.Node.WalletDirectory = filepath.Join(c.Node.DataDirectory, "wallets")
	} else {
		c.Node.WalletDirectory = replaceHome(c.Node.WalletDirectory, home)
	}

	if c.Node.DBPath == "" {
		c.Node.DBPath = filepath.Join(c.Node.DataDirectory, "data.db")
	} else {
		c.Node.DBPath = replaceHome(c.Node.DBPath, home)
	}

	if c.Node.RunMaster {
		// Run in arbitrating mode if the node is master
		c.Node.Arbitrating = true
	}

	userAgentData := useragent.Data{
		Coin:    c.Node.CoinName,
		Version: c.Build.Version,
		Remark:  c.Node.UserAgentRemark,
	}

	if _, err := userAgentData.Build(); err != nil {
		return err
	}

	c.Node.userAgent = userAgentData

	apiSets, err := buildAPISets(c.Node)
	if err != nil {
		return err
	}

	// Don't open browser to load wallets if wallet apis are disabled.
	c.Node.enabledAPISets = apiSets
	if _, ok := c.Node.enabledAPISets[api.EndpointsWallet]; !ok {
		c.Node.EnableGUI = false
		c.Node.LaunchBrowser = false
	}

	if c.Node.EnableGUI {
		c.Node.GUIDirectory = file.ResolveResourceDirectory(c.Node.GUIDirectory)
	}

	if c.Node.DisableDefaultPeers {
		c.Node.DefaultConnections = nil
	}

	if c.Node.HostWhitelist != "" {
		c.Node.hostWhitelist = strings.Split(c.Node.HostWhitelist, ",")
	}

	httpAuthEnabled := c.Node.WebInterfaceUsername != "" || c.Node.WebInterfacePassword != ""
	if httpAuthEnabled && !c.Node.WebInterfaceHTTPS && !c.Node.WebInterfacePlaintextAuth {
		return errors.New("Web interface auth enabled but HTTPS is not enabled. Use -web-interface-plaintext-auth=true if this is desired")
	}

	return nil
}

// buildAPISets builds the set of enable APIs by the following rules:
// * If EnableAll, all API sets are added
// * For each api set in EnabledAPISets, add
// * For each api set in DisabledAPISets, remove
func buildAPISets(c NodeConfig) (map[string]struct{}, error) {
	enabledAPISets := strings.Split(c.EnabledAPISets, ",")
	if err := validateAPISets("-enable-api-sets", enabledAPISets); err != nil {
		return nil, err
	}

	disabledAPISets := strings.Split(c.DisabledAPISets, ",")
	if err := validateAPISets("-disable-api-sets", disabledAPISets); err != nil {
		return nil, err
	}

	apiSets := make(map[string]struct{})

	allAPISets := []string{
		api.EndpointsRead,
		api.EndpointsStatus,
		api.EndpointsWallet,
		api.EndpointsTransaction,
		// Do not include insecure or deprecated API sets, they must always
		// be explicitly enabled through -enable-api-sets
	}

	if c.EnableAllAPISets {
		for _, s := range allAPISets {
			apiSets[s] = struct{}{}
		}
	}

	// Add the enabled API sets
	for _, k := range enabledAPISets {
		apiSets[k] = struct{}{}
	}

	// Remove the disabled API sets
	for _, k := range disabledAPISets {
		delete(apiSets, k)
	}

	return apiSets, nil
}

func validateAPISets(opt string, apiSets []string) error {
	for _, k := range apiSets {
		k = strings.ToUpper(strings.TrimSpace(k))
		switch k {
		case api.EndpointsRead,
			api.EndpointsStatus,
			api.EndpointsTransaction,
			api.EndpointsWallet,
			api.EndpointsInsecureWalletSeed,
			api.EndpointsDeprecatedWalletSpend:
		case "":
			continue
		default:
			return fmt.Errorf("Invalid value in %s: %q", opt, k)
		}
	}
	return nil
}

// RegisterFlags binds CLI flags to config values
func (c *NodeConfig) RegisterFlags() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX, "disable PEX peer discovery")
	flag.BoolVar(&c.DownloadPeerList, "download-peerlist", c.DownloadPeerList, "download a peers.txt from -peerlist-url")
	flag.StringVar(&c.PeerListURL, "peerlist-url", c.PeerListURL, "with -download-peerlist=true, download a peers.txt file from this url")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing", c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming", c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking", c.DisableNetworking, "Disable all network activity")
	flag.BoolVar(&c.EnableGUI, "enable-gui", c.EnableGUI, "Enable GUI")
	flag.BoolVar(&c.EnableUnversionedAPI, "enable-unversioned-api", c.EnableUnversionedAPI, "Enable the deprecated unversioned API endpoints without /api/v1 prefix")
	flag.BoolVar(&c.DisableCSRF, "disable-csrf", c.DisableCSRF, "disable CSRF check")
	flag.BoolVar(&c.DisableCSP, "disable-csp", c.DisableCSP, "disable content-security-policy in http response")
	flag.StringVar(&c.Address, "address", c.Address, "IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")

	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert", c.WebInterfaceCert, "skycoind.cert file for web interface HTTPS. If not provided, will autogenerate or use skycoind.cert in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key", c.WebInterfaceKey, "skycoind.key file for web interface HTTPS. If not provided, will autogenerate or use skycoind.key in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https", c.WebInterfaceHTTPS, "enable HTTPS for web interface")
	flag.StringVar(&c.HostWhitelist, "host-whitelist", c.HostWhitelist, "Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.")
	flag.StringVar(&c.EnabledAPISets, "enable-api-sets", c.EnabledAPISets, "enable API set. Options are READ, STATUS, WALLET, INSECURE_WALLET_SEED, DEPRECATED_WALLET_SPEND. Multiple values should be separated by comma")
	flag.StringVar(&c.DisabledAPISets, "disable-api-sets", c.DisabledAPISets, "disable API set. Options are READ, STATUS, WALLET, INSECURE_WALLET_SEED, DEPRECATED_WALLET_SPEND. Multiple values should be separated by comma")
	flag.BoolVar(&c.EnableAllAPISets, "enable-all-api-sets", c.EnableAllAPISets, "enable all API sets, except for deprecated or insecure sets. This option is applied before -disable-api-sets.")
	flag.StringVar(&c.WebInterfaceUsername, "web-interface-username", c.WebInterfaceUsername, "username for the web interface")
	flag.StringVar(&c.WebInterfacePassword, "web-interface-password", c.WebInterfacePassword, "password for the web interface")
	flag.BoolVar(&c.WebInterfacePlaintextAuth, "web-interface-plaintext-auth", c.WebInterfacePlaintextAuth, "allow web interface auth without https")

	flag.BoolVar(&c.RPCInterface, "rpc-interface", c.RPCInterface, "enable the deprecated JSON 2.0 RPC interface")

	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address", c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
	flag.StringVar(&c.DBPath, "db-path", c.DBPath, "path of database file (defaults to ~/.skycoin/data.db)")
	flag.BoolVar(&c.DBReadOnly, "db-read-only", c.DBReadOnly, "open bolt db read-only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	flag.StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.BoolVar(&c.DisablePingPong, "no-ping-log", c.DisablePingPong, `disable "reply to ping" and "received pong" debug log messages`)
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory, "static content directory for the HTML interface")

	flag.BoolVar(&c.VerifyDB, "verify-db", c.VerifyDB, "check the database for corruption")
	flag.BoolVar(&c.ResetCorruptDB, "reset-corrupt-db", c.ResetCorruptDB, "reset the database if corrupted, and continue running instead of exiting")

	flag.BoolVar(&c.DisableDefaultPeers, "disable-default-peers", c.DisableDefaultPeers, "disable the hardcoded default peers")
	flag.StringVar(&c.CustomPeersFile, "custom-peers-file", c.CustomPeersFile, "load custom peers from a newline separate list of ip:port in a file. Note that this is different from the peers.json file in the data directory")

	flag.StringVar(&c.UserAgentRemark, "user-agent-remark", c.UserAgentRemark, "additional remark to include in the user agent sent over the wire protocol")

	// Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster, "run the daemon as blockchain master server")

	flag.StringVar(&c.BlockchainPubkeyStr, "master-public-key", c.BlockchainPubkeyStr, "public key of the master chain")
	flag.StringVar(&c.BlockchainSeckeyStr, "master-secret-key", c.BlockchainSeckeyStr, "secret key, set for master")

	flag.StringVar(&c.GenesisAddressStr, "genesis-address", c.GenesisAddressStr, "genesis address")
	flag.StringVar(&c.GenesisSignatureStr, "genesis-signature", c.GenesisSignatureStr, "genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp, "genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory, "location of the wallet files. Defaults to ~/.skycoin/wallet/")
	flag.IntVar(&c.MaxOutgoingConnections, "max-outgoing-connections", c.MaxOutgoingConnections, "The maximum outgoing connections allowed")
	flag.IntVar(&c.MaxDefaultPeerOutgoingConnections, "max-default-peer-outgoing-connections", c.MaxDefaultPeerOutgoingConnections, "The maximum default peer outgoing connections allowed")
	flag.IntVar(&c.PeerlistSize, "peerlist-size", c.PeerlistSize, "The peer list size")
	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate", c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly, "Run on localhost and only connect to localhost peers")
	flag.BoolVar(&c.Arbitrating, "arbitrating", c.Arbitrating, "Run node in arbitrating mode")
	flag.StringVar(&c.WalletCryptoType, "wallet-crypto-type", c.WalletCryptoType, "wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305")
	flag.BoolVar(&c.Version, "version", false, "show node version")
}

func (c *NodeConfig) applyConfigMode(configMode string) {
	if runtime.GOOS == "windows" {
		c.ColorLog = false
	}
	switch configMode {
	case "":
	case "STANDALONE_CLIENT":
		c.EnableAllAPISets = true
		c.EnabledAPISets = "INSECURE_WALLET_SEED"
		c.EnableGUI = true
		c.LaunchBrowser = true
		c.DisableCSRF = false
		c.DisableCSP = false
		c.DownloadPeerList = true
		c.RPCInterface = false
		c.WebInterface = true
		c.LogToFile = false
		c.ResetCorruptDB = true
		c.WebInterfacePort = 0 // randomize web interface port
	default:
		panic("Invalid ConfigMode")
	}
}

func panicIfError(err error, msg string, args ...interface{}) { // nolint: unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
