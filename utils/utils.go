package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Options ...
type Options []Option

// Option ...
type Option struct {
	Direction string
}

// getDefaultOrder возвращает дефолтное значение порядка чтения байтов
func getDefaultOrder(order binary.ByteOrder) binary.ByteOrder {
	if order == nil {
		return binary.BigEndian
	}

	return order
}

// Int8 ...
func Int8(d []byte, order binary.ByteOrder) (int8, error) {
	var y int8
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// UInt8 ...
func UInt8(d []byte, order binary.ByteOrder) (uint8, error) {
	var y uint8
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// Int16 ...
func Int16(d []byte, order binary.ByteOrder) (int16, error) {
	var y int16
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// UInt16 ...
func UInt16(d []byte, order binary.ByteOrder) (uint16, error) {
	var y uint16
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// Int32 ...
func Int32(d []byte, order binary.ByteOrder) (int32, error) {
	var y int32
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// UInt32 ...
func UInt32(d []byte, order binary.ByteOrder) (uint32, error) {
	var y uint32
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// Int64 ...
func Int64(d []byte, order binary.ByteOrder) (int64, error) {
	var y int64
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// UInt64 ...
func UInt64(d []byte, order binary.ByteOrder) (uint64, error) {
	var y uint64
	err := binary.Read(bytes.NewReader(d), getDefaultOrder(order), &y)
	return y, err
}

// GetReflectValue привеодит значение к соответствующему типу
func GetReflectValue(v interface{}, t reflect.Kind) interface{} {
	switch t {
	case reflect.Int:
		return v.(int)
	case reflect.Int8:
		return int(t)
	case reflect.Int16:
		return int16(t)
	case reflect.Int32:
		return int32(t)
	case reflect.Int64:
		return int64(t)
	case reflect.Uint:
		return uint(t)
	case reflect.Uint8:
		return uint(t)
	case reflect.Uint16:
		return uint16(t)
	case reflect.Uint32:
		return uint32(t)
	case reflect.Uint64:
		return uint64(v.(int))
	case reflect.Float32:
		return float32(t)
	case reflect.Float64:
		return float64(t)
	}

	return nil
}

// GetLengthByType ...
func GetLengthByType(t string, len uint8) (uint8, error) {

	if strings.HasPrefix(t, "buf") {
		bufLenString := strings.ReplaceAll(t, "buf", "")
		bufLen, err := strconv.ParseInt(bufLenString, 10, 32)
		return uint8(bufLen), err
	}

	if strings.HasPrefix(t, "string") {
		if len == 0 {
			return 0, fmt.Errorf("undefined len for string type")
		}
		return len, nil
	}

	switch t {
	case "int8", "uint8":
		return 1, nil
	case "int16", "uint16":
		return 2, nil
	case "int24", "uint24":
		return 3, nil
	case "int32", "uint32":
		return 4, nil
	case "int64", "uint64":
		return 8, nil
	default:
		return 0, fmt.Errorf("undefined type")
	}
}

// GetValueByType ...
func GetValueByType(t string, len uint8, value []byte, order binary.ByteOrder) (interface{}, error) {

	if strings.HasPrefix(t, "buf") {
		bufLenString := strings.ReplaceAll(t, "buf", "")
		_, err := strconv.ParseInt(bufLenString, 10, 32)
		return nil, err
	}

	if strings.HasPrefix(t, "string") {
		//v, err := String(value, order)
		return string(value), nil
	}

	switch t {
	case "int8":
		return Int8(value, order)
	case "uint8":
		return UInt8(value, order)
	case "int16":
		return Int16(value, order)
	case "uint16":
		return UInt16(value, order)
	case "int32":
		return Int8(value, order)
	case "uint32":
		return UInt32(value, order)
	case "int64":
		return Int8(value, order)
	case "uint64":
		return UInt8(value, order)
	}

	return nil, fmt.Errorf("undefined type")
}

// TODO Переделать в дженерики
func CheckBit16(v uint16, pos uint8) bool {
	return v&(1<<pos) != 0
}

func CheckBit8(v uint8, pos uint8) bool {
	return v&(1<<pos) != 0
}
