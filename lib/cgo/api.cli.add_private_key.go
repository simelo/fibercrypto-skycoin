package main

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

//export SKY_cli_AddPrivateKey
func SKY_cli_AddPrivateKey(_wlt *C.Handle, _key string) (____error_code uint32) {
	//TODO: Wallet must be Handle
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	obj, ok := lookupHandleObj(Handle(*_wlt))
	____error_code = SKY_ERROR
	if ok {
		if wlt, isWallet := (obj).(*wallet.Wallet); isWallet {
			____return_err := cli.AddPrivateKey(wlt, _key)
			____error_code = libErrorCode(____return_err)
		}
	}
	return
}

//export SKY_cli_AddPrivateKeyToFile
func SKY_cli_AddPrivateKeyToFile(_walletFile, _key string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	walletFile := _walletFile
	key := _key
	____return_err := cli.AddPrivateKeyToFile(walletFile, key)
	____error_code = libErrorCode(____return_err)
	return
}
