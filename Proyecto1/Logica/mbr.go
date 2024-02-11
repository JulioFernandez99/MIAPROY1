package Logica

import (
	"time"
)

type MBR struct {
	Mbr_tamanio int64
	Mbr_fecha_creacion [19]byte
	Mbr_disk_signature int64
	Mdsk_fit byte
	Mbr_partition [4]byte
}

func GetCurrentTime() string {
    // Obtener la hora actual
    currentTime := time.Now()

    // Formatear la hora actual como "Hora:Minuto:Segundo Fecha"
    formattedTime := currentTime.Format("02-01-2006 15:04:05")

    return formattedTime
}