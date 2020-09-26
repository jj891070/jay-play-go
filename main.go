package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/internal"
)

type Pooler interface {
	NewConn() (*Conn, error)
	CloseConn(*Conn) error

	Get() (*Conn, error)
	Put(*Conn)
	Remove(*Conn)

	Len() int
	IdleLen() int
	Stats() *Stats

	Close() error
}

type ConnPool struct {
	opt *Options //初始化的配置项

	dialErrorsNum uint32 // atomic 连接错误次数

	lastDialError   error //连接错误的最后一次的错误类型
	lastDialErrorMu sync.RWMutex

	queue chan struct{} //池子里面空闲的conn的同步channel

	connsMu sync.Mutex
	conns   []*Conn //活跃的active conns

	idleConnsMu sync.RWMutex
	idleConns   []*Conn //空闲的idle conns

	stats Stats

	_closed uint32 // atomic  //池子是否关闭标签
}

var _ Pooler = (*ConnPool)(nil) //接口检查

func NewConnPool(opt *Options) *ConnPool {
	p := &ConnPool{
		opt: opt,

		queue:     make(chan struct{}, opt.PoolSize), //同步用的
		conns:     make([]*Conn, 0, opt.PoolSize),
		idleConns: make([]*Conn, 0, opt.PoolSize),
	}

	if opt.IdleTimeout > 0 && opt.IdleCheckFrequency > 0 {
		go p.reaper(opt.IdleCheckFrequency) //定时任务，清理过期的conn
	}

	return p
}

//reaper字面意思为收割者,为清理的意思
func (p *ConnPool) reaper(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		if p.closed() {
			break
		}
		//定时清理无用的conns
		n, err := p.ReapStaleConns()
		if err != nil {
			internal.Logf("ReapStaleConns failed: %s", err)
			continue
		}
		atomic.AddUint32(&p.stats.StaleConns, uint32(n))
	}
}

func (p *ConnPool) ReapStaleConns() (int, error) {
	var n int
	for {
		//往channel里面写入一个,表示占用一个任务
		p.getTurn()

		p.idleConnsMu.Lock()
		cn := p.reapStaleConn()
		p.idleConnsMu.Unlock()

		if cn != nil {
			p.removeConn(cn)
		}

		//处理完了，释放占用的channel的位置
		p.freeTurn()

		if cn != nil {
			p.closeConn(cn)
			n++
		} else {
			break
		}
	}
	return n, nil
}

func (p *ConnPool) reapStaleConn() *Conn {
	if len(p.idleConns) == 0 {
		return nil
	}

	//取第一个空闲conn
	cn := p.idleConns[0]
	//判断是否超时没有人用
	if !cn.IsStale(p.opt.IdleTimeout) {
		return nil
	}

	//超时没有人用则从空闲列表里面移除
	p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)

	return cn
}

//判断是否超时的处理
func (cn *Conn) IsStale(timeout time.Duration) bool {
	return timeout > 0 && time.Since(cn.UsedAt()) > timeout
}

//移除连接就是一个很简单的遍历
func (p *ConnPool) removeConn(cn *Conn) {
	p.connsMu.Lock()
	for i, c := range p.conns {
		if c == cn {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			break
		}
	}
	p.connsMu.Unlock()
}

func main() {
	way = make(map[int]string, 5)
	way[0] = "fmt.Sprintf"
	way[1] = "+"
	way[2] = "strings.Join"
	way[3] = "bytes.Buffer"

	k := 4
	d := [5]time.Duration{}
	for i := 0; i < k; i++ {
		d[i] = benchmarkStringFunction(10000, i)
	}
}
