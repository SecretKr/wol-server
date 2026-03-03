package main

import (
	"net"

	"github.com/mdlayher/wol"
)

func WakeDevice(macAddr string) error {
	target, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Wake("10.0.1.255:9", target)
	if err != nil {
		return err
	}

	return nil
}