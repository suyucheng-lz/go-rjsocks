package rjsocks

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"os/exec"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ListNetworkDev() ([]string, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, ifc := range interfaces {
		ret = append(ret, ifc.Description)
	}
	return ret, nil
}

func ListNetworkAdapter() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, ifc := range interfaces {
		ret = append(ret, ifc.Name)
	}
	return ret, nil
}

func SelectNetworkDev(dev string) (*pcap.Interface, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	for _, ifc := range interfaces {
		if ifc.Description == dev {
			return &ifc, nil
		}
	}
	return nil, errors.New("未发现硬件设备" + dev)
}

func SelectNetworkAdapter(adapter string) (net.HardwareAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifc := range interfaces {
		if ifc.Name == adapter {
			return ifc.HardwareAddr, nil
		}
	}
	return nil, errors.New("无法获取对应网卡")
}

func Symmetric(data []byte) {
	for i := 0; i < 4; i++ {
		data[i] = byteReverse(data[i])
	}
}

func byteReverse(base byte) byte {
	result := byte(0)
	for i := 0; i < 8; i++ {
		result <<= 1
		result |= base & 0x01
		base >>= 1
	}
	return ^result
}

var table = []byte{
	0x00, 0x00, 0x21, 0x10, 0x42, 0x20, 0x63, 0x30, 0x84, 0x40, 0xA5, 0x50, 0xC6, 0x60, 0xE7, 0x70,
	0x08, 0x81, 0x29, 0x91, 0x4A, 0xA1, 0x6B, 0xB1, 0x8C, 0xC1, 0xAD, 0xD1, 0xCE, 0xE1, 0xEF, 0xF1,
	0x31, 0x12, 0x10, 0x02, 0x73, 0x32, 0x52, 0x22, 0xB5, 0x52, 0x94, 0x42, 0xF7, 0x72, 0xD6, 0x62,
	0x39, 0x93, 0x18, 0x83, 0x7B, 0xB3, 0x5A, 0xA3, 0xBD, 0xD3, 0x9C, 0xC3, 0xFF, 0xF3, 0xDE, 0xE3,
	0x62, 0x24, 0x43, 0x34, 0x20, 0x04, 0x01, 0x14, 0xE6, 0x64, 0xC7, 0x74, 0xA4, 0x44, 0x85, 0x54,
	0x6A, 0xA5, 0x4B, 0xB5, 0x28, 0x85, 0x09, 0x95, 0xEE, 0xE5, 0xCF, 0xF5, 0xAC, 0xC5, 0x8D, 0xD5,
	0x53, 0x36, 0x72, 0x26, 0x11, 0x16, 0x30, 0x06, 0xD7, 0x76, 0xF6, 0x66, 0x95, 0x56, 0xB4, 0x46,
	0x5B, 0xB7, 0x7A, 0xA7, 0x19, 0x97, 0x38, 0x87, 0xDF, 0xF7, 0xFE, 0xE7, 0x9D, 0xD7, 0xBC, 0xC7,
	0xC4, 0x48, 0xE5, 0x58, 0x86, 0x68, 0xA7, 0x78, 0x40, 0x08, 0x61, 0x18, 0x02, 0x28, 0x23, 0x38,
	0xCC, 0xC9, 0xED, 0xD9, 0x8E, 0xE9, 0xAF, 0xF9, 0x48, 0x89, 0x69, 0x99, 0x0A, 0xA9, 0x2B, 0xB9,
	0xF5, 0x5A, 0xD4, 0x4A, 0xB7, 0x7A, 0x96, 0x6A, 0x71, 0x1A, 0x50, 0x0A, 0x33, 0x3A, 0x12, 0x2A,
	0xFD, 0xDB, 0xDC, 0xCB, 0xBF, 0xFB, 0x9E, 0xEB, 0x79, 0x9B, 0x58, 0x8B, 0x3B, 0xBB, 0x1A, 0xAB,
	0xA6, 0x6C, 0x87, 0x7C, 0xE4, 0x4C, 0xC5, 0x5C, 0x22, 0x2C, 0x03, 0x3C, 0x60, 0x0C, 0x41, 0x1C,
	0xAE, 0xED, 0x8F, 0xFD, 0xEC, 0xCD, 0xCD, 0xDD, 0x2A, 0xAD, 0x0B, 0xBD, 0x68, 0x8D, 0x49, 0x9D,
	0x97, 0x7E, 0xB6, 0x6E, 0xD5, 0x5E, 0xF4, 0x4E, 0x13, 0x3E, 0x32, 0x2E, 0x51, 0x1E, 0x70, 0x0E,
	0x9F, 0xFF, 0xBE, 0xEF, 0xDD, 0xDF, 0xFC, 0xCF, 0x1B, 0xBF, 0x3A, 0xAF, 0x59, 0x9F, 0x78, 0x8F,
	0x88, 0x91, 0xA9, 0x81, 0xCA, 0xB1, 0xEB, 0xA1, 0x0C, 0xD1, 0x2D, 0xC1, 0x4E, 0xF1, 0x6F, 0xE1,
	0x80, 0x10, 0xA1, 0x00, 0xC2, 0x30, 0xE3, 0x20, 0x04, 0x50, 0x25, 0x40, 0x46, 0x70, 0x67, 0x60,
	0xB9, 0x83, 0x98, 0x93, 0xFB, 0xA3, 0xDA, 0xB3, 0x3D, 0xC3, 0x1C, 0xD3, 0x7F, 0xE3, 0x5E, 0xF3,
	0xB1, 0x02, 0x90, 0x12, 0xF3, 0x22, 0xD2, 0x32, 0x35, 0x42, 0x14, 0x52, 0x77, 0x62, 0x56, 0x72,
	0xEA, 0xB5, 0xCB, 0xA5, 0xA8, 0x95, 0x89, 0x85, 0x6E, 0xF5, 0x4F, 0xE5, 0x2C, 0xD5, 0x0D, 0xC5,
	0xE2, 0x34, 0xC3, 0x24, 0xA0, 0x14, 0x81, 0x04, 0x66, 0x74, 0x47, 0x64, 0x24, 0x54, 0x05, 0x44,
	0xDB, 0xA7, 0xFA, 0xB7, 0x99, 0x87, 0xB8, 0x97, 0x5F, 0xE7, 0x7E, 0xF7, 0x1D, 0xC7, 0x3C, 0xD7,
	0xD3, 0x26, 0xF2, 0x36, 0x91, 0x06, 0xB0, 0x16, 0x57, 0x66, 0x76, 0x76, 0x15, 0x46, 0x34, 0x56,
	0x4C, 0xD9, 0x6D, 0xC9, 0x0E, 0xF9, 0x2F, 0xE9, 0xC8, 0x99, 0xE9, 0x89, 0x8A, 0xB9, 0xAB, 0xA9,
	0x44, 0x58, 0x65, 0x48, 0x06, 0x78, 0x27, 0x68, 0xC0, 0x18, 0xE1, 0x08, 0x82, 0x38, 0xA3, 0x28,
	0x7D, 0xCB, 0x5C, 0xDB, 0x3F, 0xEB, 0x1E, 0xFB, 0xF9, 0x8B, 0xD8, 0x9B, 0xBB, 0xAB, 0x9A, 0xBB,
	0x75, 0x4A, 0x54, 0x5A, 0x37, 0x6A, 0x16, 0x7A, 0xF1, 0x0A, 0xD0, 0x1A, 0xB3, 0x2A, 0x92, 0x3A,
	0x2E, 0xFD, 0x0F, 0xED, 0x6C, 0xDD, 0x4D, 0xCD, 0xAA, 0xBD, 0x8B, 0xAD, 0xE8, 0x9D, 0xC9, 0x8D,
	0x26, 0x7C, 0x07, 0x6C, 0x64, 0x5C, 0x45, 0x4C, 0xA2, 0x3C, 0x83, 0x2C, 0xE0, 0x1C, 0xC1, 0x0C,
	0x1F, 0xEF, 0x3E, 0xFF, 0x5D, 0xCF, 0x7C, 0xDF, 0x9B, 0xAF, 0xBA, 0xBF, 0xD9, 0x8F, 0xF8, 0x9F,
	0x17, 0x6E, 0x36, 0x7E, 0x55, 0x4E, 0x74, 0x5E, 0x93, 0x2E, 0xB2, 0x3E, 0xD1, 0x0E, 0xF0, 0x1E,
}

func checkSum(buf []byte) {
	var b0, b1 byte
	var index int32
	for i := 0; i < 0x15; i++ {
		index = int32(b0 ^ buf[i])
		b0 = b1 ^ table[index*2+1]
		b1 = table[index*2]
	}
	buf[0x15], buf[0x16] = b0, b1
	for i := 0; i < 0x17; i++ {
		buf[i] = byteReverse(buf[i])
	}
}

func reNewIP(adapter string) {
	cmd := exec.Command("ipconfig", "/renew", adapter)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	go cmd.Run()
}

var fillbuf = []byte{
	// dhcp layer
	0xff, 0xff, 0x37, 0x77, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xfd, 0x36,
	// padding layer
	0x00, 0x00, 0x13, 0x11, 0x38, 0x30, 0x32, 0x31, 0x78, 0x2e, 0x65, 0x78, 0x65, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x04, 0x0a, 0x00, 0x02, 0x00, 0x00, 0x00, 0x13, 0x11, 0x01, 0x8c, 0x1a,
	0x28, 0x00, 0x00, 0x13, 0x11, 0x17, 0x22, 0x36, 0x38, 0x44, 0x43, 0x31, 0x32, 0x33, 0x42, 0x37,
	0x45, 0x42, 0x32, 0x33, 0x39, 0x46, 0x32, 0x33, 0x41, 0x38, 0x43, 0x30, 0x30, 0x30, 0x33, 0x38,
	0x38, 0x34, 0x39, 0x38, 0x36, 0x33, 0x39, 0x1a, 0x0c, 0x00, 0x00, 0x13, 0x11, 0x18, 0x06, 0x00,
	0x00, 0x00, 0x00, 0x1a, 0x0e, 0x00, 0x00, 0x13, 0x11, 0x2d, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x1a, 0x08, 0x00, 0x00, 0x13, 0x11, 0x2f, 0x02, 0x1a, 0x09, 0x00, 0x00, 0x13, 0x11, 0x35,
	0x03, 0x01, 0x1a, 0x18, 0x00, 0x00, 0x13, 0x11, 0x36, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1a, 0x18, 0x00, 0x00, 0x13, 0x11,
	0x38, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, 0x86,
	0x13, 0x4c, 0x1a, 0x88, 0x00, 0x00, 0x13, 0x11, 0x4d, 0x82, 0x36, 0x38, 0x64, 0x63, 0x31, 0x32,
	0x33, 0x62, 0x30, 0x37, 0x65, 0x62, 0x32, 0x33, 0x39, 0x66, 0x32, 0x33, 0x61, 0x38, 0x30, 0x64,
	0x63, 0x66, 0x32, 0x35, 0x38, 0x37, 0x35, 0x64, 0x30, 0x35, 0x37, 0x37, 0x30, 0x63, 0x37, 0x32,
	0x31, 0x65, 0x34, 0x35, 0x36, 0x34, 0x35, 0x65, 0x35, 0x33, 0x37, 0x61, 0x62, 0x33, 0x35, 0x31,
	0x62, 0x62, 0x36, 0x33, 0x31, 0x35, 0x35, 0x61, 0x65, 0x31, 0x36, 0x32, 0x36, 0x31, 0x36, 0x37,
	0x65, 0x62, 0x30, 0x39, 0x32, 0x32, 0x33, 0x65, 0x32, 0x61, 0x30, 0x61, 0x37, 0x38, 0x30, 0x33,
	0x31, 0x31, 0x36, 0x31, 0x61, 0x63, 0x30, 0x39, 0x64, 0x61, 0x32, 0x64, 0x63, 0x30, 0x37, 0x33,
	0x36, 0x39, 0x33, 0x61, 0x34, 0x66, 0x35, 0x61, 0x32, 0x39, 0x32, 0x38, 0x36, 0x37, 0x35, 0x31,
	0x66, 0x39, 0x37, 0x66, 0x34, 0x64, 0x30, 0x34, 0x36, 0x38, 0x1a, 0x28, 0x00, 0x00, 0x13, 0x11,
	0x39, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x1a, 0x48, 0x00, 0x00, 0x13, 0x11, 0x54, 0x42, 0x48, 0x55, 0x53, 0x54, 0x4d, 0x4f,
	0x4f, 0x4e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1a, 0x08, 0x00, 0x00, 0x13, 0x11,
	0x55, 0x02, 0x1a, 0x09, 0x00, 0x00, 0x13, 0x11, 0x62, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var fillLayer = gopacket.Payload(fillbuf)
/*
type RawLayer struct {
	RawBytes []byte
}

func (r *RawLayer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.AppendBytes(len(r.RawBytes))
	if err != nil {
		return err
	}
	copy(bytes, r.RawBytes)
	return nil
}
*/