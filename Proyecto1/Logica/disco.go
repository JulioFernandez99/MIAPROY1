package Logica

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/binary"
)


type Disco struct {
	Size []byte
	FileName string
	File *os.File
}


//crea el file con el tamaño especificado
func (d *Disco) CreateDskFIle() {
	d.CreateFIle()
	d.OpenDskFIle()
	d.InitDisk(d.File, d.Size, 0)
	d.CloseDskFIle()
}


//crea el archivo .disk pero sin el tamaño
func (d *Disco) CreateFIle() {
	// Create bin file
	if err := CreateFile(d.GetFileName()); err != nil {
		return
	}
	
} 

//abre el archivo .disk
func (d *Disco) OpenDskFIle() {
	// Open bin file
	file, err := OpenFile(d.GetFileName())
	if err != nil {	
		return
	}
	d.File = file
}

//cierra el archivo .disk
func (d *Disco) CloseDskFIle() {
	// Close bin file
	if err := d.File.Close(); err != nil {
		fmt.Println("Err CloseDskFIle==", err)
		return	
	}
}

//crea el archivo .disk
func CreateFile(name string) error {
	
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile Create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

//abre el archivo .disk
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}



func (d *Disco)  InitDisk(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err InitDisk==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func (d *Disco) ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}


// InsertObject inserta un objeto en el disco desde el inicio del archivo
func (d *Disco) InsertObject(data interface{}) error {
	// Mover el puntero de escritura al inicio del archivo
	d.File.Seek(0, 0)
	
	// Escribir el objeto en el disco
	err := binary.Write(d.File, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Error al insertar el objeto:", err)
		return err
	}
	
	return nil
}




//Metodo get para el nombre del disco
func (d *Disco) GetFileName() string {
	return d.FileName
}

//set para el nombre del disco
func (d *Disco) SetFileName(FileName string ){
	d.FileName = FileName
}

//set para el tamaño del disco
func (d *Disco) SetSize(Size []byte) {
	d.Size = Size
}

//get para el tamaño del disco
func (d *Disco) GetSize() []byte {
	return d.Size
}

//set para el ajuste del disco
