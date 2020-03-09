package main

// valid
// var _foo int
// var _123 int
// var fooBar int

// not recommand
// var foo_bar int

// syntax error
// var 123 int

// may cause type ambiguous, int is not a type
var int int
