package main

import (
	"fmt"
)

const (
	SKY_OK    = 0
	SKY_ERROR = 0xFFFFFFFF
)

func libErrorCode(err error) uint32 {
	if err == nil {
		return SKY_OK
	}
	// TODO: Implement error codes
	return SKY_ERROR
}

func catchApiPanic(err interface{}) (retVal uint32) {
	retVal = SKY_OK
	if err != nil {
		fmt.Printf("API panic detected : %v\n", err)
		// TODO: Fix to be like retVal = libErrorCode(err)
		retVal = SKY_ERROR
	}
	return
}