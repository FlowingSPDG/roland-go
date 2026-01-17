package roland

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"

	"golang.org/x/xerrors"
)

// DeviceType represents the Roland device model
type DeviceType string

const (
	DeviceV160HD  DeviceType = "V-160HD"
	DeviceV80HD   DeviceType = "V-80HD"
	DeviceVR120HD DeviceType = "VR-120HD"
	DeviceVR6HD   DeviceType = "VR-6HD"
)

// DeviceLimits contains parameter limits for each device type
type DeviceLimits struct {
	VideoInputSources []string
	HDMIInputs        []string
	SDIInputs         []string
	StillImages       []string
	InputChannels     []string
	AuxBuses          []string
	PinPCount         int
	DSKCount          int
	SplitCount        int
	MemoryCount       int
	CameraCount       int
	PresetCount       int
}

// Device limits mapping
var deviceLimits = map[DeviceType]DeviceLimits{
	DeviceV160HD: {
		VideoInputSources: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6", "HDMI7", "HDMI8",
			"SDI1", "SDI2", "SDI3", "SDI4", "SDI5", "SDI6", "SDI7", "SDI8",
			"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16",
			"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8",
			"INPUT9", "INPUT10", "INPUT11", "INPUT12", "INPUT13", "INPUT14", "INPUT15", "INPUT16",
			"INPUT17", "INPUT18", "INPUT19", "INPUT20"},
		HDMIInputs: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6", "HDMI7", "HDMI8"},
		SDIInputs:  []string{"SDI1", "SDI2", "SDI3", "SDI4", "SDI5", "SDI6", "SDI7", "SDI8"},
		StillImages: []string{"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16"},
		InputChannels: []string{"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8",
			"INPUT9", "INPUT10", "INPUT11", "INPUT12", "INPUT13", "INPUT14", "INPUT15", "INPUT16",
			"INPUT17", "INPUT18", "INPUT19", "INPUT20"},
		AuxBuses:    []string{"AUX1", "AUX2", "AUX3"},
		PinPCount:   4,
		DSKCount:    4,
		SplitCount:  4,
		MemoryCount: 30,
		CameraCount: 4,
		PresetCount: 16,
	},
	DeviceV80HD: {
		VideoInputSources: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4",
			"SDI1", "SDI2", "SDI3", "SDI4",
			"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16",
			"STILL17", "STILL18", "STILL19", "STILL20", "STILL21", "STILL22", "STILL23", "STILL24",
			"STILL25", "STILL26", "STILL27", "STILL28", "STILL29", "STILL30", "STILL31", "STILL32",
			"V.PLAYER", "SRT",
			"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8",
			"INPUT9", "INPUT10", "INPUT11", "INPUT12", "INPUT13", "INPUT14", "INPUT15", "INPUT16"},
		HDMIInputs: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4"},
		SDIInputs:  []string{"SDI1", "SDI2", "SDI3", "SDI4"},
		StillImages: []string{"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16",
			"STILL17", "STILL18", "STILL19", "STILL20", "STILL21", "STILL22", "STILL23", "STILL24",
			"STILL25", "STILL26", "STILL27", "STILL28", "STILL29", "STILL30", "STILL31", "STILL32"},
		InputChannels: []string{"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8",
			"INPUT9", "INPUT10", "INPUT11", "INPUT12", "INPUT13", "INPUT14", "INPUT15", "INPUT16"},
		AuxBuses:    []string{"AUX1", "AUX2"},
		PinPCount:   2,
		DSKCount:    2,
		SplitCount:  2,
		MemoryCount: 32,
		CameraCount: 2,
		PresetCount: 16,
	},
	DeviceVR120HD: {
		VideoInputSources: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6",
			"SDI1", "SDI2", "SDI3", "SDI4", "SDI5", "SDI6",
			"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16",
			"V.PLAYER", "SRT",
			"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8"},
		HDMIInputs: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6"},
		SDIInputs:  []string{"SDI1", "SDI2", "SDI3", "SDI4", "SDI5", "SDI6"},
		StillImages: []string{"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16"},
		InputChannels: []string{"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6", "INPUT7", "INPUT8"},
		AuxBuses:      []string{"AUX1", "AUX2", "AUX3"},
		PinPCount:     2,
		DSKCount:      2,
		SplitCount:    2,
		MemoryCount:   32,
		CameraCount:   2,
		PresetCount:   16,
	},
	DeviceVR6HD: {
		VideoInputSources: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6",
			"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16",
			"V.PLAYER", "SRT",
			"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6"},
		HDMIInputs: []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "HDMI5", "HDMI6"},
		SDIInputs:  []string{},
		StillImages: []string{"STILL1", "STILL2", "STILL3", "STILL4", "STILL5", "STILL6", "STILL7", "STILL8",
			"STILL9", "STILL10", "STILL11", "STILL12", "STILL13", "STILL14", "STILL15", "STILL16"},
		InputChannels: []string{"INPUT1", "INPUT2", "INPUT3", "INPUT4", "INPUT5", "INPUT6"},
		AuxBuses:      []string{"AUX1"},
		PinPCount:     1,
		DSKCount:      1,
		SplitCount:    1,
		MemoryCount:   32,
		CameraCount:   1,
		PresetCount:   16,
	},
}

// Common constants
const (
	// States
	StateOFF = "OFF"
	StateON  = "ON"

	// Transition types
	TransitionMIX  = "MIX"
	TransitionWIPE = "WIPE"

	// Camera controls
	CameraLEFT  = "LEFT"
	CameraRIGHT = "RIGHT"
	CameraSTOP  = "STOP"
	CameraUP    = "UP"
	CameraDOWN  = "DOWN"

	// Zoom controls
	ZoomWIDE_FAST = "WIDE_FAST"
	ZoomWIDE_SLOW = "WIDE_SLOW"
	ZoomTELE_SLOW = "TELE_SLOW"
	ZoomTELE_FAST = "TELE_FAST"

	// Focus controls
	FocusNEAR = "NEAR"
	FocusFAR  = "FAR"

	// Stream states
	StreamSTOP     = "STOP"
	StreamSTART    = "START"
	StreamSTOPPING = "STOPPING"
	StreamSTARTING = "STARTING"
	StreamONAIR    = "ONAIR"

	// Player controls
	PlayerSTOP   = "STOP"
	PlayerPLAY   = "PLAY"
	PlayerPAUSE  = "PAUSE"
	PlayerRESUME = "RESUME"

	// Scan modes
	ScanNORMAL  = "NORMAL"
	ScanREVERSE = "REVERSE"
	ScanRANDOM  = "RANDOM"
)

type Client struct {
	conn       net.Conn
	mu         sync.Mutex // protects concurrent access to connection
	deviceType DeviceType
	limits     DeviceLimits
}

// NewClient connects to a Roland switcher-compatible device (e.g., V-160HD, V-80HD, VR-120HD, VR-6HD)
// via Telnet over TCP (typically port 8023).
// Commands are sent in ASCII format without stx (02H) prefix as per Telnet specification.
func NewClient(ipAddress string, port string) (*Client, error) {
	return NewClientWithDeviceType(ipAddress, port, DeviceV160HD)
}

// NewClientWithDeviceType connects to a Roland switcher-compatible device with specific device type
// via Telnet over TCP. The connection uses ASCII command format without stx (02H) prefix.
func NewClientWithDeviceType(ipAddress string, port string, deviceType DeviceType) (*Client, error) {
	if ipAddress == "" {
		return nil, errors.New("ipAddress must not be empty")
	}
	if port == "" {
		return nil, errors.New("port must not be empty")
	}

	limits, exists := deviceLimits[deviceType]
	if !exists {
		return nil, fmt.Errorf("unsupported device type: %s", deviceType)
	}

	// Connect via Telnet (TCP connection)
	address := net.JoinHostPort(ipAddress, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, xerrors.Errorf("dial telnet tcp %s: %w", address, err)
	}

	return &Client{
		conn:       conn,
		deviceType: deviceType,
		limits:     limits,
	}, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return xerrors.Errorf("close tcp connection: %w", err)
	}
	return nil
}

// send writes a Telnet command (ASCII format without stx prefix) and waits for response, ensuring thread safety.
// Command format: "COMMAND:param1,param2;" (stx 02H is omitted for Telnet as per specification)
func (c *Client) send(command string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(command) == 0 || command[len(command)-1] != ';' {
		command = command + ";"
	}
	// Send command via Telnet (ASCII format, stx 02H is omitted)
	if _, err := c.conn.Write([]byte(command)); err != nil {
		return xerrors.Errorf("write telnet command %q: %w", command, err)
	}

	// Wait for response after sending command (read until LF)
	reader := bufio.NewReader(c.conn)
	_, err := reader.ReadString('\n')
	if err != nil {
		return xerrors.Errorf("read telnet response for command %q: %w", command, err)
	}
	return nil
}

// SendRaw sends a prepared Telnet command body (e.g., "PGM:INPUT1;") and waits for response.
// The command should be in ASCII format without stx (02H) prefix.
func (c *Client) SendRaw(raw string) error { return c.send(raw) }

// ReadResponse reads Telnet response until LF following the final ';' (for manual response reading)
func (c *Client) ReadResponse() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	reader := bufio.NewReader(c.conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", xerrors.Errorf("read telnet response: %w", err)
	}
	return line, nil
}

// GetDeviceType returns the device type of the client
func (c *Client) GetDeviceType() DeviceType {
	return c.deviceType
}

// GetDeviceLimits returns the device limits of the client
func (c *Client) GetDeviceLimits() DeviceLimits {
	return c.limits
}

// validateVideoInputSource validates if the input source is valid for the device
func (c *Client) validateVideoInputSource(input string) error {
	for _, validInput := range c.limits.VideoInputSources {
		if validInput == input {
			return nil
		}
	}
	return fmt.Errorf("invalid video input source %s for device %s", input, c.deviceType)
}

// validateAuxBus validates if the aux bus is valid for the device
func (c *Client) validateAuxBus(aux string) error {
	for _, validAux := range c.limits.AuxBuses {
		if validAux == aux {
			return nil
		}
	}
	return fmt.Errorf("invalid aux bus %s for device %s", aux, c.deviceType)
}

// validatePinP validates if the PinP number is valid for the device
func (c *Client) validatePinP(pinp string) error {
	if c.limits.PinPCount == 0 {
		return fmt.Errorf("PinP not supported on device %s", c.deviceType)
	}
	for i := 1; i <= c.limits.PinPCount; i++ {
		if fmt.Sprintf("PinP%d", i) == pinp {
			return nil
		}
	}
	return fmt.Errorf("invalid PinP %s for device %s (max: PinP%d)", pinp, c.deviceType, c.limits.PinPCount)
}

// validateDSK validates if the DSK number is valid for the device
func (c *Client) validateDSK(dsk string) error {
	if c.limits.DSKCount == 0 {
		return fmt.Errorf("DSK not supported on device %s", c.deviceType)
	}
	for i := 1; i <= c.limits.DSKCount; i++ {
		if fmt.Sprintf("DSK%d", i) == dsk {
			return nil
		}
	}
	return fmt.Errorf("invalid DSK %s for device %s (max: DSK%d)", dsk, c.deviceType, c.limits.DSKCount)
}

// validateSplit validates if the Split number is valid for the device
func (c *Client) validateSplit(split string) error {
	if c.limits.SplitCount == 0 {
		return fmt.Errorf("Split not supported on device %s", c.deviceType)
	}
	for i := 1; i <= c.limits.SplitCount; i++ {
		if fmt.Sprintf("SPLIT%d", i) == split {
			return nil
		}
	}
	return fmt.Errorf("invalid Split %s for device %s (max: SPLIT%d)", split, c.deviceType, c.limits.SplitCount)
}

// validateMemory validates if the memory number is valid for the device
func (c *Client) validateMemory(memory string) error {
	for i := 1; i <= c.limits.MemoryCount; i++ {
		if fmt.Sprintf("MEMORY%d", i) == memory {
			return nil
		}
	}
	return fmt.Errorf("invalid memory %s for device %s (max: MEMORY%d)", memory, c.deviceType, c.limits.MemoryCount)
}

// validateCamera validates if the camera number is valid for the device
func (c *Client) validateCamera(camera string) error {
	if c.limits.CameraCount == 0 {
		return fmt.Errorf("Camera control not supported on device %s", c.deviceType)
	}
	for i := 1; i <= c.limits.CameraCount; i++ {
		if fmt.Sprintf("CAMERA%d", i) == camera {
			return nil
		}
	}
	return fmt.Errorf("invalid camera %s for device %s (max: CAMERA%d)", camera, c.deviceType, c.limits.CameraCount)
}

// validatePreset validates if the preset number is valid for the device
func (c *Client) validatePreset(preset string) error {
	for i := 1; i <= c.limits.PresetCount; i++ {
		if fmt.Sprintf("PRESET%d", i) == preset {
			return nil
		}
	}
	return fmt.Errorf("invalid preset %s for device %s (max: PRESET%d)", preset, c.deviceType, c.limits.PresetCount)
}

// VIDEO Commands

// PGM selects program channel by numeric index (legacy helper). Prefer SetProgram("INPUT1").
func (c *Client) PGM(channel int) error { return c.send(fmt.Sprintf("PGM:%d;", channel)) }

// SetProgram sets the program channel with validation
func (c *Client) SetProgram(input string) error {
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("PGM:%s;", input))
}

// SetPreset sets the preset channel with validation
func (c *Client) SetPreset(input string) error {
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("PST:%s;", input))
}

// QueryVideoInputSource queries video input source
func (c *Client) QueryVideoInputSource(index int) error {
	return c.send(fmt.Sprintf("QVISRC:%d;", index))
}

// QueryVideoInputStatus queries video input status
func (c *Client) QueryVideoInputStatus(input string) error {
	// Validate HDMI or SDI input
	valid := false
	for _, hdmi := range c.limits.HDMIInputs {
		if hdmi == input {
			valid = true
			break
		}
	}
	if !valid {
		for _, sdi := range c.limits.SDIInputs {
			if sdi == input {
				valid = true
				break
			}
		}
	}
	if !valid {
		return fmt.Errorf("invalid HDMI/SDI input %s for device %s", input, c.deviceType)
	}
	return c.send(fmt.Sprintf("QVIST:%s;", input))
}

// SetVideoFaderLevel sets video fader level (0-2047)
func (c *Client) SetVideoFaderLevel(level int) error {
	if level < 0 || level > 2047 {
		return fmt.Errorf("video fader level must be between 0 and 2047, got %d", level)
	}
	return c.send(fmt.Sprintf("VFL:%d;", level))
}

// QueryVideoFaderLevel queries video fader level
func (c *Client) QueryVideoFaderLevel() error {
	return c.send("QVFL;")
}

// QueryProgram queries current program channel
func (c *Client) QueryProgram() error { return c.send("QPGM;") }

// QueryPreset queries current preset channel
func (c *Client) QueryPreset() error { return c.send("QPST;") }

// SetAux sets AUX bus output with validation
func (c *Client) SetAux(aux string, input string) error {
	if err := c.validateAuxBus(aux); err != nil {
		return err
	}
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("AUX:%s,%s;", aux, input))
}

// QueryAux queries AUX bus output
func (c *Client) QueryAux(aux string) error {
	if err := c.validateAuxBus(aux); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QAUX:%s;", aux))
}

// AutoTransition performs auto transition
func (c *Client) AutoTransition() error { return c.send("ATO;") }

// AutoTo performs auto transition to specific input with optional time
func (c *Client) AutoTo(input string, tenths int) error {
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	if tenths < -1 || tenths > 40 {
		return fmt.Errorf("transition time must be between -1 and 40, got %d", tenths)
	}
	return c.send(fmt.Sprintf("ATO:%s,%d;", input, tenths))
}

// CutTransition performs cut transition
func (c *Client) CutTransition() error { return c.send("CUT;") }

// CutTo performs cut transition to specific input
func (c *Client) CutTo(input string) error {
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("CUT:%s;", input))
}

// QueryAutoTransition queries auto transition status
func (c *Client) QueryAutoTransition() error {
	return c.send("QATG;")
}

// SetFreeze sets freeze state
func (c *Client) SetFreeze(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("freeze state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("FRZ:%s;", state))
}

// ToggleFreeze toggles freeze state
func (c *Client) ToggleFreeze() error {
	return c.send("FRZ;")
}

// QueryFreeze queries freeze state
func (c *Client) QueryFreeze() error {
	return c.send("QFRZ;")
}

// SetOutputFade sets output fade state
func (c *Client) SetOutputFade(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("output fade state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("FTB:%s;", state))
}

// ToggleOutputFade toggles output fade state
func (c *Client) ToggleOutputFade() error {
	return c.send("FTB;")
}

// QueryOutputFade queries output fade state
func (c *Client) QueryOutputFade() error {
	return c.send("QFTB;")
}

// SetTransitionType sets transition type
func (c *Client) SetTransitionType(transitionType string) error {
	if transitionType != TransitionMIX && transitionType != TransitionWIPE {
		return fmt.Errorf("transition type must be MIX or WIPE, got %s", transitionType)
	}
	return c.send(fmt.Sprintf("TRS:%s;", transitionType))
}

// QueryTransitionType queries transition type
func (c *Client) QueryTransitionType() error {
	return c.send("QTRS;")
}

// SetTransitionTime sets transition time
func (c *Client) SetTransitionTime(transitionType string, time int) error {
	if time < 0 || time > 40 {
		return fmt.Errorf("transition time must be between 0 and 40, got %d", time)
	}
	return c.send(fmt.Sprintf("TIM:%s,%d;", transitionType, time))
}

// QueryTransitionTime queries transition time
func (c *Client) QueryTransitionTime(transitionType string) error {
	return c.send(fmt.Sprintf("QTIM:%s;", transitionType))
}

// SetPinPSource sets PinP source with validation
func (c *Client) SetPinPSource(pinp string, input string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("PIS:%s,%s;", pinp, input))
}

// QueryPinPSource queries PinP source
func (c *Client) QueryPinPSource(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QPIS:%s;", pinp))
}

// SetPinPPGM sets PinP PGM state
func (c *Client) SetPinPPGM(pinp string, state string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("PinP PGM state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("PPS:%s,%s;", pinp, state))
}

// TogglePinPPGM toggles PinP PGM state
func (c *Client) TogglePinPPGM(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("PPS:%s;", pinp))
}

// QueryPinPPGM queries PinP PGM state
func (c *Client) QueryPinPPGM(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QPPS:%s;", pinp))
}

// SetPinPPVW sets PinP PVW state
func (c *Client) SetPinPPVW(pinp string, state string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("PinP PVW state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("PPW:%s,%s;", pinp, state))
}

// TogglePinPPVW toggles PinP PVW state
func (c *Client) TogglePinPPVW(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("PPW:%s;", pinp))
}

// QueryPinPPVW queries PinP PVW state
func (c *Client) QueryPinPPVW(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QPPW:%s;", pinp))
}

// SetPinPPosition sets PinP window position
func (c *Client) SetPinPPosition(pinp string, h, v int) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	if h < -1000 || h > 1000 {
		return fmt.Errorf("PinP horizontal position must be between -1000 and 1000, got %d", h)
	}
	if v < -1000 || v > 1000 {
		return fmt.Errorf("PinP vertical position must be between -1000 and 1000, got %d", v)
	}
	return c.send(fmt.Sprintf("PIP:%s,%d,%d;", pinp, h, v))
}

// QueryPinPPosition queries PinP window position
func (c *Client) QueryPinPPosition(pinp string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QPIP:%s;", pinp))
}

// SetDSK sets DSK state with validation
func (c *Client) SetDSK(dsk string, state string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("DSK state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("DSK:%s,%s;", dsk, state))
}

// ToggleDSK toggles DSK state
func (c *Client) ToggleDSK(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("DSK:%s;", dsk))
}

// QueryDSK queries DSK state
func (c *Client) QueryDSK(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QDSK:%s;", dsk))
}

// SetDSKPVW sets DSK PVW state
func (c *Client) SetDSKPVW(dsk string, state string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("DSK PVW state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("DVW:%s,%s;", dsk, state))
}

// ToggleDSKPVW toggles DSK PVW state
func (c *Client) ToggleDSKPVW(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("DVW:%s;", dsk))
}

// QueryDSKPVW queries DSK PVW state
func (c *Client) QueryDSKPVW(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QDVW:%s;", dsk))
}

// SetDSKFillSource sets DSK fill source
func (c *Client) SetDSKFillSource(dsk string, input string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if err := c.validateVideoInputSource(input); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("DSS:%s,%s;", dsk, input))
}

// QueryDSKFillSource queries DSK fill source
func (c *Client) QueryDSKFillSource(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QDSS:%s;", dsk))
}

// SetDSKLevel sets DSK level (0-255)
func (c *Client) SetDSKLevel(dsk string, level int) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if level < 0 || level > 255 {
		return fmt.Errorf("DSK level must be between 0 and 255, got %d", level)
	}
	return c.send(fmt.Sprintf("KYL:%s,%d;", dsk, level))
}

// QueryDSKLevel queries DSK level
func (c *Client) QueryDSKLevel(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QKYL:%s;", dsk))
}

// SetDSKGain sets DSK gain (0-255)
func (c *Client) SetDSKGain(dsk string, gain int) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if gain < 0 || gain > 255 {
		return fmt.Errorf("DSK gain must be between 0 and 255, got %d", gain)
	}
	return c.send(fmt.Sprintf("KYG:%s,%d;", dsk, gain))
}

// QueryDSKGain queries DSK gain
func (c *Client) QueryDSKGain(dsk string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QKYG:%s;", dsk))
}

// SetSplit sets Split state
func (c *Client) SetSplit(split string, state string) error {
	if err := c.validateSplit(split); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("Split state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("SPS:%s,%s;", split, state))
}

// ToggleSplit toggles Split state
func (c *Client) ToggleSplit(split string) error {
	if err := c.validateSplit(split); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("SPS:%s;", split))
}

// QuerySplit queries Split state
func (c *Client) QuerySplit(split string) error {
	if err := c.validateSplit(split); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QSPS:%s;", split))
}

// SetSplitPosition sets Split position
func (c *Client) SetSplitPosition(split string, pgmCenter, pstCenter int) error {
	return c.SetSplitPositionWithCenter(split, pgmCenter, pstCenter, 0)
}

// SetSplitPositionWithCenter sets Split position with center position
func (c *Client) SetSplitPositionWithCenter(split string, pgmCenter, pstCenter, centerPos int) error {
	if err := c.validateSplit(split); err != nil {
		return err
	}
	if pgmCenter < -500 || pgmCenter > 500 {
		return fmt.Errorf("PGM center must be between -500 and 500, got %d", pgmCenter)
	}
	if pstCenter < -500 || pstCenter > 500 {
		return fmt.Errorf("PST center must be between -500 and 500, got %d", pstCenter)
	}
	if centerPos < -500 || centerPos > 500 {
		return fmt.Errorf("center position must be between -500 and 500, got %d", centerPos)
	}
	return c.send(fmt.Sprintf("SPT:%s,%d,%d,%d;", split, pgmCenter, pstCenter, centerPos))
}

// QuerySplitPosition queries Split position
func (c *Client) QuerySplitPosition(split string) error {
	if err := c.validateSplit(split); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QSPT:%s;", split))
}

// SetStillOutput sets still image output
func (c *Client) SetStillOutput(still string) error {
	if still != StateOFF && still != "V.PLAYER" && still != "SRT" {
		// Validate still image number
		valid := false
		for _, validStill := range c.limits.StillImages {
			if validStill == still {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid still output %s for device %s", still, c.deviceType)
		}
	}
	return c.send(fmt.Sprintf("STO:%s;", still))
}

// QueryStillOutput queries still image output
func (c *Client) QueryStillOutput() error {
	return c.send("QSTO;")
}

// AUDIO Commands

// SetAudioOutputAssign sets audio output assign
func (c *Client) SetAudioOutputAssign(output string, assign string) error {
	return c.send(fmt.Sprintf("AOS:%s,%s;", output, assign))
}

// QueryAudioOutputAssign queries audio output assign
func (c *Client) QueryAudioOutputAssign(output string) error {
	return c.send(fmt.Sprintf("QAOS:%s;", output))
}

// SetInputLevel sets audio input level with validation
func (c *Client) SetInputLevel(input string, deciDb int) error {
	if deciDb < -800 || deciDb > 100 {
		return fmt.Errorf("input level must be between -800 and 100, got %d", deciDb)
	}
	return c.send(fmt.Sprintf("IAL:%s,%d;", input, deciDb))
}

// QueryInputLevel queries audio input level
func (c *Client) QueryInputLevel(input string) error {
	return c.send(fmt.Sprintf("QIAL:%s;", input))
}

// SetInputMute sets audio input mute with validation
func (c *Client) SetInputMute(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("input mute state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("IAM:%s,%s;", input, state))
}

// ToggleInputMute toggles audio input mute
func (c *Client) ToggleInputMute(input string) error {
	return c.send(fmt.Sprintf("IAM:%s;", input))
}

// QueryInputMute queries audio input mute
func (c *Client) QueryInputMute(input string) error {
	return c.send(fmt.Sprintf("QIAM:%s;", input))
}

// SetInputSolo sets audio input solo with validation
func (c *Client) SetInputSolo(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("input solo state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("IAS:%s,%s;", input, state))
}

// ToggleInputSolo toggles audio input solo
func (c *Client) ToggleInputSolo(input string) error {
	return c.send(fmt.Sprintf("IAS:%s;", input))
}

// QueryInputSolo queries audio input solo
func (c *Client) QueryInputSolo(input string) error {
	return c.send(fmt.Sprintf("QIAS:%s;", input))
}

// SetInputDelayTime sets audio input delay time
func (c *Client) SetInputDelayTime(input string, delay int) error {
	if delay < 0 || delay > 5000 {
		return fmt.Errorf("delay time must be between 0 and 5000, got %d", delay)
	}
	return c.send(fmt.Sprintf("ADT:%s,%d;", input, delay))
}

// QueryInputDelayTime queries audio input delay time
func (c *Client) QueryInputDelayTime(input string) error {
	return c.send(fmt.Sprintf("QADT:%s;", input))
}

// SetHPF sets high pass filter with validation
func (c *Client) SetHPF(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("HPF state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("HPF:%s,%s;", input, state))
}

// QueryHPF queries high pass filter
func (c *Client) QueryHPF(input string) error {
	return c.send(fmt.Sprintf("QHPF:%s;", input))
}

// SetGate sets gate with validation
func (c *Client) SetGate(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("gate state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("GATE:%s,%s;", input, state))
}

// QueryGate queries gate
func (c *Client) QueryGate(input string) error {
	return c.send(fmt.Sprintf("QGATE:%s;", input))
}

// SetStereoLink sets stereo link with validation
func (c *Client) SetStereoLink(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("stereo link state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("STLK:%s,%s;", input, state))
}

// ToggleStereoLink toggles stereo link
func (c *Client) ToggleStereoLink(input string) error {
	return c.send(fmt.Sprintf("STLK:%s;", input))
}

// QueryStereoLink queries stereo link
func (c *Client) QueryStereoLink(input string) error {
	return c.send(fmt.Sprintf("QSTLK:%s;", input))
}

// SetVoiceChanger sets voice changer with validation
func (c *Client) SetVoiceChanger(input string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("voice changer state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("VOCH:%s,%s;", input, state))
}

// ToggleVoiceChanger toggles voice changer
func (c *Client) ToggleVoiceChanger(input string) error {
	return c.send(fmt.Sprintf("VOCH:%s;", input))
}

// QueryVoiceChanger queries voice changer
func (c *Client) QueryVoiceChanger(input string) error {
	return c.send(fmt.Sprintf("QVOCH:%s;", input))
}

// SetOutputLevel sets audio output level with validation
func (c *Client) SetOutputLevel(output string, deciDb int) error {
	if deciDb < -800 || deciDb > 100 {
		return fmt.Errorf("output level must be between -800 and 100, got %d", deciDb)
	}
	return c.send(fmt.Sprintf("OAL:%s,%d;", output, deciDb))
}

// QueryOutputLevel queries audio output level
func (c *Client) QueryOutputLevel(output string) error {
	return c.send(fmt.Sprintf("QOAL:%s;", output))
}

// SetOutputMute sets audio output mute with validation
func (c *Client) SetOutputMute(output string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("output mute state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("OAM:%s,%s;", output, state))
}

// QueryOutputMute queries audio output mute
func (c *Client) QueryOutputMute(output string) error {
	return c.send(fmt.Sprintf("QOAM:%s;", output))
}

// SetReverb sets reverb with validation
func (c *Client) SetReverb(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("reverb state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("RVB:%s;", state))
}

// ToggleReverb toggles reverb
func (c *Client) ToggleReverb() error {
	return c.send("RVB;")
}

// QueryReverb queries reverb
func (c *Client) QueryReverb() error {
	return c.send("QRVB;")
}

// SetAudioAutoMixing sets audio auto mixing with validation
func (c *Client) SetAudioAutoMixing(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("audio auto mixing state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("ATM:%s;", state))
}

// ToggleAudioAutoMixing toggles audio auto mixing
func (c *Client) ToggleAudioAutoMixing() error {
	return c.send("ATM;")
}

// QueryAudioAutoMixing queries audio auto mixing
func (c *Client) QueryAudioAutoMixing() error {
	return c.send("QATM;")
}

// METER Commands

// SetMeterAutoSend sets meter auto send with validation
func (c *Client) SetMeterAutoSend(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("meter auto send state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("MTRSW:%s;", state))
}

// QueryMeterAutoSend queries meter auto send
func (c *Client) QueryMeterAutoSend() error {
	return c.send("QMTRSW;")
}

// QueryMeterPFL queries meter PFL
func (c *Client) QueryMeterPFL() error {
	return c.send("MTRLV:PFL;")
}

// QueryMeterAFL queries meter AFL
func (c *Client) QueryMeterAFL() error {
	return c.send("MTRLV:AFL;")
}

// QueryMeterChannel queries meter channel
func (c *Client) QueryMeterChannel(channel int) error {
	return c.send(fmt.Sprintf("MTRCH:%d;", channel))
}

// SetCompGRAutoSend sets comp GR auto send with validation
func (c *Client) SetCompGRAutoSend(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("comp GR auto send state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("GRSW:%s;", state))
}

// QueryCompGRAutoSend queries comp GR auto send
func (c *Client) QueryCompGRAutoSend() error {
	return c.send("QGRSW;")
}

// QueryCompGRLevel queries comp GR level
func (c *Client) QueryCompGRLevel() error {
	return c.send("GRLV;")
}

// QueryCompGRChannel queries comp GR channel
func (c *Client) QueryCompGRChannel(channel int) error {
	return c.send(fmt.Sprintf("GRCH:%d;", channel))
}

// SetAutoMixingAutoSend sets auto mixing auto send with validation
func (c *Client) SetAutoMixingAutoSend(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("auto mixing auto send state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("AMSW:%s;", state))
}

// QueryAutoMixingAutoSend queries auto mixing auto send
func (c *Client) QueryAutoMixingAutoSend() error {
	return c.send("QAMSW;")
}

// QueryAutoMixingLevel queries auto mixing level
func (c *Client) QueryAutoMixingLevel() error {
	return c.send("AMLV;")
}

// QueryAutoMixingChannel queries auto mixing channel
func (c *Client) QueryAutoMixingChannel(channel int) error {
	return c.send(fmt.Sprintf("AMCH:%d;", channel))
}

// SetSigPeakAutoSend sets sig/peak auto send with validation
func (c *Client) SetSigPeakAutoSend(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("sig/peak auto send state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("SPSW:%s;", state))
}

// QuerySigPeakAutoSend queries sig/peak auto send
func (c *Client) QuerySigPeakAutoSend() error {
	return c.send("QSPSW;")
}

// QuerySigPeakLevel queries sig/peak level
func (c *Client) QuerySigPeakLevel() error {
	return c.send("SPLV;")
}

// QuerySigPeakChannel queries sig/peak channel
func (c *Client) QuerySigPeakChannel(channel int) error {
	return c.send(fmt.Sprintf("SPCH:%d;", channel))
}

// SetAuxAutoSend sets aux auto send with validation
func (c *Client) SetAuxAutoSend(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("aux auto send state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("AUXSW:%s;", state))
}

// QueryAuxAutoSend queries aux auto send
func (c *Client) QueryAuxAutoSend() error {
	return c.send("QAUXSW;")
}

// QueryAuxLevel queries aux level
func (c *Client) QueryAuxLevel(aux string) error {
	if err := c.validateAuxBus(aux); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("AUXLV:%s;", aux))
}

// QueryAuxChannel queries aux channel
func (c *Client) QueryAuxChannel(channel int) error {
	return c.send(fmt.Sprintf("AUXCH:%d;", channel))
}

// CONTROL Commands

// RecallMemory recalls memory with validation
func (c *Client) RecallMemory(memory string) error {
	if err := c.validateMemory(memory); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("MEM:%s;", memory))
}

// QueryMemory queries current memory
func (c *Client) QueryMemory() error {
	return c.send("QMEM;")
}

// SetGPO sets GPO output with validation
func (c *Client) SetGPO(gpo string, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("GPO state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("GPO:%s,%s;", gpo, state))
}

// ToggleGPO toggles GPO output
func (c *Client) ToggleGPO(gpo string) error {
	return c.send(fmt.Sprintf("GPO:%s;", gpo))
}

// QueryGPO queries GPO state
func (c *Client) QueryGPO(gpo string) error {
	return c.send(fmt.Sprintf("QGPO:%s;", gpo))
}

// QueryTally queries tally state
func (c *Client) QueryTally() error {
	return c.send("TLY;")
}

// SetAutoSwitching sets auto switching with validation
func (c *Client) SetAutoSwitching(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("auto switching state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("ASW:%s;", state))
}

// ToggleAutoSwitching toggles auto switching
func (c *Client) ToggleAutoSwitching() error {
	return c.send("ASW;")
}

// QueryAutoSwitching queries auto switching
func (c *Client) QueryAutoSwitching() error {
	return c.send("QASW;")
}

// ExecuteInputScan executes input scan
func (c *Client) ExecuteInputScan(mode string) error {
	if mode != ScanNORMAL && mode != ScanREVERSE && mode != ScanRANDOM {
		return fmt.Errorf("input scan mode must be NORMAL, REVERSE, or RANDOM, got %s", mode)
	}
	return c.send(fmt.Sprintf("INSC:%s;", mode))
}

// ExecuteMemoryScan executes memory scan
func (c *Client) ExecuteMemoryScan(mode string) error {
	if mode != ScanNORMAL && mode != ScanREVERSE && mode != ScanRANDOM {
		return fmt.Errorf("memory scan mode must be NORMAL, REVERSE, or RANDOM, got %s", mode)
	}
	return c.send(fmt.Sprintf("MEMSC:%s;", mode))
}

// ExecutePinPSourceScan executes PinP source scan
func (c *Client) ExecutePinPSourceScan(pinp string, mode string) error {
	if err := c.validatePinP(pinp); err != nil {
		return err
	}
	if mode != ScanNORMAL && mode != ScanREVERSE && mode != ScanRANDOM {
		return fmt.Errorf("PinP source scan mode must be NORMAL, REVERSE, or RANDOM, got %s", mode)
	}
	return c.send(fmt.Sprintf("PPSC:%s,%s;", pinp, mode))
}

// ExecuteDSKSourceScan executes DSK source scan
func (c *Client) ExecuteDSKSourceScan(dsk string, mode string) error {
	if err := c.validateDSK(dsk); err != nil {
		return err
	}
	if mode != ScanNORMAL && mode != ScanREVERSE && mode != ScanRANDOM {
		return fmt.Errorf("DSK source scan mode must be NORMAL, REVERSE, or RANDOM, got %s", mode)
	}
	return c.send(fmt.Sprintf("DSKSC:%s,%s;", dsk, mode))
}

// SetStreamType sets stream type
func (c *Client) SetStreamType(streamType string) error {
	if streamType != "RTMP" && streamType != "SRT" {
		return fmt.Errorf("stream type must be RTMP or SRT, got %s", streamType)
	}
	return c.send(fmt.Sprintf("STRTY:%s;", streamType))
}

// QueryStreamType queries stream type
func (c *Client) QueryStreamType() error {
	return c.send("QSTRTY;")
}

// StreamStart starts stream
func (c *Client) StreamStart() error {
	return c.send("STROA:START;")
}

// StreamStop stops stream
func (c *Client) StreamStop() error {
	return c.send("STROA:STOP;")
}

// QueryStreamStatus queries stream status
func (c *Client) QueryStreamStatus() error {
	return c.send("QSTRST;")
}

// QueryStreamTime queries stream time
func (c *Client) QueryStreamTime() error {
	return c.send("QSTRTM;")
}

// SetSRTOutStart starts SRT out
func (c *Client) SetSRTOutStart() error {
	return c.send("SRTOSS:START;")
}

// SetSRTOutStop stops SRT out
func (c *Client) SetSRTOutStop() error {
	return c.send("SRTOSS:STOP;")
}

// SetSafetyImage sets safety image with validation
func (c *Client) SetSafetyImage(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("safety image state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("SFIM:%s;", state))
}

// QuerySafetyImage queries safety image
func (c *Client) QuerySafetyImage() error {
	return c.send("QSFIM;")
}

// SetExternalRecControl sets external rec control with validation
func (c *Client) SetExternalRecControl(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("external rec control state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("EXRC:%s;", state))
}

// QueryExternalRecControl queries external rec control
func (c *Client) QueryExternalRecControl() error {
	return c.send("QEXRC;")
}

// SetAudioPlayer sets audio player with validation
func (c *Client) SetAudioPlayer(file int, action string) error {
	if action != PlayerSTOP && action != PlayerPLAY && action != PlayerPAUSE && action != PlayerRESUME {
		return fmt.Errorf("audio player action must be STOP, PLAY, PAUSE, or RESUME, got %s", action)
	}
	return c.send(fmt.Sprintf("APL:%d,%s;", file, action))
}

// QueryAudioPlayerStatus queries audio player status
func (c *Client) QueryAudioPlayerStatus(file int) error {
	return c.send(fmt.Sprintf("QAPS:%d;", file))
}

// QueryAudioPlayerImport queries audio player import
func (c *Client) QueryAudioPlayerImport(file int) error {
	return c.send(fmt.Sprintf("QAPE:%d;", file))
}

// QueryAudioPlayerDuration queries audio player duration
func (c *Client) QueryAudioPlayerDuration(file int) error {
	return c.send(fmt.Sprintf("QAPD:%d;", file))
}

// QueryAudioPlayerCurrentTime queries audio player current time
func (c *Client) QueryAudioPlayerCurrentTime(file int) error {
	return c.send(fmt.Sprintf("QAPC:%d;", file))
}

// SetAudioPlayerRepeat sets audio player repeat with validation
func (c *Client) SetAudioPlayerRepeat(file int, state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("audio player repeat state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("APR:%d,%s;", file, state))
}

// QueryAudioPlayerRepeat queries audio player repeat
func (c *Client) QueryAudioPlayerRepeat(file int) error {
	return c.send(fmt.Sprintf("QAPR:%d;", file))
}

// SetVideoPlayerType sets video player type
func (c *Client) SetVideoPlayerType(playerType string) error {
	if playerType != "V.PLAYER" && playerType != "SRT" {
		return fmt.Errorf("video player type must be V.PLAYER or SRT, got %s", playerType)
	}
	return c.send(fmt.Sprintf("PBTY:%s;", playerType))
}

// QueryVideoPlayerType queries video player type
func (c *Client) QueryVideoPlayerType() error {
	return c.send("QPBTY;")
}

// SetVideoPlayer sets video player with validation
func (c *Client) SetVideoPlayer(action string) error {
	if action != PlayerSTOP && action != PlayerPLAY && action != PlayerPAUSE && action != PlayerRESUME {
		return fmt.Errorf("video player action must be STOP, PLAY, PAUSE, or RESUME, got %s", action)
	}
	return c.send(fmt.Sprintf("VPL:%s;", action))
}

// QueryVideoPlayerStatus queries video player status
func (c *Client) QueryVideoPlayerStatus() error {
	return c.send("QVPS;")
}

// QueryVideoPlayerImport queries video player import
func (c *Client) QueryVideoPlayerImport() error {
	return c.send("QVPE;")
}

// QueryVideoPlayerDuration queries video player duration
func (c *Client) QueryVideoPlayerDuration() error {
	return c.send("QVPD;")
}

// QueryVideoPlayerCurrentTime queries video player current time
func (c *Client) QueryVideoPlayerCurrentTime() error {
	return c.send("QVPC;")
}

// SetSRTInStart starts SRT in
func (c *Client) SetSRTInStart() error {
	return c.send("SRTISS:START;")
}

// SetSRTInStop stops SRT in
func (c *Client) SetSRTInStop() error {
	return c.send("SRTISS:STOP;")
}

// QuerySRTInStatus queries SRT in status
func (c *Client) QuerySRTInStatus() error {
	return c.send("QSRTIST;")
}

// QueryStreamStability queries stream stability
func (c *Client) QueryStreamStability() error {
	return c.send("STRSTB;")
}

// QueryVideoPlayerStability queries video player stability
func (c *Client) QueryVideoPlayerStability() error {
	return c.send("VPSTB;")
}

// QuerySRTInStability queries SRT in stability
func (c *Client) QuerySRTInStability() error {
	return c.send("SRTISTB;")
}

// CAMERA Commands (PTZ)

// CameraPanTilt sets camera pan/tilt with validation
func (c *Client) CameraPanTilt(camera string, horiz string, vert string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if horiz != CameraLEFT && horiz != CameraSTOP && horiz != CameraRIGHT {
		return fmt.Errorf("horizontal direction must be LEFT, STOP, or RIGHT, got %s", horiz)
	}
	if vert != CameraDOWN && vert != CameraSTOP && vert != CameraUP {
		return fmt.Errorf("vertical direction must be DOWN, STOP, or UP, got %s", vert)
	}
	return c.send(fmt.Sprintf("CAMPT:%s,%s,%s;", camera, horiz, vert))
}

// CameraPanTiltSpeed sets camera pan/tilt speed
func (c *Client) CameraPanTiltSpeed(camera string, speed int) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if speed < 1 || speed > 24 {
		return fmt.Errorf("pan/tilt speed must be between 1 and 24, got %d", speed)
	}
	return c.send(fmt.Sprintf("CAMPTS:%s,%d;", camera, speed))
}

// QueryCameraPanTiltSpeed queries camera pan/tilt speed
func (c *Client) QueryCameraPanTiltSpeed(camera string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QCAMPTS:%s;", camera))
}

// CameraZoom sets camera zoom with validation
func (c *Client) CameraZoom(camera string, speed string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if speed != ZoomWIDE_FAST && speed != ZoomWIDE_SLOW && speed != CameraSTOP &&
		speed != ZoomTELE_SLOW && speed != ZoomTELE_FAST {
		return fmt.Errorf("zoom speed must be WIDE_FAST, WIDE_SLOW, STOP, TELE_SLOW, or TELE_FAST, got %s", speed)
	}
	return c.send(fmt.Sprintf("CAMZM:%s,%s;", camera, speed))
}

// CameraZoomReset resets camera zoom
func (c *Client) CameraZoomReset(camera string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("CAMZMR:%s;", camera))
}

// CameraFocus sets camera focus with validation
func (c *Client) CameraFocus(camera string, direction string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if direction != FocusNEAR && direction != CameraSTOP && direction != FocusFAR {
		return fmt.Errorf("focus direction must be NEAR, STOP, or FAR, got %s", direction)
	}
	return c.send(fmt.Sprintf("CAMFC:%s,%s;", camera, direction))
}

// CameraAutoFocus sets camera auto focus with validation
func (c *Client) CameraAutoFocus(camera string, state string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("auto focus state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("CAMAFC:%s,%s;", camera, state))
}

// QueryCameraAutoFocus queries camera auto focus
func (c *Client) QueryCameraAutoFocus(camera string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QCAMAFC:%s;", camera))
}

// CameraAutoExposure sets camera auto exposure with validation
func (c *Client) CameraAutoExposure(camera string, state string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if state != StateOFF && state != StateON {
		return fmt.Errorf("auto exposure state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("CAMAEP:%s,%s;", camera, state))
}

// QueryCameraAutoExposure queries camera auto exposure
func (c *Client) QueryCameraAutoExposure(camera string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QCAMAEP:%s;", camera))
}

// CameraPreset sets camera preset with validation
func (c *Client) CameraPreset(camera string, preset string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	if err := c.validatePreset(preset); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("CAMPR:%s,%s;", camera, preset))
}

// QueryCameraPreset queries camera preset
func (c *Client) QueryCameraPreset(camera string) error {
	if err := c.validateCamera(camera); err != nil {
		return err
	}
	return c.send(fmt.Sprintf("QCAMPR:%s;", camera))
}

// SYSTEM Commands

// SendACK sends ACK command
func (c *Client) SendACK() error {
	return c.send("ACS;")
}

// QueryVersion queries version
func (c *Client) QueryVersion() error {
	return c.send("VER;")
}

// QueryBusy queries busy state
func (c *Client) QueryBusy() error {
	return c.send("QBSY;")
}

// SetHDCP sets HDCP with validation
func (c *Client) SetHDCP(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("HDCP state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("HDCP:%s;", state))
}

// QueryHDCP queries HDCP
func (c *Client) QueryHDCP() error {
	return c.send("QHDCP;")
}

// TestPattern sets test pattern with validation
func (c *Client) TestPattern(pattern string) error {
	validPatterns := []string{
		"OFF", "COLORBAR75", "COLORBAR100", "RAMP", "STEP", "HATCH", "DIAMOND", "CIRCLE",
		"COLORBAR75-SP", "COLORBAR100-SP", "RAMP-SP", "STEP-SP", "HATCH-SP",
	}
	for _, valid := range validPatterns {
		if valid == pattern {
			return c.send(fmt.Sprintf("TPT:%s;", pattern))
		}
	}
	return fmt.Errorf("invalid test pattern %s", pattern)
}

// QueryTestPattern queries test pattern
func (c *Client) QueryTestPattern() error {
	return c.send("QTPT;")
}

// TestTone sets test tone with validation
func (c *Client) TestTone(level string) error {
	return c.TestToneWithFrequency(level, "1k", "1k")
}

// TestToneWithFrequency sets test tone with frequency
func (c *Client) TestToneWithFrequency(level, freqL, freqR string) error {
	validLevels := []string{"OFF", "-20", "-10", "0dB"}
	validFreqs := []string{"500", "1k", "2kHz"}

	levelValid := false
	for _, valid := range validLevels {
		if valid == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		return fmt.Errorf("invalid test tone level %s", level)
	}

	freqLValid := false
	for _, valid := range validFreqs {
		if valid == freqL {
			freqLValid = true
			break
		}
	}
	if !freqLValid {
		return fmt.Errorf("invalid test tone frequency L %s", freqL)
	}

	freqRValid := false
	for _, valid := range validFreqs {
		if valid == freqR {
			freqRValid = true
			break
		}
	}
	if !freqRValid {
		return fmt.Errorf("invalid test tone frequency R %s", freqR)
	}

	return c.send(fmt.Sprintf("TTN:%s,%s,%s;", level, freqL, freqR))
}

// QueryTestTone queries test tone
func (c *Client) QueryTestTone() error {
	return c.send("QTTN;")
}

// MACRO, SEQUENCER, GRAPHICS PRESENTER Commands

// ExecuteMacro executes macro
func (c *Client) ExecuteMacro(macro int) error {
	if macro < 1 || macro > 100 {
		return fmt.Errorf("macro number must be between 1 and 100, got %d", macro)
	}
	return c.send(fmt.Sprintf("MCREX:%d;", macro))
}

// QueryMacroStatus queries macro status
func (c *Client) QueryMacroStatus(macro int) error {
	if macro < 1 || macro > 100 {
		return fmt.Errorf("macro number must be between 1 and 100, got %d", macro)
	}
	return c.send(fmt.Sprintf("QMCRST:%d;", macro))
}

// SetSequencer sets sequencer with validation
func (c *Client) SetSequencer(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("sequencer state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("SEQSW:%s;", state))
}

// QuerySequencer queries sequencer
func (c *Client) QuerySequencer() error {
	return c.send("QSEQSW;")
}

// SetSequencerAutoSequence sets sequencer auto sequence with validation
func (c *Client) SetSequencerAutoSequence(state string) error {
	if state != StateOFF && state != StateON {
		return fmt.Errorf("sequencer auto sequence state must be OFF or ON, got %s", state)
	}
	return c.send(fmt.Sprintf("SEQAS:%s;", state))
}

// QuerySequencerAutoSequence queries sequencer auto sequence
func (c *Client) QuerySequencerAutoSequence() error {
	return c.send("QSEQAS;")
}

// SequencerPrevious goes to previous sequence
func (c *Client) SequencerPrevious() error {
	return c.send("SEQPV;")
}

// SequencerPreviousToStart goes to start of sequence
func (c *Client) SequencerPreviousToStart() error {
	return c.send("SEQPV:1;")
}

// SequencerNext goes to next sequence
func (c *Client) SequencerNext() error {
	return c.send("SEQNX;")
}

// SequencerJump jumps to specific sequence
func (c *Client) SequencerJump(sequence string) error {
	if sequence != "START" {
		// Validate sequence number
		if len(sequence) < 4 || sequence[:3] != "SEQ" {
			return fmt.Errorf("sequence must be START or SEQ1-SEQ1000, got %s", sequence)
		}
	}
	return c.send(fmt.Sprintf("SEQJP:%s;", sequence))
}

// GraphicsPresenterNextContent selects next content
func (c *Client) GraphicsPresenterNextContent() error {
	return c.send("GPNC;")
}

// GraphicsPresenterSelectContent selects specific content
func (c *Client) GraphicsPresenterSelectContent(content int) error {
	if content < 1 || content > 124 {
		return fmt.Errorf("content number must be between 1 and 124, got %d", content)
	}
	return c.send(fmt.Sprintf("GPSC:CONTENT%d;", content))
}

// GraphicsPresenterHideFront hides front content
func (c *Client) GraphicsPresenterHideFront() error {
	return c.send("GPHF;")
}

// GraphicsPresenterHideBackground hides background content
func (c *Client) GraphicsPresenterHideBackground() error {
	return c.send("GPHB;")
}

// GraphicsPresenterOnAir toggles on air
func (c *Client) GraphicsPresenterOnAir() error {
	return c.send("GPOA;")
}
