package packet

import "unsafe"

// BytesWBytes 从data 字节数组指定的位置中写入字符串
//返回　字节数组最后位置
func BytesWBytes(buffer []byte, data []byte, start int) int {

	copy(buffer[start:], data)
	return start + len(data)
}

// BytesWString 从data 字节数组指定的位置中写入字符串
//返回　字节数组最后位置
func BytesWString(buffer []byte, data string, start int) int {
	return BytesWBString(buffer, []byte(data), start)

}

// BytesWBString 从data 字节数组指定的位置中写入字符串字节数组
//返回　字节数组最后位置
func BytesWBString(buffer []byte, data []byte, start int) int {

	var cnt = len(data)
	buffer[start] = byte(cnt >> 8)
	buffer[start+1] = byte(cnt & 0xff)
	copy(buffer[start+2:], data)
	return start + cnt + 2

}

// BytesWUint32 从data 字节数组指定的位置中写入uint32整数
//返回　字节数组最后位置
func BytesWUint32(buffer []byte, data uint32, start int) int {
	buffer[start] = byte(data >> 24)
	buffer[start+1] = byte(data >> 16)
	buffer[start+2] = byte(data >> 8)
	buffer[start+3] = byte(data)
	return start + 4
}

// BytesWUint16 从data 字节数组指定的位置中写入uint16整数
//返回　字节数组最后位置
func BytesWUint16(buffer []byte, data uint16, start int) int {
	buffer[start] = byte(data >> 8)
	buffer[start+1] = byte(data)
	return start + 2
}

// BytesWUint8 从data 字节数组指定的位置中写入uint8整数
//返回　字节数组最后位置
func BytesWUint8(buffer []byte, data uint8, start int) int {
	buffer[start] = byte(data)
	return start + 1
}

// BytesWByte 从data 字节数组指定的位置中写入uint8整数
//返回　字节数组最后位置
func BytesWByte(buffer []byte, data byte, start int) int {
	buffer[start] = data
	return start + 1
}

// ByteRBool 字节转换布尔型
func ByteRBool(b uint8, pos uint8) bool {
	var t uint8 = (1 << pos)
	return (b & t) == t
}

// BytesRByte 从data 字节数组指定位置读取一个uint8
//返回　读取的uint8整数和字节数组最后位置
func BytesRByte(data []byte, start int) (byte, int) {
	return data[start], start + 1

}

// BytesRUint16 从data 字节数组指定位置读取一个uint16
//返回　读取的uint16整数和字节数组最后位置
func BytesRUint16(data []byte, start int) (uint16, int) {
	return uint16(data[start])<<8 | uint16(data[start+1]), start + 2

}

// BytesRUint8 从data 字节数组指定位置读取一个uint8
//返回　读取的uint8整数和字节数组最后位置
func BytesRUint8(data []byte, start int) (uint8, int) {
	return uint8(data[start]), start + 1

}

// BytesRString 从data 字节数组指定的位置中读取字符串
//返回　读取的字符串和字节数组最后位置
func BytesRString(data []byte, start int) (string, int) {
	var l = (int(data[start]) << 8) | int(data[start+1])
	end := start + 2 + l
	return byte2str(data[start+2 : end]), end

}

// BytesRBString 从data 字节数组指定的位置中读取字符串字节数组
//返回　读取的字符串和字节数组最后位置
func BytesRBString(data []byte, start int) ([]byte, int) {
	var l = (int(data[start]) << 8) | int(data[start+1])
	end := start + 2 + l
	return data[start+2 : end], end

}

// BytesRBytes 从data 字节数组指定的位置中读取字节数组
//返回　读取的字节数组核最后位置
func BytesRBytes(data []byte, start int, cnt int) (string, int) {
	var end = start + cnt
	return byte2str(data[start:end]), end
}

// BytesRBytesEnd 从data 字节数组指定的位置中读取字节数组到结束
//返回　读取的字节数组
func BytesRBytesEnd(data []byte, start int) []byte {
	return data[start:]
}

// String2Bytes 字符串转换字节数组 ０　，１字节为字符串长度
func String2Bytes(msg string) []byte {
	var tmp = []byte(msg)
	var cnt = len(tmp)
	var out = make([]byte, cnt+2)
	out[0] = byte(cnt >> 8)
	out[1] = byte(cnt & 0xff)
	for i := range tmp {
		out[i+2] = tmp[i]
	}
	return out
}

// Bytes2Remlen 将字节数组转换为剩余长度
//返回剩余长度和读取了字节数量
func Bytes2Remlen(data []byte) (int32, int) {
	var r int32
	var d = int32(data[0])
	r += d & 0x7f
	if (d >> 7) == 0 {
		return r, 1
	}
	d = int32(data[1])
	r |= d & 0x7f << 7
	if (d >> 7) == 0 {
		return r, 2
	}
	d = int32(data[2])
	r |= d & 0x7f << 14
	if (d >> 7) == 0 {
		return r, 3
	}
	d = int32(data[3])
	r |= d & 0x7f << 21
	return r, 4
}

func readRemlen(data []byte) (int32, []byte, bool) {
	if len(data) == 0 {
		return 0, data, false
	}
	var r int32
	var d = int32(data[0])
	r |= d & 0x7f

	if (d >> 7) == 0 {
		return r, data[1:], true
	}
	if len(data) == 1 {
		return 0, data, false
	}
	d = int32(data[1])
	r |= d & 0x7f << 7
	if (d >> 7) == 0 {
		return r, data[2:], true
	}
	if len(data) == 2 {
		return 0, data, false
	}
	d = int32(data[2])
	r |= d & 0x7f << 14
	if (d >> 7) == 0 {
		return r, data[3:], true
	}
	if len(data) == 3 {
		return 0, data, false
	}
	d = int32(data[3])
	r |= d & 0x7f << 21
	return r, data[4:], true
}

// Remlen2Bytes 剩余长度转换字节数组
func Remlen2Bytes(remlen int32) []byte {
	var tmp int32 = 0
	var r = [4]byte{0, 0, 0, 0}
	r[0] = byte(remlen & 0x7f)
	tmp = remlen >> 7
	if tmp > 0 {
		r[0] += 0x80
		r[1] = byte(tmp & 0x7f)
		tmp = tmp >> 7
		if tmp > 0 {
			r[1] += 0x80
			r[2] = byte(tmp & 0x7f)
			tmp = tmp >> 7
			if tmp > 0 {
				r[2] += 0x80
				r[3] = byte(tmp & 0x7f)
				return r[0:4]
			}
			return r[0:3]
		}
		return r[0:2]
	}
	return r[0:1]

}
func byte2str(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
