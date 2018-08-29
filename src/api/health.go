package api

import (
	"fmt"
	"net/http"
	"time"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
)

// BlockchainMetadata extends visor.BlockchainMetadata to include the time since the last block
type BlockchainMetadata struct {
	*visor.BlockchainMetadata
	TimeSinceLastBlock wh.Duration `json:"time_since_last_block"`
}

// HealthResponse is returned by the /health endpoint
type HealthResponse struct {
	BlockchainMetadata    BlockchainMetadata `json:"blockchain"`
	Version               visor.BuildInfo    `json:"version"`
	OpenConnections       int                `json:"open_connections"`
	Uptime                wh.Duration        `json:"uptime"`
	CSRFEnabled           bool               `json:"csrf_enabled"`
	CSPEnabled            bool               `json:"csp_enabled"`
	WalletAPIEnabled      bool               `json:"wallet_api_enabled"`
	GUIEnabled            bool               `json:"gui_enabled"`
	UnversionedAPIEnabled bool               `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool               `json:"json_rpc_enabled"`
}

// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, csrfStore *CSRFStore, gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		health, err := gateway.GetHealth()
		if err != nil {
			err = fmt.Errorf("gateway.GetHealth failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		elapsedBlockTime := time.Now().UTC().Unix() - int64(health.BlockchainMetadata.Head.Time)
		timeSinceLastBlock := time.Second * time.Duration(elapsedBlockTime)

		wh.SendJSONOr500(logger, w, HealthResponse{
			BlockchainMetadata: BlockchainMetadata{
				BlockchainMetadata: health.BlockchainMetadata,
				TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
			},
			Version:               health.Version,
			OpenConnections:       health.OpenConnections,
			Uptime:                wh.FromDuration(health.Uptime),
			CSRFEnabled:           csrfStore.Enabled,
			CSPEnabled:            !c.disableCSP,
			UnversionedAPIEnabled: c.enableUnversionedAPI,
			GUIEnabled:            c.enableGUI,
			JSON20RPCEnabled:      c.enableJSON20RPC,
			WalletAPIEnabled:      gateway.IsWalletAPIEnabled(),
		})
	}
}
