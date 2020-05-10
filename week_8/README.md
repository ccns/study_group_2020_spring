讀書會 Week8
===

# Package
- 依照目前的軟體規模，可能一個工具就有上萬的函式，但我們可以透過分類包裝成 Package，讓他能更簡單、更有組織性的重複使用。
- Go 自帶了萬用的 Client 端工具 `go`，除了可以用來執行程式和打包執行檔外，也能用來管理套件

## 10.1 Introduction
- Go 語言的 Package 就是模組化開發的實現。
- 每個 Package 都是獨立的 namespace，避免糾結於命名的衝突。
- Package 同時也有封裝的功能，控制個個物件 Public 或 Private，藉以規範安全的使用方式。
- Go 的 Package 在改動後都需要重新 compile 才會生效，不過針對 Package import 的導入非常快，主要原因有三
  - 需在檔案開頭即聲明導入的 Package，立刻觸發的編譯
  - 導入 Package 之間需為有向無環圖，讓它們能各自獨立編譯甚至並行化
  - 自身編譯好之後，可以提供給上層多個相依 Package，無須重複編譯

## 10.2 Import path
- Go 的原生語法並沒有特別規範 import 字串中特別的意義，因此都是由 `go` 這個 Client 端工具去解讀的。
```go
import (

  "fmt"
  "math/rand"
  "encoding/json"

  "golang.org/x/net/html"
  "github.com/go-sql-driver/mysql"
)
```

- 對於一個要使用的 Package，他的 import path 必須要式 globally unique 的，否則可能會造成導入錯誤。
- 若是網路上的 Package，就要提供正確的 hostname 以及 route.

## 10.3 Package declaration
- 每個 Go file 必定由 Package 宣告開始
- 該 Package Name 也是在其他 Package 中進行存取的 identifier，稱作 last segment convention.
```go=
package main

import (
  
  "fmt"
  "math/rand"
)

func main() {

  fmt.Println(rand.Int())
}
```
- 比如透過 rand 去存取 rand.Int()
- 但 last segment convention 有三個例外情況
  - 第一是 main package，他應作為包成持行檔的程式入口，而不被用於 import
  - 第二是以 \_test 為 suffix 的 package，他主要是給 `go test` 指令使用的，這部分會在後面詳細介紹
  - 第三是帶有版本號的 suffix，比如 "gopkg.in/yaml.v2" 就可以用 yaml 作為 identifier

## 10.4 Import Declaration
- 可以分開寫，也可以寫在一起，通常會寫在一起
```go=
import "fmt"
import "os"
```

```go=
import (
  
  "fmt"
  "os"
)
```

- 可以用空行去分群，但如果空太多行可能你的 Go Compiler 會有意見
```go=
import (

  "fmt"
  "html/template"
  "os"

  "golang.org/x/net/html"
  "golang.org/x/net/ipv4"
)

```

- 如果 last segment 有重複，或是太醜太長你想要簡化，可以同時宣告他的 identifier
```go=
import (

  "crypto/rand"
  mrand "math/rand" // alternative name mrand avoids conflict
)

```

- 要注意如果 Circular import 的話 compiler 會報錯

## 10.5 Blank imports
- 有時候我們需要驅動一個 Package 的初始化，但並不會存取其內部物件
- 為避免 Compiler 判定為 Package imported but not used，就會借助 blank import 來迴避
```go=
// The jpeg command reads a PNG image from the standard input
// and writes it as a JPEG image to the standard output. package main
import (

  "fmt"
  "image"
  "image/jpeg"
  _"image/png" // register PNG decoder
  "io"
  "os"
)

func main() {
  
  if err := toJPEG(os.Stdin, os.Stdout); err != nil {
  
    fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
    os.Exit(1)
  }
}

func toJPEG(in io.Reader, out io.Writer) error {

  img, kind, err := image.Decode(in)
  if err != nil {
    
    return err
  }
  
  fmt.Fprintln(os.Stderr, "Input format =", kind)
  return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

```
- 比如並未使用到 image/png 的內部物件，但仍需要他初始化 png decoder 協助讀檔
```go=
package png // image/png

func Decode(r io.Reader) (image.Image, error)
func DecodeConfig(r io.Reader) (image.Config, error)

func init() {

  const pngHeader = "\x89PNG\r\n\x1a\n"   
  image.RegisterFormat("png", pngHeader, Decode, DecodeConfig)
}
```

- 其他像是 Database 相關的幾個 Package 也有類似的機制
```go=
import (

  "database/mysql"
  _"github.com/lib/pq" // enable support for Postgres 
  _"github.com/go-sql-driver/mysql" // enable support for MySQL
)

db, err = sql.Open("postgres", dbname) // OK
db, err = sql.Open("mysql", dbname) // OK
db, err = sql.Open("sqlite3", dbname) // returns error: unknown driver "sqlite3"
```

## 10.6 Packages and Naming
- Package 命名盡量以簡潔但不至於無法理解的單詞為主
- 也不要與常用的變數名稱重複
- 其實也沒強迫啦

## 10.7 The Go tool
- `go` 是個萬用工具，功能包含 Downloading, formatting, building, testing 等等
```bash
$ go
...
  build compile packages and dependencies
  clean remove object files
  doc show documentation for package or symbol
  env print Go environment information
  fmt run gofmt on package sources
  get download and install packages and dependencies
  install compile and install packages and dependencies
  list list packages
  run compile and run Go program
  test test packages
  version print Go version
  vet run go tool vet on packages

Use "go help [command]" for more information about a command.
...
```

### 10.7.1 Workspace Organization
- 他的 Workspace 大概只有考慮到學校作業等級，像下面這樣
- GO 1.11 前確實需要像他這樣整個專案掛在 GOPATH 下，但現在已經沒人這樣幹了
```
GOPATH/
  src/
    gopl.io/
      .git/
      ch1/
        helloworld/
          main.go
        dup/
          main.go
        ...
    golang.org/x/net/
      .git/
      html/
        parse.go
        node.go
        ...
  bin/
    helloworld
    dup
  pkg/
    darwin_amd64/ 
```
- 比較實用又不會太浮誇的架構可以參考 AppleBOY 在 ModernWeb 的 Session
{%slideshare appleboy/golang-project-layout-and-practice-167597350 %}

### 10.7.2 Downloading Packages
- 最簡單的方法，就是 `go get`，可以加上 -u flag 拉取最新版本
- 除了 Go 官方自己的 host，也有支援 GitHub 和 Launchpad 等等熱門原始碼的網站 package 下載
```bash
$ go get github.com/golang/lint/golint 
```
- 他會拉下來一個 git client repo，所以會有 remote origin 去比較版本差異
```bash
$ cd $GOPATH/src/golang.org/x/net
$ git remote -v
origin https://go.googlesource.com/net (fetch)
origin https://go.googlesource.com/net (push)
```

### 10.7.3 Building Package
- `go build` 可用來檢測目標 Package 是否有編譯錯誤
- 但只有在 Package 為 main 時，才會生成和 .go file 檔名相同的可執行檔
```go
$ cat quoteargs.go
package main

import (

  "fmt"
  "os"
)

func main() {

  fmt.Printf("%q\n", os.Args[1:])
}

$ go build quoteargs.go
$ ./quoteargs one "two three" four\ five
["one" "two three" "four five"]
```
- 或者想要直接執行，沒有要特地留下可執行檔，就用 `go run`
```bash
$ go run quoteargs.go one "two three" four\ five
["one" "two three" "four five"]
```
- 若針對不同編譯情況要用不同檔案的話，可以使用 Build tag，要寫在 Package declaration 前面
- 比如說 linux 和 Mac only
```go
// +build linux darwin
```
- 或者是直接忽略，這很常用於 test scripts
```go
// +build ignore
```

### 10.7.4 Documenting Package
- <del>程式寫的爛才需要註解，優秀的程式能夠超越註解</del>
- 可以用 `go doc` 查看 Package, Member 或 Method 的註解
```bash
$ go doc time
package time // import "time"

Package time provides functionality for measuring and displaying time.

const Nanosecond Duration = 1 ...
func After(d Duration) <-chan Time
func Sleep(d Duration)
func Since(t Time) Duration
func Now() Time
type Duration int64
type Time struct { ... }
... many more ...
```

```bash
$ go doc json.decode
func (dec *Decoder) Decode(v interface{}) error

Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
```

- 看起來他還提供了 Local host 一個 doc host 的功能，但我沒用過，也聽別人提到過，有興趣可以試試看
```
$ godoc -http :8000
```

### 10.7.5 Internal Package
- 看他的描述應該就是封測之類的功能，有些 Package 再正式釋出給社群使用之前，希望透過架構來限制他的影響力
- 有點類似 Chaos Engineering 中的 blast radius limitation
- 在 internal 之下的 Package，將會帶有 import 權限限制
```bash=
net/http
net/http/internal/chunked
net/http/httputil
net/url
```
- `net/http/internal/chunked` 能被 `net/http` 和 `net/http/httputil` import，但 `net/url` 無法

### 10.7.6 Querying Packages
- `go list` 直接查詢 Package 目前是否在 workspace 中
```bash
$ go list github.com/go-sql-driver/mysql
github.com/go-sql-driver/mysql
```
- 或用萬用字符 `.`來列出 workspace 中所有 Package
```bash
$ go list ...
archive/tar
archive/zip
bufio
bytes
cmd/addr2line
cmd/api
... many more ...
```
- 也可以 traverse 一個 Package subtree
```bash
$ go list gopl.io/ch3/...
gopl.io/ch3/basename1
gopl.io/ch3/basename2
gopl.io/ch3/comma
gopl.io/ch3/mandelbrot
gopl.io/ch3/netflag
gopl.io/ch3/printints
gopl.io/ch3/surface
```
- 或是找關鍵字，哦幹這奇葩 convention 也太猛了吧
```bash
$ go list ...xml...
encoding/xml
gopl.io/ch7/xmlselect
```
- 他也能查 Package Hash，如果你的套件行為和你想的不一樣，記得去看看 Hash
```
$ go list -jsonhash
{
  "Dir": "/home/gopher/go/src/hash",
  "ImportPath": "hash",
  "Name": "hash",
  "Doc": "Package hash provides interfaces for hash functions.",
  "Target": "/home/gopher/go/pkg/darwin_amd64/hash.a",
  "Goroot": true,
  "Standard": true,
  "Root": "/home/gopher/go",
  "GoFiles": [
    "hash.go"
  ],
  "Imports": [
    "io"
  ],
  "Deps": [
    "errors",
    "io",
    "runtime",
    "sync",
    "sync/atomic",
    "unsafe"
  ]
}
```
- 而且還支援自訂的 formatting，對 Artifact management 超友善 r
```bash
$ go list -f '{{.ImportPath}} -> {{join .Imports " "}}' 
compress/... compress/bzip2 -> bufio io sort
compress/flate -> bufio fmt io math sort strconv 
compress/gzip -> bufio compress/flate errors fmt hash hash/crc32 io time
compress/lzw -> bufio errors fmt io
compress/zlib -> bufio compress/flate errors fmt hash hash/adler32 io
```

# Testing
- 許多現代軟體已經複雜到難以理解了，我們需要自動化測試確保他們如預期中的運作
- `go test` 就是為了做這件事

## 11.1 The `go test` Tool
- 相關的程式碼要寫在 _test 為 suffix 的 .go file 中，這些將不會被 `go build` 打包進入執行檔
- 其中會被 `go test` 觸發的物件有三個，將會在後續段落獨立詳述
  - Tests，以結果作正確性判斷
  - Benchmark，做效能評估
  - Example
- `go test` 將會蒐集全部 *_test.go files，生成合適的 main Package 編譯並運行，並在回報結果後自行清除

## 11.2 Test Function
- Test Function 必須以 Test 為 prefix，並透過參數 *testing.T 回報結果和 Log
```go=
// Package word provides utilities for word games.
package word

// IsPalindrome reports whether s reads the same forward and backward.
// (Our first attempt.)
func IsPalindrome(s string) bool {

  for i := range s {
    
    if s[i] != s[len(s)-1-i] {
       
      return false
    }
  }
  return true
}
```
```go=
package word

import "testing"

func TestPalindrome(t *testing.T) {

  if !IsPalindrome("detartrated") { 
    
    t.Error(`IsPalindrome("detartrated") = false`)
  }
  if !IsPalindrome("kayak") { 
      
    t.Error(`IsPalindrome("kayak") = false`)
  }
}

func TestNonPalindrome(t *testing.T) {

  if IsPalindrome("palindrome") { 
    
    t.Error(`IsPalindrome("palindrome") = true`)
  }
}

```
- 執行結果
```bash
$ cd $GOPATH/src/gopl.io/ch11/word1
$ go test
ok gopl.io/ch11/word1 0.008s

```
- 然後來試試在 Test case 中放入未定義字元進去觸發錯誤
```go
func TestFrenchPalindrome(t *testing.T) {

  if !IsPalindrome("été") {
    
    t.Error(`IsPalindrome("été") = false`)
  }
}

func TestCanalPalindrome(t *testing.T) {

  input := "A man, a plan, a canal: Panama"
  if !IsPalindrome(input) {
  
    t.Errorf(`IsPalindrome(%q) = false`,input)
  }
}
```
- 可以加上 -v tag 輸出詳細 Log
```bash
$ go test -v
=== RUN TestPalindrome
--- PASS: TestPalindrome (0.00s)
=== RUN TestNonPalindrome
--- PASS: TestNonPalindrome (0.00s)
=== RUN TestFrenchPalindrome
--- FAIL: TestFrenchPalindrome (0.00s)
    word_test.go:28: IsPalindrome("été") = false
=== RUN TestCanalPalindrome
--- FAIL: TestCanalPalindrome (0.00s)
    word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
FAIL
exit status 1
FAIL gopl.io/ch11/word1 0.017s
```
- 或是加上 -run flag 透過 RegExp 指定測試群組
```bash
$ go test -v -run="French|Canal"
=== RUN TestFrenchPalindrome
--- FAIL: TestFrenchPalindrome (0.00s)
    word_test.go:28: IsPalindrome("été") = false
=== RUN TestCanalPalindrome
--- FAIL: TestCanalPalindrome (0.00s)
    word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
FAIL
exit status 1
FAIL gopl.io/ch11/word1 0.014s
```
- 我們根據測試的錯誤，改進功能性函數
```go
// Package word provides utilities for word games.
package word

import "unicode"
// IsPalindrome reports whether s reads the same forward and backward.
// Letter case is ignored, as are non-letters.
func IsPalindrome(s string) bool {

  var letters []rune
  for _, r := range s {
  
    if unicode.IsLetter(r) {
    
      letters = append(letters, unicode.ToLower(r))
    }
  }
  for i := range letters {
  
    if letters[i] != letters[len(letters)-1-i] {
    
      return false
    }
  }
  return true
}
```
- 再用更高標準的測試資料進行測試，Table driven 的測試因為擴展性高所以蠻實用的
```go
func TestIsPalindrome(t *testing.T) {

  var tests = []struct {
    
    input string
    want bool
  }{
    {"", true},
    {"a", true},
    {"aa", true},
    {"ab", false},
    {"kayak", true},
    {"detartrated", true},
    {"A man, a plan, a canal: Panama", true},
    {"Evil I did dwell; lewd did I live.", true},
    {"Able was I ere I saw Elba", true},
    {"été", true},
    {"Et se resservir, ivresse reste.", true},
    {"palindrome", false}, // non-palindrome
    {"desserts", false}, // semi-palindrome
  }
  
  for _, test := range tests {
  
    if got := IsPalindrome(test.input); got != test.want {
    
      t.Errorf("IsPalindrome(%q) = %v", test.input, got)
    }
  }
}

```

### 11.2.1 Randomized Testing
- 有時手寫 Test case 曠日廢時，還會有開發者偏誤，不如用隨機生成
```go=
import "math/rand"

func randomPalindrome(rng *rand.Rand) string {

  n := rng.Intn(25)
  runes := make([]rune, n)
  for i := 0; i < (n+1)/2; i++ {
  
    r := rune(rng.Intn(0x1000))
    runes[i] = r runes[n-1-i] = r
  }
  return string(runes)
}

func TestRandomPalindromes(t *testing.T) {

  seed := time.Now().UTC().UnixNano()
  t.Logf("Random seed: %d", seed)
  rng := rand.New(rand.NewSource(seed))
  for i := 0; i < 1000; i++ {
  
    p := randomPalindrome(rng)
    if !IsPalindrome(p) {
    
      t.Errorf("IsPalindrome(%q) = false", p)
    }
  }
}

```

### 11.2.2 Testing a Command
- 基本上就是上面的東西換個例子，所以先跳過

### 11.2.3 White-Box Testing
- 白盒測試和黑盒測試的主要差別就是和功能性程式碼放在同個 Package，因此能有權限存取非公開內容
```go=
package storage

import (

  "fmt"
  "log"
  "net/smtp"
)

func bytesInUse(username string) int64 { return 0 /* ... */ }

const sender = "notifications@example.com"
const password = "correcthorsebatterystaple"
const hostname = "smtp.example.com"
const template = `Warning: you are using %d bytes of storage, %d%% of your quota.`

var notifyUser = func(username, msg string) {

  auth := smtp.PlainAuth("", sender, password, hostname)
  err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
  if err != nil {
  
    log.Printf("smtp.SendEmail(%s) failed: %s", username, err)
  }
}

func CheckQuota(username string) {

  used := bytesInUse(username)
  const quota = 1000000000
  percent := 100 * used / quota
  if percent < 90 {
  
    return
  }
  msg := fmt.Sprintf(template, used, percent) notifyUser(username, msg)
}
```
```go=
package storage

import (

  "strings"
  "testing"
)

func TestCheckQuotaNotifiesUser(t *testing.T) {

  var notifiedUser, notifiedMsg string
  notifyUser = func(user, msg string) {
  
    notifiedUser, notifiedMsg = user, msg
  }
  // ...simulate a 980MB-used condition...
  const user = "joe@example.org"
  CheckQuota(user)
  if notifiedUser == "" && notifiedMsg == "" {
  
    t.Fatalf("notifyUser not called")
  }
  if notifiedUser != user {
  
    t.Errorf("wrong user (%s) notified, want %s", notifiedUser, user)
  }
  
  const wantSubstring = "98% of your quota"
  if !strings.Contains(notifiedMsg, wantSubstring) { 
  
    t.Errorf("unexpected notification message <<%s>>, "+ "want substring %q", notifiedMsg, wantSubstring)
  }
}
```
- 白盒測試可能會汙染功能性 Package 的內容，因此要把 clean 做好

### 12.2.4 External Test Package
- 這部分在講如何用 *_test Package 解開 circular import
- 以及黑盒測試造成的存取受限問題
- 基本上都能很直覺的解開，加上它文字內容太抽象了所以先跳過

### 11.2.5 Writing EffectiveTests
- Go 沒有 assert，它們覺得這是開發者自己該 Handle 的部分
- 一個良好的 asserting 機制能提供更多有用的資訊
```go=
func TestSplit(t *testing.T) {

  s, sep := "a:b:c", ":"
  words := strings.Split(s, sep)
  if got, want := len(words), 3; got != want {
  
    t.Errorf("Split(%q, %q) returned %d words, want %d", s, sep, got, want)
  }
}
```

### 11.2.6 Avoiding Brittle Tests
- 若你的測試只要有一點點的功能微調就會出錯，那對維護者的體驗也蠻糟的，會消耗大量時間在處理測試上
- 不需要地毯性的檢測每個值，只需要測試會有影響的關鍵部分和步驟即可，因為中間實作和輸入結構都是很可能頻繁微調的

## 11.3 Coverage
- 測試的覆蓋率是評估測試水平的一個指標，可以用 `go test` 的 -coverprofile 產出覆蓋率報告
- 覆蓋率不代表一切，測試的品質及彈性也在開發過程中有很大的影響

```bash=
$ go test -run=Coverage -coverprofile=c.out 
gopl.io/ch7/eval ok gopl.io/ch7/eval 0.032s coverage: 68.5% of statements
```
- 然後可以套入 html 模板呈現
```bash=
$ go tool cover -html=c.out
```
![](https://i.imgur.com/9WeVVAg.png)

## 11.4 Benchmark Functions
- Benchmark 可用於測試效能，透過參數 \*testing.B 蒐集資料並回報
```go=
import "testing"

func BenchmarkIsPalindrome(b *testing.B) {

  for i := 0; i < b.N; i++ {
  
    IsPalindrome("A man, a plan, a canal: Panama")
  }
}
```
- 比較特別的是，一定要用 -bench flag 以 RegExp 指定執行的 benchmark，不然他預設不會執行任何 benchmark
```bash=
$ cd $GOPATH/src/gopl.io/ch11/word2
$ go test -bench=.
PASS BenchmarkIsPalindrome-8 1000000 1035 ns/op
ok gopl.io/ch11/word2 2.179s
```
- 加上 -benchmem flag 顯示記憶體的開銷
```bash=
$ go test -bench=. -benchmem
PASS BenchmarkIsPalindrome 2000000 807 ns/op 128 B/op 1 allocs/op
```
- 若想要一次比較多個執行次數，比較時間和記憶體開銷的增長曲線，可以這樣寫
```go=
func benchmark(b *testing.B, size int) { /* ... */ }
func Benchmark10(b *testing.B) { benchmark(b, 10) }
func Benchmark100(b *testing.B) { benchmark(b, 100) }
func Benchmark1000(b *testing.B) { benchmark(b, 1000) }
```

## 11.5 Profiling
- 當你測試後發現你的程式效能很糟，然後思考要從哪裡下手優化
![](https://i.imgur.com/Cqt6oPv.jpg)
- Profiling 提供了整體程式執行上，更詳細的時間和資源開銷
- 他也是整合進了 `go test` 裡面，主要分為三種
  - cpu profile 用於檢查運算時間開銷
  - heap profile 用於檢查記憶體開銷
  - blocking profile 用於檢查 goroutine 堵塞的時間
```bash=
$ go test -cpuprofile=cpu.out
$ go test -blockprofile=block.out
$ go test -memprofile=mem.out
```
- 而解析 profile 結果的是另個工具 `pprof`，可以把 `go test` 產生的 log 留下來丟給他
```bash=
$ go test -run=NONE -bench=ClientServerParallelTLS64 \
  -cpuprofile=cpu.log net/http
PASS 
BenchmarkClientServerParallelTLS64-8 1000
3141325 ns/op 143010 B/op 1747 allocs/op
ok net/http 3.395s

$ go tool pprof -text -nodecount=10 ./http.test cpu.log
2570ms of 3590ms total (71.59%)
Dropped 129 nodes (cum <= 17.95ms)
Showing top 10 nodes out of 166 (cum >= 60ms)
  flat flat% sum% cum cum%
  1730ms 48.19% 48.19% 1750ms 48.75% crypto/elliptic.p256ReduceDegree
  230ms 6.41% 54.60% 250ms 6.96% crypto/elliptic.p256Diff
  120ms 3.34% 57.94% 120ms 3.34% math/big.addMulVVW
  110ms 3.06% 61.00% 110ms 3.06% syscall.Syscall
  90ms 2.51% 63.51% 1130ms 31.48% crypto/elliptic.p256Square
  70ms 1.95% 65.46% 120ms 3.34% runtime.scanobject
  60ms 1.67% 67.13% 830ms 23.12% crypto/elliptic.p256Mul 
  60ms 1.67% 68.80% 190ms 5.29% math/big.nat.montgomery 
  50ms 1.39% 70.19% 50ms 1.39% crypto/elliptic.p256ReduceCarry
  50ms 1.39% 71.59% 60ms 1.67% crypto/elliptic.p256Sum
```
- -nodecount flag 指定 query 的數量，然後你發現 crypto 套件裡的東西就是拖慢效能的主因
- 從效能占比最大的部分開始著手優化，能有更好的投資報酬率

## 11.6 Example Functions
- 透過 Example prefix 可以在上面 10.7.4 提到的 Doc host 中直接把程式碼引入，和功能性函數放在一起
- R 不過我好像真的沒看其他人用過這個 @@
```go=
func ExampleIsPalindrome() {

  fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
  fmt.Println(IsPalindrome("palindrome"))
  // Output:
  // true
  // false
}
```