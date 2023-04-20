/*
 *    ______                 __
 *   /\__  _\               /\ \
 *   \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *      \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *       \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *       /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *       \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 * @Author       : ICE
 * @Date         : 2022-03-04 11:53:43
 * @LastEditTime : 2022-03-31 15:57:38
 * @LastEditors  : ICE
 * Copyright (c) 2022 ICE, All Rights Reserved.
 * @Description  :
 */

package httputils

import (
	"encoding/json"
)

const (
	JSON_PARSE_ERR = 1001
)

type JSONResult struct {
	Code    int         `json:"code" `
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Err_Handler(code int, err error) *JSONResult {
	return &JSONResult{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
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
			}
		}
	} else {
		write_data = nil
	}
	return &JSONResult{
		Code:    200,
		Message: "",
		Data:    string(write_data),
	}

}
