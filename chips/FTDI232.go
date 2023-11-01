package chips

import (
	"dmxongo/objects"
	"fmt"

	"github.com/albenik/go-serial/v2"

	"sync"
)

var conn DMX
var mutex = &sync.Mutex{}

type Device struct {
	VID          int
	PID          int
	SerialNumber string
	Name         string
}

// DMX represents the DMX controller.
type DMX struct {
	Data          []byte
	NumOfChannels int
	Ser           *serial.Port
	Device        *Device
	StartByte     []byte
	EndByte       []byte
}

func NewDMX(numOfChannels int, serialNumber string) (*DMX, error) {
	dmx := &DMX{
		Data:          make([]byte, numOfChannels),
		NumOfChannels: numOfChannels,
	}

	var err error

	dmx.StartByte = []byte{0x7E, 0x06, 0x01, 0x02, 0x00}
	dmx.EndByte = []byte{0xE7}
	dmx.Ser, err = nil, nil

	port, err := serial.Open("/dev/ttyUSB0",
		serial.WithBaudrate(250000),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(1000),
		serial.WithWriteTimeout(1000),
		serial.WithHUPCL(false),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("Opened serial port")

	dmx.Ser = port

	return dmx, nil
}

func (d *DMX) Send() {
	fmt.Println("Sending DMX")
	data := append(append(d.StartByte, d.Data...), d.EndByte...)
	_, err := fmt.Fprint(d.Ser, data)
	if err != nil {
		panic(err)
	}
}

func Init() {
	connection, err := NewDMX(512, "")
	if err != nil {
		panic(err)
	}

	conn = *connection
}

func SendUniverse(u objects.Universe) {
	mutex.Lock()
	defer mutex.Unlock()

	conn.Data = make([]byte, len(u.ChannelValues))
	for i, v := range u.ChannelValues {
		conn.Data[i] = byte(v)
	}
	conn.Send()
}
