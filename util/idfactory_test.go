package util

//"time"

//func TestCreateid(t *testing.T) {
//	runtime.GOMAXPROCS(1)
//	quit := make(chan int,1)
//	idfactory := NewPacketIdFactory()
//	for j := 0; j < 10; j++ {
//		go func(j int, idfactory *PacketIdFactory) {
//			for i := 0; i < 65540; i++ {
//				id := idfactory.CreateId()
//				fmt.Printf("%d线程:%d \n", j, id)
//				idfactory.ReleaseId(id)
//				//runtime.Gosched()
//				//time.Sleep(1000 * time.Millisecond)
//			}
//			quit <- 1
//		}(j, idfactory)
//	}
//	for i := 0; i < 10; i++ {
//		<-quit
//	}
//}

//func BenchmarkCreateidM(b *testing.B) {
//	runtime.GOMAXPROCS(8)

//	quit := make(chan int)
//	for j := 0; j < 2; j++ {

//		go func(b *testing.B) {
//			idfactory := NewPacketIdFactory()
//			for i := 0; i < b.N; i++ {
//				idfactory.ReleaseId(idfactory.CreateId())
//			}
//			quit <- 1
//		}(b)
//	}
//	for i := 0; i < 2; i++ {
//		<-quit
//	}
//}

//func BenchmarkCreateid(b *testing.B) {
//	idfactory := NewPacketIdFactory()
//	for i := 0; i < b.N; i++ {
//		idfactory.ReleaseId(idfactory.CreateId())
//	}
//}

//func BenchmarkCreateOldid(b *testing.B) {
//	idfactory := NewIdFactory()
//	for i := 0; i < b.N; i++ {
//		idfactory.ReleaseId(idfactory.CreateId())
//	}
//}
