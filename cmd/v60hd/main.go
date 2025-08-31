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

	// Create client with specific device type (V-160HD in this example)
	c, err := roland.NewClientWithDeviceType(ipAddress, port, roland.DeviceV160HD)
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer c.Close()

	fmt.Printf("Connected to %s device\n", c.GetDeviceType())

	// Example 1: Basic video switching
	fmt.Println("\n=== Video Switching Examples ===")

	// Set program to HDMI1
	if err := c.SetProgram("HDMI1"); err != nil {
		fmt.Println("SetProgram error:", err)
	} else {
		fmt.Println("✓ Set program to HDMI1")
	}

	// Set preset to SDI1
	if err := c.SetPreset("SDI1"); err != nil {
		fmt.Println("SetPreset error:", err)
	} else {
		fmt.Println("✓ Set preset to SDI1")
	}

	// Auto transition to HDMI2 with 2.0 second duration
	if err := c.AutoTo("HDMI2", 20); err != nil {
		fmt.Println("AutoTo error:", err)
	} else {
		fmt.Println("✓ Auto transition to HDMI2 (2.0s)")
	}

	// Example 2: PinP (Picture in Picture) control
	fmt.Println("\n=== PinP Control Examples ===")

	// Set PinP1 source to HDMI3
	if err := c.SetPinPSource("PinP1", "HDMI3"); err != nil {
		fmt.Println("SetPinPSource error:", err)
	} else {
		fmt.Println("✓ Set PinP1 source to HDMI3")
	}

	// Enable PinP1 on program
	if err := c.SetPinPPGM("PinP1", roland.StateON); err != nil {
		fmt.Println("SetPinPPGM error:", err)
	} else {
		fmt.Println("✓ Enabled PinP1 on program")
	}

	// Set PinP1 position (center)
	if err := c.SetPinPPosition("PinP1", 0, 0); err != nil {
		fmt.Println("SetPinPPosition error:", err)
	} else {
		fmt.Println("✓ Set PinP1 position to center")
	}

	// Example 3: DSK (Downstream Keyer) control
	fmt.Println("\n=== DSK Control Examples ===")

	// Set DSK1 fill source to STILL1
	if err := c.SetDSKFillSource("DSK1", "STILL1"); err != nil {
		fmt.Println("SetDSKFillSource error:", err)
	} else {
		fmt.Println("✓ Set DSK1 fill source to STILL1")
	}

	// Set DSK1 level to 128 (50%)
	if err := c.SetDSKLevel("DSK1", 128); err != nil {
		fmt.Println("SetDSKLevel error:", err)
	} else {
		fmt.Println("✓ Set DSK1 level to 128")
	}

	// Enable DSK1
	if err := c.SetDSK("DSK1", roland.StateON); err != nil {
		fmt.Println("SetDSK error:", err)
	} else {
		fmt.Println("✓ Enabled DSK1")
	}

	// Example 4: Audio control
	fmt.Println("\n=== Audio Control Examples ===")

	// Set XLR1 input level to -60dB
	if err := c.SetInputLevel("XLR1", -60); err != nil {
		fmt.Println("SetInputLevel error:", err)
	} else {
		fmt.Println("✓ Set XLR1 input level to -60dB")
	}

	// Enable high pass filter on XLR1
	if err := c.SetHPF("XLR1", roland.StateON); err != nil {
		fmt.Println("SetHPF error:", err)
	} else {
		fmt.Println("✓ Enabled HPF on XLR1")
	}

	// Set main output level to -20dB
	if err := c.SetOutputLevel("MAIN", -20); err != nil {
		fmt.Println("SetOutputLevel error:", err)
	} else {
		fmt.Println("✓ Set main output level to -20dB")
	}

	// Example 5: Camera control (PTZ)
	fmt.Println("\n=== Camera Control Examples ===")

	// Pan camera left
	if err := c.CameraPanTilt("CAMERA1", roland.CameraLEFT, roland.CameraSTOP); err != nil {
		fmt.Println("CameraPanTilt error:", err)
	} else {
		fmt.Println("✓ Panning camera left")
	}

	// Set camera zoom to wide
	if err := c.CameraZoom("CAMERA1", roland.ZoomWIDE_SLOW); err != nil {
		fmt.Println("CameraZoom error:", err)
	} else {
		fmt.Println("✓ Zooming camera wide")
	}

	// Recall camera preset 1
	if err := c.CameraPreset("CAMERA1", "PRESET1"); err != nil {
		fmt.Println("CameraPreset error:", err)
	} else {
		fmt.Println("✓ Recalled camera preset 1")
	}

	// Example 6: Stream control
	fmt.Println("\n=== Stream Control Examples ===")

	// Set stream type to RTMP
	if err := c.SetStreamType("RTMP"); err != nil {
		fmt.Println("SetStreamType error:", err)
	} else {
		fmt.Println("✓ Set stream type to RTMP")
	}

	// Start streaming
	if err := c.StreamStart(); err != nil {
		fmt.Println("StreamStart error:", err)
	} else {
		fmt.Println("✓ Started streaming")
	}

	// Example 7: Memory and scene control
	fmt.Println("\n=== Memory Control Examples ===")

	// Recall memory 1
	if err := c.RecallMemory("MEMORY1"); err != nil {
		fmt.Println("RecallMemory error:", err)
	} else {
		fmt.Println("✓ Recalled memory 1")
	}

	// Example 8: System queries
	fmt.Println("\n=== System Query Examples ===")

	// Query device version
	if err := c.QueryVersion(); err != nil {
		fmt.Println("QueryVersion error:", err)
	} else {
		fmt.Println("✓ Queried device version")
	}

	// Query current program
	if err := c.QueryProgram(); err != nil {
		fmt.Println("QueryProgram error:", err)
	} else {
		fmt.Println("✓ Queried current program")
	}

	// Query tally state
	if err := c.QueryTally(); err != nil {
		fmt.Println("QueryTally error:", err)
	} else {
		fmt.Println("✓ Queried tally state")
	}

	fmt.Println("\n=== All examples completed ===")
	fmt.Println("Note: Some commands may not work depending on device configuration and current state.")
}
