//
// Author       : ICE
// Date         : 2022-06-13 18:14:50
// LastEditTime : 2022-06-13 18:46:18
// LastEditors  : ICE
// Copyright (c) 2022 ICE, All Rights Reserved.
// Description  :  int <--> bytes
//
package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// IntToBytes 将int类型的数转化为字节并以小端存储
func IntToBytes(intNum int) []byte {
	uint16Num := uint16(intNum)
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.LittleEndian, uint16Num)
	return buf.Bytes()
}

// BytesToUint
func BytesToUint(bytesArr []byte) uint64 {
	parseUint, _ := strconv.ParseUint(fmt.Sprintf("%x", binary.LittleEndian.Uint32(bytesArr)), 16, 32)
	return parseUint
}
