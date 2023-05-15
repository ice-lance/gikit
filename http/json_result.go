/***
 * @Author       : ICE
 * @Date         : 2023-04-20 19:23:14
 * @LastEditTime : 2023-04-21 16:35:11
 * @LastEditors  : ICE
 * @Copyright (c) 2023 ICE, All Rights Reserved.
 * @Description  : http返回处理
 */

package http

import (
	"encoding/json"
	"time"
)

const (
	JSON_PARSE_ERR = 1001
)

type JSONResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Time    int64       `json:"time"`
}

func Err_Handler(code int, err error) *JSONResult {
	return &JSONResult{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
		Time:    time.Now().UnixMilli(),
	}
}

func Write_Handler(data interface{}) *JSONResult {
	var write_data []byte
	var err error
	if data != nil {
		write_data, err = json.Marshal(data)
		if err != nil {
			return &JSONResult{
				Code:    JSON_PARSE_ERR,
				Message: err.Error(),
				Data:    nil,
				Time:    time.Now().UnixMilli(),
			}
		}
	} else {
		write_data = nil
	}
	return &JSONResult{
		Code:    200,
		Message: "",
		Data:    string(write_data),
		Time:    time.Now().UnixMilli(),
	}

}
