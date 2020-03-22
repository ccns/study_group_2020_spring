package main

func main() {

	type Weekday int

	const (
		Sunday Weekday = iota
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
