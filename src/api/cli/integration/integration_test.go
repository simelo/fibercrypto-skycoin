package cli_integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"io"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	binaryName = "skycoin-cli"
	walletName = "integration_test.wlt"

	testModeStable = "stable"
	testModeLive   = "live"

	// Number of random transactions of live transaction test.
	randomLiveTransactionNum = 500
)

var (
	binaryPath string
	walletDir  string
)

var (
	update     = flag.Bool("update", false, "update golden files")
	liveTxFull = flag.Bool("live-tx-full", false, "run live transaction test against full blockchain")
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestGenerateAddresses(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "generateAddresses").CombinedOutput()
	require.NoError(t, err)
	o := strings.Trim(string(output), "\n")
	require.Equal(t, "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC", o)

	wltPath := filepath.Join(walletDir, walletName)
	var w wallet.ReadableWallet
	loadJSON(t, wltPath, &w)

	golden := filepath.Join("testdata", "generateAddresses.golden")
	if *update {
		writeJSON(t, golden, w)
	}

	var expect wallet.ReadableWallet
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, w)
}

func TestVerifyAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	tt := []struct {
		name   string
		addr   string
		err    error
		errMsg string
	}{
		{
			"valid skycoin address",
			"2Kg3eRXUhY6hrDZvNGB99DKahtrPDQ1W9vN",
			nil,
			"",
		},
		{
			"invalid skycoin address",
			"2KG9eRXUhx6hrDZvNGB99DKahtrPDQ1W9vn",
			errors.New("exit status 1"),
			"Invalid version",
		},
		{
			"invalid bitcoin address",
			"1Dcb9gpaZpBKmjqjCsiBsP3sBW1md2kEM2",
			errors.New("exit status 1"),
			"Invalid version",
		},
	}

	for _, tc := range tt {
		output, err := exec.Command(binaryPath, "verifyAddress", tc.addr).CombinedOutput()
		if err != nil {
			require.Equal(t, tc.err.Error(), err.Error())
			require.Equal(t, tc.errMsg, strings.Trim(string(output), "\n"))
			return
		}

		require.Empty(t, output)
	}
}

func TestDecodeRawTransaction(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	rawTx := `2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1ccb6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a00000000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e40d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0fc9b01000000001b12000000000000`
	output, err := exec.Command(binaryPath, "decodeRawTransaction", rawTx).CombinedOutput()
	require.NoError(t, err)

	var txn visor.TransactionJSON
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&txn)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "decodeRawTransaction.golden")
	if *update {
		writeJSON(t, golden, txn)
	}

	var expect visor.TransactionJSON
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, txn)
}

func TestAddressGen(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressGen").CombinedOutput()
	require.NoError(t, err)

	var wltAddress wallet.ReadableWallet
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddress)
	require.NoError(t, err)
}

func TestStableListWallets(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "listWallets").CombinedOutput()
	require.NoError(t, err)

	var wlts struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "listWallets.golden")
	if *update {
		writeJSON(t, golden, wlts)
	}

	var expect struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, wlts)
}

func TestLiveListWallets(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "listWallets").CombinedOutput()
	require.NoError(t, err)

	var wlts struct {
		Wallets []cli.WalletEntry `json:"wallets"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wlts)
	require.NoError(t, err)
}

func TestStableListAddress(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "listAddresses").CombinedOutput()
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "listAddresses.golden")
	if *update {
		writeJSON(t, golden, wltAddresses)
	}

	var expect struct {
		Addresses []string `json:"addresses"`
	}
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, wltAddresses)
}

func TestLiveListAddresses(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "listAddresses").CombinedOutput()
	require.NoError(t, err)

	var wltAddresses struct {
		Addresses []string `json:"addresses"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltAddresses)
	require.NoError(t, err)
}

func TestStableAddressBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "addressBalance.golden")
	if *update {
		writeJSON(t, golden, addrBalance)
	}

	var expect cli.BalanceResult
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, addrBalance)
}

func TestLiveAddressBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressBalance", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrBalance)
	require.NoError(t, err)
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "walletBalance").CombinedOutput()
	require.NoError(t, err)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "walletBalance.golden")
	if *update {
		writeJSON(t, golden, wltBalance)
	}

	var expect cli.BalanceResult
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, wltBalance)
}

func TestLiveWalletBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "walletBalance").CombinedOutput()
	require.NoError(t, err)

	var wltBalance cli.BalanceResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltBalance)
	require.NoError(t, err)
}

func TestStableWalletOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "walletOutputs").CombinedOutput()
	require.NoError(t, err)

	var wltOutput webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "walletOutputs.golden")
	if *update {
		writeJSON(t, golden, wltOutput)
	}

	var expect webrpc.OutputsResult
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, wltOutput)
}

func TestLiveWalletOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "walletOutputs").CombinedOutput()
	require.NoError(t, err)

	var wltOutput webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&wltOutput)
	require.NoError(t, err)
}

func TestStableAddressOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrOutputs webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
	require.NoError(t, err)

	golden := filepath.Join("testdata", "addressOutputs.golden")
	if *update {
		writeJSON(t, golden, addrOutputs)
	}

	var expect webrpc.OutputsResult
	loadJSON(t, golden, &expect)
	require.Equal(t, expect, addrOutputs)
}

func TestLiveAddressOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "addressOutputs", "2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt").CombinedOutput()
	require.NoError(t, err)

	var addrOutputs webrpc.OutputsResult
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&addrOutputs)
	require.NoError(t, err)
}

func TestStableStatus(t *testing.T) {
	if !doStable(t) {
		return
	}

	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)
	var ret struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)

	// TimeSinceLastBlock is not stable
	ret.TimeSinceLastBlock = ""

	var expect struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	golden := filepath.Join("testdata", "status.golden")
	if *update {
		writeJSON(t, golden, ret)
	}

	loadJSON(t, golden, &expect)
	require.Equal(t, expect, ret)
}

func TestLiveStatus(t *testing.T) {
	if !doLive(t) {
		return
	}

	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)

	var ret struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}

	err = json.NewDecoder(bytes.NewReader(output)).Decode(&ret)
	require.NoError(t, err)
	require.True(t, ret.Running)
	require.Equal(t, ret.RPCAddress, rpcAddress())
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		args       []string
		err        error
		errMsg     string
		goldenFile string
	}{
		{
			"invalid txid",
			[]string{"abcd"},
			errors.New("exit status 1"),
			"invalid txid\n",
			"",
		},
		{
			"not exist",
			[]string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			errors.New("exit status 1"),
			"transaction doesn't exist [code: -32600]\n",
			"",
		},
		{
			"empty txid",
			[]string{""},
			errors.New("exit status 1"),
			"txid is empty\n",
			"",
		},
		{
			"genesis transaction",
			[]string{"d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add"},
			nil,
			"",
			"./testdata/genesisTransaction.golden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"transaction"}, tc.args...)
			o, err := exec.Command(binaryPath, args...).CombinedOutput()
			if err != nil {
				require.Equal(t, tc.err.Error(), err.Error())
				require.Equal(t, tc.errMsg, string(o))
				return
			}

			// Decode the output into visor.TransactionJSON
			var tx webrpc.TxnResult
			err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
			require.NoError(t, err)

			if tc.goldenFile != "" && *update {
				writeJSON(t, tc.goldenFile, tx)
			}
			var expect webrpc.TxnResult
			loadJSON(t, tc.goldenFile, &expect)

			require.Equal(t, expect, tx)
		})
	}

	scanTransactions(t, true)
}

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	o, err := exec.Command(binaryPath, "transaction", "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add").CombinedOutput()
	require.NoError(t, err)
	var tx webrpc.TxnResult
	err = json.NewDecoder(bytes.NewReader(o)).Decode(&tx)
	require.NoError(t, err)

	var expect webrpc.TxnResult

	golden := filepath.Join("testdata", "genesisTransaction.golden")
	if *update {
		writeJSON(t, golden, tx)
	}

	loadJSON(t, golden, &expect)
	require.Equal(t, expect.Transaction.Transaction, tx.Transaction.Transaction)

	scanTransactions(t, *liveTxFull)

	// scan pending transactions
	scanPendingTransactions(t)
}

func scanPendingTransactions(t *testing.T) {
}

// scanTransactions scans transactions against blockchain.
// If fullTest is true, scan the whole blockchain, and test every transactions,
// otherwise just test random transactions.
func scanTransactions(t *testing.T, fullTest bool) {
	// Gets blockchain height through "status" command
	output, err := exec.Command(binaryPath, "status").CombinedOutput()
	require.NoError(t, err)
	var status struct {
		webrpc.StatusResult
		RPCAddress string `json:"webrpc_address"`
	}
	err = json.NewDecoder(bytes.NewReader(output)).Decode(&status)
	require.NoError(t, err)

	txids := getTxids(t, status.BlockNum)

	l := len(txids)
	if !fullTest && l > randomLiveTransactionNum {
		txidMap := make(map[string]struct{})
		var ids []string
		for len(txidMap) < randomLiveTransactionNum {
			// get random txid
			txid := txids[rand.Intn(l)]
			if _, ok := txidMap[txid]; !ok {
				ids = append(ids, txid)
				txidMap[txid] = struct{}{}
			}
		}

		// reassign the txids
		txids = ids
	}

	checkTransctions(t, txids)
}

func checkTransctions(t *testing.T, txids []string) {
	// Start goroutines to check transactions
	var wg sync.WaitGroup
	txC := make(chan string, 500)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case txid, ok := <-txC:
					if !ok {
						return
					}

					t.Run(fmt.Sprintf("%v", txid), func(t *testing.T) {
						o, err := exec.Command(binaryPath, "transaction", txid).CombinedOutput()
						require.NoError(t, err)
						var txRlt webrpc.TxnResult
						err = json.NewDecoder(bytes.NewReader(o)).Decode(&txRlt)
						require.NoError(t, err)
						require.Equal(t, txid, txRlt.Transaction.Transaction.Hash)
						require.True(t, txRlt.Transaction.Status.Confirmed)
					})
				}
			}
		}()
	}

	for _, txid := range txids {
		txC <- txid
	}
	close(txC)

	wg.Wait()
}

func getTxids(t *testing.T, blockNum uint64) []string {
	// p represents the number of blocks that each time we query,
	// do not get all blocks in one query, which might run out of
	// memory when blockchain becomes very huge.
	p := 500
	n := int(blockNum / uint64(p))

	// Collects all transactions' id
	var txids []string
	for i := 0; i < int(n); i++ {
		txids = append(txids, getTxidsInBlocks(t, i*p+1, (i+1)*p)...)
	}

	if (blockNum % uint64(p)) > 0 {
		txids = append(txids, getTxidsInBlocks(t, n*p+1, int(blockNum)-1)...)
	}

	return txids
}

func getTxidsInBlocks(t *testing.T, start, end int) []string {
	s := strconv.Itoa(start)
	e := strconv.Itoa(end)
	o, err := exec.Command(binaryPath, "blocks", s, e).CombinedOutput()
	require.NoError(t, err)
	var blocks visor.ReadableBlocks
	err = json.NewDecoder(bytes.NewReader(o)).Decode(&blocks)
	require.NoError(t, err)
	require.Len(t, blocks.Blocks, end-start+1)

	var txids []string
	for _, b := range blocks.Blocks {
		for _, tx := range b.Body.Transactions {
			txids = append(txids, tx.Hash)
		}
	}
	return txids
}

// Do setup and teardown here.
func TestMain(m *testing.M) {
	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("get binary name absolute path failed: %v\n", err))
		os.Exit(1)
	}

	binaryPath = abs

	// Build cli binary file.
	args := []string{"build", "-o", binaryPath, "../../../../cmd/cli/cli.go"}
	if err := exec.Command("go", args...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Make %v binary failed: %v\n", binaryName, err))
		os.Exit(1)
	}

	dir, clean, err := createTempWalletFile(filepath.Join("testdata", "integration_test.wlt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer clean()

	walletDir = dir

	os.Setenv("WALLET_DIR", dir)
	os.Setenv("WALLET_NAME", walletName)

	ret := m.Run()

	// Remove the generated cli binary file.
	if err := os.Remove(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Delete %v failed: %v", binaryName, err))
		os.Exit(1)
	}

	os.Exit(ret)
}

// createTempWalletFile creates a temporary dir, and copy the 'from' file to dir.
// returns the temporary dir path, cleanup callback function, and error if any.
func createTempWalletFile(from string) (string, func(), error) {
	dir, err := ioutil.TempDir("", "integration_test")
	if err != nil {
		return "", nil, fmt.Errorf("Get temporary dir failed: %v", err)
	}

	// Copy the  the temporary dir.
	wltPath := filepath.Join(dir, walletName)
	f, err := os.Create(wltPath)
	if err != nil {
		return "", nil, fmt.Errorf("Create temporary file: %v failed: %v", wltPath, err)
	}

	defer f.Close()

	rf, err := os.Open(from)
	if err != nil {
		return "", nil, fmt.Errorf("Open %v failed: %v", from, err)
	}

	defer rf.Close()
	io.Copy(f, rf)

	fun := func() {
		// Delete the temporary dir
		os.RemoveAll(dir)
	}

	return dir, fun, nil
}

func loadJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(obj)
	require.NoError(t, err)
}

func writeJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	require.NoError(t, enc.Encode(obj))
}

func mode(t *testing.T) string {
	mode := os.Getenv("SKYCOIN_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive, testModeStable:
	default:
		t.Fatal("Invalid test mode, must be stable or live")
	}
	return mode
}

func enabled() bool {
	return os.Getenv("SKYCOIN_INTEGRATION_TESTS") == "1"
}

func doStable(t *testing.T) bool {
	if enabled() && mode(t) == testModeStable {
		return true
	}

	t.Skip("Stable tests disabled")
	return false
}

func doLive(t *testing.T) bool {
	if enabled() && mode(t) == testModeLive {
		return true
	}

	t.Skip("Live tests disabled")
	return false
}

func doLiveOrStable(t *testing.T) bool {
	if enabled() {
		switch mode(t) {
		case testModeStable, testModeLive:
			return true
		}
	}

	t.Skip("Live and stable tests disabled")
	return false
}

func rpcAddress() string {
	rpcAddr := os.Getenv("RPC_ADDR")
	if rpcAddr == "" {
		rpcAddr = "127.0.0.1:6430"
	}

	return rpcAddr
}
