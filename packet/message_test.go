package packet

import (
	"fmt"
	"testing"
)

func _test(n int32) {
	var r = Remlen2Bytes(n)
	//var k, _ = Bytes2Remlen(r)
	var k, _, _ = readRemlen(r)
	fmt.Printf("原:%d %x 新 %d %X  还原:%d\n", n, n, r, r, k)
}

func ExampleFixedHeader() {
	_test(0x7f)
	_test(128)
	_test(16383)
	_test(16384)
	_test(2097151)
	_test(2097152)
	_test(268435455)
	_test(268435454)
	// Output:原:127 7f 新 [127] 7F  还原:127
	//原:128 80 新 [128 1] 8001  还原:128
	//原:16383 3fff 新 [255 127] FF7F  还原:16383
	//原:16384 4000 新 [128 128 1] 808001  还原:16384
	//原:2097151 1fffff 新 [255 255 127] FFFF7F  还原:2097151
	//原:2097152 200000 新 [128 128 128 1] 80808001  还原:2097152
	//原:268435455 fffffff 新 [255 255 255 127] FFFFFF7F  还原:268435455
	//原:268435454 ffffffe 新 [254 255 255 127] FEFFFF7F  还原:268435454

}
func _testb() {
	var test int32 = 0x7f
	var r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:127 7f 新 [127] 7F
	test = 128
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:128 80 新 [128 1] 8001
	test = 16383
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:16383 3fff 新 [255 127] FF7F
	test = 16384
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:16384 4000 新 [128 128 1] 808001

	test = 2097151
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:2097151 1fffff 新 [255 255 127] FFFF7F
	test = 2097152
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:2097152 200000 新 [128 128 128 1] 80808001
	test = 268435455
	r = Remlen2Bytes(test)
	fmt.Printf("原:%d %x 新 %d %X\n", test, test, r, r)
	// Output:原:268435455 fffffff 新 [255 255 255 127] FFFFFF7F
}

//func TestFixedHeader(t *testing.T) {
//	_testb()
//
//	//0 (0x00)
//	//127 (0x7F)
//	//128 (0x80, 0x01)
//	//16383 (0xFF, 0x7F)
//	//16384 (0x80, 0x80, 0x01)
//	//2097151 (0xFF, 0xFF, 0x7F)
//	// 2097152 (0x80, 0x80, 0x80, 0x01)
//	//268435455 (0xFF, 0xFF, 0xFF, 0x7F)
//10000000
//}

func BenchmarkFixedHeader(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		//_test(int32(i))
		Bytes2Remlen(Remlen2Bytes(int32(i)))
	}
}
