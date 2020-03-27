package main

import (
	"fmt"
	"github.com/sony/sonyflake"
	"os"
	"time"
)

//sonyflake

//MachineID 可以由用户自定义的函数，如果用户不定义的话，会默认将本机 ip 的低 16 位作为 machine id。
func getMachineID() (uint16, error) {
	var machineID uint16
	var err error
	machineID = readMachineIDFromLocalFile() //获取机器id
	if machineID == 0 {
		machineID, err = generateMachineID()
		if err != nil {
			return 0, err
		}
	}
	return machineID, nil
}

//是由用户提供的检查 MachineID 是否冲突的函数
func checkMachineID(machineID uint16) bool {
	saddResult, err := saddMachineIDToRedisSet() //用 Redis 的 set 来检查冲突。
	if err != nil || saddResult == 0 {
		return true
	}
	err := saveMachineIDToLocalFile(machineID)
	if err != nil {
		return true
	}
	return false
}
func main() {
	t, _ := time.Parse("2006-01-02", "2020-01-01")
	settings := sonyflake.Settings{
		StartTime:      t,
		MachineID:      getMachineID,
		CheckMachineID: checkMachineID,
	}
	sf := sonyflake.NewSonyflake(settings)
	id, err := sf.NextID()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(id)
}
