# Ch7 Interfaces
from [HackMD](https://hackmd.io/@HexRabbit/HylNA8Dw8)

## 7.1 Interfaces as Contracts

到目前為止書中所介紹的型別都是 *concrete* types，而 interface 則是一種 *abstract* type，他只能抽象的描述物件能做甚麼，而無法得知實際的實作細節 (更明白地說，你只能得知 interface 有甚麼 Method 可以使用)

```go
package fmt

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)

func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}

func Sprintf(format string, args ...interface{}) string {
    var buf bytes.Buffer
    Fprintf(&buf, format, args...)
    return buf.String()
}
```
> 我們晚點再來談 `interface{}`
> [color=green]

上方 `Fprintf` 的第一個參數其實不是 `bytes.Buffer` 這個型別，而是一個 interface `io.Writer`，來看他的宣告
```go
package io

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    // Write writes len(p) bytes from p to the underlying data stream.
    // It returns the number of bytes written from p (0 <= n <= len(p))
    // and any error encountered that caused the write to stop early.
    // Write must return a non-nil error if it returns n < len(p).
    // Write must not modify the slice data, even temporarily.
    //
    // Implementations must not retain p.
    Write(p []byte) (n int, err error)
}
```

這裡宣告了 `Writer` interface，規範任何能被該 interface 所接受的型別提供了 `Write` 這個 Method，

我們也可以把這句話反過來說，任何提供 `Write` method 的型別都將滿足這個 interface

這像是在函數與呼叫者之間綁定一紙合約，對呼叫者來說他只能傳入型別符合規範的物件，對函數本身來說不需要去考慮被傳入的物件是如何實做的，只要符合規範函數就會正確運作

這可以帶來一些有趣的應用，例如以下
```go
type ByteCounter int
func (c *ByteCounter) Write(p []byte) (int, error) {
    *c += ByteCounter(len(p)) // convert int to ByteCounter
    return len(p), nil
}
```
```go
var c ByteCounter

c.Write([]byte("hello"))
fmt.Println(c) // "5", = len("hello")

c = 0 // reset the counter

var name = "Dolly"

fmt.Fprintf(&c, "hello, %s", name)
fmt.Println(c) // "12", = len("hello, Dolly")
```

還有我們熟悉的 `String()` method，只要型別提供能夠被 `fmt.Stringer` 接受，`fmt.Print` 就會使用我們定義的 `String()` method 來印出物件的內容
```go
package fmt
// The String method is used to print values passed
// as an operand to any format that accepts a string
// or to an unformatted printer such as Print.
type Stringer interface {
    String() string
}
```
> 至於 fmt 中是如何得知變數有沒有符合 `fmt.Stringer` 這個 interface，在後面的章節會提到
 
## 7.2 Interface Types

```go
package io

type Reader interface {
    Read(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}
```
當然也可以做 embedding

```go
type ReadWriter interface {
    Reader
    Writer
}
type ReadWriteCloser interface {
    Reader
    Write(p []byte) (n int, err error)
    Close() error
}
```

## 7.3 Interface Satisfaction

當一個型別 *satisfies* 一個 interface，這表示該型別提供所有在 interface 中宣告的 method，在 Go 中我們用 "is" 來表達型別與 interface 之間的關係，而能不能滿足該條件是在 compile time 被決定

根據 *assignability rule* (§2.4.2)，只有滿足該 interface 的型別的變數才能被賦值
```go
var w io.Writer
w = os.Stdout // OK: *os.File has Write method
w = new(bytes.Buffer) // OK: *bytes.Buffer has Write method
w = time.Second // compile error: time.Duration lacks Write method

var rwc io.ReadWriteCloser
rwc = os.Stdout // OK: *os.File has Read, Write, Close methods
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method
```

一些初學者常犯的錯誤
```go
type IntSet struct { /* ... */ }
func (*IntSet) String() string
var _ = IntSet{}.String() // compile error: String requires *IntSet receiver

var s IntSet
var _ = s.String() // OK: s is a variable and &s has a String method

var _ fmt.Stringer = &s // OK
var _ fmt.Stringer = s // compile error: IntSet lacks String method
```

同時就算原先的型別提供多個 method，interface 會限制使用者只能存取 interface 中宣告的 method

```go
os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
os.Stdout.Close() // OK: *os.File has Close method

var w io.Writer
w = os.Stdout

w.Write([]byte("hello")) // OK: io.Writer has Write method
w.Close() // compile error: io.Writer lacks Close method
```

那如果我宣告一個不包含任何 method 的 interface 呢?
```go
var any interface{}
any = true
any = 12.34
any = "hello"
any = map[string]int{"one": 1}
any = new(bytes.Buffer)
```

看起來好像相當無用，雖然我們可以放入任何值到該變數中，但因為 `interface{}` 沒有包含任何 method，所以基本上甚麼操作都做不了

不過可以看到 `Printf` 中為了要處理使用者餵進的各式各樣的型別，利用空 interface `interface{}` 去接參數，所以這肯定是有用的

> 當然不可能接完參數不做事，Go 也提供了方法將其中的型別抽出來，不過這讓我們留到後面再談

另外，interface 也不失為一個有效的方式去分類型別，考慮用來整理各種數位媒體的一支程式

```go
/*
Album
Book
Movie
Magazine
Podcast
TVEpisode
Track
*/

type Artifact interface {
    Title() string
    Creators() []string
    Created() time.Time
}

type Text interface {
    Pages() int
    Words() int
    PageSize() int
}

type Audio interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string // e.g., "MP3", "WAV"
}

type Video interface {
    Stream() (io.ReadCloser, error)
    RunningTime() time.Duration
    Format() string // e.g., "MP4", "WMV"
    Resolution() (x, y int)
}
```

> 回想我們剛剛用 "is a" 去表達型別與 interface 的關係
 
## 7.4 Parsing Flags with `flag.Value`
這章主要介紹 `flag` 這個 package 中的各種好用 function 與 interface

```go
var period = flag.Duration("period", 1*time.Second, "sleep period")
func main() {
    flag.Parse()
    fmt.Printf("Sleeping for %v...", *period)
    time.Sleep(*period)
    fmt.Println()
}
```
```
$ ./sleep -period 50ms
Sleeping for 50ms...
$ ./sleep -period 2m30s
Sleeping for 2m30s...
$ ./sleep -period 1.5h
Sleeping for 1h30m0s...
$ ./sleep -period "1 day"
invalid value "1 day" for flag -period: time: invalid duration 1 day
```

看起來很神奇對吧，而且我們也可以註冊自定義的 commandline flag，只要提供 `flag.CommandLine.Var` 符合 `flag.Value` 這個 interface 的型別即可

```go
package flag
// Value is the interface to the value stored in a flag.
type Value interface {
    String() string
    Set(string) error
}
```

```go
// *celsiusFlag satisfies the flag.Value interface.
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
    var unit string
    var value float64
    fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
    
    switch unit {
    case "C", "°C":
        f.Celsius = Celsius(value)
        return nil
    case "F", "°F":
        f.Celsius = FToC(Fahrenheit(value))
        return nil
    }
    return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag defines a Celsius flag with the specified name,
// default value, and usage, and returns the address of the flag variable.
// The flag argument must have a quantity and a unit, e.g., "100C".
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
    f := celsiusFlag{value}
    flag.CommandLine.Var(&f, name, usage)
    return &f.Celsius
}

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")
func main() {
    flag.Parse()
    fmt.Println(*temp)
}
```

## 7.5 Interface Values

概念上，所謂 *interface value* 可以被拆成一個型別以及其對應的值，分別是: *dynamic type*, *dynamic value*
若要我來說的話，interface 變數本身就像是一個容器，可以乘載各種型別與數值

考慮宣告 `w` 這個 interface 變數，其 type 與 value 初始值都會是 nil
```go
var w io.Writer
```

![](https://i.imgur.com/puK19Ry.png)

因為這時候 type 的部分是 nil，我們可以說該 interface 是空(nil) 的，可以用 `w == nil` 檢驗

對 nil interface 呼叫 method 會觸發 panic
```go
w.Write([]byte("hello")) // panic: nil pointer dereference
```

若是我們對他賦值則會改變 type 及 value
```go
w = os.Stdout
```

![](https://i.imgur.com/4UHDHuJ.png)

當然若是呼叫 `Write` method，則會印出 "hello"
```go
w.Write([]byte("hello")) // "hello"
```

看到這大家可能有點疑惑 value 部分是不是都存著指標指向物件，實際上 value 可以存下任意大小的值，例如說
```go
var x interface{} = time.Now()
```
![](https://i.imgur.com/hjyxhJy.png)

interface 之間也是可以用 `a == b` 比較的，會先比較 *dynamic type* 再來才比較 *dynamic value*，若是兩個屬性皆相同(且都可以比較)才回傳 `true`

```go
var x interface{} = []int{1, 2, 3}
fmt.Println(x == x) // panic: comparing uncomparable type []int
```

可以在 `fmt.Printf` 中用 "%T" 印出 interface 的 *dynamic type* 
```go
var w io.Writer
fmt.Printf("%T\n", w) // "<nil>"

w = os.Stdout
fmt.Printf("%T\n", w) // "*os.File"

w = new(bytes.Buffer)
fmt.Printf("%T\n", w) // "*bytes.Buffer"
```

### 7.5.1 Caveat: An Interface Containing a Nil Pointer Is Non-Nil

以下是個錯誤示範

```go
const debug = true
func main() {
    var buf *bytes.Buffer
    
    if debug {
        buf = new(bytes.Buffer) // enable collection of output
    }
    
    f(buf) // NOTE: subtly incorrect!
    
    if debug {
    // ...use buf...
    }
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
// ...do something...
    if out != nil {
        out.Write([]byte("done!\n"))
    }
}
```

因為 `f()` 中的變數 `out` 帶有 `*bytes.Buffer` 這個型別，所以 `out != nil` 依然是 `true`
![](https://i.imgur.com/Kz0qxaa.png)

> 有夠他媽混淆的

至於修正的方法則是宣告 buf 成 `io.Writer`，這樣
```go
var buf io.Writer
if debug {
    buf = new(bytes.Buffer) // enable collection of output
}
f(buf) // OK
```

## 7.6 Sorting with sort.Interface

和其他語言相同，Go 也提供了通用的 sorting function
只要滿足 `sort.Interface` 就可以使用

```go
package sort
type Interface interface {
    Len() int
    Less(i, j int) bool // i, j are indices of sequence elements
    Swap(i, j int)
}
```
以 `[]string` 為範例
```go
type StringSlice []string
func (p StringSlice) Len() int { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
```
```go
sort.Sort(StringSlice(names))
```

接著是一個我覺得還不錯的範例
```go
type Track struct {
    Title string
    Artist string
    Album string
    Year int
    Length time.Duration
}

var tracks = []*Track{
    {"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
    {"Go", "Moby", "Moby", 1992, length("3m37s")},
    {"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
    {"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
    d, err := time.ParseDuration(s)
    if err != nil {
        panic(s)
    }
    return d
}

func printTracks(tracks []*Track) {
    const format = "%v\t%v\t%v\t%v\t%v\t\n"
    tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
    fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
    fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
    for _, t := range tracks {
        fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
    }
    tw.Flush() // calculate column widths and print table
}
```
```go
type byArtist []*Track
func (x byArtist) Len() int { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
```
```go
sort.Sort(byYear(tracks))
```
```
Title Artist Album Year Length
----- ------ ----- ---- ------
Go Ahead Alicia Keys As I Am 2007 4m36s
Go Delilah From the Roots Up 2012 3m38s
Ready 2 Go Martin Solveig Smash 2011 4m24s
Go Moby Moby 1992 3m37s
```

當然對於比較常見的資料型態如 `[]int` sort package 中就有提供對應的函數
```go
values := []int{3, 1, 4, 1}
fmt.Println(sort.IntsAreSorted(values)) // "false"

sort.Ints(values)
fmt.Println(values) // "[1 1 3 4]"
fmt.Println(sort.IntsAreSorted(values)) // "true"

sort.Sort(sort.Reverse(sort.IntSlice(values)))
fmt.Println(values) // "[4 3 1 1]"
fmt.Println(sort.IntsAreSorted(values)) // "false"
```

## 7.7 The http.Handler Interface

這章也是 Go 的火力展示，主要著重在介紹 `http.Handler` 的各種使用姿勢
礙於篇幅我只會擷取部分比較有學習價值的程式碼

```go
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}
func ListenAndServe(address string, h Handler) error
```

幫一個 db 加上 `ServeHTTP` method 就可以被拿來回應 http request 了LOL
```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32
type database map[string]dollars

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price) // http.ResponseWriter satisfies io.Writer
    }
}
```

若是想要做出更複雜的 routing 設定的話，可以利用 `net/http` 裡的 `ServeMux`
```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    mux := http.NewServeMux()
    mux.Handle("/list", http.HandlerFunc(db.list)) // 注意這邊
    mux.Handle("/price", http.HandlerFunc(db.price))
    log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price)
    }
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("item")
    price, ok := db[item]
    if !ok {
        w.WriteHeader(http.StatusNotFound) // 404
        fmt.Fprintf(w, "no such item: %q\n", item)
        return
    }
    fmt.Fprintf(w, "%s\n", price)
}
```

因為 `db.list` 及 `db.price` 兩者的型別是 `func(http.ResponseWriter, *http.Request)`，但是 `mux.Handle` 要求的是滿足 `http.Handler` interface 的型別，所以我們勢必需要將想執行的 callback function 利用 `http.HandlerFunc(db.list)` 包裝起來，讓他成為滿足 `http.Handler` 的型別

大家到這邊可以先暫停一下，稍微思考這個 wrapper 可能是怎麼實作的
> 提供一個想法，例如說利用 struct 去存 callback，然後再加上 method

以下公布答案

.

.

.

.

.

.

.

.

.

.

.

.

.

.

.

.

.

.

```go
package http
type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```
> The expression `http.HandlerFunc(db.list)` is a conversion, not a function call, since `http.HandlerFunc` is a type.

![](https://i.imgur.com/sf8kt7u.png)

## 7.8 The error Interface

從書的第一章開始我們便開始使用 error 這個神奇的 type，現在是時候來介紹他了

其實就只是個單純的 interface，提供一個 `Error()` 用以輸出錯誤訊息
```go
type error interface {
    Error() string
}
```
可以用 `errors.New()` 取得一個 error
```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct { text string }
func (e *errorString) Error() string { return e.text }
```

更常用的是 fmt package 裡的 `Errorf` wrapper function
```go
package fmt
import "errors"

func Errorf(format string, args ...interface{}) error {
    return errors.New(Sprintf(format, args...))
}
```

## 7.9 Example: Expression Evaluator

沒啥東西，跳過

## 7.10 Type Assertions

*type assertion* 是一個作用在 interface 的操作，語法上可以表示成 `x.(T)`，其中 `x` 是 interface 變數，`T` 則是一個型別

當 `T` 是 *concrete type* 時，
這個操作將會檢查 `x` 的 *dynamic type* 是否與 `T` 一致，若成功則回傳 `x` 的 *dynamic value* (型別當然為 `T`)

這通常被用來判斷 interface 裡存放的 type
```go
var w io.Writer
w = os.Stdout
f := w.(*os.File) // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer
```

當 `T` 是 *interface type* 時，
這個操作會檢查 `x` 的 *dynamic type* 是否能夠滿足 `T`，若成功則回傳型別為 `T` 的 interface，其中 *dynamic type* 和 *dynamic value* 皆與原本相同

這通常被用來新增受 interface 限制而無法使用的 method
```go
var w io.Writer
w = os.Stdout
rw := w.(io.ReadWriter) // success: *os.File has both Read and Write
w = new(ByteCounter)
rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method
```

也有失敗時不會 panic 的寫法
```go
var w io.Writer = os.Stdout
f, ok := w.(*os.File) // success: ok, f == os.Stdout
b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil
```

記得使用簡潔的 if 寫法
```go
if f, ok := w.(*os.File); ok {
// ...use f...
}
```

## 7.11 Discriminating Errors with Type Assertions

透過 *type assertion* 我們便可以判斷各種不同 error

```go
package os
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op string
    Path string
    Err error
}
func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

這也是另一種符合 error interface 的實作
```go
package syscall
type Errno uintptr // operating system error code
func (err *Errno) Error() string {
    ...
}
```

透過 *type assertion* 確定型別，並將內部的 error code 取出來
```go
import (
    "errors"
    "syscall"
)

var ErrNotExist = errors.New("file does not exist")

// IsNotExist returns a boolean indicating whether the error is known to
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.
func IsNotExist(err error) bool {
    if pe, ok := err.(*PathError); ok {
        err = pe.Err
    }
    return err == syscall.ENOENT || err == ErrNotExist
}
```

## 7.12 Querying Behaviors with Interface Type Assertions

*interface type assertion* 常會被用在確認 interface 中存放的物件是否有額外的 method 可以使用
```go
// writeString writes s to w.
// If w has a WriteString method, it is invoked instead of w.Write.
func writeString(w io.Writer, s string) (n int, err error) {
    type stringWriter interface {
        WriteString(string) (n int, err error)
    }
    if sw, ok := w.(stringWriter); ok {
        return sw.WriteString(s) // avoid a copy
    }
    return w.Write([]byte(s)) // allocate temporary copy
}

func writeHeader(w io.Writer, contentType string) error {
    if _, err := writeString(w, "Content-Type: "); err != nil {
        return err
    }
    if _, err := writeString(w, contentType); err != nil {
        return err
    }
    // ...
}
```
這裡其實就是透過定義 `stringWriter` interface 額外做了一個檢查: 「`w` 是否有額外包含 `WriteString` 這個 method？」

## 7.13 Type Switches

為了讓大家寫 type assertion 時不需要不斷的透過 if...else 去判斷型別(尤其是 `interface{}`)，可以利用 `switch x.(type) { /* ... */ }` 語法建立一個 type switch 來更方便的處理
```go
func sqlQuote(x interface{}) string {
    switch x := x.(type) {
    case nil:
        return "NULL"
    case int, uint:
        return fmt.Sprintf("%d", x)
    case bool:
        if x {
            return "TRUE"
        }
        return "FALSE"
    case string:
        return sqlQuoteString(x) // (not shown)
    default:
        panic(fmt.Sprintf("unexpected type %T: %v", x, x))
    }
}
```

## 7.15. A Few Words of Advice

作者給開發者的一些建議
- 不要宣告只被一種實作滿足的 interface
- 上述的唯一例外是需要 [decouple two packages](https://stackoverflow.com/questions/30784215/what-does-decoupling-two-classes-at-the-interface-level-mean) 的情況
- 若是要限制存取請愛用大小寫區分的 export 手法
- 一個好的 interface 設計理念: ask only for what you need
- 不是甚麼都需要 OOP，請用腦思考你要什麼
