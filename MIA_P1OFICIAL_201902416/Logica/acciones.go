package Logica

import (
	"MIA_P1OFICIAL_201902416/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)


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
func ReadObject(name string,data interface{}) error {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		//fmt.Println("Err OpenFile==",err)
		return err
	}

	file.Seek(0, 0)
	err = binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		//fmt.Println("Err ReadObject==",err)
		return err
	}
	return err
}


func GetMBR(path string) (*Structs.MBR, error) {
	var TempData Structs.MBR
	err := ReadObject(path, &TempData);
	if err != nil {
		return &TempData,err
	}
	return &TempData,err
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

