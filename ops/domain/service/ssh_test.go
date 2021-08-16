package service

import (
	"fmt"
	"log"
	"testing"
)

func TestConnection_SendCommands(t *testing.T) {

	// hostInfoSh := `cat /proc/cpuinfo | grep 'cpu cores'; cat /proc/meminfo; fdisk -l |grep Disk`
	processInfoSh := `ps -ef | awk '{print $8}' | sed '1d'`

	conn, err := Connect("47.241.12.170:22", "root", "YCJ19910421x")
	if err != nil {
		log.Fatal(err)
	}

	// output, err := conn.SendCommands(hostInfoSh)
	output, err := conn.SendCommands(processInfoSh)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}
