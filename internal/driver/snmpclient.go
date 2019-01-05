// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 Dell Technologies
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"errors"
	g "github.com/soniah/gosnmp"
)

// SNMPClient represents the SNMP device is used for getting and setting SNMP device via OID
type SNMPClient struct {
	ipAddr string
	ipPort uint16
}

func NewSNMPClient(addr string, port uint16) SNMPClient {
	return SNMPClient{
		ipAddr: addr,
		ipPort: port,
	}
}

type DeviceCommand struct {
	operation string
	value     int
}

func NewGetDeviceCommand(op string) DeviceCommand {
	return DeviceCommand{
		operation: op,
		// for Gets, value is not used
		value: 0,
	}
}

func NewSetDeviceCommand(op string, val int) DeviceCommand {
	return DeviceCommand{
		operation: op,
		value:     val,
	}
}

func (c *SNMPClient) GetValues(commands []DeviceCommand) ([]int, error) {

	var results []int
	var oids []string

	for _, command := range commands {
		if command.operation == "" {
			return results, errors.New("Unknown operation: " + command.operation)
		}
		oids = append(oids, command.operation)
	}
	g.Default.Target = c.ipAddr
	g.Default.Port = c.ipPort
	err := g.Default.Connect()
	if err != nil {
		return results, err
	}
	defer g.Default.Conn.Close()

	packets, err2 := g.Default.Get(oids)
	if err2 != nil {
		return results, err2
	}

	var temp uint64
	for _, variable := range packets.Variables {
		temp = g.ToBigInt(variable.Value).Uint64()
		results = append(results, int(temp))
	}
	return results, nil
}

func (c *SNMPClient) GetValue(command DeviceCommand) (int, error) {
	commands := []DeviceCommand{command}
	results, err := c.GetValues(commands)
	if err != nil {
		return 0, err
	}
	return results[0], nil
}

func (c *SNMPClient) SetValues(commands []DeviceCommand) ([]int, error) {

	var results []int
	var pdus []g.SnmpPDU

	for _, command := range commands {
		if command.operation == "" {
			return results, errors.New("Unknown operation: " + command.operation)
		}
		// TODO pass in logger
		pdu := g.SnmpPDU{Name: command.operation, Type: g.Integer, Value: command.value, Logger: nil}
		pdus = append(pdus, pdu)
	}
	g.Default.Target = c.ipAddr
	g.Default.Port = c.ipPort
	g.Default.Community = COMMUNITY_ACCESS
	err := g.Default.Connect()
	if err != nil {
		return results, err
	}
	defer g.Default.Conn.Close()

	packets, err2 := g.Default.Set(pdus)
	if err2 != nil {
		return results, err2
	}

	var temp uint64
	for _, variable := range packets.Variables {
		temp = g.ToBigInt(variable.Value).Uint64()
		results = append(results, int(temp))
	}
	return results, nil
}
