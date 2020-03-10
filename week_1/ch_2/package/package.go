// skip may cause :
// expected 'package', found 'func'
package main

// duplicate may cause :
// syntax error: non-declaration statement outside function body
// package sub

// skip may cause :
// function main is undeclared in the main package
func main() { /* Do something*/ }
