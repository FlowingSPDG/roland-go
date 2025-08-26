package main

import (
	"fmt"
	"os"

	"github.com/FlowingSPDG/roland-go/roland"
)

func main() {
	// Device IP address and port (default TCP API port is 8023)
	ipAddress := "192.168.2.254"
	port := "8023"

	// Example: program selection
	channel := 1

	c, err := roland.NewClient(ipAddress, port)
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer c.Close()

	// Legacy helper using numeric channel
	if err := c.PGM(channel); err != nil {
		fmt.Println("PGM command error:", err)
		os.Exit(1)
	}

	// Modern helper using label (uncomment for usage)
	// if err := c.SetProgram("INPUT1"); err != nil {
	// 	fmt.Println("SetProgram error:", err)
	// 	os.Exit(1)
	// }

	fmt.Println("Done")
}
