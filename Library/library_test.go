package Library

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewBook(t *testing.T) {
	newBook, err := NewBook("Книга", "","Издательство", "Вид", "Размер", "")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Printf("new book : %#v", newBook)
}

func TestNewBook_ErrorEmptyName(t *testing.T) {
	_, err := NewBook("", "Author","Издательство", "Вид", "Размер", "")
	if err == nil {
		t.Log("Failed: the function doesn't return an error")
		t.Fail()
	}
}

func TestNewBook_ErrorEmptyAuthor(t *testing.T) {
	_, err := NewBook("Name", "","Издательство", "Вид", "Размер", "")
	if err == nil {
		t.Log("Failed: the function doesn't return an error")
		t.Fail()
	}
}

func TestNewBook_ErrorEmptyArgs(t *testing.T) {
	if _, err := NewBook("Книга", "Автор", "", "", "", ""); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNewBook_ErrorBadConstructor(t *testing.T) {
	newBook, err := NewBook("nf","nf", "nf", "nf", "nf", "nf")
	if err != nil {
		t.Log("Unhandled error")
		t.Fail()
	}
	if newBook == nil || newBook.Name != "nf" || newBook.Author != "nf" || newBook.Publisher !=  "nf" ||
		newBook.Kind != "nf" || newBook.Size != "nf" || newBook.HoldingShelf != "nf" {
		t.Log("Failed. Expected full 'nf', received ", newBook)
		t.Fail()
	}
}

func TestStrings_Split(t *testing.T) {
	str := strings.Replace("-f -a -t", " ", "", -1)
	str = strings.Replace(str, "-", "|", -1)
	fmt.Printf("%#v", strings.Split(str, "|")[1:])
}

func TestBooks_SortBooksWithParams(t *testing.T) {
	var db Database
	db = &FileIO{}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	var library Books
	//var lib *Books
	err = db.LoadBooks(&library)
	if err != nil {
		log.Fatal(err)
	}

	Book{}.GetStringTableTitle()
	for _, b := range library {
		fmt.Println(b.GetStringTableItem())
	}
	fmt.Printf("\nStart sorting\n\n")

	err = library.SortBooksWithParams("name", "-r")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(library.GetBooksTableWithTitle("Sorting books test"))
}

func TestBooks_GetBooksTable(t *testing.T) {
	var db Database
	db = &FileIO{}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	var library Books
	//var lib *Books
	err = db.LoadBooks(&library)
	if err != nil {
		log.Fatal(err)
	}

	var book Book
	err = db.LoadBook(&book, "Совершенный код")

	for i := 0; i < 100; i++ {
		book.ID = 4 + i
		book.Name = "Совершенный код" + " ч. " + strconv.Itoa(i)
		book.AddedAt = book.AddedAt.Add(time.Hour * 24)
		library = append(library, book)
	}

	fmt.Println(library.GetBooksTableWithTitle("Find books with params - name=Объектно flags=\"\""))
}

func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var db Database
		db = &FileIO{}

		err := db.Connect()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Disconnect()

		var library Books
		err = db.LoadBooks(&library)
		if err != nil {
			log.Fatal(err)
		}

		err = library.SortBooksWithParams("name", "-r")
		if err != nil {
			log.Fatal(err)
		}
	}
}