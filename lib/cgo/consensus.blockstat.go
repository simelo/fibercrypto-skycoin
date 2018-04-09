package main

import "unsafe"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"
*/
import "C"

// export SKY_consensus_BlockStat_Init
func SKY_consensus_BlockStat_Init(_self *C.BlockStat) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStat)(unsafe.Pointer(_self))
	self.Init()
	return
}

// export SKY_consensus_BlockStat_GetSeqno
func SKY_consensus_BlockStat_GetSeqno(_self *C.BlockStat, _arg0 *uint64) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStat)(unsafe.Pointer(_self))
	__arg0 := self.GetSeqno()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_BlockStat_Clear
func SKY_consensus_BlockStat_Clear(_self *C.BlockStat) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStat)(unsafe.Pointer(_self))
	self.Clear()
	return
}

// export SKY_consensus_BlockStat_GetBestHashPubkeySig
func SKY_consensus_BlockStat_GetBestHashPubkeySig(_self *C.BlockStat, _arg0 *C.SHA256, _arg1 *C.PubKey, _arg2 *C.Sig) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStat)(unsafe.Pointer(_self))
	__arg0, __arg1, __arg2 := self.GetBestHashPubkeySig()
	return
}

// export SKY_consensus_BlockStat_Print
func SKY_consensus_BlockStat_Print(_self *C.BlockStat) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStat)(unsafe.Pointer(_self))
	self.Print()
	return
}

// export SKY_consensus_PriorityQueue_Len
func SKY_consensus_PriorityQueue_Len(_pq *C.PriorityQueue, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*PriorityQueue)(unsafe.Pointer(_pq))
	__arg0 := pq.Len()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_PriorityQueue_Less
func SKY_consensus_PriorityQueue_Less(_pq *C.PriorityQueue, _i int, _j int, _arg2 *bool) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*PriorityQueue)(unsafe.Pointer(_pq))
	i := _i
	j := _j
	__arg2 := pq.Less(i, j)
	*_arg2 = __arg2
	return
}

// export SKY_consensus_PriorityQueue_Swap
func SKY_consensus_PriorityQueue_Swap(_pq *C.PriorityQueue, _i int, _j int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := *(*PriorityQueue)(unsafe.Pointer(_pq))
	i := _i
	j := _j
	pq.Swap(i, j)
	return
}

// export SKY_consensus_PriorityQueue_Push
func SKY_consensus_PriorityQueue_Push(_pq *C.PriorityQueue, _x interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := (*PriorityQueue)(unsafe.Pointer(_pq))
	pq.Push(x)
	return
}

// export SKY_consensus_PriorityQueue_Pop
func SKY_consensus_PriorityQueue_Pop(_pq *C.PriorityQueue, _arg0 interface{}) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	pq := (*PriorityQueue)(unsafe.Pointer(_pq))
	__arg0 := pq.Pop()
	return
}

// export SKY_consensus_BlockStatQueue_Len
func SKY_consensus_BlockStatQueue_Len(_self *C.BlockStatQueue, _arg0 *int) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStatQueue)(unsafe.Pointer(_self))
	__arg0 := self.Len()
	*_arg0 = __arg0
	return
}

// export SKY_consensus_BlockStatQueue_Print
func SKY_consensus_BlockStatQueue_Print(_self *C.BlockStatQueue) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	self := (*BlockStatQueue)(unsafe.Pointer(_self))
	self.Print()
	return
}
