package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

// 自旋锁
type Spinlock struct {
	state *int32
}

const free = int32(0)

func (l *Spinlock) Lock() {
	for !atomic.CompareAndSwapInt32(l.state, free, 42) {
		runtime.Gosched() //释放当前协程的cpu资源，如果没有得到锁的话
	}
}
func (l *Spinlock) Unlock() {
	atomic.StoreInt32(l.state, free) //将锁重置为未被获取的状态
}
//票务的原子操作
type TicketStore struct {
	ticket *uint64
	done   *uint64
	slots  []string
}
func (ts *TicketStore)Put(s string){
	t:=atomic.AddUint64(ts.ticket,1)-1
	ts.slots[t]=s
	for !atomic.CompareAndSwapUint64(ts.done,t,t+1){
		runtime.Gosched()
	}
}
func (ts *TicketStore)GetDone()[]string{
	return ts.slots[:atomic.LoadUint64(ts.done)+1]
}
//管道的csp
func main() {
	ch := make(chan int, 1)
	// close(ch)
	// fmt.Println(<-ch)
	// ch<-1
	for {

		data, more, ok := TryRecive(ch, time.Second)
		fmt.Println(data, more, ok)
	}
}
func TryRecive(ch <-chan int, duration time.Duration) (data int, more, ok bool) {
	select {
	case data, more := <-ch:
		return data, more, true
	case <-time.After(duration):
		return 0, true, false
	}
}
