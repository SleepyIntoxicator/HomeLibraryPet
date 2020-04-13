package Library

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type Commits map[int]Commit
type Commit struct {
	ID int
	changes Books
}

type HoldingShelf struct {
	Short string
	Description string
}

//DONE: added_at timestamp
type Book struct {
	ID int
	Name string
	Author string
	Publisher string
	Kind string
	Size string
	HoldingShelf string
	AddedAt time.Time `json:"AddedAt string"`
}

type BookJSON struct {
	ID int
	Name string
	Author string
	Publisher string
	Kind string
	Size string
	HoldingShelf string
	AddedAt string
}

func (Bjson BookJSON) convertFromBookJSON() (Book, error) {
	loadTime, err := time.Parse(dateRFClex, Bjson.AddedAt)
	if err != nil {
		err = errors.New("unsupported time format when converting: " + err.Error())
		return Book{}, err
	}

	book, err := NewBook(Bjson.Name, Bjson.Author, Bjson.Publisher, Bjson.Kind, Bjson.Size, Bjson.HoldingShelf)
	if err != nil {
		return Book{}, err
	}
	book.ID = Bjson.ID
	book.AddedAt = loadTime
	return *book, nil
}

func (Bjson *BookJSON) convertFromBook(book Book){
	*Bjson = BookJSON{book.ID, book.Name, book.Author, book.Publisher, book.Kind, book.Size, book.HoldingShelf, "N\\A"}
	Bjson.AddedAt = fmt.Sprintf("%.2d.%.2d.%.2d", book.AddedAt.Day(), book.AddedAt.Month(), book.AddedAt.Year())
}

type Books []Book

//Создаёт новую книгу с заданными параметрами, но без полки хранения. Полку можно указать позже.
func NewBook(name, author, publisher, kind, size, shelf string) (book *Book, err error){
	book = &Book{-1, "", "", publisher, kind,size, "none", time.Now()}
	book.ID = -1
	if name == "" {
		book.Name = "Error: in library.NewBook(). Name param is null. Check details."
		err = fmt.Errorf("empty book name. Book can't have no name")
		return book, err
	} else {
		book.Name = name
	}
	if author == "" {
		book.Author = "Error: in NewBook(). Author param is null. Check details."
		err = fmt.Errorf("empty book author. Book can't have no author")
		return book, err
	} else {
		book.Author = author
	}
	if publisher == "" {	//TODO: rewrite logging errors
		log.Print("Warning: in NewBook(...). Publisher param is null. Check details.")
	}
	if kind == "" || size == "" {
		log.Print("Warning: in NewBook(...). Kind or size param is null. Check details.")
	}
	if shelf != "" {
		book.HoldingShelf = shelf
	}
	timestamp := time.Now()
	book.AddedAt = time.Date(timestamp.Year(),timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.Local)
	return book, nil
}

func (book *Book) SetShelf(newShelf string) {
	if newShelf != "" {
		book.HoldingShelf = newShelf
	}
}

func (book *Book) IsEmpty() bool {
	return book.ID == -1 || (book.ID == 0 && book.Name == "" && book.Author == "" && book.Publisher == "")
}

func (book *Book) GetDateStr() string {
	return fmt.Sprintf("%.2d.%.2d.%.4d", book.AddedAt.Day(), int(book.AddedAt.Month()), book.AddedAt.Year())
}

func (bs *Books) String() string {
	return bs.GetBooksTable()
}

func (book Book) getBooksHeaderWithTitleF(j []int, name string) string {

	//TODO: ISSUE #3	CLOSED
	// Add the handling of 3 situations
	// Summary of table width and nameLength may be even
	// *
	// n = 19 w = 161		r = 0		| bWR = 0
	// n = 20 w = 161		r = 0		| bWR = 1
	// n = 19 w = 162		r = -1		| bWR = 0
	// n = 20 w = 162		r = +1		| bWR = 1

	var width int	//ширина всей таблицы
	width = len(j) * 2 - 1 //одно боковое деление не отрисовывается
	for _, indexW := range j {
		width += indexW
	}
	nameLength := utf8.RuneCountInString(name)
	balanceWeight := 0 //To balance width when nameLength is even
	if  (width - nameLength) % 2 == 1 {
		balanceWeight = 1
	}

	//----------------
	var str string
	str += "┏" + strings.Repeat("━", width) + "┓\n"
	str += "┃" + strings.Repeat(" ", (width - nameLength)/2) + name + strings.Repeat(" ", (width - nameLength)/2+balanceWeight) + "┃\n"
	//----------------
	str += "┣"
	for i, b := range j {
		if i != 0 {
			str += "┳"
		}
		str += strings.Repeat("━", b + 1)
	}
	str += "┫\n"
	for _, b := range j {
		str += "┃ %-"
		str += fmt.Sprintf("%ds", b)
	}
	str += "┃\n"
	//----------------
	str += "┡"
	for i, b := range j {
		if i != 0 {
			str += "╇"
		}
		str += strings.Repeat("━", b + 1)
	}
	str += "┩\n"
	return str
}

func (book Book) GetStringTableTitle() string {
	str := "┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━┓\n"
	str += fmt.Sprintf("┃ №    %-48s\t ┃ %-20s\t┃ %-18s\t\t ┃ %-16s\t\t ┃ %s\t\t┃ %s\t\t┃ %s┃\n", "Название", "Автор", "Издательство", "Тип", "Размер", "Полка хранения", "Дата добавления")
	str += "┡━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━┩\n"
	return str
}

func (book *Book) GetStringTableItem() string {
	if book != nil {
		return fmt.Sprintf("[%d]  %-48s\t  %-20s\t %-18s\t\t  %-16s\t\t  %s\t\t%-15s\t%s", book.ID, book.Name, book.Author, book.Publisher, book.Kind, book.Size, book.HoldingShelf, book.GetDateStr())
	}
	return "Error. The book is empty or not exist"
}

func (book Book) getBooksHeaderF(j []int) string {
	str := "┏"
	for i, b := range j {
		if i != 0 {
			str += "┳"
		}
		str += strings.Repeat("━", b + 1)
	}
	str += "┓\n"
	//----------------
	for _, b := range j {
		str += "┃ %-"
		str += fmt.Sprintf("%ds", b)
	}
	str += "┃\n"
	//----------------
	str += "┡"
	for i, b := range j {
		if i != 0 {
			str += "╇"
		}
		str += strings.Repeat("━", b + 1)
	}
	str += "┩\n"
	return str
}

func (book Book) getBooksItemF(j []int) string {
	var str string

	for i, b := range j {
		if i == 0 {
			str += "│ %-"
			str += fmt.Sprintf("%dd", b)
		} else {
			str += "│ %-"
			str += fmt.Sprintf("%ds", b)
		}
	}
	str += "│\n"
	return str
}

func (book Book) getBooksEnd(j []int) string {
	str := "└"
	for i, b := range j {
		if i != 0 {
			str += "┴"
		}
		str += strings.Repeat("─", b + 1)
	}
	str += "┘\n"
	return str
}

func (bs *Books) GetBooksTable() string {
	var str string
	literal := []int{4, 48, 20, 18, 16, 10, 15, 15}
	str = fmt.Sprintf(Book{}.getBooksHeaderF(literal), "№", "Название", "Автор", "Издательство", "Тип", "Размер", "Полка хранения", "Дата добавления")
	for _, book := range *bs {
		str += fmt.Sprintf(book.getBooksItemF(literal), book.ID, book.Name, book.Author, book.Publisher, book.Kind, book.Size, book.HoldingShelf, book.GetDateStr())
	}
	str += fmt.Sprintf(Book{}.getBooksEnd(literal))
	return str
}

func (bs *Books) GetBooksTableWithTitle(title string) string {
	var str string
	//№, Название, Автор, Изд, Тип, Разм, Полка, Дата
	literal := []int{4, 48, 21, 18, 17, 10, 15, 15}
	str = fmt.Sprintf(Book{}.getBooksHeaderWithTitleF(literal, title), "№", "Название", "Автор", "Издательство", "Тип", "Размер", "Полка хранения", "Дата добавления")
	for _, book := range *bs {
		str += fmt.Sprintf(book.getBooksItemF(literal), book.ID, book.Name, book.Author, book.Publisher, book.Kind, book.Size, book.HoldingShelf, book.GetDateStr())
	}
	str += fmt.Sprintf(Book{}.getBooksEnd(literal))
	return str
}

func (book *Book) ChangeBookShelf(shelf string) {
	if shelf != "" {
		book.HoldingShelf = shelf
	}
}

/*
Добавляет новую книгу
*/
func (bs *Books) len() int {
	return len(*bs)
}
func (bs *Books) swap(i, j int) {
	(*bs)[i], (*bs)[j] = (*bs)[j], (*bs)[i]
}
func (bs *Books) lessByID(i, j int) bool {
	return (*bs)[i].ID < (*bs)[j].ID
}
func (bs *Books) lessByName(i, j int) bool {
	return (*bs)[i].Name < (*bs)[j].Name
}
func (bs *Books) lessByAuthor(i, j int) bool {
	return (*bs)[i].Author < (*bs)[j].Author
}
func (bs *Books) lessByPublisher(i, j int) bool {
	return (*bs)[i].Publisher < (*bs)[j].Publisher
}
func (bs *Books) lessByKind(i, j int) bool {
	return (*bs)[i].Kind < (*bs)[j].Kind
}
func (bs *Books) lessBySize(i, j int) bool {
	return (*bs)[i].Size < (*bs)[j].Size
}
func (bs *Books) lessByShelf(i, j int) bool {
	return (*bs)[i].HoldingShelf < (*bs)[j].HoldingShelf
}
func (bs *Books) lessByDate(i, j int) bool {
	return (*bs)[i].AddedAt.Before((*bs)[j].AddedAt)
}


func (bs *Books) AddBook(book Book) {
	*bs = append(*bs, book)
}

func (bs *Books) GetBookByName(name string) (number int, err error){
	for i, item := range *bs {
		if strings.Contains(item.Name, name) {
			return i, nil
		}
	}
	return -1, errors.New("there is no books with this name")
}

//Сортирует список книг по param. flag=="-r" - разворачивает список
func (bs *Books) SortBooksWithParams(param string, flag string) error {
	newBooks := *bs
	switch param{
	case "id":
		sort.Slice(newBooks, newBooks.lessByID)
	case "name":
		sort.Slice(newBooks, newBooks.lessByName)
	case "author":
		sort.Slice(newBooks, newBooks.lessByAuthor)
	case "publisher":
		sort.Slice(newBooks, newBooks.lessByPublisher)
	case "kind":
		sort.Slice(newBooks, newBooks.lessByKind)
	case "size":
		sort.Slice(newBooks, newBooks.lessBySize)
	case "shelf":
		sort.Slice(newBooks, newBooks.lessByShelf)
	case "date":
		sort.Slice(newBooks, newBooks.lessByDate)
	default:
		return errors.New("unhandled sorting param")
	}
	booksLen := len(newBooks)
	*bs = make(Books, booksLen)
	for i, b := range newBooks {
		if flag == "-r" {	//Reverse sort
			(*bs)[booksLen-i-1] = b
		} else {
			(*bs)[i] = b
		}
	}
	return nil
}

/*	Возвращает книги у которых в param содержится value
	Список возможных параметров
	id			int			-range ([2]int)
	name		string		- флаги
	author		string		-
	publisher	string		-
	kind		string		-
	size		string		-
	shelf		string		-
	date		time.Time	-a (dd.mm.yyyy) -b (dd.mm.yyyy) -n
*/
func (bs *Books) FindBooksWithParams(param string, value interface{}, flags string) (*Books, error) {
	if len(*bs) == 0 {
		return nil, errors.New("empty book bs")
	}

	founded := make(Books, 0)
	for _, book := range *bs {
		switch param {
		case "id":
			if strings.Contains(flags, "-range") {
				left, right := value.([2]int)[0], value.([2]int)[1]
				if book.ID >= left && book.ID <= right {
					founded = append(founded, book)
				}
			} else {
				if book.ID == value.(int) {
					founded = append(founded, book)
				}
			}
		case "name":
			if strings.Contains(strings.ToLower(book.Name), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "author":
			if strings.Contains(strings.ToLower(book.Author), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "publisher":
			if strings.Contains(strings.ToLower(book.Publisher), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "kind":
			if strings.Contains(strings.ToLower(book.Kind), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "size":
			if strings.Contains(strings.ToLower(book.Size), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "shelf":
			if strings.Contains(strings.ToLower(book.HoldingShelf), strings.ToLower(value.(string))) {
				founded = append(founded, book)
			}
		case "date":
			if strings.Contains(flags, "-n") && book.AddedAt.Year() == time.Now().Year()  && book.AddedAt.Month() == time.Now().Month() && book.AddedAt.Day() == time.Now().Day() {
				founded = append(founded, book)
				break
			}
			if strings.Contains(flags, "-a") || strings.Contains(flags, "-b") {
				var valueTime time.Time
				var d, m, y int
				if value != "" {
					_, err := fmt.Sscanf(value.(string),"%d.%d.%d", &d, &m, &y)
					if err != nil {
						return nil, err/*errors.New(err.Error() + " With -a or -b flags must be use value=\"dd.mm.yyyy\"")*/
					}
				} else {
					return nil, errors.New("failed. With -a or -b flags must be use value=\"dd.mm.yyyy\"")
				}
				valueTime = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
				if strings.Contains(flags, "-a") && book.AddedAt.After(valueTime) {
					founded = append(founded, book)
				} else if strings.Contains(flags, "-b") && book.AddedAt.Before(valueTime) {
					founded = append(founded, book)
				}

			}
		default:
			return nil, errors.New("unhandled param")
		}
	}
	return &founded, nil
}