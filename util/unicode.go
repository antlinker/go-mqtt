package util

import (
	"unicode/utf8"
)

//校验uincode0的字符串，如果带有0则返回false否则返回true
func VerifyZeroByMqtt(data []rune) bool {
	for _, c := range data {
		if c == 0 {
			return false
		}
	}
	return true
}

//校验字符串符合mqtt字符串格式规范
func VerifyUnicodeByMqtt(data []rune) bool {
	for _, c := range data {
		if !verifyUnicodeByMqttByOne(c) {
			return false
		}
	}
	return true
}

func verifyUnicodeByMqttByOne(c rune) bool {
	//// U+0001 和 U+001F U+0000
	if c <= 0x001F {
		return false
	}
	if c == utf8.RuneError {
		return false
	}
	// U+D800 和 U+DFFF
	if c >= 0xD800 && c <= 0xDFFF {
		return false
	}
	//U+007F 和 U+009F
	if c >= 0x007F && c <= 0x009F {
		return false
	}
	if c == 0x0FFFF {
		return false
	}
	return true
}
