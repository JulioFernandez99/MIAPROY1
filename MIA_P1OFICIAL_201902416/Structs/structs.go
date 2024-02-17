package Structs

type MBR struct {
	Mbr_tamanio int64
	Mbr_fecha_creacion int64
	Mbr_disk_signature int64
	Mdsk_fit [1]byte	
	Mbr_partition [4]Partition
}

type Partition struct {
	Part_status byte //indica si la particion esta montada o no
	Part_tyPe byte //indica el tipo de particion
	Part_fit byte //indica el ajuste de la particion
	Part_start int64 //indica el byte de inicio de la particion
	Part_s int64 //indica el tamanio de la particion
	Part_name [16]byte	//nombre de la particion
	Part_correlative int64	//numero de particion
	Part_id [4]byte	//identificador unico de la particion
}


