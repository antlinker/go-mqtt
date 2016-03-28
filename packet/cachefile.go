package packet

import (
	"fmt"
	//"io"
	"os"
	//"strconv"
	"sync"
)

var cacheFactory *CacheFactory = nil
var cacheLock = &sync.Mutex{}

//缓存工厂
type CacheFactory struct {
	sync.Mutex
	cur      uint64
	cachemap map[uint64]*CacheFile
}

//生成缓存管理工厂，单例模式
func NewCacheFactory() *CacheFactory {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	if cacheFactory == nil {
		tmp := &CacheFactory{}
		tmp.cur = 0
		tmp.cachemap = make(map[uint64]*CacheFile)
	}
	return cacheFactory
}

//获取缓存
func (f *CacheFactory) GetCache(id uint64) *CacheFile {
	f.Lock()
	defer f.Unlock()
	cf := f.cachemap[id]

	if cf == nil {
		return f.CreateCache()
	}
	return cf
}

//创建缓存文件
func (f *CacheFactory) CreateCache() *CacheFile {
	f.Lock()
	defer f.Unlock()
	id := f.createId()
	tmp := NewCacheFile()
	tmp.id = id
	tmp.create()
	f.cachemap[id] = tmp
	return tmp
}

//释放文件缓存
func (f *CacheFactory) ReleaseCache(id uint64) {
	f.Lock()
	f.cachemap[id].deleteFile()
	delete(f.cachemap, id)
	f.Unlock()
}

//缓存文件
type CacheFile struct {
	//唯一标识
	id uint64
	//缓存文件
	cacheFile *os.File

	filesize int64
}

func NewCacheFile() *CacheFile {
	return &CacheFile{}
}
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//创建文件
func (c *CacheFile) create() error {
	src := PUBLISH_CACHE_DIR + "/"
	if !IsExist(src) {
		if err := os.MkdirAll(src, 0777); err != nil {
			if os.IsPermission(err) {
				fmt.Println("你不够权限创建文件")
			}
			return err
		}
	}
	src += fmt.Sprint("%d", c.id)
	tmp, err := os.Create(src)
	if err != nil {
		return err
	}
	c.cacheFile = tmp
	return nil
}

//获取缓存文件
func (c *CacheFile) GetFileSize() int64 {
	return c.filesize
}

//获取缓存文件
func (c *CacheFile) GetFile() *os.File {
	return c.cacheFile
}

//删除文件
func (c *CacheFile) deleteFile() {
	c.cacheFile.Close()
	os.Remove(c.cacheFile.Name())
}

//追加数据
func (c *CacheFile) Append(data []byte) error {
	_, err := c.cacheFile.Write(data)
	if err != nil {
		return err
	}
	c.filesize += int64(len(data))
	return nil
}
func (c *CacheFile) Read(data []byte) (int, error) {
	return c.cacheFile.Read(data)
}
func (c *CacheFile) ReadAt(data []byte, off int64) (int, error) {
	return c.cacheFile.ReadAt(data, off)
}

//缓存id发生器

func (p *CacheFactory) createId() uint64 {
	var cur = p.cur
	var idmap = p.cachemap
	for {
		if idmap[cur] != nil {
			cur += 1
			continue
		}
		p.cur = cur + 1
		p.Unlock()
		return cur
	}
}
