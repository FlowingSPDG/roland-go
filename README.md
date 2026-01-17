# Roland Go Control Library

A comprehensive Go library for controlling Roland video switchers and AV mixers via Telnet (LAN) interface.

The library communicates with Roland devices using Telnet over TCP, sending commands in ASCII format without the stx (02H) prefix as specified in the Roland control documentation.

## Supported Devices

- **V-160HD** (STREAMING VIDEO SWITCHER) Ver.3.3+
- **V-80HD** (DIRECT STREAMING VIDEO SWITCHER) Ver.1.1+
- **VR-120HD** (DIRECT STREAMING AV MIXER) Ver.2.1+
- **VR-6HD** (DIRECT STREAMING AV MIXER) Ver.2.1+

## Features

- **Device-specific validation**: Automatic parameter validation based on device capabilities
- **Comprehensive command coverage**: All commands from the Roland LAN/RS-232 manual
- **Thread-safe operations**: Concurrent access protection
- **Error handling**: Detailed error messages with validation
- **Easy to use**: High-level API with intuitive method names

## Installation

```bash
go get github.com/FlowingSPDG/roland-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/FlowingSPDG/roland-go/roland"
)

func main() {
    // Connect to V-160HD device
    client, err := roland.NewClientWithDeviceType("192.168.2.254", "8023", roland.DeviceV160HD)
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer client.Close()

    // Set program to HDMI1
    if err := client.SetProgram("HDMI1"); err != nil {
        log.Printf("Failed to set program: %v", err)
    }

    // Auto transition to HDMI2 with 2-second duration
    if err := client.AutoTo("HDMI2", 20); err != nil {
        log.Printf("Failed to auto transition: %v", err)
    }
}
```

## API Overview

### Connection

```go
// Basic connection (defaults to V-160HD)
client, err := roland.NewClient("192.168.2.254", "8023")

// Connection with specific device type
client, err := roland.NewClientWithDeviceType("192.168.2.254", "8023", roland.DeviceV80HD)
```

### Video Commands

```go
// Program and Preset control
client.SetProgram("HDMI1")
client.SetPreset("SDI1")
client.AutoTo("HDMI2", 20)  // 2.0 second transition
client.CutTo("HDMI3")

// PinP (Picture in Picture) control
client.SetPinPSource("PinP1", "HDMI3")
client.SetPinPPGM("PinP1", roland.StateON)
client.SetPinPPosition("PinP1", 0, 0)  // Center position

// DSK (Downstream Keyer) control
client.SetDSKFillSource("DSK1", "STILL1")
client.SetDSKLevel("DSK1", 128)
client.SetDSK("DSK1", roland.StateON)

// Split control
client.SetSplit("SPLIT1", roland.StateON)
client.SetSplitPosition("SPLIT1", 0, 0)
```

### Audio Commands

```go
// Input control
client.SetInputLevel("XLR1", -60)  // -60dB
client.SetInputMute("XLR1", roland.StateOFF)
client.SetHPF("XLR1", roland.StateON)

// Output control
client.SetOutputLevel("MAIN", -20)  // -20dB
client.SetOutputMute("MAIN", roland.StateOFF)

// Effects
client.SetReverb(roland.StateON)
client.SetAudioAutoMixing(roland.StateON)
```

### Camera Control (PTZ)

```go
// Pan/Tilt control
client.CameraPanTilt("CAMERA1", roland.CameraLEFT, roland.CameraSTOP)
client.CameraPanTiltSpeed("CAMERA1", 12)

// Zoom control
client.CameraZoom("CAMERA1", roland.ZoomWIDE_SLOW)
client.CameraZoomReset("CAMERA1")

// Focus control
client.CameraFocus("CAMERA1", roland.FocusNEAR)
client.CameraAutoFocus("CAMERA1", roland.StateON)

// Preset recall
client.CameraPreset("CAMERA1", "PRESET1")
```

### Stream Control

```go
// Stream configuration
client.SetStreamType("RTMP")
client.StreamStart()
client.StreamStop()

// SRT control
client.SetSRTOutStart()
client.SetSRTInStart()
```

### Memory and Scene Control

```go
// Memory recall
client.RecallMemory("MEMORY1")

// Sequencer control
client.SetSequencer(roland.StateON)
client.SequencerNext()
client.SequencerPrevious()
```

### System Commands

```go
// Device information
client.QueryVersion()
client.QueryBusy()

// Test patterns and tones
client.TestPattern("COLORBAR75")
client.TestTone("0dB")
```

## Device-Specific Features

The library automatically validates parameters based on the connected device:

- **V-160HD**: 8 HDMI, 8 SDI, 16 Still, 20 Input, 4 PinP, 4 DSK, 4 Split, 4 Camera
- **V-80HD**: 4 HDMI, 4 SDI, 32 Still, 16 Input, 2 PinP, 2 DSK, 2 Split, 2 Camera
- **VR-120HD**: 6 HDMI, 6 SDI, 16 Still, 8 Input, 2 PinP, 2 DSK, 2 Split, 2 Camera
- **VR-6HD**: 6 HDMI, 0 SDI, 16 Still, 6 Input, 1 PinP, 1 DSK, 1 Split, 1 Camera

## Constants

The library provides constants for common values:

```go
// States
roland.StateOFF
roland.StateON

// Transition types
roland.TransitionMIX
roland.TransitionWIPE

// Camera controls
roland.CameraLEFT
roland.CameraRIGHT
roland.CameraSTOP
roland.CameraUP
roland.CameraDOWN

// Zoom controls
roland.ZoomWIDE_FAST
roland.ZoomWIDE_SLOW
roland.ZoomTELE_SLOW
roland.ZoomTELE_FAST

// Player controls
roland.PlayerSTOP
roland.PlayerPLAY
roland.PlayerPAUSE
roland.PlayerRESUME
```

## Error Handling

All methods return errors with detailed validation messages:

```go
err := client.SetProgram("INVALID_INPUT")
// Error: invalid video input source INVALID_INPUT for device V-160HD

err := client.SetPinP("PinP5", "HDMI1")
// Error: invalid PinP PinP5 for device V-160HD (max: PinP4)
```

## Examples

See the `cmd/v60hd/main.go` file for comprehensive usage examples covering all major features.

## Building

```bash
# Build the example
go build -o v60hd cmd/v60hd/main.go

# Run the example
./v60hd
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## References

- [Roland LAN/RS-232 簡易制御コマンドマニュアル](https://www.roland.com/global/support/by_product/v-160hd/owners_manuals/)
- [Roland V-160HD Product Page](https://www.roland.com/global/products/v-160hd/)
