# Ch6 Methods
from [HackMD](https://hackmd.io/@HexRabbit/ByObRGcUL)

## 6.1 Method Declarations
Go 的 method 宣告方法與一些常見的語言如 C++、Java、Python 不同，一般 method 可能會被宣告在 class 之中，而 Go 的作法則是在宣告函數時多加入一個 *receiver* 來指定這是對應哪個 **Type** 的 method
```go
package geometry

import "math"
type Point struct{ X, Y float64 }

// traditional function
func Distance(p, q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```
`p Point` 被稱為該 method 的 *receiver*，`p` 在這裡其實就是個變數名稱，當然你要取成 `this`、`self` 也可以，只是書中提到 Go 不推薦這樣做的原因就是 `p` 極有可能在函數中出現多次，所以當然越短越省時間


除了 interface 與 pointer (下方會提到)，在 Go 中你幾乎可以為任意型別宣告 method
可以注意到下方 Path 的型別只是個 Slice (!)
```go
// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
    sum := 0.0
    for i := range path {
        if i > 0 {
            sum += path[i-1].Distance(path[i])
        }
    }
    return sum
}
```

雖然在宣告時 Distance 這個函數名稱被宣告了兩次，但因為對應的型別不同所以是能夠正常 compile 的

## 6.2 Methods with a Pointer Receiver

在 Go 中，你當然也可以為指標型別宣告 methods，不過在語言設計上為了避免混淆，要求使用者不能為指標型別的 *named types* 宣告 methods

```go
type P *int
func (P) f() { /* ... */ } // compile error: invalid receiver type
```

而正確的寫法則是

```go
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}
```

在呼叫指標型別的 method 時，因為編譯器很聰明，所以基本上不需要顯式的取址/取值之後再呼叫

```go
pptr := &Point{1, 2}
p := Point{1, 2}

pptr.ScaleBy(2)
pptr.Distance(q) // implicit (*pptr)
fmt.Println(*pptr) // "{2, 4}"

p.ScaleBy(2) // implicit (&p)
p.Distance(q)
fmt.Println(p) // "{2, 4}"
```

但有些例外，例如編譯器沒辦法幫 temporary value 取地址 (卻可以幫右值取地址??)

```go
Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal
(&Point{1, 2}).ScaleBy(2) // well, you can do this
```

此外就算指標變數的值為 `nil`，你也可以呼叫該變數的 method，
書中提到若是 method 接受 `nil` 則最好在文件中明確註記
```go
// An IntList is a linked list of integers.
// A nil *IntList represents the empty list.
type IntList struct {
    Value int
    Tail *IntList
}

// Sum returns the sum of the list elements.
func (list *IntList) Sum() int {
    if list == nil {
        return 0
    }
    return list.Value + list.Tail.Sum()
}
```

## 6.3 Composing Types by Struct Embedding

本章節是 4.4.3 中 struct embedding 的進階用法，對語法不熟的話請回去複習一下

```go
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
    Point
    Color color.RGBA
}
```

考慮上方這個 embedded structure，因為 Point 被 embedded 在 ColoredPoint 裡面，
所有 Point 的 method 會被 *promote* 至 ColoredPoint，也就是 method receiver type 從 Point 改為 ColoredPoint，這讓我們可以直接使用 ColoredPoint 呼叫那些定義於 Point 中的 method


當然若是要使用 embedded 在其中的 Point，需要顯式將他取出來編譯器才知道你要使用該變數 (下方的 `q.Point`)
> 在 Go 裡的 struct 彼此之間只會有 "has a" 的關係

```go
red := color.RGBA{255, 0, 0, 255}
blue := color.RGBA{0, 0, 255, 255}

var p = ColoredPoint{Point{1, 1}, red}
var q = ColoredPoint{Point{5, 4}, blue}

fmt.Println(p.Distance(q.Point)) // "5"

p.ScaleBy(2) 
q.ScaleBy(2)

fmt.Println(p.Distance(q.Point)) // "10"
```

也可以 embed 一個指標型別進去

```go
type ColoredPoint struct {
    *Point
    Color color.RGBA
}

p := ColoredPoint{&Point{1, 1}, red}
q := ColoredPoint{&Point{5, 4}, blue}
fmt.Println(p.Distance(*q.Point)) // "5"

q.Point = p.Point // p and q now share the same Point

p.ScaleBy(2)
fmt.Println(*p.Point, *q.Point) // "{2 2} {2 2}"
```

Go 允許使用者在 struct 內 embed 多個 anonymous field，由於變數名稱可能會衝突，compiler 在 resolve method/field (如 p.ScaleBy) 的時候會順著階層找，先從該 struct 的 method/field 找起，接著是被 *promote* 一次的，再來是被 *promote* 兩次的...以此類推

一旦在同一個階層裡找到相同名稱的 field/method 編譯器就會報錯，如下
```go
package main
import (
    "fmt"
)

type A struct { }
type B struct { }
type C struct {
    A
    B
}

func (a A) Print() {
    fmt.Println("A")
}
func (b B) Print() {
    fmt.Println("B")
}

func main() {
    c := C{
        A: A{},
        B: B{},
    }
    c.Print()
}
```
```
$ go run test1.go
./test1.go:25:4: ambiguous selector c.Print
```
或是 field 與 method 名稱重複也會報錯
```
$ go run test2.go
./test2.go:15:6: type A has both field and method named Print
```
> 有趣的是只要不用到重複的 field/method 編譯器就不吭聲
> [color=orange][name=HexRabbit]
 
書中還提到一個有趣的 trick，利用 *unnamed struct* 與 *struct embedding* 來提高 code 可讀性
~~其實就是 inheritance 嘛~~

```go
var (
    mu sync.Mutex // guards mapping
    mapping = make(map[string]string)
)

func Lookup(key string) string {
    mu.Lock()
    v := mapping[key]
    mu.Unlock()
    return v
}
```

```go
var cache = struct {
    sync.Mutex
    mapping map[string]string
} {
    mapping: make(map[string]string),
}

func Lookup(key string) string {
    cache.Lock()
    v := cache.mapping[key]
    cache.Unlock()
    return v
}
```

## 6.4 Method Values and Expressions

由於 Function 在 Go 中是 first-class citizen，我們其實可以把呼叫 method 這個動作拆開來成為: 選擇 method(`var.Method`) 及 呼叫(`(...)`)
`var.Method` 會回傳一個 *method value*，是一個將 `Type.Method` 綁定在該 *receiver* `p` 上的 Function
```go
p := Point{1, 2}
q := Point{4, 6}

distanceFromP := p.Distance // method value
fmt.Println(distanceFromP(q)) // "5"
var origin Point // {0, 0}
fmt.Println(distanceFromP(origin)) // "2.23606797749979", ;5

scaleP := p.ScaleBy // method value
scaleP(2) // p becomes (2, 4)
scaleP(3) // then (6, 12)
scaleP(10) // then (60, 120)
```

這提供一些好處，例如說最明顯的就是縮短程式碼

```go
type Rocket struct { /* ... */ }
func (r *Rocket) Launch() { /* ... */ }
r := new(Rocket)

time.AfterFunc(10 * time.Second, func() { r.Launch() })
time.AfterFunc(10 * time.Second, r.Launch) // shorter
```

此外，若是不想使用 OOP 的方式呼叫 Method 的話，也可以將 Method 當作一般 Function 使用，只需要在最一開始提供要操作的物件作為參數即可

```go
p := Point{1, 2}
q := Point{4, 6}

distance := Point.Distance // method expression
fmt.Println(distance(p, q)) // "5"
fmt.Printf("%T\n", distance) // "func(Point, Point) float64"

scale := (*Point).ScaleBy // supply '*' to select pointer version
scale(&p, 2)
fmt.Println(p) // "{2 4}"
fmt.Printf("%T\n", scale) // "func(*Point, float64)"
```

## 6.5 Example: Bit Vector Type

沒啥好說的，請各位回去看完本章 (照理來說應該都看過了) 之後，把後面的 Exercise 當回家作業
有時間的話可以五題都做一下，沒空就做 Exercise 6.1, 6.2, 6.4

## 6.6 Encapsulation

封裝是 OOP 相當重要的一項特性，透過 expose/hide Methods 讓使用者可以專注在 API 的功能上，同時隱藏了實作細節讓開發者能夠有更大的彈性。在 Go 裡的區分相當簡單： 首字母大寫的將會被 exposed 給使用者，其餘隱藏起來 (也就是使用者完全無法存取)，這不僅適用在 Field 上，同時 Method 也是一樣

不過正因為這個設計，Go 的封裝技巧只適用在 struct 上

```go
type Counter struct { n int }
func (c *Counter) N() int { return c.n }
func (c *Counter) Increment() { c.n++ }
func (c *Counter) Reset() { c.n = 0 }
```

> 講得好像很厲害可是沒啥有趣的內容lol
 
本章最後提到，雖然 Method 對 OOP 來說相當重要，但這也只是 half the picture 罷了，重頭戲是下周的 Interface 章節，還請靜待下回分曉！
