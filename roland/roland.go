package roland

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"golang.org/x/xerrors"
)

type Client struct {
	conn net.Conn
}

const (
	stx = "\x02"
)

// NewClient connects to a Roland switcher-compatible device (e.g., V-160HD, V-80HD, VR-120HD, VR-6HD)
// via TCP port (typically 8023).
func NewClient(ipAddress string, port string) (*Client, error) {
	if ipAddress == "" {
		return nil, errors.New("ipAddress must not be empty")
	}
	if port == "" {
		return nil, errors.New("port must not be empty")
	}
	address := net.JoinHostPort(ipAddress, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, xerrors.Errorf("dial tcp %s: %w", address, err)
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return xerrors.Errorf("close tcp connection: %w", err)
	}
	return nil
}

// send writes a command ensuring it is framed as STX + body + ';'
func (c *Client) send(command string) error {
	if len(command) == 0 || command[len(command)-1] != ';' {
		command = command + ";"
	}
	payload := fmt.Sprintf("%s%s", stx, command)
	if _, err := c.conn.Write([]byte(payload)); err != nil {
		return xerrors.Errorf("write command %q: %w", command, err)
	}
	return nil
}

// SendRaw sends a prepared command body (e.g., "PGM:INPUT1;")
func (c *Client) SendRaw(raw string) error { return c.send(raw) }

// ReadResponse reads until LF following the final ';'
func (c *Client) ReadResponse() (string, error) {
	reader := bufio.NewReader(c.conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", xerrors.Errorf("read response: %w", err)
	}
	return line, nil
}

// VIDEO Commands

// PGM selects program channel by numeric index (legacy helper). Prefer SetProgram("INPUT1").
func (c *Client) PGM(channel int) error { return c.send(fmt.Sprintf("PGM:%d;", channel)) }

func (c *Client) SetProgram(input string) error { return c.send(fmt.Sprintf("PGM:%s;", input)) }
func (c *Client) SetPreset(input string) error  { return c.send(fmt.Sprintf("PST:%s;", input)) }

func (c *Client) AutoTransition() error { return c.send("ATO;") }
func (c *Client) CutTransition() error  { return c.send("CUT;") }
func (c *Client) AutoTo(input string, tenths int) error {
	return c.send(fmt.Sprintf("ATO:%s,%d;", input, tenths))
}

func (c *Client) QueryProgram() error { return c.send("QPGM;") }
func (c *Client) QueryPreset() error  { return c.send("QPST;") }

func (c *Client) SetPinP(pinp string, input string) error {
	return c.send(fmt.Sprintf("PIS:%s,%s;", pinp, input))
}
func (c *Client) SetDSK(dsk string, state string) error {
	return c.send(fmt.Sprintf("DSK:%s,%s;", dsk, state))
}

// AUDIO Commands
func (c *Client) SetInputLevel(input string, deciDb int) error {
	return c.send(fmt.Sprintf("IAL:%s,%d;", input, deciDb))
}
func (c *Client) SetInputMute(input string, state string) error {
	return c.send(fmt.Sprintf("IAM:%s,%s;", input, state))
}
func (c *Client) SetOutputMute(output string, state string) error {
	return c.send(fmt.Sprintf("OAM:%s,%s;", output, state))
}
func (c *Client) SetHPF(input string, state string) error {
	return c.send(fmt.Sprintf("HPF:%s,%s;", input, state))
}
func (c *Client) SetGate(input string, state string) error {
	return c.send(fmt.Sprintf("GATE:%s,%s;", input, state))
}

// METER Commands
func (c *Client) QueryMeterPFL() error                { return c.send("MTRLV:PFL;") }
func (c *Client) SetMeterAutoSend(state string) error { return c.send(fmt.Sprintf("MTRSW:%s;", state)) }

// CONTROL Commands
func (c *Client) RecallMemory(memory string) error { return c.send(fmt.Sprintf("MEM:%s;", memory)) }
func (c *Client) QueryMemory() error               { return c.send("QMEM;") }
func (c *Client) QueryTally() error                { return c.send("TLY;") }
func (c *Client) StreamStart() error               { return c.send("STROA:START;") }
func (c *Client) StreamStop() error                { return c.send("STROA:STOP;") }
func (c *Client) QueryStreamStatus() error         { return c.send("QSTRST;") }

// CAMERA Commands (PTZ)
func (c *Client) CameraPanTilt(camera string, horiz string, vert string) error {
	return c.send(fmt.Sprintf("CAMPT:%s,%s,%s;", camera, horiz, vert))
}
func (c *Client) CameraZoom(camera string, speed string) error {
	return c.send(fmt.Sprintf("CAMZM:%s,%s;", camera, speed))
}
func (c *Client) CameraPreset(camera string, preset string) error {
	return c.send(fmt.Sprintf("CAMPR:%s,%s;", camera, preset))
}

// SYSTEM Commands
func (c *Client) QueryVersion() error              { return c.send("VER;") }
func (c *Client) TestPattern(pattern string) error { return c.send(fmt.Sprintf("TPT:%s;", pattern)) }
func (c *Client) SetHDCP(state string) error       { return c.send(fmt.Sprintf("HDCP:%s;", state)) }
