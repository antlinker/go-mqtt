package packet

//从data 字节数组指定的位置中写入字符串
//返回　字节数组最后位置
func BytesWBytes(buffer []byte, data []byte, start int) int {

	copy(buffer[start:], data)
	return start + len(data)
}

//从data 字节数组指定的位置中写入字符串
//返回　字节数组最后位置
func BytesWString(buffer []byte, data string, start int) int {
	return BytesWBString(buffer, []byte(data), start)

}

//从data 字节数组指定的位置中写入字符串字节数组
//返回　字节数组最后位置
func BytesWBString(buffer []byte, data []byte, start int) int {

	var cnt = len(data)
	buffer[start] = byte(cnt >> 8)
	buffer[start+1] = byte(cnt & 0xff)
	copy(buffer[start+2:], data)
	return start + cnt + 2

}

//从data 字节数组指定的位置中写入uint32整数
//返回　字节数组最后位置
func BytesWUint32(buffer []byte, data uint32, start int) int {
	buffer[start] = byte(data >> 24)
	buffer[start+1] = byte(data >> 16)
	buffer[start+2] = byte(data >> 8)
	buffer[start+3] = byte(data)
	return start + 4
}

//从data 字节数组指定的位置中写入uint16整数
//返回　字节数组最后位置
func BytesWUint16(buffer []byte, data uint16, start int) int {
	buffer[start] = byte(data >> 8)
	buffer[start+1] = byte(data)
	return start + 2
}

//从data 字节数组指定的位置中写入uint8整数
//返回　字节数组最后位置
func BytesWUint8(buffer []byte, data uint8, start int) int {
	buffer[start] = byte(data)
	return start + 1
}

//从data 字节数组指定的位置中写入uint8整数
//返回　字节数组最后位置
func BytesWByte(buffer []byte, data byte, start int) int {
	buffer[start] = data
	return start + 1
}

func ByteRBool(b uint8, pos uint8) bool {
	var t uint8 = (1 << pos)
	return (b & t) == t
}

//从data 字节数组指定位置读取一个uint8
//返回　读取的uint8整数和字节数组最后位置
func BytesRByte(data []byte, start int) (byte, int) {
	return data[start], start + 1

}

//从data 字节数组指定位置读取一个uint16
//返回　读取的uint16整数和字节数组最后位置
func BytesRUint16(data []byte, start int) (uint16, int) {
	return uint16(data[start])<<8 | uint16(data[start+1]), start + 2

}

//从data 字节数组指定位置读取一个uint8
//返回　读取的uint8整数和字节数组最后位置
func BytesRUint8(data []byte, start int) (uint8, int) {
	return uint8(data[start]), start + 1

}

//从data 字节数组指定的位置中读取字符串
//返回　读取的字符串和字节数组最后位置
func BytesRString(data []byte, start int) (string, int) {
	var l int = (int(data[start]) << 8) | int(data[start+1])
	end := start + 2 + l
	return string(data[start+2 : end]), end

}

//从data 字节数组指定的位置中读取字符串字节数组
//返回　读取的字符串和字节数组最后位置
func BytesRBString(data []byte, start int) ([]byte, int) {
	var l int = (int(data[start]) << 8) | int(data[start+1])
	end := start + 2 + l
	return data[start+2 : end], end

}

//从data 字节数组指定的位置中读取字节数组
//返回　读取的字节数组核最后位置
func BytesRBytes(data []byte, start int, cnt int) (string, int) {
	var end = start + cnt
	return string(data[start:end]), end
}

//从data 字节数组指定的位置中读取字节数组到结束
//返回　读取的字节数组
func BytesRBytesEnd(data []byte, start int) []byte {
	return data[start:]
}

//字符串转换字节数组 ０　，１字节为字符串长度
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

//将字节数组转换为剩余长度
//返回剩余长度和读取了字节数量
func Bytes2Remlen(data []byte) (int32, int) {
	var r int32 = 0
	var d int32 = int32(data[0])
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
	var r int32 = 0
	var d int32 = int32(data[0])
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

////循环方式效率低
//func Actlen2Remlen(data []byte) int32 {
//	var r int32 = 0
//	var i int = 0
//	for i = range data {
//		var d int32 = int32(data[i])
//		r += d & 0x7f << uint(i*7)
//		if (d >> 7) == 0 {
//			return r
//		}
//	}
//	return r
//}
//加if效率不高
//func Remlen2actlen(remlen int32) []byte {
//	if remlen > 16383 {
//		return Remlen2actlen1(remlen)
//	} else {
//		return Remlen2actlen2(remlen)
//	}
//}

////remlen>16383时速度较快
//func Remlen2actlen1(remlen int32) []byte {
//	var flag int32 = 0
//	var cnt = 1
//	var r = [4]byte{0, 0, 0, 0}
//	r[3] = byte(remlen >> 21 & 0x7f)
//	if r[3] > 0 {
//		cnt = 4
//		flag = 0x80
//	}
//	r[2] = byte((remlen >> 14 & 0x7f) + flag)
//	if r[2] > 0 {
//		if cnt < 3 {
//			cnt = 3
//		}
//		flag = 0x80
//	}
//	r[1] = byte((remlen >> 7 & 0x7f) + flag)
//	if r[1] > 0 {
//		if cnt < 2 {
//			cnt = 2
//		}
//		flag = 0x80
//	}
//	r[0] = byte((remlen & 0x7f) + flag)
//	return r[:cnt]
//}

//剩余长度转换字节数组
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
			} else {
				return r[0:3]
			}
		} else {
			return r[0:2]
		}
	} else {
		return r[0:1]
	}
}

//
////切片实现速度很慢
//func Remlen2actlen3(remlen int32) []byte {
//	var tmp int32 = 0
//	var r = make([]byte, 1, 4)
//	r = append(r, byte(remlen&0x7f))
//	tmp = remlen >> 7
//	if tmp > 0 {
//		r[0] += 0x80
//	} else {
//		return r[0:1]
//	}
//	r = append(r, byte(tmp&0x7f))
//	tmp = tmp >> 7
//	if tmp > 0 {
//		r[1] += 0x80
//	} else {
//		return r
//	}
//	r = append(r, byte(tmp&0x7f))
//	tmp = tmp >> 7
//	if tmp > 0 {
//		r[2] += 0x80
//	} else {
//		return r
//	}
//	r = append(r, byte(tmp&0x7f))
//	tmp = tmp >> 7
//	if tmp > 0 {
//		r[3] += 0x80
//	} else {
//		return r
//	}
//	return r
//}

////打包
//func (fh *FixedHeader) packet() *byte[]{
//
//}
