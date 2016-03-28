/**
作者：ybq
时间:2015-09-08
**/
package packet

import "io"

//const (
//	FLAG_CONNECT uint = uint8(iota + 1)
//	FLAG_CONNACK
//	FLAG_PUBLISH
//	FLAG_PUBACK
//	FLAG_PUBREC
//	FLAG_PUBREL
//	FLAG_PUBCOMP
//	FLAG_SUBSCRIBE
//	FLAG_SUBACK
//	FLAG_UNSUBSCRIBE
//	FLAG_UNSUBACK
//	FLAG_PINGREQ
//	FLAG_PINGRESP
//	FLAG_DISCONNECT
//)

//固定报头
type FixedHeader struct {
	hdata   byte
	remlen  int //剩余长度
	totalen int
	//ispacket   bool
	//packetbuff *bytes.Buffer
}

func (fh *FixedHeader) GetRemlen() int {
	return fh.remlen
}

//类型和标志合并后的字节
func (fh *FixedHeader) SetHeaderTypeAndFlag(hdata byte) {
	fh.hdata = hdata
	//	fh.ispacket = false
}

////类型和标志合设置
//func (fh *FixedHeader) SetHeaderTypeAndFlag2(t uint8, flag uint8) {
//	fh.hdata = (t << 4) | (flag & 0xF)
//}
////类型和标志设置
//func (fh *FixedHeader) SetHeaderTypeAndFlag3(t uint8, flag1 bool, flag2 bool, flag3 bool, flag4 bool) {
//	fh.hdata = (t << 4)
//	if flag1 {
//		fh.hdata += 1
//	}
//	if flag2 {
//		fh.hdata += 2
//	}
//	if flag3 {
//		fh.hdata += 4
//	}
//	if flag4 {
//		fh.hdata += 8
//	}
//}

//设置类型
func (fh *FixedHeader) SetHeaderType(t uint8) {
	fh.hdata = (t << 4) | (fh.hdata & 0xf)
	//fh.ispacket = false
}

//设置标志
func (fh *FixedHeader) SetHeaderFlag(t uint8) {
	fh.hdata = (fh.hdata & 0xf0) | (t & 0xf)
	//fh.ispacket = false
}

//获取所有固定头标志
func (fh *FixedHeader) GetHeaderFlag() uint8 {
	return fh.hdata & 0xf
}

//获取类型
func (fh *FixedHeader) GetHeaderType() uint8 {
	return fh.hdata >> 4
}

//获取类型和标志
func (fh *FixedHeader) GetHeaderTypeFlag() uint8 {
	return fh.hdata
}

//从文件解码 主要为PUBLIST类型报文使用，大数据传输缓存中读取数据其他报文不需要重载
func (fh *FixedHeader) UnPacketFile(head byte, varheader []byte, in io.Reader) error {
	return nil
}

//打包到，大数据传输缓存中写入数据 除PUBLIST外其他报文不需要重载
func (fh *FixedHeader) PacketTo(out io.Writer, size int) error {
	return nil
}
func (fh *FixedHeader) Packet() []byte {

	return nil
}

//将信息发送到一个Writer中
func (fh *FixedHeader) WriteTo(out io.Writer) (int64, error) {
	n, err := out.Write(fh.Packet())
	return int64(n), err
}

const (
	_LEN3 int = 2 << 21
	_LEN2 int = 2 << 14
	_LEN1 int = 2 << 7
)

//获取数据长度在调用过Packet或UnPacket后生效否则计算可能出错,如果返回0则还未经过计算不能生效
func (fh *FixedHeader) GetDataLen() (sum int) {
	if fh.totalen <= 0 {
		remlen := fh.remlen
		sum = 2 + remlen
		if remlen >= _LEN3 {
			sum += 3
		} else if remlen >= _LEN2 {
			sum += 2
		} else if remlen >= _LEN1 {
			sum += 1
		}
		fh.totalen = sum
		return
	} else {
		return fh.totalen
	}

}
func (fh *FixedHeader) setDataLen(sum int) {
	fh.totalen = sum
}

// func (fh *FixedHeader) setBuffByRemdata(remdata []byte) {
// 	rs := Remlen2Bytes(int32(fh.remlen))
// 	buffer := new(bytes.Buffer)
// 	//var buff = make([]byte, 1, 1+len(rs)+fh.remlen)
// 	buffer.WriteByte(fh.hdata)
// 	//buff[0] = fh.hdata
// 	buffer.Write(rs)
// 	buffer.Write(remdata)
// 	//buff = append(buff, rs...)
// 	//buff = append(buff, remdata...)
// 	fh.packetbuff = buffer
// 	fh.ispacket = true
// 	fh.totalen = buffer.Len()
// }
