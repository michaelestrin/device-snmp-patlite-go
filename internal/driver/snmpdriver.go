// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 Dell Technologies
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	ds_models "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
	"github.com/edgexfoundry/edgex-go/pkg/models"
	"sync"
	"time"
)

type SNMPDriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *ds_models.AsyncValues
	switchButton bool
}

var client *SNMPClient
//Used to avoid get/set at the same time. If this happens simultaneously, state
//of the device can get out of sync with command actuation result
var mu sync.Mutex

// DisconnectDevice handles protocol-specific cleanup when a device
// is removed.
func (s *SNMPDriver) DisconnectDevice(address *models.Addressable) error {
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *SNMPDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *ds_models.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *SNMPDriver) HandleReadCommands(addr *models.Addressable, reqs []ds_models.CommandRequest) (res []*ds_models.CommandValue, err error) {

	var commands []DeviceCommand
	for _, req := range reqs {
		s.lc.Debug(fmt.Sprintf("SNMPDriver.HandleReadCommand: device: %s operation: %v attributes: %v", addr.Name, req.RO.Operation, req.DeviceObject.Attributes))
		oid := req.DeviceObject.Attributes["oid"].(string)
		commands = append(commands, NewGetDeviceCommand(oid))
	}

	port := uint16(addr.Port)
	if addr.Port == 0 {
		port = DEFAULT_PORT
	}
	mu.Lock()
	if client == nil {
		client = NewSNMPClient(addr.Address, port)
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	vals, err2 := client.GetValues(commands)
	mu.Unlock()

	if err2 != nil {
		s.lc.Error(fmt.Sprintf("SNMPDriver.HandleReadCommands; %s", err2))
		return
	}

	for i, val := range vals {
		s.lc.Debug(fmt.Sprintf("SNMPDriver.HandleReadCommands; value of %v is: %d", reqs[i].RO, val))
		cv, _ := ds_models.NewInt32Value(&reqs[i].RO, now, int32(val))
		res = append(res, cv)
	}

	return
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource (aka DeviceObject).
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *SNMPDriver) HandleWriteCommands(addr *models.Addressable, reqs []ds_models.CommandRequest,
	params []*ds_models.CommandValue) error {

	var commands []DeviceCommand
	for i, req := range reqs {
		s.lc.Debug(fmt.Sprintf("SNMPDriver.HandleWriteCommands: device: %s operation: %v attributes: %v", addr.Name, req.RO.Operation, req.DeviceObject.Attributes))
		oid := req.DeviceObject.Attributes["oid"].(string)
		val, err := (params[i]).Int32Value()
		if err != nil {
			return err
		}
		commands = append(commands, NewSetDeviceCommand(oid, int(val)))
	}

	port := uint16(addr.Port)
	if addr.Port == 0 {
		port = DEFAULT_PORT
	}

	mu.Lock()
	if client == nil {
		client = NewSNMPClient(addr.Address, port)
	}
	_, err2 := client.SetValues(commands)
	mu.Unlock()

	if err2 != nil {
		s.lc.Error(fmt.Sprintf("SNMPDriver.HandleWriteCommands; %s", err2))
		return err2
	}
	return nil

}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *SNMPDriver) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("SNMPDriver.Stop called: force=%v", force))
	client.Disconnect()
	return nil
}
