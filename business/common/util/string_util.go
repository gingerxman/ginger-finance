package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func ToCamelString(str string) string {
	if str == "" {
		return ""
	}

	if !strings.Contains(str, "_"){
		return str
	}

	temp := strings.Split(str, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') {
				vv[0] -= 32
			}
			s += string(vv)
		}
	}
	return s
}

// return if s in slice
func StringInSlice(s string, slice []string) bool{
	if len(slice) == 0{
		return false
	}

	for _, ss := range slice{
		if ss == s{
			return true
		}
	}

	return false
}

func GetMD5(str string) string{
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}