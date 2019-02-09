package api

// APIs for network-related information

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// connectionHandler returns a specific connection
// URI: /api/v1/network/connection
// Method: GET
// Args:
//	addr - An IP:Port string
func connectionHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/connection networkConnection
	//
	// This endpoint returns a specific connection.
	//
	// ---
	//
	//
	// produces:
	// - application/json
	// parameters:
	// - name: addr
	//   in: query
	//   description: Address port
	//   required: true
	//   type: string
	//
	// security:
	// - csrfAuth: []
	//
	// responses:
	//   200:
	//     description: This endpoint return a connection struct
	//     schema:
	//       type: object
	//       properties:
	//         id:
	//           type: integer
	//           format: int64
	//         address:
	//           type: string
	//         state:
	//           type: string
	//           enum: [pending, connected, introduced]
	//         user_agent:
	//           type: string
	//         last_sent:
	//           type: integer
	//           format: int64
	//         last_received:
	//           type: integer
	//           format: int64
	//         connected_at:
	//           type: integer
	//           format: int64
	//         outgoing:
	//           type: boolean
	//         mirror:
	//           type: integer
	//           format: int32
	//         listen_port:
	//           type: integer
	//           format: int32
	//         height:
	//           type: integer
	//           format: int64
	//         is_trusted_peer:
	//           type: boolean
	//         unconfirmed_verify_transaction:
	//           description: Represent unconfirmed transactions
	//           type: object
	//           properties:
	//             burn_factor:
	//               type: integer
	//               format: int32
	//             max_transaction_size:
	//               type: integer
	//               format: int32
	//             max_decimals:
	//               type: integer
	//               format: int32
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("addr")
		if addr == "" {
			wh.Error400(w, "addr is required")
			return
		}

		c, err := gateway.GetConnection(addr)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		if c == nil {
			wh.Error404(w, "")
			return
		}

		wh.SendJSONOr500(logger, w, readable.NewConnection(c))
	}
}

// Connections wraps []Connection
type Connections struct {
	Connections []readable.Connection `json:"connections"`
}

// NewConnections copies []daemon.Connection to a struct with json tags
func NewConnections(dconns []daemon.Connection) Connections {
	conns := make([]readable.Connection, len(dconns))
	for i, dc := range dconns {
		conns[i] = readable.NewConnection(&dc)
	}

	return Connections{
		Connections: conns,
	}
}

// connectionsHandler returns all outgoing connections
// URI: /api/v1/network/connections
// Method: GET
// Args:
//	states: [optional] comma-separated list of connection states ("pending", "connected" or "introduced"). Defaults to "connected,introduced"
//  direction: [optional] "outgoing" or "incoming". If not provided, both are included.
func connectionsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/connections networkConnections
	//
	// This endpoint returns all outgoings connections.
	//
	// ---
	// produces:
	// - application/json
	//
	// security:
	// - csrfAuth: []
	//
	// parameters:
	// - name: states
	//   in: query
	//   description: Connection status.
	//   required: false
	//   enum: [pending, connected, introduced]
	//   type: string
	// - name: direction
	//   in: query
	//   description: Direction of the connection.
	//   required: false
	//   enum: [connected, introduced]
	//   type: string
	// responses:
	//   200:
	//     description: This endpoint return a connection struct
	//     schema:
	//       type: array
	//       items:
	//         properties:
	//           id:
	//             type: integer
	//             format: int64
	//           address:
	//             type: string
	//           state:
	//             type: string
	//             enum: [pending, connected, introduced]
	//           user_agent:
	//             type: string
	//           last_sent:
	//             type: integer
	//             format: int64
	//           last_received:
	//             type: integer
	//             format: int64
	//           connected_at:
	//             type: integer
	//             format: int64
	//           outgoing:
	//             type: boolean
	//           mirror:
	//             type: integer
	//             format: int32
	//           listen_port:
	//             type: integer
	//             format: int32
	//           height:
	//             type: integer
	//             format: int64
	//           is_trusted_peer:
	//             type: boolean
	//           unconfirmed_verify_transaction:
	//             description: Represent unconfirmed transactions
	//             type: object
	//             properties:
	//               burn_factor:
	//                 type: integer
	//                 format: int32
	//               max_transaction_size:
	//                 type: integer
	//                 format: int32
	//               max_decimals:
	//                 type: integer
	//                 format: int32
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		formStates := r.FormValue("states")
		statesMap := make(map[daemon.ConnectionState]struct{}, 3)
		if formStates != "" {
			states := strings.Split(formStates, ",")
			for _, s := range states {
				switch daemon.ConnectionState(s) {
				case daemon.ConnectionStatePending,
					daemon.ConnectionStateConnected,
					daemon.ConnectionStateIntroduced:
					statesMap[daemon.ConnectionState(s)] = struct{}{}
				default:
					wh.Error400(w, fmt.Sprintf("Invalid state in states. Valid states are %q, %q or %q", daemon.ConnectionStatePending, daemon.ConnectionStateConnected, daemon.ConnectionStateIntroduced))
					return
				}
			}
		}

		// "connected" and "introduced" are the defaults, if not specified
		if len(statesMap) == 0 {
			statesMap[daemon.ConnectionStateConnected] = struct{}{}
			statesMap[daemon.ConnectionStateIntroduced] = struct{}{}
		}

		direction := r.FormValue("direction")
		switch direction {
		case "incoming", "outgoing", "":
		default:
			wh.Error400(w, "Invalid direction. Valid directions are \"outgoing\" or \"incoming\"")
			return
		}

		conns, err := gateway.GetConnections(func(c daemon.Connection) bool {
			switch direction {
			case "outgoing":
				if !c.Outgoing {
					return false
				}
			case "incoming":
				if c.Outgoing {
					return false
				}
			}

			if _, ok := statesMap[c.State]; !ok {
				return false
			}

			return true
		})

		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, NewConnections(conns))
	}
}

// defaultConnectionsHandler returns the list of default hardcoded bootstrap addresses.
// They are not necessarily connected to.
// URI: /api/v1/network/defaultConnections
// Method: GET
func defaultConnectionsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/defaultConnections defaultConnections
	//
	// defaultConnectionsHandler returns the list of default hardcoded bootstrap addresses.\n They are not necessarily connected to.
	//
	// ---
	// produces:
	// - application/json
	//
	// security:
	// - csrfAuth: []
	//
	// responses:
	//   200:
	//     description: This endpoint return an list of default connections.
	//     schema:
	//       type: array
	//       items:
	//         type: string
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetDefaultConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

// trustConnectionsHandler returns all trusted connections
// In the default configuration, these will be a subset of the default hardcoded bootstrap addresses
// URI: /api/v1/network/trust
// Method: GET
func trustConnectionsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/connections/trust networkConnectionsTrust
	//
	// trustConnectionsHandler returns all trusted connections.\n They are not necessarily connected to. In the default configuration, these will be a subset of the default hardcoded bootstrap addresses.
	//
	// ---
	// produces:
	// - application/json
	//
	// security:
	// - csrfAuth: []
	//
	// responses:
	//   200:
	//     description: This endpoint return a list of trusted connections.
	//     schema:
	//       type: array
	//       items:
	//         type: string
	//   default:
	//     $ref: '#/responses/genericError'

	// TODO Add URI: /api/v1/network/trust to README.md or Most be URI: /api/v1/network/connections/trust

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetTrustConnections()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

// exchgConnectionsHandler returns all connections found through peer exchange
// URI: /api/v1/network/exchange
// Method: GET
func exchgConnectionsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/connections/exchange networkConnectionsExchange
	//
	// This endpoint returns all connections found through peer exchange
	//
	// ---
	// produces:
	// - application/json
	//
	// security:
	// - csrfAuth: []
	//
	// responses:
	//   200:
	//     description: This endpoint return a list of all connections found through peer exchange.
	//     schema:
	//       type: array
	//       items:
	//         type: string
	//   default:
	//     $ref: '#/responses/genericError'

	// TODO Add URI: /api/v1/network/exchange to README.md or Most be URI: /api/v1/network/connections/exchange

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		conns := gateway.GetExchgConnection()
		sort.Strings(conns)

		wh.SendJSONOr500(logger, w, conns)
	}
}

// disconnectHandler disconnects a connection by ID or address
// URI: /api/v1/network/connection/disconnect
// Method: POST
// Args:
//	id: ID of the connection
func disconnectHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/network/connection/disconnect networkConnectionsDisconnect
	//
	// This endpoint disconnects a connection by ID or address
	//
	// ---
	//
	// parameters:
	// - name: id
	//   in: query
	//   description: Address id.
	//   required: true
	//   type: string
	//
	// security:
	// - csrfAuth: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   200:
	//     description: This endpoint
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		formID := r.FormValue("id")
		if formID == "" {
			wh.Error400(w, "id is required")
			return
		}

		id, err := strconv.ParseUint(formID, 10, 64)
		if err != nil || id == 0 { // gnet IDs are non-zero
			wh.Error400(w, "invalid id")
			return
		}

		if err := gateway.Disconnect(uint64(id)); err != nil {
			switch err {
			case daemon.ErrConnectionNotExist:
				wh.Error404(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, struct{}{})
	}
}
