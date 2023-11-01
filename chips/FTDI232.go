package chips

import (
	"dmxongo/objects"
	"fmt"
	"log"

	"github.com/tarm/serial"
	secondSerial "go.bug.st/serial"
)

var conn DMX
var previousUniverse objects.Universe

type Device struct {
	VID          int
	PID          int
	SerialNumber string
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

var DUMMY = Device{VID: 0, PID: 0, SerialNumber: ""}

var EUROLITEUSBDMX512PROCABLEINTERFACE = Device{VID: 1027, PID: 24577, SerialNumber: ""}

var DEVICE_LIST = []Device{DUMMY, EUROLITEUSBDMX512PROCABLEINTERFACE}

func NewDMX(numOfChannels int, serialNumber string) (*DMX, error) {
	dmx := &DMX{
		Data:          make([]byte, numOfChannels),
		NumOfChannels: numOfChannels,
	}

	var err error
	dmx.Device, err = findDevice(serialNumber)
	if err != nil {
		return nil, err
	}

	if dmx.Device.VID == EUROLITEUSBDMX512PROCABLEINTERFACE.VID &&
		dmx.Device.PID == EUROLITEUSBDMX512PROCABLEINTERFACE.PID {
		dmx.StartByte = []byte{0x7E, 0x06, 0x01, 0x02, 0x00}
		dmx.EndByte = []byte{0xE7}
		dmx.Ser, err = setupSerialPort(dmx.Device.SerialNumber)
		if err != nil {
			return nil, err
		}
	} else {
		dmx.StartByte = []byte{0x00}
		dmx.EndByte = []byte{}
		dmx.Ser, err = setupSerialPort(dmx.Device.SerialNumber)
		if err != nil {
			return nil, err
		}
	}

	return dmx, nil
}

func (d *DMX) SetData(channelID int, data int, autoSend bool) {
	if channelID < 1 || channelID > 512 {
		log.Fatal("Channel ID must be between 1 and 512.")
	}
	if data < 0 || data > 255 {
		log.Fatal("Data must be between 0 and 255.")
	}
	if channelID > d.NumOfChannels {
		log.Fatal("Channel ID was not reserved. Please set the number of channels first.")
	}

	d.Data[channelID-1] = byte(data)

	if autoSend {
		d.Send()
	}
}

func (d *DMX) Send() {
	data := append(append(d.StartByte, d.Data...), d.EndByte...)
	_, err := d.Ser.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func findDevice(serialNumber string) (*Device, error) {
	ports, err := secondSerial.GetPortsList()
	if err != nil {
		return nil, err
	}

	for _, port := range ports {
		for _, device := range DEVICE_LIST {
			for _, knownDevice := range DEVICE_LIST {
				if device.VID == knownDevice.VID &&
					device.PID == knownDevice.PID &&
					device.SerialNumber == "" {
					device.SerialNumber = port
				}
			}
		}
	}

	return nil, fmt.Errorf("Could not find the RS-485 interface.")
}

func setupSerialPort(serialNumber string) (*serial.Port, error) {
	c := &serial.Config{Name: serialNumber, Baud: 250000}
	ser, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func Init() {
	connection, err := NewDMX(512, "")
	if err != nil {
		fmt.Println(err)
	}

	conn = *connection
}

func SendUniverse(u objects.Universe) {
	for channel, value := range u.ChannelValues {
		if previousUniverse.ChannelValues != nil && value != previousUniverse.ChannelValues[channel] {
			conn.SetData(channel, value, false)
		}
	}
	previousUniverse = u
	conn.Send()
}
