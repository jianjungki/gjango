package main

import (
	"fmt"

	"tiktok_tools"
)

func main() {
	tiktok_tools.New().
		WithRoutes(&MyServices{}).
		Run()
}

// MyServices implements tiktok_tools/route.ServicesI
type MyServices struct{}

// SetupRoutes is our implementation of custom routes
func (s *MyServices) SetupRoutes() {
	fmt.Println("set up our custom routes!")
}
