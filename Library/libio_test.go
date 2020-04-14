package Library

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestFileIO_Connect(t *testing.T) {
	var db Database
	db = &FileIO{}
	err := db.Connect()
	if err != nil {
		log.Print(err)
		t.Fail()
	}
}

func TestFileIO_LoadBookSuccess(t *testing.T) {
	var db Database
	db = &FileIO{}
	err := db.Connect()
	if err != nil {
		t.Fail()
	}
	var book Book
	err = db.LoadBook(&book, "Название")
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	if book.IsEmpty() {
		log.Print("Error: loaded book is empty. Book not found or wasn't loaded.")
		log.Printf("Got: %#v", book)
		t.Fail()
	}
	fmt.Printf(Book.GetStringTableTitle(Book{}))
	fmt.Printf("%s\n", book.GetStringTableItem())
}

func TestFileIO_LoadBookFailed(t *testing.T) {
	var db Database
	db = &FileIO{}
	err := db.Connect()
	if err != nil {
		t.Fail()
	}
	var book Book
	err = db.LoadBook(&book, "Не существующая книга")
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	if !book.IsEmpty() {
		log.Print("Error: loaded book must be empty, т.к её не существует. Не та книга.")
		log.Printf("Got: %#v", book)
		t.Fail()
	}
}

func TestFileIO_LoadBooks(t *testing.T) {
	var db Database
	db = &FileIO{}
	err := db.Connect()
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	var books Books
	err = db.LoadBooks(&books)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	fmt.Printf(Book.GetStringTableTitle(Book{}))
	for _, b := range books {
		fmt.Printf("%s\n", b.GetStringTableItem())
	}
	yyyy, mm, dd := time.Now().Date()
	fmt.Printf("[dd.mm.yyyy] %.2d.%.2d.%d\n", dd, mm, yyyy)
}

//FAILING TEST
func TestParseJSON(t *testing.T) {
	js := `{
    "id": 0,
    "name": "Название",
    "author": "Автор",
    "publisher": "Издательство",
    "kind": "Вид",
    "size": "Размер",
    "holdingShelf": "КП-1В",
    "addedAt": "2020-03-05T00:00:00Z" 
  }`
	var book Book
	err := json.Unmarshal([]byte(js), &book)
	if err != nil {
		log.Print(err)
		log.Print(book)
		t.Fail()
	}
	fmt.Println(book)
}