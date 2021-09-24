package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch vTypes := v.(type) {
	case float64:
		return vTypes
	case string:
		vStr := vTypes
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	default:
		panic("to float64 error.")
	}
}

func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch vTypes := v.(type) {
	case string:
		vStr := vTypes
		vInt, _ := strconv.Atoi(vStr)
		return vInt
	case int:
		return vTypes
	case float64:
		vF := vTypes
		return int(vF)
	default:
		panic("to int error.")
	}
}

func ToUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}

	switch vType := v.(type) {
	case int:
		return uint64(vType)
	case float64:
		return uint64((vType))
	case string:
		uV, _ := strconv.ParseUint(vType, 10, 64)
		return uV
	default:
		panic("to uint64 error.")
	}
}

func ToInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}

	switch vType := v.(type) {
	case float64:
		return int64(vType)
	default:
		vv := fmt.Sprint(v)

		if vv == "" {
			return 0
		}

		vvv, err := strconv.ParseInt(vv, 0, 64)
		if err != nil {
			return 0
		}

		return vvv
	}
}

func FloatToString(f float64, precision int) string {
	return fmt.Sprint(FloatToFixed(f, precision))
}

func FloatToFixed(f float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(f*p) / p
}
