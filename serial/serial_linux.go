package serial

import (
	"errors"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type Termios syscall.Termios

var baud = map[int]uint32{
	0:       syscall.B0,
	50:      syscall.B50,
	75:      syscall.B75,
	110:     syscall.B110,
	134:     syscall.B134,
	150:     syscall.B150,
	200:     syscall.B200,
	300:     syscall.B300,
	600:     syscall.B600,
	1200:    syscall.B1200,
	1800:    syscall.B1800,
	2400:    syscall.B2400,
	4800:    syscall.B4800,
	9600:    syscall.B9600,
	19200:   syscall.B19200,
	38400:   syscall.B38400,
	57600:   syscall.B57600,
	115200:  syscall.B115200,
	230400:  syscall.B230400,
	460800:  syscall.B460800,
	500000:  syscall.B500000,
	576000:  syscall.B576000,
	921600:  syscall.B921600,
	1000000: syscall.B1000000,
	1152000: syscall.B1152000,
	1500000: syscall.B1500000,
	2000000: syscall.B2000000,
	2500000: syscall.B2500000,
	3000000: syscall.B3000000,
	3500000: syscall.B3500000,
	4000000: syscall.B4000000,
}

var bits = map[int]uint32{
	5: syscall.CS5,
	6: syscall.CS6,
	7: syscall.CS7,
	8: syscall.CS8,
}

// Constants not defined in syscall module
const (
	cbaud   = 0010017
	cbaudex = 0010000
	crtscts = 020000000000
	tcflsh  = 0x540B
)

func (s *Serial) tcGetAttr(cfg *Termios) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(cfg)),
	)
	if e != 0 {
		return os.NewSyscallError("tcgetattr", e)
	}
	return nil
}

func (s *Serial) tcSetAttr(cfg *Termios) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		syscall.TCSETS,
		uintptr(unsafe.Pointer(cfg)),
	)
	if e != 0 {
		return os.NewSyscallError("tcsetattr", e)
	}
	return nil
}

func (s *Serial) init() error {
	var t Termios
	t.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	t.Oflag &^= (syscall.OPOST)
	t.Lflag &^= (syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	t.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	t.Cflag |= (baud[9600] | bits[8] | syscall.CREAD | syscall.HUPCL)
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
	t.Ispeed = baud[9600]
	t.Ospeed = baud[9600]
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setSpeed(b int) error {
	var t Termios
	bb, ok := baud[b]
	if !ok {
		return errors.New("Unknown boud rate")
	}
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	t.Cflag &^= cbaud | cbaudex
	t.Cflag |= bb
	t.Ispeed = bb
	t.Ospeed = bb
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setParity(mode int) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	switch mode {
	case PAR_NONE:
		t.Cflag &^= syscall.PARENB
	case PAR_EVEN:
		t.Cflag |= syscall.PARENB
		t.Cflag &^= syscall.PARODD
	case PAR_ODD:
		t.Cflag |= syscall.PARENB
		t.Cflag |= syscall.PARODD
	default:
		return errors.New("invalid parity mode")
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setStopBits2(two bool) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if two {
		t.Cflag |= syscall.CSTOPB
	} else {
		t.Cflag &^= syscall.CSTOPB
	}

	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setHwFlowCtrl(hw bool) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if hw {
		t.Cflag |= crtscts
	} else {
		t.Cflag &^= crtscts
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setSwFlowCtrl(sw bool) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if sw {
		t.Iflag |= (syscall.IXON | syscall.IXOFF | syscall.IXANY)
	} else {
		t.Iflag &^= (syscall.IXON | syscall.IXOFF | syscall.IXANY)
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setLocal(local bool) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if local {
		t.Cflag |= syscall.CLOCAL
	} else {
		t.Cflag &^= syscall.CLOCAL
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setReadTimeout(vmin int, vtime time.Duration) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	t.Cc[syscall.VMIN] = uint8(vmin)
	t.Cc[syscall.VTIME] = uint8(vtime / (time.Second / 10))
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setHup(hup bool) error {
	var t Termios
	if err := s.tcGetAttr(&t); err != nil {
		return err
	}
	if hup {
		t.Cflag |= syscall.HUPCL
	} else {
		t.Cflag &^= syscall.HUPCL
	}
	if err := s.tcSetAttr(&t); err != nil {
		return err
	}
	return nil
}

func (s *Serial) setDTR(level bool) error {
	var cmd uintptr
	if level {
		cmd = syscall.TIOCMBIS
	} else {
		cmd = syscall.TIOCMBIC
	}
	v := syscall.TIOCM_DTR
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		cmd,
		uintptr(unsafe.Pointer(&v)),
	)
	if e != 0 {
		return os.NewSyscallError("setDTR", e)
	}
	return nil
}

func (s *Serial) setRTS(level bool) error {
	var cmd uintptr
	if level {
		cmd = syscall.TIOCMBIS
	} else {
		cmd = syscall.TIOCMBIC
	}
	v := syscall.TIOCM_RTS
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		cmd,
		uintptr(unsafe.Pointer(&v)),
	)
	if e != 0 {
		return os.NewSyscallError("setRTS", e)
	}
	return nil
}

func (s *Serial) inpWaiting() (int, error) {
	var v int
	cmd := uintptr(syscall.TIOCINQ)
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		cmd,
		uintptr(unsafe.Pointer(&v)),
	)
	if e != 0 {
		return 0, os.NewSyscallError("inpWaiting", e)
	}
	return v, nil
}

func (s *Serial) outWaiting() (int, error) {
	var v int
	cmd := uintptr(syscall.TIOCOUTQ)
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		cmd,
		uintptr(unsafe.Pointer(&v)),
	)
	if e != 0 {
		return 0, os.NewSyscallError("outWaiting", e)
	}
	return v, nil
}

func (s *Serial) flush(mode int) error {
	var v int
	cmd := uintptr(tcflsh)
	switch mode {
	case FLUSH_I:
		v = syscall.TCIFLUSH
	case FLUSH_O:
		v = syscall.TCOFLUSH
	case FLUSH_IO:
		v = syscall.TCIOFLUSH
	default:
		return errors.New("invalid flush mode")
	}
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(s.f.Fd()),
		cmd,
		uintptr(v),
	)
	if e != 0 {
		return os.NewSyscallError("flush", e)
	}
	return nil
}
