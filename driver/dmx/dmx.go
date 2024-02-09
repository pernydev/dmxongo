package dmx

import (
	"driver/serial"
	"fmt"
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
	StartByte     []byte
	EndByte       []byte
}

func NewDMX(numOfChannels int, serialNumber string) (*DMX, error) {
	dmx := &DMX{}
	dmx.Data = make([]byte, numOfChannels)
	dmx.NumOfChannels = numOfChannels
	d, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyUSB0", Baud: 250000})
	if err != nil {
		fmt.Println("Error opening serial port:", err)
		return nil, err
	}

	fmt.Println(d)

	dmx.Serial = d
	dmx.StartByte = []byte{0x7E, 0x06, 0x01, 0x02, 0x00}
	dmx.EndByte = []byte{0xE7}
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
