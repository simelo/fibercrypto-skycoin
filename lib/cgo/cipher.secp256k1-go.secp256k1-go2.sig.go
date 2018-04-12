package main

import (
	"reflect"
	"unsafe"
	secp256k1go "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_secp256k1go_Signature_Print
func SKY_secp256k1go_Signature_Print(_sig *C.Number, _lab string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	lab := _lab
	sig.Print(lab)
	return
}

// export SKY_secp256k1go_Signature_Verify
func SKY_secp256k1go_Signature_Verify(_sig *C.Number, _pubkey *C.secp256k1go__XY, _message *C.Number, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	pubkey := (*secp256k1go.XY)(unsafe.Pointer(_pubkey))
	message := (*secp256k1go.Number)(unsafe.Pointer(_message))
	__arg2 := sig.Verify(pubkey, message)
	*_arg2 = __arg2
	return
}

// export SKY_secp256k1go_Signature_Recover
func SKY_secp256k1go_Signature_Recover(_sig *C.Number, _pubkey *C.secp256k1go__XY, _m *C.Number, _recid int, _arg3 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	pubkey := (*secp256k1go.XY)(unsafe.Pointer(_pubkey))
	m := (*secp256k1go.Number)(unsafe.Pointer(_m))
	recid := _recid
	__arg3 := sig.Recover(pubkey, m, recid)
	*_arg3 = __arg3
	return
}

// export SKY_secp256k1go_Signature_Sign
func SKY_secp256k1go_Signature_Sign(_sig *C.Number, _seckey, _message, _nonce *C.Number, _recid *int, _arg2 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	seckey := (*secp256k1go.Number)(unsafe.Pointer(_seckey))
	message := (*secp256k1go.Number)(unsafe.Pointer(_message))
	nonce := (*secp256k1go.Number)(unsafe.Pointer(_nonce))
	recid := _recid
	__arg2 := sig.Sign(seckey, message, nonce, recid)
	*_arg2 = __arg2
	return
}

// export SKY_secp256k1go_Signature_ParseBytes
func SKY_secp256k1go_Signature_ParseBytes(_sig *C.Number, _v *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	v := *(*[]byte)(unsafe.Pointer(_v))
	sig.ParseBytes(v)
	return
}

// export SKY_secp256k1go_Signature_Bytes
func SKY_secp256k1go_Signature_Bytes(_sig *C.Number, _arg0 *C.GoSlice_) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	sig := (*secp256k1go.Signature)(unsafe.Pointer(_sig))
	__arg0 := sig.Bytes()
	copyToGoSlice(reflect.ValueOf(__arg0), _arg0)
	return
}
