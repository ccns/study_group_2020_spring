---
tags: CCNS
---

# 讀書會 Week 0

----

# CH1 Tutorial

----

# 1.1 Hello World

- Standard IO

```go=
package main

import "fmt"

func main() {
  fmt.Println("Hello, World")
}
```

----

# Loop

- Loop syntax

```go=
for initialization; condition; post {

  // zero or more statements
}

// a traditional "while" loop
for condition { 
  // ...
}
```

----

# 1.2 Echo

## echo_1

- Standard args read

```go=
// Echo1 prints its command-line arguments. 
package main

import (
  "fmt"
  "os"
 )

func main() {
  var s, sep string
  for i := 1; i < len(os.Args); i++ {
    s += sep + os.Args[i]
    sep = " "
  }
  fmt.Println(s)
}

```

----

## echo_2

- Standard args read with loop

```go=

// Echo2 prints its command-line arguments.
package main

import (
  "fmt"
  "os"
)

func main() {
  s, sep := "", ""
  for _, arg := range os.Args[1:] { 
    s += sep + arg
    sep = " "
  }
  fmt.Println(s)
}
```
----

# 1.3 Finding Duplicated Lines

- 事實上，我們很少拿 Golang 來進行複雜的字串操作

```go=
package main

import (
  "bufio"
  "fmt"
  "os"
)  

func main() {
  counts := make(map[string]int)
  input := bufio.NewScanner(os.Stdin)
  for input.Scan() {
    counts[input.Text()]++
  }
  // NOTE: ignoring potential errors from input.Err()
  
  for line, n := range counts {
    if n > 1 {
      fmt.Printf("%d\t%s\n", n, line)
    }
  }
}
```

----

# 1.4 Animated GIFs

- 我其實不知道它可以做 GIF
- 我也沒看過別人用這功能
- Code 基本上都是數學函數，先跳過

----

# 1.5 Fetch a URL

```go=
package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)

func main() {
  for _, url := range os.Args[1:] {
    resp, err := http.Get(url)
    if err != nil {
      fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
      os.Exit(1)
    }
    b, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
      fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err) 
      os.Exit(1)
    }
    fmt.Printf("%s", b)
  }
}
```

----

# 1.6 Fetch Concurrently

- fetch concurrently

```go=
package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "time"
)

func main() {
  start := time.Now()
  ch := make(chan string)
  for _, url := range os.Args[1:] {
    go fetch(url, ch)
    // start a goroutine
  }
  for range os.Args[1:] {
    fmt.Println(<-ch)
    // receive from channel ch
  }
  fmt.Printf("%.2fs elapsed\n",time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

```

----

# 1.7 Web Server

## server_1

```go=
package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  // each request calls handler
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

```

----

## server_2

```go=
package main
import (
  "fmt"
  "log"
  "net/http"
  "sync"
)

var mu sync.Mutex
var count int

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/count", counter)
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
  mu.Lock()
  count++
  mu.Unlock()
  fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
// func counter(w http.ResponseWriter, r *http.Request) { mu.Lock() fmt.Fprintf(w, "Count %d\n", count) mu.Unlock() }
```

----

## request content

```go=
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
  for k, v := range r.Header {
    fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
  }
  fmt.Fprintf(w, "Host = %q\n", r.Host)
  fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
  if err := r.ParseForm(); err != nil {
    log.Print(err)
  }
  
  for k, v := range r.Form {
    fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
  }
}
```

----

# 1.8 Loose Ends

- just a concept of control flow
- for example: cases switching

```go=
switch coinflip() {
  case "heads":
    heads++
  case "tails":
    tails++
  default:
    fmt.Println("landed on edge!")
}
```

---

# CH2 Program Structure

----

# 命名

----

## 規範

- [a-zA-Z_][a-zA-Z0-9_]*
- camel case
- 別用到關鍵字
![](https://i.imgur.com/NZhjiZM.png)
![](https://i.imgur.com/QT2AG3f.png)

```go=
package main

// valid
var _foo int
var _123 int
var fooBar int

// not recommended
var foo_bar int

// syntax error
var 123 int

// may cause type ambiguous, int is not a type
// var int int
```

----

## Scope

```go=
package main

import (
	"fmt"
)

var _foo int = 0

func foo() {

	// unused, because of not accessible out of func
	// var _foo int = 1
}

func main() {

	var _bar int = 2
	fmt.Println("_foo: ", _foo)
	fmt.Println("_bar: ", _bar)
}

// $ go run week_1/ch_2/variable/scope.go
// _foo:  0
// _bar:  2
```

----

# 宣告

----

## package

- Package 是 Go program 的模組化管理單位
- 任何 .go file 必須以 Package 宣告開頭
- 而後接連的是 import
- 任何 .go file 必須屬於一個 Package，且只能屬於一個
- 但一個 Package 可以包含多個 .go files

----

### main package

- 程式起點，且必須實作 main 函數
- 若缺少 main 函數則程式無法啟動
  - Go run 無法執行
  - Go build 不會包東西出來

```go=
// skip may cause :
// expected 'package', found 'func'
package main

// duplicate may cause :
// syntax error: non-declaration statement outside function body
// package sub

// skip may cause :
// function main is undeclared in the main package
func main() { /* Do something*/ }
```

----

## import

```go=
package main

// single package import
import "fmt"

// multiple package import
import (
  "fmt"
  "os"
)

// preferred arrangement
import (
  "standard packages"
  
  "third party packages"
  
  "in-project local packages"
)
```

----

## var
  - default value, no un-inited var
  - auto type map with list or func call
  - init 時機
    - 執行前
    - 被 import 前
  - short declare(a declaration, = is assignment), error if re-declare

```go=
package main

import (
	"fmt"
	"reflect"
)

var foo, bar, mur = 0.2, true, "Hi"

func declare() {

	unique := 0
	unique, sub := 1, 2
    unique, sub = sub, unique

	// re-declare error
	// unique := 3

	fmt.Println("declare: ", unique, sub)
}

func main() {

	declare()
	fmt.Println(foo, bar, mur)
	fmt.Println(reflect.TypeOf(foo), reflect.TypeOf(bar), reflect.TypeOf(mur))
}

// $ go run week_1/ch_2/variable/declare.go
// declare:  2 1
// 0.2 true Hi
// float64 bool string

```
----

- pointer
    - operator order (pointer first)
    - new 沒啥人用，跳過
    - lifetime, GC 太複雜，跳過

```go=
package main

import (
	"fmt"
)

func incr(p *int) int {
	*p++
	return *p
}

func main() {

	v := 1
	incr(&v)
	fmt.Println(incr(&v))
}

// $ go run week_1/ch_2/pointer/order.go
// 3
```

----

## assign
  - assign with operator
  - tuple assign
  - unwanted assign

```go=
package main

import (
	"fmt"
)

func multipleReturn() (string, int) {

	return "first", 0
}

func main() {

	_, value := multipleReturn()
	value++
	fmt.Println("value: ", value)

	medals := []string{"gold", "silver", "bronze"}
	medals[0] = "Platinum"
	fmt.Println("medals: ", medals)
}

// $ go run week_1/ch_2/assign/assign.go
// value:  1
// medals:  [Platinum silver bronze]
```

----

## type
  - struct declaration 十分重要，為程式架構核心能力
  - 留意 comparable

```go=
package main

import (
	"fmt"
)

type price float64
type height float64
type house struct {
	Price  price
	Height height
}

const (
	house_price  price  = 4000000.0
	house_height height = 35.5
)

func main() {

	// type mismatch may cause :
	// invalid operation: house_height < house_price (mismatched types height and price)
	// fmt.Println(house_height < house_price)

	fmt.Println(price(house_height) < house_price)
}

// $ go run week_1/ch_2/type/declare.go
// true
```

---

# Practice
- [Reverse Integer](https://leetcode.com/problems/reverse-integer/)
- [Container with most water](https://leetcode.com/problems/container-with-most-water/)
- [search insert position](https://leetcode.com/problems/search-insert-position/)
