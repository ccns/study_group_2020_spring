# 讀書會 Week 3

## Chapter.5 Functions

### 5.1 Function Declarations
#### Declarations
```go
func name(parameter-list) (result-list) {
    body
}
```

```go
* Equivalent

func f(i, j, k int, s, t string) { body }
func f(i int, j int, k int, s string, t string) { body }

```

* Arguments are **passed by value**
    * If the args contain some kind of reference, like
        * pointer
        * slice
        * map
        * func
        * channel
    * Caller may be affected by any modifications the function makes to variable

* If a function is declared without a body
    * It indicates that the function is implemented in lang other than Go
    ```
    package math
    func Sin(x float64) float64
    ```
    Reference: [go/sin_s390x.s](https://github.com/golang/go/blob/master/src/math/sin_s390x.s)
    
### 5.2 Recursion
:::info
`go get golang.org/x/net/html`
:::
* Reference: [gopl.io/ch5/findlinks1](https://github.com/adonovan/gopl.io/tree/master/ch5/findlinks1)
* Reference: [gopl.io/ch1/fetch](https://github.com/adonovan/gopl.io/tree/master/ch1/fetch)

* Typical Go implementations use variable-size stacks, which start small and grow up as needed up to a gigabyte

### 5.3 Multiple Return Values
* Reference: [gopl.io/ch5/findlinks2](https://github.com/adonovan/gopl.io/tree/master/ch5/findlinks2)
```go
Equivalent, but the first one is rarely used in production code

log.Println(findLinks(url))

links, err := findLinks(url)
log.Println(links, err)
```

* Some good ways to define your function
```go
func Size(rect image.Rectangle) (width, height int)
func Split(path string) (dir, file string)
func HourMinSec(t time.Time) (hour, minute, second int)
```

* Bare return
```go
func CountWordsandImages(url string) (words, images int, err error) {
    resp, error := http.Get(url)
    if err != nil {
        return
    }
    doc, err := html.Parse(resp.Body)
    resp.Body.Close()
    if err != nil {
        err = fmt.Errorf("parsing HTML: %s", err)
        return
    }
    words, images = countWordsAndImages(doc)
    return
}
func countWordsAndImages (n *html.Node) (words, images int) { /* .... */ }
```
Every return above is equivalent to `return words, images, err`

### 5.4 Errors
* Builtin-type **error** is an interface type
    * May be *nil* or *non-nil*
    * Can be obtained by calling **Error** method or
      `fmt.Println(err)` / `fmt.Printf("%v", err)`
* Golang do have *exceptions*
    * But it's only for **unexpected errors**
    * Error is for routine errors, in order to maintain a robust program
* `fmt.Errorf`
    * Uses `fmt.Sprintf` and returns a new error value
    * Should provide a clear causal chain from the root problem to overall failure
    * Message strings should not be capitalized and newlines should be avoided
        * They will be self-contained when found by tools like grep
* Error handling strategy
    * Propagate the error
    * Retry the failed operation
    * Print the error and stop the program gracefully
    * Log the error and continue
    * Ignore the error entirely
* Get into the habit of considering errors after every function call
  When you ignore one, document your intention clearly
* EOF
    * `errors.New("EOF")`

### 5.5 Function Values
```go
func square(n int) int {return n*n}
func negative(n int) int {return -n}
func product(m, n int) int {return m*n}

f := square
f(3) // "9"

f = negative
f(3) // "-3"

g := product // will work
f = product // compile error: can't assign f(int, int) int to f(int) int
```

```go
func add1(r rune) rune {return r + 1}

fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
```
Reference: [strings - #Map](https://golang.org/pkg/strings/#Map)

### 5.6 Anonymous Functions
```go
func squares() func() int {
    var x int
    return func() int {
        x++
        return x * x
    }
}
func main(){
    f := squares()
    fmt.Println(f()) // '1'
    fmt.Println(f()) // '4'
    fmt.Println(f()) // '9'
}
```

```go
var rmdirs []func()
for _, d := range tempDirs() {
    dir := d
    os.MkdirAll(dir, 0755)
    rmdirs = append(rmdirs, func(){
        os.RemoveAll(dir)
    })
}

// ....

for _, rmdir := range rmdirs {
    rmdir()
}
```

```go
var rmdirs []func()
for _, dir := range tempDirs() {
    os.MkdirAll(dir, 0755)
    rmdirs = append(rmdirs, func(){
        os.RemoveAll(dir)
    })
}
```
* A loop "captures" and shares the same variable — an addressable storage location, not its value at the particular moment

### 5.7 Variadic Functions
* **...** symbol
```go
func sum(vals ...int) int {
    total := 0
    for _, val := range vals {
        total += val
    }
    return total
}
fmt.Println(sum()) // '0'
fmt.Println(sum(3)) // '3'
fmt.Println(sum(1,2,3,4)) // '10'

values := []int{1,2,3,4}
fmt.Println(sum(values...)) // '10'
```

### 5.8 Deferred Function Calls
* `defer`
```go
func main(){
    defer fmt.Println("foo")
    fmt.Println("bar")
}
```
```
bar
foo
```
* Useful to "restore" something
```go
func double(x int) (result int){
    defer func() { fmt.Printf("double(%d) = %d\n", x, result)}()
    return x + x
}
_ = double(4) // '8'
```

```go
func triple(x int) (result int){
    defer func() {result += x}()
    return double(x)
}
_ = triple(4) // '12'
```

```go
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // WARNING: may run out of file descriptor
    .....
}
```

```go
for _, filename := range filenames {
    if err := doFile(filename); err != nil {
        return err
    }
}
func doFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    .....
}
```

### 5.9 Panic
* panic
    * Normal execution stops
    * All deferred function calls executed
    * Program crashes with log message
        * panic value
        * stack trace for each goroutine
    * Can be called manually
        * Takes any values as argument
* Deferred functions are run in reverse order, starting with the one of the topmost function on the stack and proceeding up to main

### 5.10 Recover
* recover
    * Ends the current state of panic and returns the panic value

```go
func Parse(input string) (s *Syntax, err error){
    defer func(){
        if p := recover(); p != nil {
            err = fmt.Errorf("internal error: %v", p)
        }
    }()
    .....
}
```

## Exercises
* [Reverse Linked List](https://leetcode.com/problems/reverse-linked-list/)
* [N-th Tribonacci Number](https://leetcode.com/problems/n-th-tribonacci-number/)
* [Decode String](https://leetcode.com/problems/decode-string/)
