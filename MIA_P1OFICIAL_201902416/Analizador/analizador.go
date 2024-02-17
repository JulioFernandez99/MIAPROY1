package Analizador

import (
	"MIA_P1OFICIAL_201902416/Logica"
	"MIA_P1OFICIAL_201902416/Structs"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ----------------------------------flags mkdisk----------------------------------
var (
	sizeDisk = 0
	fitDisk  = "f"
	unitDisk = "m"

	sizeFDisk        = 0 		//obligatorio
	driveletterFDisk = "A"		//obligatorio
	nameFDisk        = "path"	//obligatorio
	unitFDisk        = "k"		//Opcional	
	typeFDisk        = "p"		//Opcional
	fitFDisk         = "f"		//Opcional
	deleteFDisk      = "false"  //Opcional
	addFDisk         = "0"	//Opcional

	cadenaByte= ""
	

	reDisk  = regexp.MustCompile(`(?:mkdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reRdisk = regexp.MustCompile(`(?:rmdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reFDisk = regexp.MustCompile(`(?:fdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
)

//-----------------------------------fin glags mkdisk------------------------------------

//Run es la funcion principal que se encarga de leer la entrada del usuario
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

//Esta funcion lee la entrada del usuario
func ReadLine() (string, error) {
	fmt.Print(">>")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	return text, err
}

//Esta funcion procesa todos los comandos
func ProcessInput(input string) {
	matchMkDisk := reDisk.FindAllStringSubmatch(input, -1)
	comandDisk := strings.Trim(matchMkDisk[0][0], " ")

	matchRmDisk := reRdisk.FindAllStringSubmatch(input, -1)
	comandRmDisk := strings.Trim(matchRmDisk[0][0], " ")

	matchFDisk := reFDisk.FindAllStringSubmatch(input, -1)
	comandFDisk := strings.Trim(matchFDisk[0][0], " ")

	if comandDisk == "mkdisk" && len(matchMkDisk) > 0 {
		processMkDisk(matchMkDisk[1:])
	} else if comandRmDisk == "rmdisk" && len(matchRmDisk) > 0 {
		processRmDisk(matchRmDisk[1:])
	} else if comandFDisk == "fdisk" && len(matchFDisk) > 0 {
		processFDisk(matchFDisk[1:])
	}
}

//Esta funcion procesa el comando fdisk
func processFDisk(match [][]string) {
	//Este for es para iterar sobre los flags y signarles su valor
	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]

		if flagName == "size" {
			value, _ := strconv.Atoi(flagValue)
			sizeFDisk = value
		} else if flagName == "driveletter" {
			driveletterFDisk = "DISCOS/" + flagValue + ".disk"
		} else if flagName == "name" {
			flagValue = strings.Trim(flagValue, "\"")
			flagValue = strings.ToLower(flagValue)
			nameFDisk = strings.Trim(flagValue, " ")
		} else if flagName == "unit" {
			flagValue = strings.ToLower(flagValue)
			unitFDisk = flagValue
		} else if flagName == "type" {
			typeFDisk = flagValue
		} else if flagName == "fit" {
			flagValue = flagValue[:1]
			flagValue = strings.ToLower(flagValue)
			fitFDisk = flagValue
		} else if flagName == "delete" {
			deleteFDisk = flagValue
		} else if flagName == "add" {
			addFDisk = flagValue
		} else {
			fmt.Println("Error: Flag not found", flagName, flagValue)

		}

	}
	


	
	//Aqui obtengo el mbr del disco
	mbrPrueba, err := Logica.GetMBR(driveletterFDisk)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", driveletterFDisk,",verifique si existe.")
		return
	}

	

	//Este for es para iterar sobre las particiones y buscar una libre
	for i := range mbrPrueba.Mbr_partition {
		//Aqui obtengo el nombre de la particion
		cadenaByte = string(mbrPrueba.Mbr_partition[i].Part_name[:])
		cadenaByte = strings.ReplaceAll(cadenaByte, "”", "")
		cadenaByte = strings.Trim(cadenaByte, "\x00")
		addValue, _ := strconv.ParseInt(addFDisk, 10, 64)
		
		//Este if es para verificar si la particion esta libre
		if mbrPrueba.Mbr_partition[i].Part_status == 0  && deleteFDisk == "false" && addFDisk == "0"{
			fmt.Println("Particion libre")
			if unitFDisk == "k" {
				sizeFDisk = sizeFDisk * 1000
			}else if unitFDisk == "m" {
				sizeFDisk = sizeFDisk * 1000000
			}
			mbrPrueba.Mbr_partition[i].Part_status = 1
			mbrPrueba.Mbr_partition[i].Part_tyPe = typeFDisk[0]
			mbrPrueba.Mbr_partition[i].Part_fit = fitFDisk[0]
			mbrPrueba.Mbr_partition[i].Part_s = int64(sizeFDisk)
			copy(mbrPrueba.Mbr_partition[i].Part_name[:], nameFDisk)
			mbrPrueba.Mbr_partition[i].Part_correlative = int64(i + 1)
			copy(mbrPrueba.Mbr_partition[i].Part_id[:], "0000")
			fmt.Println("Particion creada con exito", mbrPrueba.Mbr_partition[i])
			break
		}else if deleteFDisk == "full" && cadenaByte == nameFDisk && addFDisk == "0" { 
			//Este if es para eliminar una particion del mbr
			fmt.Println("Eliminará la partición", deleteFDisk)
			mbrPrueba.Mbr_partition[i]= Structs.Partition{}
			fmt.Println("Particion eliminada con exito")
			break
		}else if addFDisk != "0" && cadenaByte == nameFDisk {
			//Este if es para aumentar el tamaño de una particion del mbr
			if unitFDisk == "k" {
				addValue = addValue * 1000
			}else if unitFDisk == "m" {
				addValue = addValue * 1000000
			}
			mbrPrueba.Mbr_partition[i].Part_s = mbrPrueba.Mbr_partition[i].Part_s + addValue
			fmt.Println("Particion aumentada con exito")
			break
		}else if i == 3 && (deleteFDisk == "full" || addFDisk != "0") {
			fmt.Println("No se encontro una particion con el nombre:", nameFDisk)
		}else if i == 3 && deleteFDisk == "false" {
			fmt.Println("Todas las particiones están ocupadas")
			
		}
	}

	//Aqui guardo el mbr con la particion creada
	Logica.SaveMBR(driveletterFDisk, *mbrPrueba)

}

//Esta funcion procesa el comando rmdisk
func processRmDisk(match [][]string) {
	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]

		switch flagName {
		case "driveletter":
			path := path.Join("DISCOS", flagValue+".disk")
			err := os.Remove(path)
			if err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println("Error: Flag not found", flagName, flagValue)
		}
	}
}

func processMkDisk(match [][]string) {

	// Obtener el próximo nombre de disco disponible
	nombreDisco, err := ObtenerProximoNombreDisco("DISCOS")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	FileName := "DISCOS/" + nombreDisco

	if len(match) == 1 && match[0][1] == "size" {
		//Si sOlo tiene un parametro
		value, _ := strconv.Atoi(match[0][2])
		sizeDisk = value * 1000000
	} else if len(match) == 2 && match[1][1] == "fit" {
		//Si tiene dos parametros y el segundo es fit
		fmt.Println("entro")
		value, _ := strconv.Atoi(match[0][2])
		sizeDisk = value * 1000000
	} else {
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
				sizeDisk = sizeValue
			case "fit":
				flagValue = flagValue[:1]
				flagValue = strings.ToLower(flagValue)
				fitDisk = flagValue
			case "unit":
				flagValue = strings.ToLower(flagValue)
				if flagValue == "m" {
					//fmt.Println("Mega")
					sizeDisk = sizeDisk * (1000000)
				}
				if flagValue == "k" {
					sizeDisk = sizeDisk * (1000)

				}
				unitDisk = flagValue
			default:
				fmt.Println("Error: Flag not found", flagName, flagValue)
			}
		}
	}
	Logica.CreateDisk(sizeDisk, unitDisk, fitDisk, FileName)
	Logica.InsertMBR(FileName, sizeDisk, unitDisk, fitDisk)
	// mbrPrueba:=Logica.GetMBR(FileName)
	// fmt.Print("Disco creado con exito\n")
	// fmt.Println(mbrPrueba)

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
