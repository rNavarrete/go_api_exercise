package main

func main() {
	a := App{}
	// you need to set your Username and password here
	a.Initialize("root", "godfather", "rest_api_example")
	a.Run(":8080")
}
