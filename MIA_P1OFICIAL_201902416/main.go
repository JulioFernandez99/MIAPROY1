package main

import (
	"MIA_P1OFICIAL_201902416/Analizador"
	
)
//mkdisk -size=12 -unit=K -fit="BF" A
//mkdisk -size=6 -unit=M B
//mkdisk -size=412 C
//mkdisk -size=12 -unit=K -fit="FF" D
//rmdisk -driveletter=A
//fdisk -size=300 -driveletter=A -name=Particion4 
//iterar sobre las particiones
//fmt.Println("Particiones",mbrPrueba.Mbr_partition)
//fdisk -name=Particion3 -delete=full -driveletter=A
//fmt.Println("Size particion",mbrPrueba.Mbr_partition[0].Part_s)
	
//mbrPrueba.Mbr_partition[0].Part_s = mbrPrueba.Mbr_partition[0].Part_s + addValue
//fmt.Println("ADDFDISK",addValue)
//fmt.Println("Size particion",mbrPrueba.Mbr_partition[0].Part_s)
//fdisk -size=300 -driveletter=B -name=Particion2
//fdisk -name=Particion2 -delete=full -driveletter=A
//fdisk -add=-500 -size=10 -unit=K -driveletter=A -name=”Particion5”
//fdisk -add=1 -unit=M -driveletter=B -name="Particion2"
	


func main() {
	Analizador.Run()	
	//Logica.PrintMBR("DISCOS/A.disk")
	//Logica.PrintMBR("DISCOS/B.disk")
	//Logica.PrintMBR("DISCOS/C.disk")
}
