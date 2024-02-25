package Logica

import (
	"MIA_P1OFICIAL_201902416/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	
	"time"
)
var data Structs.EBR
var tempStart int64
var tempSize int64
var start int64
var next int64


func CreateDisk(size int, unit string, fit string, path string)  {
	
	// Create bin file
	if err := CreateFile(path); err != nil {
		return 
	}

	// Open bin file
	file, err := OpenFileD(path)
	if err != nil {
		return 
	}

	sizeDisk := make([]byte, size)//size es el tamaño del disco
	
	if err := WriteObject(file,sizeDisk,0); err != nil {
	 	return 
	}
}

func InsertMBR(path string, size int, unit string, fit string) {
	// Create MBR
	var fitChar [1]byte
	copy(fitChar[:], fit)
	mbr := Structs.MBR{}
	mbr.Mbr_tamanio = int64(size)
	mbr.Mbr_fecha_creacion = GetCurrentTime().Unix()
	mbr.Mbr_disk_signature = 0
	mbr.Mdsk_fit = fitChar
	
	// Open bin file
	file, err := OpenFileD(path)
	if err != nil {
		return 
	}
	// Write MBR
	if err := WriteObject(file, mbr, 0); err != nil {
		return
	}

	defer file.Close()


}

// Funtion to create bin file
func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==",err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==",err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFileD(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("ERROR AL ABRIR EL ARCHIVO!")
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position  int64) error {
	
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==",err)
		return err
	}
	return nil
}


//-------------------------------------Devuelve el MBR del disco-------------------------------------
func ReadObject(name string,data interface{},posicion int64) error {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		//fmt.Println("Err OpenFile==",err)
		return err
	}

	file.Seek(posicion, 0)
	err = binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		//fmt.Println("Err ReadObject==",err)
		return err
	}
	return err
}



func GetMBR(path string) (*Structs.MBR, error) {
	var TempData Structs.MBR
	err := ReadObject(path, &TempData,0);
	if err != nil {
		return &TempData,err
	}
	return &TempData,err
}


//funcion para verificar el tipo de dato en una posicion del disco
func InsertEbr(path string,position int64,ebrInsert Structs.EBR) (bool) {
	

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==",err)
		return false
	}

	//fmt.Println("Posicion",position)
	file.Seek(position, 0)
	

	
	
	err = binary.Read(file, binary.LittleEndian, &data)

	if err != nil {
		fmt.Println("Err ReadObject9==",err)
		return false
	}
	


	if data.Part_next == 0 {
		//crear un for
		
		fmt.Println("Es un cero ",ebrInsert)
		// ebrInsert.Part_name = data.Part_name
		//en position del file.seek se debe de poner el valor de la posicion del ebr
		WriteObject(file, ebrInsert, position)
		return true
		
	}else {
		
		fmt.Println("Ya hay un ebr ",data.Part_next)

		i:=0
		for data.Part_next != 0 {
			i++
			//fmt.Println("Estamos en el next: ",data.Part_start)
			//fmt.Println("Tenemos que ir al next: ",data.Part_next)
			tempStart = data.Part_next
			tempSize = data.Part_s

			start = tempStart + int64(tempSize) + int64(binary.Size(data))
			//fmt.Println("Start-",start)
			next = start + int64(binary.Size(ebrInsert))

			data.Part_start = start
			data.Part_next = next

			
			file.Seek(start, 0)
			err = binary.Read(file, binary.LittleEndian, &data)
			fmt.Println("Data del ebr en el next" ,tempStart,"->",data)
			if err != nil {
				fmt.Println("Err ReadObject==",err)
				
			}
		}

		//fmt.Println("Estamos en el next->",tempStart)
		//fmt.Println("IMPRIMIENDO EL EBR",data)
		//fmt.Println("IMPRESION DEL SIZE",tempSize)
		fmt.Println("Entro",i,"veces")
		//file.Seek(3000, 0)
	
		//fmt.Println("Next-",next)

		ebrInsert.Part_start = start
		ebrInsert.Part_next = next		
		ebrInsert.Part_s = tempSize 	
		
		//fmt.Println("quiero ir a la posicion",start)
		
		// ir al file en la posicion start
		
		

		
		

		//escribir el ebr en la posicion start
		WriteObject(file, ebrInsert, start)

		
		

		
		// fmt.Println("Next del ebr",data.Part_next)
		
		// var newEbr Structs.EBR
		// newEbr.Part_name = data.Part_name
		

		// fmt.Println("Next del new Ebr",newEbr.Part_next)
		// Insert(path, newEbr, newEbr.Part_next)
		return false
	}
	
	
	//Validando si es un cero o un ebr lo que viene en la posicion
	// if data.Part_next == 0 {

	// 	fmt.Println("Es un cero ",data.Part_name)
	// 	// 
	// 	ebrInsert.Part_name = data.Part_name
	// 	Insert(path, ebrInsert, position)
	// 	// data.Part_next = 666
	// 	// Insert(path, data, position)
		
	// 	return true
		
	// }else {

	// 	fmt.Println("Es un ebr ",data.Part_name)
		

	// 	//InsertEbr(path,data.Part_next+100,ebrInsert)
	// 	// file.Seek(position, 0)
	// 	// var data Structs.EBR
	// 	// err = binary.Read(file, binary.LittleEndian, &data)
	// 	// if err != nil {
	// 	// 	fmt.Println("Err ReadObject==",err)
	// 	// 	return false
	// 	// }
	// 	// fmt.Println("Data del ebr",data.Part_next)
		
	// 	return false
	// }
	
}


func SaveMBR(path string, mbr Structs.MBR)  {
	// Open bin file
	file, err := OpenFileD(path)
	if err != nil {
		fmt.Println("Err OpenFileD==",err)
		return
	}
	// Write MBR
	if err := WriteObject(file, mbr, 0); err != nil {
		fmt.Println("Err WriteObject==",err)
		return
	}

	defer file.Close()
}


func PrintMBR(path string) {
	TempData,err := GetMBR(path)
	if err != nil {
		fmt.Println("Err GetMBR==",err)
		return
	}
	// Print object
	fmt.Println("MBR ",path," ----->",TempData)
	fmt.Println("MBR==",TempData)
	fmt.Println("TempData==",TempData.Mbr_tamanio)
	fmt.Println("TempData==",time.Unix(TempData.Mbr_fecha_creacion, 0))
	fmt.Println("TempData==",TempData.Mbr_disk_signature)
	fmt.Println("TempData==", string(rune(int(TempData.Mdsk_fit[0]))))
	fmt.Println("TempData==",TempData.Mbr_partition[0])
}

func GetCurrentTime() time.Time {
	currentTime := time.Now()
	return currentTime
}



