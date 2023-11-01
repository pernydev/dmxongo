package dmx

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/tarm/serial"
)

type Device struct {
	Vid          int
	Pid          int
	SerialNumber string
}

// DMX represents a DMX controller
type DMX struct {
	Data          []byte
	NumOfChannels int
	Serial        *serial.Port
	Device        *Device
	StartByte     []byte
	EndByte       []byte
}

var Dummy = Device{0, 0, ""}
var EUROLITEUSBDMX512PROCABLEINTERFACE = Device{1027, 24577, ""}

var DeviceList = []Device{Dummy, EUROLITEUSBDMX512PROCABLEINTERFACE}

func NewDMX(numOfChannels int, serialNumber string) (*DMX, error) {
	dmx := &DMX{}
	dmx.Data = make([]byte, numOfChannels)
	dmx.NumOfChannels = numOfChannels

	for _, device := range DeviceList {
		ports, err := serial.GetPortsList()
		if err != nil {
			return nil, err
		}
		for _, port := range ports {
			d, err := serial.OpenPort(&serial.Config{Name: port, Baud: 250000})
			if err != nil {
				continue
			}
			if reflect.DeepEqual(dmx.Device, device) && (serialNumber == "" || serialNumber == device.SerialNumber) {
				dmx.Serial = d
				dmx.Device = &device
				break
			}
			d.Close()
		}
		if dmx.Device != nil {
			break
		}
	}
	if dmx.Device == nil {
		return nil, fmt.Errorf("Could not find the RS-485 interface.")
	}
	if dmx.Device.Vid == EUROLITEUSBDMX512PROCABLEINTERFACE.Vid && dmx.Device.Pid == EUROLITEUSBDMX512PROCABLEINTERFACE.Pid {
		dmx.StartByte = []byte{0x7E, 0x06, 0x01, 0x02, 0x00}
		dmx.EndByte = []byte{0xE7}
	} else {
		dmx.StartByte = []byte{0x00}
		dmx.EndByte = []byte{}
	}
	return dmx, nil
}

func (dmx *DMX) SetData(channelID int, data byte) {
	if channelID >= 1 && channelID <= 512 {
		dmx.Data[channelID-1] = data
	}
}

func (dmx *DMX) Send() {
	data := append(append(dmx.StartByte, dmx.Data...), dmx.EndByte...)
	dmx.Serial.Write(data)
	dmx.Serial.Flush()
}

func (dmx *DMX) Close() {
	dmx.NumOfChannels = 512
	dmx.Data = make([]byte, dmx.NumOfChannels)
	dmx.Send()
	dmx.Serial.Close()
}

func main() {
	dmx, err := NewDMX(512, "")
	if err != nil {
		log.Fatal(err)
	}
	defer dmx.Close()

	for {
		for i := 0; i < 255; i += 5 {
			dmx.SetData(1, byte(i))
			dmx.SetData(2, byte(i))
			dmx.SetData(3, byte(i))
			dmx.SetData(4, byte(i))
			dmx.Send()
			time.Sleep(10 * time.Millisecond)
		}

		for i := 255; i > 0; i -= 5 {
			dmx.SetData(1, byte(i))
			dmx.SetData(2, byte(i))
			dmx.SetData(3, byte(i))
			dmx.SetData(4, byte(i))
			dmx.Send()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
