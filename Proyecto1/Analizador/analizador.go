package Analizador

import (
	"MIA_P1_201902416/Logica"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"path/filepath"
	"strings"
)


//----------------------------------flags mkdisk----------------------------------
var (
	comandDisk = flag.String("comandDisk", "0", "comandDisk")
	sizeDisk = flag.Int("size", 0, "Tamaño")
	fitDisk = flag.String("fit", "f", "Ajuste")
	unitDisk = flag.String("unit", "m", "Unidad")
	reDisk = regexp.MustCompile(`(?:mkdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reRdisk = regexp.MustCompile(`(?:rmdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	delateDisk = flag.String("driveletter", "A", "driveletter")
)
//-----------------------------------fin glags mkdisk------------------------------------

func Run() {
	flag.Parse()
	for {
		input, err := ReadLine()
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}
		ProcessInput(input)
	}
}

func ReadLine() (string,error) {
	fmt.Print(">>")
	reader := bufio.NewReader(os.Stdin)
	text,err := reader.ReadString('\n')
	return text,err
}

func ProcessInput(input string) {
	matchMkDisk := reDisk.FindAllStringSubmatch(input, -1)
	comandDisk :=strings.Trim(matchMkDisk[0][0], " ")

	matchRmDisk := reRdisk.FindAllStringSubmatch(input, -1)
	comandRmDisk :=strings.Trim(matchRmDisk[0][0], " ")
	
	fmt.Println("matchRmDisk",matchRmDisk)
	
	if comandDisk == "mkdisk" && len(matchMkDisk) > 0{
		disco:=Logica.Disco{}
		processMkDisk(matchMkDisk[1:],disco)
		
	}else if comandRmDisk == "rmdisk" && len(matchRmDisk) > 0{
		processRmDisk(matchRmDisk[1:])
	}
	

}

func processRmDisk(match [][]string) {
	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]

		switch flagName {
			case "driveletter":
				path := path.Join("DISCOS", flagValue+".disk")
				fmt.Println("Path:", path)
				err := os.Remove(path)
				if err != nil {
					fmt.Println("Error:", err)
				}
			default:
				fmt.Println("Error: Flag not found",flagName,flagValue)
		}
	}
}

func processMkDisk(match [][]string, disco Logica.Disco) {
	// Obtener el próximo nombre de disco disponible
    nombreDisco, err := ObtenerProximoNombreDisco("DISCOS")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	fmt.Println("Nombre del disco:", nombreDisco)
	disco.FileName = "DISCOS/"+nombreDisco
	
	if len(match) == 1 && match[0][1] == "size"{
		//Si sOlo tiene un parametro
		value ,_:= strconv.Atoi(match[0][2])
		*sizeDisk = value*1000000
	}else if len(match) == 2 && match[1][1] == "fit" {
		//Si tiene dos parametros y el segundo es fit
		value ,_:= strconv.Atoi(match[0][2])
		*sizeDisk = value*1000000	
	}else {
		//Si tiene mas de dos parametros
		for _, match := range match {	
			flagName := match[1]
			flagValue := match[2]

			
			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
				case "size":
					sizeValue := 0
					fmt.Sscanf(flagValue, "%d", &sizeValue)
					*sizeDisk = sizeValue
				case "fit":
					flagValue = flagValue[:1]
					flagValue = strings.ToLower(flagValue)
					*fitDisk = flagValue
				case "unit":
					flagValue = strings.ToLower(flagValue)
					if flagValue =="m"{
						fmt.Println("Mega")
						*sizeDisk = *sizeDisk*(1000000)
					}
					if flagValue =="k"{
						*sizeDisk = *sizeDisk*(1000)
						
					}
					*unitDisk = flagValue
				default:
					fmt.Println("Error: Flag not found",flagName,flagValue)
			}
		}
	}
	disco.SetSize(make([]byte, *sizeDisk))

	// bytes := []byte(*fitDisk)
	// // Obtener el primer byte (primer carácter) del slice de bytes
	// fitdiskChar := bytes[0]
	// mbr := Logica.MBR{
	// 	Mbr_tamanio:        int64(*sizeDisk),
	// 	Mbr_fecha_creacion: [19]byte{},
	// 	Mbr_disk_signature: 0,
	// 	Mdsk_fit:           fitdiskChar,
	// 	Mbr_partition:      [4]byte{},
	// }
	// bytesFechaCreacion := []byte(Logica.GetCurrentTime())
	// copy(mbr.Mbr_fecha_creacion[:], bytesFechaCreacion)
	// disco.InsertObject(mbr)
	// fmt.Println("MBR:",mbr)
	disco.CreateDskFIle()
}


//funcion que verifica los flags de mkdisk
func checkFLagsMkDisK(match []string, size *int, fit *string, unit *string,disco *Logica.Disco) {
	flagName := match[1]
			flagValue := match[2]
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
				case "size":       
					sizeValue := 0
					fmt.Sscanf(flagValue, "%d", &sizeValue)
					size:=make([]byte, sizeValue)
					disco.SetSize(size)
					
				case "fit":
					flagValue = flagValue[:1]
					flagValue = strings.ToLower(flagValue)
					
				case "unit":
					flagValue = strings.ToLower(flagValue)
					
				default:
					fmt.Println("Error: Flag not found",flagName,flagValue)
			}
}



// ObtenerProximoNombreDisco busca el próximo nombre de archivo disponible en la carpeta especificada
func ObtenerProximoNombreDisco(carpeta string) (string, error) {
    // Obtener la lista de nombres de archivos existentes en la carpeta
    archivosExistente, err := obtenerArchivosExistente(carpeta)
    if err != nil {
        return "", err
    }

    // Iterar sobre las letras del alfabeto
    for letra := 'A'; letra <= 'Z'; letra++ {
        nombre := string(letra) + ".disk"

        // Verificar si el nombre de archivo ya existe
        if _, existe := archivosExistente[nombre]; !existe {
            return nombre, nil // Devolver el nombre si no existe
        }
    }

    return "", fmt.Errorf("no se encontraron nombres disponibles en la secuencia")
}

func obtenerArchivosExistente(carpeta string) (map[string]struct{}, error) {
    archivosExistente := make(map[string]struct{})

    // Obtener la lista de archivos en la carpeta
    err := filepath.Walk(carpeta, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".disk") {
            archivosExistente[info.Name()] = struct{}{}
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    return archivosExistente, nil
}