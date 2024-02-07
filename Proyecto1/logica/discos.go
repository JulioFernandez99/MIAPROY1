
package logica

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/binary"
)

var discos Disco
var disckCreated bool = false

type Disco struct {
	Size [1024]byte
	Create bool
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
	d.SetCreate(true)
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

//Metodo set para el estado del disco
func (d *Disco) GetCreate() bool {
	return d.Create
}

//Metodo get para el nombre del disco
func (d *Disco) GetFileName() string {
	return d.FileName
}

//set para el nombre del disco
func (d *Disco) SetFileName(FileName string ){
	d.FileName = FileName
}

//Metodo get para estado del disco
func (d *Disco) SetCreate(Create bool) {
	d.Create = Create
}

