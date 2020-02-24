package core

import (
	"errors"
)

var ERROR_NO_KEY_FOUND = errors.New("no key found")
var ERROR_NO_GET_WRONG_TYPE = errors.New("get wrong type")

type TtfInfo map[string]interface{}

func (t TtfInfo) PushString(key string, val string) {
	t[key] = val
}

func (t TtfInfo) PushBytes(key string, val []byte) {
	t[key] = val
}

func (t TtfInfo) PushInt64(key string, val int64) {
	t[key] = val
}

func (t TtfInfo) PushInt(key string, val int) {
	t[key] = val
}

func (t TtfInfo) PushUInt64(key string, val uint) {
	t[key] = val
}

func (t TtfInfo) PushBool(key string, val bool) {
	t[key] = val
}

func (t TtfInfo) PushInt64s(key string, val []int) {
	t[key] = val
}

func (t TtfInfo) PushMapIntInt64(key string, val map[int]int) {
	t[key] = val
}

func (t TtfInfo) GetBool(key string) (bool, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(bool); ok {
			/* act on str */
			return m, nil
		} else {
			return false, ERROR_NO_GET_WRONG_TYPE
		}
	} else {
		return false, ERROR_NO_KEY_FOUND
	}
}

func (t TtfInfo) GetString(key string) (string, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(string); ok {
			/* act on str */
			return m, nil
		} else {
			return "", ERROR_NO_GET_WRONG_TYPE
		}
	} else {
		return "", ERROR_NO_KEY_FOUND
	}
}

func (t TtfInfo) GetInt64(key string) (int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(int); ok {
			/* act on str */
			return m, nil
		} else {
			return 0, ERROR_NO_GET_WRONG_TYPE
		}
	} else {
		return 0, ERROR_NO_KEY_FOUND
	}
}

func (t TtfInfo) GetInt64s(key string) ([]int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.([]int); ok {
			/* act on str */
			return m, nil
		} else {
			return nil, ERROR_NO_GET_WRONG_TYPE
		}
	} else {
		return nil, ERROR_NO_KEY_FOUND
	}
}

func (t TtfInfo) GetMapIntInt64(key string) (map[int]int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(map[int]int); ok {
			/* act on str */
			return m, nil
		} else {
			return nil, ERROR_NO_GET_WRONG_TYPE
		}
	} else {
		return nil, ERROR_NO_KEY_FOUND
	}
}

func NewTtfInfo() TtfInfo {
	info := make(TtfInfo)
	return info
}
