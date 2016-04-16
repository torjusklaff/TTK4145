package main

import (
	"./com"
	"./driver"
	"./elevator"
	"./logger"
	"./master"
	"./network"
	"./slave"
	"flag"
)

func main() {
	var startAsMaster, selfAsBackup bool
	flag.BoolVar(&startAsMaster, "master", false, "Start as master")
	flag.BoolVar(&selfAsBackup, "selfAsBackup", false, "Can use self as backup")
	flag.Parse()

	var elevatorEvents com.ElevatorEvent
	elevatorEvents.FloorReached = make(chan int)
	elevatorEvents.NewTargetFloor = make(chan int)

	var slaveEvents com.SlaveEvent
	slaveEvents.CompletedFloor = make(chan int)
	slaveEvents.MissedDeadline = make(chan bool)
	slaveEvents.ButtonPressed = make(chan driver.OrderButton)
	slaveEvents.FromMaster = make(chan network.UDPMessage)
	slaveEvents.ToMaster = make(chan network.UDPMessage)

	var masterEvents com.MasterEvent
	masterEvents.ToSlaves = make(chan network.UDPMessage)
	masterEvents.FromSlaves = make(chan network.UDPMessage)

	driver.ElevInit()

	go driver.EventListener(
		slaveEvents.ButtonPressed,
		elevatorEvents.FloorReached)

	elevLogger := logger.NewLogger("ELEV")
	go elevator.Init(
		slaveEvents.CompletedFloor,
		slaveEvents.MissedDeadline,
		elevatorEvents.FloorReached,
		elevatorEvents.NewTargetFloor,
		elevLogger)

	networkLogger := logger.NewLogger("NETWORK")

	if startAsMaster {
		go network.UDPInit(true, masterEvents.ToSlaves, masterEvents.FromSlaves, networkLogger)
		masterLogger := logger.NewLogger("MASTER")
		go master.InitMaster(masterEvents, nil, nil, masterLogger, selfAsBackup)
	}
	go network.UDPInit(false, slaveEvents.ToMaster, slaveEvents.FromMaster, networkLogger)
	slaveLogger := logger.NewLogger("SLAVE")
	slave.InitSlave(slaveEvents, masterEvents, elevatorEvents, slaveLogger)
}
