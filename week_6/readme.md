---
tags: CCNS
---

# 讀書會 Week 6

# CH8 Goroutines and Channels

* 這一章主要介紹Go的concurrent programming特性
* 在下一章才會提到concurrent programming的壞處和陷阱
* 很多神奇的魔法
* 比較相關的前置任務:
  * Sync vs Async
  * Block vs unblock 
  * [Parallelism vs Concurrency](https://medium.com/golang-%E7%AD%86%E8%A8%98/golang-concurrency-note-0-c4a5b489edaa)
  * [Scheduling in Go](https://medium.com/random-technical-note/scheduling-in-go-727c9b88c93a)

# 8.1 Goroutines

* 每一個concurrently executing activity都稱作**goroutine**
* a goroutine is similar to a thread
* thread和goroutine的差別在下一章會談
* goroutine用 **go** statement 創造
```go=
f()    // call f(); wait for it to return
go f() // create a new goroutine that calls f(); don't wait
```
---
* 下例的spinner用來顯示程式還在跑但program同時在算數列
* 當main func return的時候所有的goroutine都會中斷
  * 除此之外沒辦法寫扣讓他停下來... ~~止まるんじゃねぇぞ~~
  ![](https://i.imgur.com/DXxSJCE.jpg)
* 但用一些方法叫他自己停下來
```go=
package main

import (
        "fmt" 
        "time" 
)
func main() {
     go spinner(100 * time.Millisecond) //goroutine
     const n = 45
     fibN := fib(n) // slow
     fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}
func spinner(delay time.Duration) {
     for {
          for _, r := range `-\|/` {
              fmt.Printf("\r%c", r)
              time.Sleep(delay)
          }
     }
}
func fib(x int) int {
     if x < 2 {
         return x
     }
     return fib(x-1) + fib(x-2)
}
```
# 8.2. Example: Concurrent Clock Server

* 一個把自己目前的時間傳給client的clock server

```go=
// clock1.go
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn) // handle one connection at a time
	}
}
// net.Conn satisfies the io.Writer interface
func handleConn(c net.Conn) {
        // close conn when exit
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
```
output:
```
$ ./clock1 &
$ nc localhost 8000
13:58:54
13:58:55
13:58:56
13:58:57
^C
```
* 另一個實作方法，但會在第一個conn斷掉的時候開另一個conn
```go=
// netcat1.go
// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
```
output:
```
$ ./netcat1
13:58:54            $ ./netcat1
13:58:55
13:58:56
^C
                    13:58:57
                    13:58:58
                    13:58:59
                    ^C
$ killall clock1
```

---

* 把clock1中的handleConn宣告成go statement就變成concurrent形式了
* 太神奇了8傑克
```go=
// clock2.go
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // 魔法在這裡  就這麼神奇
	}	
}
```
output:
```go=
$ go build gopl.io/ch8/clock2
$ ./clock2 &
$ go build gopl.io/ch8/netcat1
$ ./netcat1
14:02:54            $ ./netcat1
14:02:55            14:02:55
14:02:56            14:02:56
14:02:57            ^C
14:02:58
14:02:59            $ ./netcat1
14:03:00            14:03:00
14:03:01            14:03:01
^C                  14:03:02
                    ^C
$ killall clock2
```

# 8.3. Example: Concurrent Echo Server

* 剛剛的例子是one goroutine per connection
* 這節則是multiple goroutines per connection的例子
* Example的目標是讓client傳一串要求，接著server回一個，一直loop到結束

echo server code:
```go=
// Reverb1.go
// Reverb1 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
```
client code: 他把input copy到goroutine中，即使input已經停止，但background goroutine其實還在執行.....等等會解決這個問題
```go=
//netcat2.go
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}


func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
```

output:
```go=
$ go build gopl.io/ch8/reverb1
$ ./reverb1 &
$ go build gopl.io/ch8/netcat2
$ ./netcat2
Hello?  
  HELLO?
  Hello?
  hello?
Is there anybody there?
  IS THERE ANYBODY THERE?
Yooo-hooo!
  Is there anybody there?
  is there anybody there?
  YOOO-HOOO!
  Yooo-hooo!
  yooo-hooo!
^D
$ killall reverb1
```
* 需要更多的goroutine來處理　所以把echo也宣告成goroutine
* 這樣不只是connect的時候是concurrent 連echo都是

```go=
// Reverb2 is a TCP server that simulates an echo.

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
           // 新加的goroutine
           go echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
```
output:
```go=
$ go build gopl.io/ch8/reverb2
$ ./reverb2 &
$ ./netcat2
Is there anybody there?
   IS THERE ANYBODY THERE?
Yooo-hooo!
   Is there anybody there?
   YOOO-HOOO!
   is there anybody there?
   Yooo-hooo!
   yooo-hooo!
^D
$ killall reverb2
```

但要注意concurrency safety 下一章會提到

# 8.4. Channels

* If goroutines are the activities of a concurrent Go program, channels are the connections between them

* channel會有一個type叫做**channel's element type**去代表他所傳輸的資料是什麼資料型別。

下面的例子使用make來創建channel，順帶一提channel只能用make來創建
```go=
ch := make(chan int) // ch has type 'chan int'
```

這邊可以複習一下make可以拿來創建的東西有哪些? 跟new的差別是?

![](https://i.imgur.com/raHL9kx.png)
[source](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-make-and-new/)

作為溝通的管道，當然要知道他是如何send & receive。
channel是使用箭頭<-來操作，第三種操作是把channel用close()關掉

```go=
ch <- x   // a send statement，把x賦值給channel
x=<-ch    // a receive expression，從ch接收資料並且賦值給x
<-ch      // a receive statement; result is discarded	 
close(ch) // close channel
```
* 對已經close的channel send資料會run-time panic
* 可以多宣告一個ok來看channel是否已經close
* channel可以用`==`來compare也可以確認是否是nil
* 送資料到nil channel會被block而送到closed channel則是會直接Return 0
* 後面會提到如何偵測closed channel

如果只是單單的創造一個channel但不給capacity的話他是沒有暫存空間的，兩者其實都各有其用途，我們在8.4.4會再回來談buffered vs unbuffered channels

```go=
ch = make(chan int)    // unbuffered channel
ch = make(chan int, 0) // unbuffered channel
ch = make(chan int, 3) // buffered channel with capacity 3
```

## 8.4.1. Unbuffered Channels

* **channels是FIFO**，這個雖然書這邊沒講到但我覺得對理解channels來說很重要
* unbuffered的用途在於保證goroutine內的讀寫都有完成之後才會結束main func
* unbuffered channels sometimes also called **synchronous channels**

> 你覺得能一次處理很多需求的是同步(synchronous)還是非同步(asynchronous)呢?

以send & receive來說，若要send一個資料給unbuffered chan，他在確定receiver可以接收之前都會把資料block住，也就是在等另一個goroutine執行receive之前都會block

* 其實就是因為你沒有buff，channel只能做一件事情，又因為FIFO的關係所以得先block一開始的行為 (先out first in的)

* 接下來他談了一段有關concurrency的概念，還蠻有趣的

  * **x happens before y**並不是只代表x在y之前會先執行，而是代表x happens before y這件事情已經被保證會發生
  * 若x不是在y之前發生也不是在y之後發生，那我們就會說x concurrent with y
  * 有趣的是即使在可以說兩者是concurrent的之後，兩者也未必會必須是同步進行的
  ![](https://i.imgur.com/vAgQ9Sm.jpg)
  * 下一章會談更多有關order的問題，這問題很重要

---

* Synchronize the two goroutines
* 這邊利用channels回來解決8.3中的client code問題，我們要讓main func等goroutines執行完之後再Exit

```go=
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
	    log.Fatal(err)	
	}
	done := make(chan struct{})
	// go statement是call匿名函式
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}
```

* Line 14做的事情就是強制讓程式必須等到接收到done這個channel才能exit，因此在exit之前都會log "done" message，這是只有unbuffered channels可以做到的事情

* 會把channels的element type宣告成struct{}{}是因為他只有synchronization這個目的而已，但宣告成bool or int更常見就是了
  * 這邊我其實沒被說服為啥要用struct (p.227)
> 這裡最後還說 done <- 1 is shorter than <- struct{}{} so more common.....
> 所以到底為啥要用structㄋ


* remove error logging的原因是因為當input關掉的時候，mustCopy會return並close connection。這會導致你的background goroutine在call io.Copy的時候會傳入conn這個已經關掉的東西，所以會報錯"read from closed connection"

* Exercise 8.3針對這個問題建議了更好的解決辦法

## 8.4.2. Pipelines

* 當用channel使一個goroutine的input是另一個goroutine的output的時候就稱為**pipeline**

* 產生數字 $\Rightarrow$ 平方 $\Rightarrow$ 印出
![](https://i.imgur.com/DpThHAX.jpg)

```go=
// Pipeline1.go
package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()
	
	// Squarer
	go func() {
		for {
			x := <-naturals
			squares <- x * x
		}
	}()
	
	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}
```
因為這樣會一直執行這個pipeline到天荒地老
所以小明想要用剛學到的close去關掉
這時候如果加上
> close(naturals)

後面的channels會發生什麼咧
答案是會傳入0 然後一直輸出0到世界末日
~~所以天荒地老跟世界末日哪個會先發生咧~~


* go在檢查channel是否close上沒有直接的做法 所以只能找個變數一起接收
* 傳統上把這個變數叫做ok
```go=
// Squarer
go func() {
  for {
    x, ok := <-naturals
    if !ok {
       break // channel was closed and drained
    }
    squares <- x * x
  }
  close(squares)
  }()
```
* 下面這個是利用for的特性所以第一個關掉後剩下兩個也會exit
```go=
// pipeline2.go
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	// Squarer
	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}
}
```

* close channel跟close file是完全不同的，不要搞混。
* 非必要可以不用關channel，因為關channel並不會幫你省下資源，但關file很重要
* 如果嘗試close一個已經closed的channel或者是想close一個nil channel都會造成panic
* 8.9有另一個利用closeing channels來做broadcast mechanism的例子

## 8.4.3. Unidirectional Channel Types

在程式寫大的時候通常會用function把他包成許多small pieces
那由於channel作為func parameter的時候通常不是要receive就是要send
go這邊就提供unidirectional channel types來讓你先定義這個channel只能做send還是只能做receive來避免誤用

* 以剛剛pipeline的例子來說
```go=
func counter(out chan int)     
func squarer(out, in chan int)
func printer(in chan int)
```

* 這個send-only or receive-only的violation會在compile的時候偵測
```go=
chan<- int  // send-only
<- chan int // receive-only
```
* 需要注意的是close的定義是讓這個channel不會再有任何的send, 所以只有send-only的channel可以call close(), 如果是receive-only的channel call close()是會造成compile-error的

* 下例把out宣告為send-only, in則是receive-only
* channel的轉換是在call function的時候會執行
* 這個轉換是不可逆的
```go=
// pipeline3.go
package main

import "fmt"

func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}
```

## 8.4.4. Buffered Channels

buffered channel就是在make channel的時候capacity大於0的channel

![](https://i.imgur.com/QwuuKLy.jpg)

在處理send & receive的時候他的處理跟**queue**一樣，也就是說他遵從FIFO

所以若宣告
```go=
ch <- "A"
ch <- "B"
ch <- "C"
```
則會得到
![](https://i.imgur.com/1JJXJ06.jpg)

當channel是滿的時候會先block掉想要傳進來的東西，直到他挪出空間之後才會繼續接收

![](https://i.imgur.com/KK3wWxd.jpg)


若想得知channel的capacity有多大可以
至於len則是回傳目前channel中有多少個element
```go=
fmt.Println(cap(ch)) // "3"
fmt.Println(len(ch)) // "2"
```

目前為止的例子都只在同一個goroutine上作，但這邊作者提醒說如果你只是想用單一個goroutine去達成queue，自己用slice寫就好，殺雞焉用牛刀

---

接下來這個例子同時丟了三個request，但只會response最快的那個，即使慢的兩個還在處理，func就會先回傳最快結束了的那個

```go=
func mirroredQuery() string {
    responses := make(chan string, 3)
    go func() { responses <- request("asia.gopl.io") }()
    go func() { responses <- request("europe.gopl.io") }()
    go func() { responses <- request("americas.gopl.io") }()
    return <-responses // return the quickest response
}
func request(hostname string) (response string) { /* ... */ }
```

* 如果今天是用unbuffered channel來處理會造成**goroutine leak**，也就是後面兩個request被block住，這個誤用很嚴重，因為go的垃圾處理機制不會幫你處理這個

* 總結就是若今天需要synchronized send和receive operation，unbuffered提供一個很好的方法，而buffered則是decouple了兩者的operation

* 書這邊還有舉一個蛋糕店的例子，這間蛋糕店裡面有三個廚師，分別做baking, icing, inscribing, 注意這個程序是有順序性的，如果今天廚房很小一次只能讓一個廚師做事，那這個情境就會比較類似unbuffered

* 如果廚房夠大(capacity夠大)，廚師就可以一起同時做事，會讓蛋糕做更快，但這邊需要考慮三個廚師各自的速度和順序，如果前面的廚師做比較慢，後面的就會一直在buffer中等待，填滿buffer。如果後面的比較快，就會變成empty buffer, 這樣開buffer就沒意義了

[Cake.go](https://github.com/adonovan/gopl.io/blob/master/ch8/cake/cake.go)

# 8.5. Looping in Parallel

這節的例子都引用了[這個package](https://github.com/adonovan/gopl.io/blob/master/ch8/thumbnail/thumbnail.go)
```go=
package thumbnail
// ImageFile reads an image from infile and writes
// a thumbnail-size version of it in the same directory.
// It returns the generated file name, e.g., "foo.thumb.jpg".
func ImageFile(infile string) (string, error)
```


這個例子是要把所有的image轉成thumbnails，這個操作的order無關緊要，他們都是互相獨立的，也就是embarrassingly parallel，而這樣的問題要做平行化是最容易的
```go=
package thumbnail_test

import (
	"log"
	"os"
	"sync"

	"gopl.io/ch8/thumbnail"
)

//!+1
// makeThumbnails makes thumbnails of the specified files.
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}
```
* 加入goroutine嘗試平行化這個處理，但這樣其實是錯的

```go=
// NOTE: incorrect!
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
	}
}
```
* 因為沒有直接的手段去等goroutine做完再exit，所以就透過剛剛看過的unbuffered channels和f來做到
* `f`是匿名函式的explicit parameter，所以會一直update
```go=
// makeThumbnails3 makes thumbnails of the specified files in parallel.
func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f)
	}

	// Wait for goroutines to complete.
	for range filenames {
		<-ch
	}
}
```
* 這邊有講如何在任何一個步驟出錯的時候return error, 我覺得沒有很難這邊就不講了，有興趣可以自己去看 (p.236)

* 直接跳來講final version 
 
  * 首先這裡的filenames不是單純的slice而是從channel傳入的
  * `sync.WaitGroup`是一個在goroutines之間作用的counter, 通常用來等goroutine做完
  * `Done` is equivalent to `Add(-1)`
  * 記得`defer`是把函式推遲到exit的時候執行, LIFO


```go=
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64) // file size
	var wg sync.WaitGroup     // number of working goroutines
	for f := range filenames {

		wg.Add(1) //before goroutine start
		// worker
		go func(f string) {
			defer wg.Done() //用來確定是否做完了
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	go func() {
		// wait會讓func block住，直到counter=0
		wg.Wait()
		close(sizes)
	}()

    // 紀錄總共處理了多少個檔案
	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
```

* thin segments代表在sleep
* thick seqments代表執行
* diagonal arrows代表goroutine的sychronize
![](https://i.imgur.com/ijO3Bw5.jpg)

* [Rain's Example](https://github.com/RainrainWu/probe/blob/master/pkg/utils/runner.go)

# 8.6. Example: Concurrent Web Crawler

* 他說5.6的時候有做過一個crawler但我沒印象了 [github](https://github.com/adonovan/gopl.io/blob/master/ch5/findlinks3/findlinks.go)
* 這一節就是要做出這個crawler的concurrent版本

* [HexRabbit's Crawler](https://github.com/HexRabbit/hydralix)

```go=
// crawl1
// This version quickly exhausts available file descriptors
// due to excessive concurrent calls to links.Extract.
// Also, it never terminates because the worklist is never closed.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//crawl
//沒變
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}


// 作用跟之前的breadthFirst一樣
// 但這次是每個crawl都有自己的goroutine
func main() {
       // worklist本來是用slice 這邊改成用channel
	worklist := make(chan []string)

	// Start with the command-line arguments.
	// 參數輸入得自己是一個goroutine不然會deadlock
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for list := range worklist {
               // link是explicit parameter
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
```
output:
```shell=
$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
http://gopl.io/
https://golang.org/help/
https://golang.org/doc/
https://golang.org/blog/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket:
                                                        too many open files
...
```
* 2 problems
  * a surprising report of a DNS lookup failure for a reliable domain，因為開太多conn, `net.Dial` GG了
    * 那就限制他R 
  * 開太多平形化處理，unbounded parallelism不被鼓勵，除非你有超級多core的CPU
    * 阿你是不會break喔

* 好 限制他
```go=
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token

	if err != nil {
		log.Print(err)
	}
	return list
}
```

* 好　break他
  * 觀察n 
```go=
func main() {
	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}
```

書本這邊還有另一個做法是開20個goroutine，然後刪掉了所有counter，可以參考
```go=
// Crawl3 crawls web links starting with the command-line arguments.
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
 main

import (
	"fmt"
	"log"
	"os"
	"gopl.io/ch5/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!+
func main() {
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() { worklist <- os.Args[1:] }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}
```

# 8.7. Multiplexing with select

* 計組很討厭的MUX
* 例子是火箭發射的倒數計時器

```go=
// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}
```

* 如果因為緊急狀況想要停止計時可以宣告一個abort channel用
* 有兩種情況的關係，不能共用一個channel，不然後者會被block
* 因此就可以使用`select`來multiplex
  * `select`就類似`switch` 
```go=
func main() {
	// ...create abort channel...
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown.  Press return to abort.")
	select {
	case <-time.After(10 * time.Second):
	// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	// 10秒經過....The World
	launch()
}
```
* 下例利用`time.Tick`
* 跑Playground
```go=

// Countdown implements the countdown for a rocket launch.
// NOTE: the ticker goroutine never terminates if the launch is aborted.
// This is a "goroutine leak".
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// ...create abort channel...
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown.  Press return to abort.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}
```
* 這邊的問題在於main exit的時候會停止從tick收**event**,但卻沒有讓ticker goroutine停止
  * **event**是指在channel上傳輸的訊息/資料
  * 因為在上面傳輸的訊息有value跟溝通發生的moment所以才會叫event
* 造成前面提過的goroutine leak
* 所以如果要使用`Tick`的話，情境最好是tick在整個program的lifetime都需要被使用
* 不然也可以像下面這樣用stop解決
```go=
ticker := time.NewTicker(1 * time.Second)
<-ticker.C // receive from the ticker's channel
ticker.Stop() // cause the ticker's goroutine to terminate
```

除此之外`select`也可以讓你在channel上send or receive的時候avoid blocking

* 有東西就收沒東西就不做事
* Doing it repeatedly is called **polling**
```go=
select {
  case <-abort:
     fmt.Printf("Launch aborted!\n")
     return
  default:
  // do nothing
}
```
* nil channel的用處?
  * 讓select裡面的case永遠不會被選到
  * enable or disable cases that correspond to features like **handling timeouts or cancellation**
  * 下一節會有例子

# 8.8. Example: Concurrent Directory Traversal

這個example是要模擬UNIX的du(**d**isk **u**sage)

* 先看WalkDir跟Dirents
* ioutil.ReadDir會returns a slice of os.FileInfo
```go=
// The du1 command computes the disk usage of the files in a directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
    
	// background goroutine
	// Traverse the file tree.
	// 在fileSizes這個channel上走訪
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// main goroutines
	// 算出接收到的filzSize有多大 然後print出來
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
// The du1 variant uses two goroutines and
// prints the total after every file is found.
```
output:
```
$ go build gopl.io/ch8/du1
$ ./du1 $HOME /usr /bin /etc
213201 files 62.7 GB
```

* 這個output要等超久才會跑出來
* 所以有進度條會比較好
* 但如果把`printDiskUsage`丟進for的話會印出超級多行
* 那就設定timer, 每過多久就輸出一次

```go=
var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// ...start background goroutine...
	
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results periodically.
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
	   
		select {
		
		case size, ok := <-fileSizes:
			// 因為沒有for fileSize所以要加終止條件
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}
```
output=
```
$ go build gopl.io/ch8/du2
$ ./du2 -v $HOME /usr /bin /etc
28608 files 8.3 GB
54147 files 10.3 GB
93591 files 15.1 GB
127169 files 52.9 GB
175931 files 62.2 GB
213201 files 62.7 GB
```

* 但效率這樣還是太低了
* 所以接下來把最費時的WalkDir改成goroutine版本

```go=
// The du3 command computes the disk usage of the files in a directory.
package main

// The du3 variant traverses all directories in parallel.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// ...determine roots...

	flag.Parse()
	
	// Determine the initial directories.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes) // final totals
	// ...select loop...
}


func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// Sync.WaitGroup跟前面說過的用法一樣，用來算現在還有多少個goroutine在運作
// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			// 每次call walkDir都會加一個goroutine去處理
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	// ...
	//!-sema

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
```

# 8.9. Cancellation

這邊講的是要如何讓goroutine停下來，一開始有提過這件事情沒辦法直接的去做，但當然還是可以透過一些方法來做到，不然真的就~~止まるんじゃねぇぞ~~了

前面8.7的rocket launch program其實就已經有用abort這個channel去做到了，但它只能一次關一個goroutine，但如果我們想要一次關多個該怎麼做咧?

直覺上會想把每個event都設置abort，但這樣會有count太大/太小的問題，太大是因為有些goroutine會自己停下來，太小則是如果goroutine內又產生了另一個goroutine就會太小，而且這樣goroutine之間無法得知彼此之間是否有abort。

所以我們的目標會是一個可以broadcast到所有event的abort機制

* 這邊利用的方法就是當channel close的時候會讓後續的receive operation直接變zero value的特性

* 這邊改得有點多 我們分段看吧
  * [du4 code](https://github.com/adonovan/gopl.io/blob/master/ch8/du4/main.go)


首先造了cancellation channel `done`，上面不會有任何數值的傳遞，但當這個channel關閉的時候就表示這個程式該停下來了

至於cannelled這個bool function則是用來check是否要cancel
```go=
var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
```

* 第二步我們造一個可讀取standard input的goroutine
* 只要讀到任何的input, 他就會把done這個channel關掉
  * 關掉done就等同於broadcast the cancellation 
```go=
	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()
```
* 這邊for這邊是在main func裡面
```go=
for {
	select {
	// 這個case是確保func在return之前會把fileSizes這個channel drain掉
	case <-done:
		// Drain fileSizes to allow existing goroutines to finish.
		for range fileSizes {
				// Do nothing.
		}
		return
	// 這個
	case size, ok := <-fileSizes:
			// ...
		if !ok {
			break loop // fileSizes was closed
		}
		nfiles++
		nbytes += size
	case <-tick:
		printDiskUsage(nfiles, nbytes)
	}
}
```

`walkDir`會polls the cancellation state
```go=
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		// ...
		}
	}
}
```
把poll放在walk Dir的loop裡面可以避免creating goroutines after the cancellation event

這樣一直遍歷雖然比較沒有效率但比較好理解，書本表示這算是做了一個trade-off


* The bottleneck是dirents中semaphore token的acquisition
```go=
var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	// ...read directory...
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()
	//...
}
```

# 8.10. Example: Chat Server

這邊會用到四個goroutine: main, broadcaster, handleConn, clientWriter

因為也蠻多code的所以分開講

[code](https://github.com/adonovan/gopl.io/blob/master/ch8/chat/chat.go)

---

* main goroutine的目標是接收client丟過來的訊息, 對每個訊息他都會創一個broadcaster goroutine

```go=
// main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
```
* Broadcaster用來處理client的連接和訊息的傳遞
* 對每個client我們只會記錄他的outgoing message
* `cli`就是server用來接收訊息的接口

```go=
// broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {

	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true
		// 聽到的如果是離開就關掉client的outgoing mess
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
```
* handleConn會幫進入的client創造一個outgoing message channel使用
* 還會跟discord一樣跟你打招呼

```go=
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	// greet
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch
	
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()
	// 沒東西讀了就關掉
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}
// 印出訊息跟講話的人是誰
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
```

output:
![](https://i.imgur.com/lprjlYv.jpg)

* 若有n個client就會有2n+2個concurrently communicating goroutines
  * 前提是沒有lockin operation (9.2會提) 

* 只有什麼是share在不同goroutine之間的呢? (這個程式的sharing varialbe是甚麼?)
  * 答案是net.Conn




