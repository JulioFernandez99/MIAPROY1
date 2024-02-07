package main

import (
	"fmt"
	"myProject1/logica"
)


func main() {

	disco:=logica.Disco{}
	
	disco.SetFileName("./disks/disk1.disk")
	disco.CreateDskFIle()
	fmt.Println(disco.GetCreate())
}



