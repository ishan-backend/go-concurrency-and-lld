package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var idCounter = 0

type Resource struct {
	id string
}

func NewResource() *Resource {
	fmt.Println("NewResource called")
	idCounter += 1
	return &Resource{fmt.Sprintf("Resource-%d", idCounter)}
}

// demo1 - we get the resource from the pool, do some work and then put it back. By doing this one step at the time in the end we'll end with just one Resource instance.
func demo1() {
	fmt.Println("demo1 start")
	resourcePool := sync.Pool{New: func() any { return NewResource() }}
	for i := 0; i < 10; i++ { // 10 requests we make sync - we get same resource
		item := resourcePool.Get().(*Resource)
		// do some work with item i.e. make some api calls with redis/db
		fmt.Println(fmt.Sprintf("%s doing work", item.id))
		resourcePool.Put(item)
	}
	fmt.Println("demo1 done")
}

// demo2 - we spawn 10 goroutines, that use the pool. Since all goroutines start roughly at the same time and require a resource to doWork, in the end the pool will have 10 Resource instances.
func demo2() {
	fmt.Println("demo2 start")
	wg := sync.WaitGroup{}
	resourcePool := sync.Pool{New: func() any { return NewResource() }}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance := resourcePool.Get().(*Resource)
			fmt.Println(fmt.Sprintf("%s doing work", instance.id))
			resourcePool.Put(instance)
		}()
	}
	wg.Wait()
	fmt.Println("demo2 done")
}

// in demo3, we are putting some random sleeps in between. The faster goroutines will also return the resource faster to the pool and slower goroutines which start at a later time will reuse the resource instead of creating a new one.
func demo3() {
	fmt.Println("demo3 start")
	wg := sync.WaitGroup{}
	resourcePool := sync.Pool{New: func() any { return NewResource() }}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(900)+100) * time.Millisecond)
			instance := resourcePool.Get().(*Resource)
			fmt.Println(fmt.Sprintf("%s doing work", instance.id))
			time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
			resourcePool.Put(instance)
		}()
	}
	wg.Wait()
	fmt.Println("demo3 done")
}

func main() {
	//demo1()
	demo2()
	//demo3()
}
