---
tags: CCNS
---

# 讀書會 Week 2

# CH3 Basic Data Type

# 3.1 Integers

## int, uint

- 都有 8, 16, 32, 64 位元版本
- 數值範圍跟 C 一樣用 2 補數去推
  - ex: int8 的數值範圍: $-2^{8-1} \sim 2^{8-1}-1= -128 \sim 127$


```go=
//overflow
package main

import (
	"fmt"
)
func main() {
     var u uint8 = 255
     fmt.Println(u, u+1, u*u) // "255 0 1"
     var i int8 = 127
     fmt.Println(i, i+1, i*i) // "127 -128 1"
}
```
## rune, byte, uintptr

- 兩者都是 data type 的別名 (alias)
- rune 代表 **int32**, 用來代表一個 Unicode code point  ~~盧恩符文是你~~
- byte 代表 **uint8**
- uintptr 拿來存 pointer, 第 13 章的時候會有用來處理 unsafe package 的 example



---

## binary operators

![](https://i.imgur.com/CU5oPOc.jpg)

```go=
//operators.go
package main

import "fmt"


var x uint8 = 1<<1 | 1<<5
var y uint8 = 1<<1 | 1<<2

func main (){
      fmt.Printf("%08b\n", x)    // "00100010", the set {1, 5}
      fmt.Printf("%08b\n", y)    // "00000110", the set {1, 2}
      fmt.Printf("%08b\n", x&y)  // "00000010", the intersection {1}
      fmt.Printf("%08b\n", x|y)  // "00100110", the union {1, 2, 5}
      fmt.Printf("%08b\n", x^y)  // "00100100", the symmetric difference {2, 5}
      fmt.Printf("%08b\n", x&^y) // "00100000", the difference {5}
	  
      for i := uint(0); i < 8; i++ {
          if x&(1<<i) != 0 { // membership test
             fmt.Println(i) // "1", "5"}
	  }
      fmt.Printf("%08b\n", x<<1) // "01000100", the set {2, 6}
      fmt.Printf("%08b\n", x>>1) // "00010001", the set {0, 4}
}
```

---

# 3.2 Floating-Point Numbers
```go=
//float.go
package main

import 	"fmt"


func main() {

	var (
		d string  = "d string."
		e float64 = 1.23
		f uint8   = 234
	)
	
	fmt.Printf("var: d = %s, e = %g, f = %d\n", d, e, f)
}


```

- printf 的 %g 會自動幫你選擇輸出浮點數最佳的精度
- 但如果是資料表，用 %e 或者是 %f 會比較好，不然效率會很差
- [fmt doc](https://golang.org/pkg/fmt/)

# 3.3 Complex Numbers

```go=
//complexnum.go

package main

import (
	"fmt"
)

func main() {
	var x complex128 = complex(1, 2) // 1+2i
	var y complex128 = complex(3, 4) // 3+4i
	fmt.Println(x * y)               // "(-5+10i)"
	fmt.Println(real(x * y))         // "-5"
	fmt.Println(imag(x * y))         // "10"
}
```
- complex 64/128 對應 float 32/64
- imag 指的是虛數 (imaginary part)

---

# 3.4 Booleans

```go=
package main

import "fmt"

func main() {

	var b bool

	b = 1       // error
	b = bool(1) // error

	b = (1 != 0) // correct
	fmt.Println(b)
}
```

- bool 不支持其他 type 的 assignment 或者是轉換

---

# 3.5 Strings

```go=
package main

import "fmt"

func main() {
		
	//ASCII 有 32 個控制字元，同字母大小寫因此差 32
	s := "CCNSccns"
	fmt.Println(len(s))     // "8"
	fmt.Println(s[0], s[4]) // "67 99" ('C' and 'c')

	
	c := s[len(s)]          // panic: index out of range
	fmt.Println(c)
	
	fmt.Println("Hello, " + s[:4]) // "Hello, CCNS"
    
	s[0] := 'A' // compile error: cannot assign to s[0]
}
```

- 好用的 len() 同 Py
- String 跟 Py 的 Tuple 一樣 immutable

---

## 3.5.1 String literals

![](https://i.imgur.com/ImEiP6P.jpg)

- regular expression

## 3.5.2 Unicode

- ASCII code to Unicode
- rune = Unicode code point

## 3.5.3 UTF-8

- 我開始覺得他只是想提到 Rob Pike & Ken Thompson (UTF-8 創始者等於 Go 語言創始者)
- Go source files are always encoded in UTF-8
- [unicode/utf-8 package](https://golang.org/pkg/unicode/utf8/) 

```go=
//rune conversion
package main

import "fmt"

func main() {
	// "program" in Japanese katakana
	s := "プログラム"
	fmt.Printf("% x\n", s) // "e3 83 97 e3 83 ad e3 82 b0 e3 83 a9 e3 83 a0"
	r := []rune(s)
	fmt.Printf("%x\n", r) // "[30d7 30ed 30b0 30e9 30e0]"

	fmt.Println(string(r))
	fmt.Println(string(65))      // "A", not "65"
	fmt.Println(string(0x4eac))  // "京"
	fmt.Println(string(1234567)) // invalid rune "�"
}

```

---

## 3.5.4 Strings and Byte Slices

- 操作字串的四個重要 packages: **bytes, strings, strconv, and unicode**

  - **strings**: 提供搜尋、取代、修剪 (trimming) 字串等等功能
    - [strings package](https://golang.org/pkg/strings/) 
  - **bytes**: 跟 string 差不多，但在要持續增加 string 長度的情況下更有效率
  - **strconv**: 將 int, float, bool 轉成 string
  - **unicode**: 用來判斷 rune 是不是 Lower, Upper 等等，也可以用來轉換

```go=
// contain.go
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Contains("seafood", "foo"))//true
	fmt.Println(strings.Contains("seafood", "bar"))//false
	fmt.Println(strings.Contains("seafood", ""))//true
	fmt.Println(strings.Contains("", ""))//true
}
```

```go=
//compare.go
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Compare("a", "b"))
	fmt.Println(strings.Compare("a", "a"))
	fmt.Println(strings.Compare("b", "a"))
}
```
## 3.5.5 Conversions between Strings and Numbers

- strconv

```go=
x := 123
y := fmt.Sprintf("%d", x)
fmt.Println(y, strconv.Itoa(x)) // "123 123"
```

---

# 3.6 Constants

- 單一宣告 vs 多重宣告
- run at compile time, not at run time

```go=
const pi = 3.141592
const e = 2.718281

const (
	e  = 2.71828
	pi = 3.14159
)
```

- constant group

```go=
const (
	a=1
	b
	c=2
	d
)
fmt.Println(a, b, c, d) // "1 1 2 2"
```

- 好潮
- 不夠方便所以有了 iota

## 3.6.1 The Constant Generator **iota**

- iota
  - 不需要把相關的每個常量一個個宣告出來就能創造出來
  - Go 實作 enum 的方法

```go=
//enum
package main

func main() {

	type Weekday int

	const (
		Sunday Weekday = iota //自動初始化為 0
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)
	println(Sunday)    //0
	println(Monday)    //1
	println(Tuesday)   //2
	println(Wednesday) //3
	println(Thursday)  //4
	println(Friday)    //5
	println(Saturday)  //6
}
```

```go=
//enum2.go
package main

import "fmt"

//Flags type
type Flags uint

//FlagUp=1, FlagBroadcast=2, FlagLoopback=4 ......
const (
	FlagUp           Flags = 1 << iota // is up
	FlagBroadcast                      // supports broadcast access capability
	FlagLoopback                       // is a loopback interface
	FlagPointToPoint                   // belongs to a point-to-point link
	FlagMulticast                      // supports multicast access capability
)


func IsUp(v Flags) bool     { return v&FlagUp == FlagUp }
func TurnDown(v *Flags)     { *v &^= FlagUp }
func SetBroadcast(v *Flags) { *v |= FlagBroadcast }
func IsCast(v Flags) bool   { return v&(FlagBroadcast|FlagMulticast) != 0 }

func main() {
	var v Flags = FlagMulticast | FlagUp
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10001 true"
	TurnDown(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10001 false"10000 false"
	SetBroadcast(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10010 false"
	fmt.Printf("%b %t\n", v, IsCast(v)) // "10010 true"
}
```


## 3.6.2 Untyped Constants

- 無法被任何型態的變數所存的常量(數字太大)

```go=
//untyped constant.go
package main

import (
	"fmt"
	"os"
)

func main() {
	const (
		_ = 1 << (10 * iota)
		KiB
		MiB
		GiB
		TiB
		PiB
		EiB
		ZiB
		YiB
	)
	fmt.Println(KiB) //1024
	fmt.Println(MiB) //1048576
	fmt.Println(GiB) //1073741824
	fmt.Println(TiB) //1099511627776
	fmt.Println(PiB) //1125899906842624
	fmt.Println(EiB) //1152921504606846976
	// fmt.Println(ZiB)   //exceeds 1 << 64 from here
	// fmt.Println(YiB)

	fmt.Println(YiB / ZiB) //1024
	os.Exit(0)
}
```

- 只有 constant 可以存成 untyped，若是變數就會自動轉換
- 為了避免麻煩還是自己寫好就好

```go=
var f float64 = 3 + 0i // untyped complex -> float64
f=2                    // untyped integer -> float64
f=1e123                // untyped floating-point -> float64
f = 'a'                // untyped rune -> float64
```

---

# CH4 Composite Types

- Ch3 is the atoms of our universe
- 喔好喔
- array, slice, maps and structs

# 4.1 Arrays

- 因為是**固定長度**的，所以在 Go 很少用
- 用 slices 更加有彈性

```go
//arrays.go
package main

import "fmt"

func main() {

	var a [5]int
	fmt.Println("Initial a:", a)
    
	a[4] = 100
	fmt.Println("Revised a:", a)
	fmt.Println("Get:", a[4])
	fmt.Println("len:", len(a))

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl b:", b)

	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
// Initial a: [0 0 0 0 0]
// Revised a: [0 0 0 0 100]
// Get: 100
// len: 5
// dcl b: [1 2 3 4 5]
// 2d:  [[0 1 2] [1 2 3]]
```

# 4.2 slices

- 長度可變
- composed by a pointer, a length, a capacity
  - pointer 指向第一個元素
  - length 代表 slice 中元素的數量, **len()**
  - capacity 是 slice 中最多能放入的元素量, **cap()**
- Slicing beyond cap(s) causes a panic, but slicing beyond len(s) extends the slice

![](https://i.imgur.com/vPFGhQF.jpg)


```go=
package main

import "fmt"

func printSlice(x []int){
   fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
}

func main() {
   var numbers = make([]int,3,5)

   printSlice(numbers)
}
//len=3 cap=5 slice=[0 0 0]
```

> append 沒有 side effect，所以要記得 assign 給自己喔 
> [name=go愛用者]
- Go 也有 new，兩者差別在 new 會回傳指標而 make 不會

```go=
//slices.go
package main

import "fmt"

func main() {
    
	//initialization
	s := make([]string, 3)
	fmt.Println("initial:", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])
	fmt.Println("len:", len(s))

	//append
	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("appd:", s)

	//copy
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("copy:", c)

	//slice 
	l := s[2:5]
	fmt.Println("slc:", l)
	
	// declare
	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)
	
	// 2-D
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}
// initial: [  ]
// set: [a b c]
// get: c
// len: 3
// appd: [a b c d e f]
// copy: [a b c d e f]
// slc: [c d e]
// dcl: [g h i]
// 2d:  [[0] [1 2] [2 3 4]]
```

## 4.2.1 In-Place Slicing Techniques

- append, pop, remove

```go=
stack = append(stack, v) // push v
top := stack[len(stack)-1] // top of stack
stack = stack[:len(stack)-1] // pop
```

```go=
//remove middle element
package main

import "fmt"

func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
func main() {
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2)) // "[5 6 8 9]"
}
```
> slice 的操作感覺和 C 的 array 很像，很底層(?)
> [name=HexRabbit]

# 4.3 Maps

- In Go, a map is a reference to a hash table, and a map type is written map[K]V, where K and V are the types of its keys and values.

```go=
ages := make(map[string]int)

ages["alice"] = 31
ages["charlie"] = 34
fmt.Println(ages["alice"])    // 3

delete(ages, "alice")         // remove element ages["alice"]

ages["bob"] = ages["bob"] + 1 // int的預設初始化是0, 所以ages["bob"]==0
fmt.Println(ages["bob"])      // 1

_ = &ages["bob"]              // compile error: cannot take address of map element
```

```go=
// enumerate all key/values pair
for name, age := range ages {
    fmt.Printf("%s\t%d\n", name, age)
}
```

- 這樣枚舉的話會 randomly 印出 map
- 若要有順序的印出需要先 sort

```go=
import "sort"

var names []string

for name := range ages {
    names = append(names, name)
}

sort.Strings(names)

for _, name := range names {

    fmt.Printf("%s\t%d\n", name, ages[name])
}
```
```go=
package main

import "fmt"

func main() {

	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13

	fmt.Println("map:", m)

	v1 := m["k1"]
	fmt.Println("v1: ", v1)

	fmt.Println("len:", len(m))

	delete(m, "k2")
	fmt.Println("map:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
// map: map[k1:7 k2:13]
// v1:  7
// len: 2
// map: map[k1:7]
// prs: false
// map: map[foo:1 bar:2]
```

# 4.4 Structs

- Each value is called a field

```go=
import "time"

type Employee struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerID int
}

var dilbert Employee

dilbert.Salary -= 5000 // demoted, for writing too few lines of code

SA := Employee{001, "SA", "ccns", 2020 ,"here", 666666, 777}

//pointer to struct
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
```
```go=
//treeSort.go
package main

import "fmt"

type node struct {
	value int
	left  *node
	right *node
}

type bst struct {
	root *node
}

func (t *bst) insert(v int) {
	if t.root == nil {
		t.root = &node{v, nil, nil}
		return
	}
	current := t.root
	for {
		if v < current.value {
			if current.left == nil {
				current.left = &node{v, nil, nil}
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &node{v, nil, nil}
				return
			}
			current = current.right
		}
	}
}

func (t *bst) inorder(visit func(int)) {
	var traverse func(*node)
	traverse = func(current *node) {
		if current == nil {
			return
		}
		traverse(current.left)
		visit(current.value)
		traverse(current.right)
	}
	traverse(t.root)
}

func (t *bst) slice() []int {
	sliced := []int{}
	t.inorder(func(v int) {
		sliced = append(sliced, v)
	})
	return sliced
}

func treesort(values []int) []int {
	tree := bst{}
	for _, v := range values {
		tree.insert(v)
	}
	return tree.slice()
}

func main() {
	fmt.Println(treesort([]int{2, 4, 3, 1, 9, 7, 8}))
}
```

## 4.4.1 Comparing Structs

```go=
type Point struct{ X, Y int }
p := Point{1, 2}
q := Point{2, 1}
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q) // "false"
```

## 4.4.2 Struct Embedding and Anonymous Fields

- Struct 裡面包另一個 struct
- 為了精簡 fields 的呼叫
  - x.d.e.f → x.f 

```go=
package main

type Point struct {
	X, Y int
}
type Circle struct {
	Center Point
	Radius int
}
type Wheel struct {
	Circle Circle
	Spokes int
}

func main() {
	var w Wheel
	w.Circle.Center.X = 8
	w.Circle.Center.Y = 8
	w.Circle.Radius = 5
	w.Spokes = 20
}

```
- 可以改寫成

```go=
package main

type Circle struct {
	Point
	Radius int
}
type Wheel struct {
	Circle
	Spokes int
}

func main() {
	var w Wheel
	w.X = 8      // equivalent to w.Circle.Point.X = 8
	w.Y = 8      // equivalent to w.Circle.Point.Y = 8
	w.Radius = 5 // equivalent to w.Circle.Radius = 5
	w.Spokes = 20
}
```

- **Circle** 和 **Point** 變成了 anonymous fields, 造成宣告上 error

```go=
//延續
w = Wheel{8, 8, 5, 20}                       // compile error: unknown fields
w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields


//--------------------要改成這樣-----------------------------
w = Wheel{
	Circle: Circle{
	Point:  Point{X: 8, Y: 8},
	Radius: 5,
    },
Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
}
```
# 4.5 JSON

- 介紹了一下 JSON
  - mapping from string to value 


- 在做 marshalling 的時候是使用 struct field name 作為 JSON 的 field name
  - 把資料轉成 JSON format 
  - json.Marshal()

- 輸出的改變是因為 field tags

- Go 的 field tags 可以是任何形式的 string，但傳統上是 key:"value"
  - 宣告之後會成為那個 field name 的 attribute 

```go=
//marshal.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`  // " " 中的是 string literal (tag) 
	Color  bool `json:"color,omitempty"` //用 Color.Tag 呼叫
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	//JSON 格式化輸出，若單純用 json.Marshal() 會全部擠在一起很噁心
	data, err := json.MarshalIndent(movies, "", " ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
}
```

- 也可以用 Unmarshal 只取自己要的資料

```go=
var titles []struct{ Title string }

if err := json.Unmarshal(data, &titles); err != nil {

        log.Fatalf("JSON unmarshalling failed: %s", err)
}

fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"
```
- 書裡後來還有提供一個 github issue tracker 的做法，挺有趣的，不過 code 太長就不放了


# 4.6 Text and HTML Templates

- 為什麼需要 text/template?
  - 因為有時候 printf 無法滿足我們的需求 (formatting)

- 用來自動產生 text, html/templates 是基於 text/templates 之上

- A template is a string or file containing one or more portions enclosed in double braces, {{...}}, called actions.

- dot(.) is a notion of the current value
  - 類似 C++ this 或是 Python self 

- 流程: Template 會先 parse 丟入的 text file 再 execute，最後回傳 response
```go=
func main() {
 //Parse
 tmpl, err := template.ParseFiles("index.html")
 //Execute (Printout)
 err := tmpl.Execute(os.Stdout, nil)
 if err != nil {
  panic(err)
 }
}
```

```go=
type Anime struct {
    Name string
    Year  int
}

func main(){
a := Anime{"ID:INVADED",2020}
    tmpl, _ := template.New("test").Parse("Name: {{.Name}}, Year: {{.Year}}")
    _ = tmpl.Execute(os.Stdout, a)
}
```

- 當然也可有迴圈、if/else.....等操作來產生
- template.Must 用來檢查模板有沒有錯誤


```go=
import "html/template"
var issueList = template.Must(template.New("issuelist").Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
	<th>#</th>
	<th>State</th>
	<th>User</th>
	<th>Title</th>
</tr>
{{range .Items}}
<tr>
	<td> <a href='{{.HTMLURL}}'>{{.Number}}</td>
	<td> {{.State}}</td>
	<td> <a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
	<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`))
```


![](https://i.imgur.com/iSfJc2r.jpg)



# Leetcode 習題

- [Sudoku solver](https://leetcode.com/problems/sudoku-solver/#) 

```go=
//sudoku solver.go
func solveSudoku(board [][]byte) {
	solve(board, 0)
}

// k 是 index
func solve(board [][]byte, k int) bool {
	if k == 81 {
		return true
	}
    //9*9
	r, c := k/9, k%9
	if board[r][c] != '.' {
		return solve(board, k+1)
	}

	// 左上角的 index
	bi, bj := r/3*3, c/3*3

	// 檢查 b 能不能放進去 board
	isValid := func(b byte) bool {
		for n := 0; n < 9; n++ {
			if board[r][n] == b ||
				board[n][c] == b ||
				board[bi+n/3][bj+n%3] == b {
				return false
			}
		}
		return true
	}
    
	for b := byte('1'); b <= '9'; b++ {
		if isValid(b) {
			board[r][c] = b
			if solve(board, k+1) {
				return true
			}
		}
	}

	board[r][c] = '.'

	return false
}

```

- [generate parentheses](https://leetcode.com/problems/generate-parentheses/)

  - 左括號用完就不能再加
  - 左右括號一樣多就不能再加右括號
  - 左右都用完就結束


```go=
//generate parentheses.go
func generateParenthesis(n int) []string {
    
	res := make([]string, 0, n*n)
	bytes := make([]byte, n*2)
	dfs(n, n, 0, bytes, &res)
	return res
}

func dfs(left, right, idx int, bytes []byte, res *[]string) {
	// 如果都加完了
	if left == 0 && right == 0 {
		*res = append(*res, string(bytes))
		return
	}

	// "(" 不用擔心配對問題
	// 只要 left > 0 就可以直接加進去
	if left > 0 {
		bytes[idx] = '('
		dfs(left-1, right, idx+1, bytes, res)
	}

    // 要加")"的時候
	// 若 left < right，
	// 而 bytes[:idx] 至少有一個 "(" 可以與這個 ")" 配對
	if right > 0 && left < right {
		bytes[idx] = ')'
		dfs(left, right-1, idx+1, bytes, res)
	}
}
```

- [rotate array](https://leetcode.com/problems/rotate-array/)
   - 沒有reverse()不會自己寫?
```go=
//rotate array.go
func rotate(nums []int, k int) {
	// if k >= 0

	n := len(nums)

	if k > n {
		k %= n
	}
	if k == 0 || k == n {
		return
	}

	reverse(nums, 0, n-1)
	reverse(nums, 0, k-1)
	reverse(nums, k, n-1)
}

func reverse(nums []int, i, j int) {
	for i < j {
		nums[i], nums[j] = nums[j], nums[i]
		i++
		j--
	}
}
```
