package util

import (
	"sync"
)

// PacketIdFactory 报文标识生产工厂
type PacketIdFactory struct {
	sync.Mutex
	cur   uint16
	idmap map[uint16]bool
}

// NewPacketIdFactory 创建报文标识生产工厂
func NewPacketIdFactory() *PacketIdFactory {
	return &PacketIdFactory{cur: 0, idmap: make(map[uint16]bool)}
}
func (p *PacketIdFactory) InitPacketIdFactory() {
	p.cur = 0
	p.idmap = make(map[uint16]bool)
}
func (p *PacketIdFactory) LockId(id uint16) {
	p.Lock()
	p.idmap[id] = true
	p.Unlock()
}
func (p *PacketIdFactory) ReleaseId(id uint16) {
	p.Lock()
	delete(p.idmap, id)
	p.Unlock()
}
func (p *PacketIdFactory) CreateId() uint16 {
	p.Lock()
	var cur = p.cur + 1
	var idmap = p.idmap
	for {
		if idmap[cur] {
			cur += 1
			continue
		}
		p.cur = cur
		idmap[cur] = true
		p.Unlock()
		return cur
	}
}
