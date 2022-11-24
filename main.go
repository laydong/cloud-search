package main

import (
	"codesearch/conf"
	"codesearch/global/glogs"
	"codesearch/global/gstore"
	"codesearch/router"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	err := conf.InitDoAfter()
	if err != nil {
		return
	}
	glogs.InitLog()
	gstore.InitMongoDb(conf.ConfInfo.MGConf.Dsn, conf.ConfInfo.MGConf.ConnMaxPoolSize, conf.ConfInfo.MGConf.ConnTimeOut)
	gstore.InitDB(conf.ConfInfo.DBConf.Dsn)
	routers := router.Routers()
	address := fmt.Sprintf("0.0.0.0:%v", conf.ConfInfo.AppConf.HttpListen)
	_ = routers.Run(address)
}

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

//func main() {
//	defer ants.Release()
//
//	runTimes := runtime.NumCPU() * 100
//
//	// Use the common pool.
//	var wg sync.WaitGroup
//	syncCalculateSum := func() {
//		demoFunc()
//		wg.Done()
//	}
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = ants.Submit(syncCalculateSum)
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", ants.Running())
//	fmt.Printf("finish all tasks.\n")
//
//	// Use the pool with a function,
//	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
//	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
//		myFunc(i)
//		wg.Done()
//	})
//	defer p.Release()
//	// Submit tasks one by one.
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = p.Invoke(int32(i))
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", p.Running())
//	fmt.Printf("finish all tasks, result is %d\n", sum)
//}
