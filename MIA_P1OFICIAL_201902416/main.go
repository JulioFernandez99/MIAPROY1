package main

import (
	"MIA_P1OFICIAL_201902416/Analizador"
	
)
//mkdisk -size=12 -unit=K -fit="BF" A
//mkdisk -size=6 -unit=M B
//mkdisk -size=412 C
//mkdisk -size=12 -unit=K -fit="FF" D
//rmdisk -driveletter=A



//fdisk -size=100 -type=L -unit=B -fit=bf -driveletter=A -name="Particion4"
//fdisk -size=300 -driveletter=A -name=Particion1
//iterar sobre las particiones
//fmt.Println("Particiones",mbrPrueba.Mbr_partition)
//fdisk -name=Particion5 -delete=full -driveletter=B
//fmt.Println("Size particion",mbrPrueba.Mbr_partition[0].Part_s)
	
//mbrPrueba.Mbr_partition[0].Part_s = mbrPrueba.Mbr_partition[0].Part_s + addValue
//fmt.Println("ADDFDISK",addValue)
//fmt.Println("Size particion",mbrPrueba.Mbr_partition[0].Part_s)

//fdisk -size=300 -driveletter=A -name=Particion1
//fdisk -name=Particion1 -delete=full -driveletter=A
//fdisk -add=-500 -size=10 -unit=K -driveletter=A -name=”Particion5”
//fdisk -add=1 -unit=M -driveletter=B -name="Particion2"

//mount -driveletter=A -name=Particion2
//fpartitions -driveletter=B   para saber si las particiones estan libres o no
//unmount -id=A119



//Creacion de particiones
// fdisk -size=1 -unit=K -driveletter=A -name=Particion5
// fdisk -type=E -driveletter=A -unit=B -name=Particion4 -size=100
// fdisk -size=20 -type=L -unit=M -fit=bf -driveletter=A -name="ParticionL1"
// fdisk -type=P -driveletter=A -name=Part2 -unit=K -size=200

// Eliminacion de particiones
// fdisk -delete=full -name="Part3" -driveletter=A
// fdisk -name=part4 -delete=full -driveletter=A

// Modificacion de particiones
// fdisk -add=200100 -size=10 -unit=K -driveletter=A -name=particion2
// fdisk -add=1 -unit=M -driveletter=A -name="Particion4"

func main() {
	Analizador.Run()	
	//Logica.PrintMBR("DISCOS/A.disk")
	//Logica.PrintMBR("DISCOS/B.disk")
	//Logica.PrintMBR("DISCOS/C.disk")
}
