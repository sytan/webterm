package serial

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Serial struct {
	f          *os.File
	LineIgnore string
	LineEnd    string
}

const (
	PAR_NONE = iota
	PAR_EVEN
	PAR_ODD
)

const (
	FLUSH_I = iota
	FLUSH_O
	FLUSH_IO
)

//Open opens serial with default params.
//	Default: 9600 8N1, soft/hard flow controll off.
func Open(name string) (*Serial, error) {
	f, err := os.OpenFile(name, os.O_RDWR|syscall.O_NOCTTY, 0600)
	if err != nil {
		return nil, err
	}
	s := &Serial{f: f, LineIgnore: "\r", LineEnd: "\n"}
	err = s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

//Close closes serial.
func (s *Serial) Close() error {
	err := s.f.Close()
	s.f = nil
	return err
}

//Read reads slice from serial.
func (s *Serial) Read(b []byte) (int, error) {
	return s.f.Read(b)
}

//WriteString writes string to serial.
func (s *Serial) WriteString(str string) (int, error) {
	return s.f.WriteString(str)
}

//Write writes byte slice to serial.
func (s *Serial) Write(b []byte) (int, error) {
	return s.f.Write(b)
}

//WriteByte writes one byte to serial.
func (s *Serial) WriteByte(c byte) error {
	_, e := s.f.Write([]byte{c})
	return e
}

//ReadByte reads one byte from serial.
func (s *Serial) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	n, e := s.f.Read(buf)
	if n == 1 {
		return buf[0], nil
	}
	return 0, e
}

//Name returns serial file name.
func (s *Serial) Name() string {
	return s.f.Name()
}

//File returns serial os.File struct.
func (s *Serial) File() *os.File {
	return s.f
}

//Fd returns serial file descriptor.
func (s *Serial) Fd() uintptr {
	return s.f.Fd()
}

//SetSpeed sets serial speed.
func (s *Serial) SetSpeed(speed int) error {
	return s.setSpeed(speed)
}

//SetHwFlowCtrl enable or disable Hardware flow control.
func (s *Serial) SetHwFlowCtrl(hw bool) error {
	return s.setHwFlowCtrl(hw)
}

//SetSwFlowCtrl enable or disable software flow control.
func (s *Serial) SetSwFlowCtrl(sw bool) error {
	return s.setSwFlowCtrl(sw)
}

//SetStopBits sets stop bits, valid values are 1 or 2.
func (s *Serial) SetStopBits(stop int) error {
	switch stop {
	case 1:
		return s.setStopBits2(false)
	case 2:
		return s.setStopBits2(true)
	default:
		return errors.New("Invalid stop bits number")
	}
}

//SetParity sets parity mode:
//	PAR_NONE
//	PAR_EVEN
//	PAR_ODD
func (s *Serial) SetParity(mode int) error {
	return s.setParity(mode)
}

//SetLocal sets local mode. In local mode, modem control lines are ignored.
func (s *Serial) SetLocal(local bool) error {
	return s.setLocal(local)
}

//SetReadTimeout sets read behavior.
//  vmin  - minimum number of characters for Read.
//  vtime - precision is 1/10 second, max timeout is 25.5 seconds.
//
//  vmin == 0 && vtime == 0 : non-blocking Read.
//  vmin == 0 && vtime > 0  : Read returns buffered characters or waits vtime for new charcters.
//  vmin > 0  && vtime > 0  : Read returns n >= vmin characters or 0 < n < vmin if vtime expires after n-th char.
//  vmin > 0  && vtime == 0 : Read returns at least vmin characters.
func (s *Serial) SetReadTimeout(vmin int, vtime time.Duration) error {
	return s.setReadTimeout(vmin, vtime)
}

//GetAttr sets Termios structure from serial attributes.
func (s *Serial) GetAttr(attr *Termios) error {
	return s.tcGetAttr(attr)
}

//SetAttr sets serial attributes from Termios structure.
func (s *Serial) SetAttr(attr *Termios) error {
	return s.tcSetAttr(attr)
}

//SetHub sets hangup mode (false -> don't reset DTR/RTS on exit).
func (s *Serial) SetHup(hup bool) error {
	return s.setHup(hup)
}

//SetDTR sets DTR signal level.
func (s *Serial) SetDTR(level bool) error {
	return s.setDTR(level)
}

//SetRTS sets RTS signal level.
func (s *Serial) SetRTS(level bool) error {
	return s.setRTS(level)
}

//InpWaiting returns number of bytes waiting on input buffer.
func (s *Serial) InpWaiting() (int, error) {
	return s.inpWaiting()
}

//OutWaiting returns number of bytes waiting on output buffer.
func (s *Serial) OutWaiting() (int, error) {
	return s.outWaiting()
}

//Flush buffers selected by mode:
//	FLUSH_I  input buffer
//	FLUSH_O  output buffer
//	FLUSH_IO input and output buffers
func (s *Serial) Flush(mode int) error {
	return s.flush(mode)
}

//ReadLine reads text line.
//Timeout is the max time to wait for line completion, when timeout occurs EOF error is returned.
//Serial.LineIgnore field has characters to be ignored (by default "\r").
//Serial.LineEnd field has end of line characters (by default "\n").
func (s *Serial) ReadLine(timeout time.Duration) (res string, err error) {
	var t Termios
	var b byte

	if err = s.GetAttr(&t); err != nil {
		return "", err
	}
	defer s.SetAttr(&t)
	s.SetReadTimeout(0, timeout)
	for {
		if b, err = s.ReadByte(); err != nil {
			return
		}
		ch := string(b)
		if strings.Contains(s.LineIgnore, ch) {
			continue
		}
		if strings.Contains(s.LineEnd, ch) {
			break
		}
		res += ch
	}
	return
}

//WaitForRe reads lines from serial and waits for line matching one regular expresion from rexp slice.
//It returns the index of rexp slice matching text line, text line itself and error != nil on timeout or I/O error.
func (s *Serial) WaitForRe(rexp []string, timeout time.Duration) (int, string, error) {
	var match string
	var err error

	for {
		if match, err = s.ReadLine(timeout); err != nil {
			return -1, "", err
		}
		for i, re := range rexp {
			ok, err := regexp.MatchString(re, match)
			if err != nil {
				return -1, "", err
			}
			if ok {
				return i, match, nil
			}
		}
	}
}
