package Analizador

import (
	"MIA_P1OFICIAL_201902416/Logica"
	"MIA_P1OFICIAL_201902416/Structs"
	"bufio"
	"encoding/binary"
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
	
	driveletterMount = "A"
	nameMount = "A"
	idMount = "A"
	cadaMount = ""

	driveletterUnMount = "A"
	idUnMount = "A"
		
	rePartitions = regexp.MustCompile(`(?:fpartitions\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reDisk  = regexp.MustCompile(`(?:mkdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reRdisk = regexp.MustCompile(`(?:rmdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reFDisk = regexp.MustCompile(`(?:fdisk\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reMount = regexp.MustCompile(`(?:mount\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)
	reUnMount = regexp.MustCompile(`(?:unmount\s+|-(mkdisk|\w+)=("[^"]+"|\S+))`)

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

	matchPartitions := rePartitions.FindAllStringSubmatch(input, -1)
	comandPartitions := strings.Trim(matchPartitions[0][0], " ")

	matchMount := reMount.FindAllStringSubmatch(input, -1)
	conmandMount := strings.Trim(matchMount[0][0], " ")

	matchUnMount := reUnMount.FindAllStringSubmatch(input, -1)
	conmandUnMount := strings.Trim(matchUnMount[0][0], " ")


	if comandDisk == "mkdisk" && len(matchMkDisk) > 0 {
		processMkDisk(matchMkDisk[1:])
	} else if comandRmDisk == "rmdisk" && len(matchRmDisk) > 0 {
		processRmDisk(matchRmDisk[1:])
	} else if comandFDisk == "fdisk" && len(matchFDisk) > 0 {
		processFDisk(matchFDisk[1:])
	}else if comandPartitions == "fpartitions" && len(matchPartitions) > 0 {
		processPartitions(matchPartitions[1:])
	}else if conmandMount == "mount" && len(matchMount) > 0 && len(matchMount) > 2{
		fmt.Println("entro mount")
		fmt.Println(len(matchMount))
		processMount(matchMount[1:]) 
	}else if conmandUnMount == "unmount" && len(matchUnMount) > 0 && len(matchUnMount) <= 2{
		fmt.Println("entro unmount")
		processUnMount(matchUnMount[1:])
	}else {
		fmt.Println("Comando no encontrado")
	}
}


func processUnMount(match [][]string) {
	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]

		if flagName == "id" {
			driveletterUnMount = "DISCOS/" + string(flagValue[0]) + ".disk"
			idUnMount = string(flagValue[1]-1)
		}
	}

	//Aqui obtengo el mbr del disco
	mbr, err := Logica.GetMBR(driveletterUnMount)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", driveletterUnMount, ",verifique si existe.")
		return
	}

	
	if err != nil {
		fmt.Println("Error al convertir idMountUnMount a int64:", err)
		return
	}

	fmt.Println("idUnMount", idUnMount)

	idUnMountInt, err := strconv.Atoi(idUnMount)
	if err != nil {
		fmt.Println("Error al convertir idUnMount a int:", err)
		return
	}

	for i := range mbr.Mbr_partition {
		if i == idUnMountInt {
			fmt.Println(i,". Particion",idUnMountInt)
			mbr.Mbr_partition[i].Part_status = 0
			fmt.Println("Particion desmontada con exito")
			
		}
	}

	//Aqui guardo el mbr con la particion creada
	Logica.SaveMBR(driveletterUnMount, *mbr)
}



func processMount(match [][]string) {
//Este for es para iterar sobre los flags y signarles su valor
	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]

		if flagName == "driveletter" {
			driveletterMount = "DISCOS/" + flagValue + ".disk"
			idMount = flagValue
		} else if flagName == "name" {
			flagValue = strings.Trim(flagValue, "\"")
			flagValue = strings.ToLower(flagValue)
			
			nameMount = strings.Trim(flagValue, " ")
		}else {
			fmt.Println("Error: Flag not found", flagName, flagValue)
		}
	}

	//Aqui obtengo el mbr del disco
	mbr, err := Logica.GetMBR(driveletterMount)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", driveletterMount, ",verifique si existe.")
		return
	}

	
	for i := range mbr.Mbr_partition {
		// convertir mbr.Mbr_partition[i].Part_name que es un array de bytes a string
		cadaMount= string(mbr.Mbr_partition[i].Part_name[:])
		if strings.Contains(cadaMount, nameMount) {
			
			mbr.Mbr_partition[i].Part_status = 1
			copy(mbr.Mbr_partition[i].Part_id[:], []byte(idMount + strconv.Itoa(i+1) + "19"))
			fmt.Println("Particion montada con exito", cadaMount)
			break
		} else if i == 3 {
			fmt.Println("No se encontro una particion con el nombre:", nameMount)
		}

		// if string(mbr.Mbr_partition[i].Part_name[:]) == nameMount {
		// 	mbr.Mbr_partition[i].Part_status = 1
		// 	copy(mbr.Mbr_partition[i].Part_id[:], []byte(idMount + strconv.Itoa(i) + "19"))
		// 	fmt.Println("Particion montada con exito", mbr.Mbr_partition[i])
		// 	break
		// } else if i == 3 {
		// 	fmt.Println("No se encontro una particion con el nombre:", nameMount)
		// }
		
	}

	//Aqui guardo el mbr con la particion creada
	Logica.SaveMBR(driveletterMount, *mbr)
	

}

func processPartitions(match [][]string) {
	path :=""

	for _, match := range match {
		flagName := match[1]
		flagValue := match[2]
		
		switch flagName {
		case "driveletter":
			path = "DISCOS/" + flagValue + ".disk"
		default:
			fmt.Println("Error: Flag not found", flagName, flagValue)
		}
	}

	mbrPrueba, err := Logica.GetMBR(path)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", driveletterFDisk,",verifique si existe.")
		return
	}
	PrintPartitionsDiponibles(*mbrPrueba)
	
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
			flagValue = flagValue[:1]
			flagValue = strings.ToLower(flagValue)
			
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

		//Aqui convierto el valor de addFDisk a entero
		addValue, _ := strconv.ParseInt(addFDisk, 10, 64)

		//Este if es para verificar si la particion esta libre
		if mbrPrueba.Mbr_partition[i].Part_s == 0  && deleteFDisk == "false" && addFDisk == "0"{

			
			//Aqui asigno los valores a la particion
			mbrPrueba.Mbr_partition[i].Part_tyPe = typeFDisk[0]
			mbrPrueba.Mbr_partition[i].Part_fit = fitFDisk[0]
			mbrPrueba.Mbr_partition[i].Part_s = GetSizeUnit(sizeFDisk, unitFDisk)
			copy(mbrPrueba.Mbr_partition[i].Part_name[:], nameFDisk)	
			mbrPrueba.Mbr_partition[i].Part_correlative = int64(i + 1)
			mbrPrueba.Mbr_partition[i].Part_start = GetSizeStart(i,*mbrPrueba)	
			fmt.Println("Particion creada con exito", mbrPrueba.Mbr_partition[i])
			
			break

		}else if deleteFDisk == "full" && cadenaByte == nameFDisk && addFDisk == "0" { 
			//Este if es para eliminar una particion del mbr
			fmt.Println("Eliminará la partición", deleteFDisk)
			mbrPrueba.Mbr_partition[i]= Structs.Partition{}
			mbrPrueba.Mbr_partition[i].Part_s = 0
			fmt.Println("Particion eliminada con exito")
			deleteFDisk = "false"
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
			addFDisk = "0"
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

func GetSizeUnit(size int, unit string) int64 {
	if unit == "k" {
		return int64(size * 1000)
	}
	if unit == "m" {
		return int64(size * 1000000)
	}
	return int64(size)
}

func PrintMBR(mbr Structs.MBR,i int) {
	//imprimir cada atributo de la particion creada
	fmt.Println("Status",mbr.Mbr_partition[i].Part_status)
	fmt.Printf("Type %c\n",mbr.Mbr_partition[i].Part_tyPe)
	fmt.Printf("Fit %c\n",mbr.Mbr_partition[i].Part_fit)
	fmt.Println("Start",mbr.Mbr_partition[i].Part_start)
	fmt.Println("Size",mbr.Mbr_partition[i].Part_s)
	fmt.Println("Name",mbr.Mbr_partition[i].Part_name)
	fmt.Println("Correlative",mbr.Mbr_partition[i].Part_correlative)
	fmt.Printf("ID %s\n",mbr.Mbr_partition[i].Part_id)
}

func GetSizeStart(i int,mbr Structs.MBR) int64 {
	if i == 0{
		return int64(binary.Size(mbr)) + 1
	}else{
		return (mbr.Mbr_partition[i-1].Part_start + mbr.Mbr_partition[i-1].Part_s)+1
	}
}

func PrintPartitionsDiponibles(mbr Structs.MBR) {
	for i := range mbr.Mbr_partition {
		if mbr.Mbr_partition[i].Part_s == 0{
			fmt.Println("-----------------Particion",i,"dispobible","-----------------")
		PrintMBR(mbr,i)
		}else{
			fmt.Println("-----------------Particion",i,"ocupada-----------------")
		PrintMBR(mbr,i)
		}
		
	}
}