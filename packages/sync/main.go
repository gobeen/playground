package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	title *string
)

func init() {
	title = flag.String("title", "once", "select example title: ||once||cond||wait||")
	flag.Parse()
	fmt.Println("Example Title:", *title)
}

func main() {
	switch *title {
	case "cond":
		cond()
	case "wait":
		waitGroup()
	default:
		once()
	}
}

func cond() {
	type button struct {
		Clicked *sync.Cond
	}
	btn := button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(btn.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(btn.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(btn.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})
	btn.Clicked.Broadcast()
	clickRegistered.Wait()
}

func once() {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func waitGroup() {
	var wg sync.WaitGroup
	dir := "./tmp"
	filename := "test-wait-group"
	start, end := 1, 3

	for i := start; i <= end; i++ {
		wg.Add(1)
		name := fmt.Sprintf("%s-%d", filename, i)
		go func(fName string) {
			defer wg.Done()
			if err := os.Mkdir(dir, 0750); err != nil {
				if !os.IsExist(err) {
					log.Fatal(err)
				}
			}
			filePath := dir + "/" + fName
			if err := os.WriteFile(filePath, []byte("Hello Gophers!"), 0666); err != nil {
				log.Fatal(err)
			}
			log.Printf("%s have been created successfully!", filePath)
		}(name)
	}
	wg.Wait()
}
