// Code generated by mockery v1.0.0. DO NOT EDIT.

package daemon

import cipher "github.com/skycoin/skycoin/src/cipher"
import coin "github.com/skycoin/skycoin/src/coin"
import gnet "github.com/skycoin/skycoin/src/daemon/gnet"
import mock "github.com/stretchr/testify/mock"
import pex "github.com/skycoin/skycoin/src/daemon/pex"
import useragent "github.com/skycoin/skycoin/src/util/useragent"
import visor "github.com/skycoin/skycoin/src/visor"

// MockDaemoner is an autogenerated mock type for the Daemoner type
type MockDaemoner struct {
	mock.Mock
}

// AddPeer provides a mock function with given fields: addr
func (_m *MockDaemoner) AddPeer(addr string) error {
	ret := _m.Called(addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddPeers provides a mock function with given fields: addrs
func (_m *MockDaemoner) AddPeers(addrs []string) int {
	ret := _m.Called(addrs)

	var r0 int
	if rf, ok := ret.Get(0).(func([]string) int); ok {
		r0 = rf(addrs)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// AnnounceAllTxns provides a mock function with given fields:
func (_m *MockDaemoner) AnnounceAllTxns() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BlockchainPubkey provides a mock function with given fields:
func (_m *MockDaemoner) BlockchainPubkey() cipher.PubKey {
	ret := _m.Called()

	var r0 cipher.PubKey
	if rf, ok := ret.Get(0).(func() cipher.PubKey); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cipher.PubKey)
		}
	}

	return r0
}

// BroadcastMessage provides a mock function with given fields: msg
func (_m *MockDaemoner) BroadcastMessage(msg gnet.Message) error {
	ret := _m.Called(msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(gnet.Message) error); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DaemonConfig provides a mock function with given fields:
func (_m *MockDaemoner) DaemonConfig() DaemonConfig {
	ret := _m.Called()

	var r0 DaemonConfig
	if rf, ok := ret.Get(0).(func() DaemonConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(DaemonConfig)
	}

	return r0
}

// Disconnect provides a mock function with given fields: addr, r
func (_m *MockDaemoner) Disconnect(addr string, r gnet.DisconnectReason) error {
	ret := _m.Called(addr, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, gnet.DisconnectReason) error); ok {
		r0 = rf(addr, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DisconnectNow provides a mock function with given fields: addr, r
func (_m *MockDaemoner) DisconnectNow(addr string, r gnet.DisconnectReason) error {
	ret := _m.Called(addr, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, gnet.DisconnectReason) error); ok {
		r0 = rf(addr, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExecuteSignedBlock provides a mock function with given fields: b
func (_m *MockDaemoner) ExecuteSignedBlock(b coin.SignedBlock) error {
	ret := _m.Called(b)

	var r0 error
	if rf, ok := ret.Get(0).(func(coin.SignedBlock) error); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMirrorPort provides a mock function with given fields: addr, mirror
func (_m *MockDaemoner) GetMirrorPort(addr string, mirror uint32) (uint16, bool) {
	ret := _m.Called(addr, mirror)

	var r0 uint16
	if rf, ok := ret.Get(0).(func(string, uint32) uint16); ok {
		r0 = rf(addr, mirror)
	} else {
		r0 = ret.Get(0).(uint16)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string, uint32) bool); ok {
		r1 = rf(addr, mirror)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetSignedBlocksSince provides a mock function with given fields: seq, count
func (_m *MockDaemoner) GetSignedBlocksSince(seq uint64, count uint64) ([]coin.SignedBlock, error) {
	ret := _m.Called(seq, count)

	var r0 []coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64, uint64) []coin.SignedBlock); ok {
		r0 = rf(seq, count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coin.SignedBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(seq, count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnconfirmedKnown provides a mock function with given fields: txns
func (_m *MockDaemoner) GetUnconfirmedKnown(txns []cipher.SHA256) (coin.Transactions, error) {
	ret := _m.Called(txns)

	var r0 coin.Transactions
	if rf, ok := ret.Get(0).(func([]cipher.SHA256) coin.Transactions); ok {
		r0 = rf(txns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coin.Transactions)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]cipher.SHA256) error); ok {
		r1 = rf(txns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnconfirmedUnknown provides a mock function with given fields: txns
func (_m *MockDaemoner) GetUnconfirmedUnknown(txns []cipher.SHA256) ([]cipher.SHA256, error) {
	ret := _m.Called(txns)

	var r0 []cipher.SHA256
	if rf, ok := ret.Get(0).(func([]cipher.SHA256) []cipher.SHA256); ok {
		r0 = rf(txns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cipher.SHA256)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]cipher.SHA256) error); ok {
		r1 = rf(txns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HeadBkSeq provides a mock function with given fields:
func (_m *MockDaemoner) HeadBkSeq() (uint64, bool, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// IncreaseRetryTimes provides a mock function with given fields: addr
func (_m *MockDaemoner) IncreaseRetryTimes(addr string) {
	_m.Called(addr)
}

// InjectTransaction provides a mock function with given fields: txn
func (_m *MockDaemoner) InjectTransaction(txn coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error) {
	ret := _m.Called(txn)

	var r0 bool
	if rf, ok := ret.Get(0).(func(coin.Transaction) bool); ok {
		r0 = rf(txn)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *visor.ErrTxnViolatesSoftConstraint
	if rf, ok := ret.Get(1).(func(coin.Transaction) *visor.ErrTxnViolatesSoftConstraint); ok {
		r1 = rf(txn)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*visor.ErrTxnViolatesSoftConstraint)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(coin.Transaction) error); ok {
		r2 = rf(txn)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// IsDefaultConnection provides a mock function with given fields: addr
func (_m *MockDaemoner) IsDefaultConnection(addr string) bool {
	ret := _m.Called(addr)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsMaxDefaultConnectionsReached provides a mock function with given fields:
func (_m *MockDaemoner) IsMaxDefaultConnectionsReached() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Mirror provides a mock function with given fields:
func (_m *MockDaemoner) Mirror() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// PexConfig provides a mock function with given fields:
func (_m *MockDaemoner) PexConfig() pex.Config {
	ret := _m.Called()

	var r0 pex.Config
	if rf, ok := ret.Get(0).(func() pex.Config); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(pex.Config)
	}

	return r0
}

// RandomExchangeable provides a mock function with given fields: n
func (_m *MockDaemoner) RandomExchangeable(n int) pex.Peers {
	ret := _m.Called(n)

	var r0 pex.Peers
	if rf, ok := ret.Get(0).(func(int) pex.Peers); ok {
		r0 = rf(n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pex.Peers)
		}
	}

	return r0
}

// RecordConnectionMirror provides a mock function with given fields: addr, mirror
func (_m *MockDaemoner) RecordConnectionMirror(addr string, mirror uint32) error {
	ret := _m.Called(addr, mirror)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint32) error); ok {
		r0 = rf(addr, mirror)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecordMessageEvent provides a mock function with given fields: m, c
func (_m *MockDaemoner) RecordMessageEvent(m AsyncMessage, c *gnet.MessageContext) error {
	ret := _m.Called(m, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(AsyncMessage, *gnet.MessageContext) error); ok {
		r0 = rf(m, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RecordPeerHeight provides a mock function with given fields: addr, height
func (_m *MockDaemoner) RecordPeerHeight(addr string, height uint64) {
	_m.Called(addr, height)
}

// RecordUserAgent provides a mock function with given fields: addr, userAgent
func (_m *MockDaemoner) RecordUserAgent(addr string, userAgent useragent.Data) error {
	ret := _m.Called(addr, userAgent)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, useragent.Data) error); ok {
		r0 = rf(addr, userAgent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveFromExpectingIntroductions provides a mock function with given fields: addr
func (_m *MockDaemoner) RemoveFromExpectingIntroductions(addr string) {
	_m.Called(addr)
}

// RequestBlocksFromAddr provides a mock function with given fields: addr
func (_m *MockDaemoner) RequestBlocksFromAddr(addr string) error {
	ret := _m.Called(addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResetRetryTimes provides a mock function with given fields: addr
func (_m *MockDaemoner) ResetRetryTimes(addr string) {
	_m.Called(addr)
}

// SendMessage provides a mock function with given fields: addr, msg
func (_m *MockDaemoner) SendMessage(addr string, msg gnet.Message) error {
	ret := _m.Called(addr, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, gnet.Message) error); ok {
		r0 = rf(addr, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHasIncomingPort provides a mock function with given fields: addr
func (_m *MockDaemoner) SetHasIncomingPort(addr string) error {
	ret := _m.Called(addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
