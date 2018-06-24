package main

func main() {
	a := App{}
	// you need to set your Username and password here
	a.Initialize("rolandonavarrete", "godfather", "rest_api_example")
	a.Run(":8080")
}
