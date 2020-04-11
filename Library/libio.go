package Library

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"time"
)

type Database interface {
	Connect() error
	LoadBook(*Book, string) error
	LoadBooks(*Books) error
	UploadBook(Book) error
	UploadBooks(Books) error
	Disconnect()
}

const postgresConnectStr = "host=10.0.0.6 user=postgres password=1 dbname=godev sslmode=disable"
const fileNameStr = "E:\\Dev\\goLibryary\\library.json"

const dateRFClex = "02.01.2006"

type Postgres struct {
	DB *sqlx.DB
	IsConnected bool
}

func (pq *Postgres) Connect() (err error) {
	pq.DB, err = sqlx.Connect("postgres", postgresConnectStr)
	if err == nil {
		pq.IsConnected = true
	}
	return err
}

func (pq *Postgres) LoadBook(book *Book, name string) (err error) {
	if !pq.IsConnected {
		return err
	}
	//TODO: some code to select from db
	return err
}

func (pq *Postgres) LoadBooks(books *Books) (err error) {
	//TODO: some code to select from db
	return err
}

func (pq *Postgres) Disconnect() {
	err := pq.DB.Close()
	if err != nil {
		panic(err)
	}
}

//FileIO read an info from file. If file was already readed (isFileCached)
//file may be contained(помещён) by cacher to be cached (fileCached)
/*
*	Файловый ввод-вывод реализующий интерфейс Database. Поддерживаемые функции:
*	LoadBook	-	загружает одну книгу с заданным параметром из json файла
*	LoadBooks	-	загружает все книги из json файла
*	UploadBook	-	считывает все книги из json файла, добавляет книгу, и загружает весь список в json файл
*	UploadBooks	-	считывает все книги из json файла, добавляет книги, и загружает весь список в json файл
 */
type FileIO struct {
	FileCached map[string]*FileCached
}


//Структура кэширования
	//Файл существует. Сомнительное поле, существование структуры уже подразумевает существование кэшированного файла
	//Файл кэширован. Сомнительное поле, существование стурктуры уже подразумевает кэширование файла.
	//Кэшированный файл в байтах.
	//Время кэширования. Если при доступе к кэшу время не совпадает со временем изменения файла, кэш обновляется.
type FileCached struct {
	IsExist bool
	IsCached bool
	Cache []byte
	CachingTimeStamp time.Time
}

//Func caching file
func (fc *FileCached) cacheFile(file []byte, timestamp time.Time) {
	if fc.IsCached == false {
		fc.Cache = file
		fc.IsCached = true
		fc.CachingTimeStamp = timestamp
	}
}

func (fc *FileCached) verifyTimeStamp(timestamp *time.Time) bool {
	return fc.CachingTimeStamp.Year() == timestamp.Year() &&
		fc.CachingTimeStamp.Month() == timestamp.Month() &&
		fc.CachingTimeStamp.Day() == timestamp.Day() &&
		fc.CachingTimeStamp.Hour() == timestamp.Hour() &&
		fc.CachingTimeStamp.Minute() == timestamp.Minute() &&
		fc.CachingTimeStamp.Second() == timestamp.Second()
}

func (fio *FileIO) freeFile(name string) {
	delete(fio.FileCached, name)
}

func (fio *FileIO) Disconnect() {
	for name := range fio.FileCached {
		delete(fio.FileCached, name)
	}
}

/*
Функция сверяет время последнего кэширования и время изменения файла.
Если время не сходится, файл кэшируется снова.
 */
/*
	Если это первый проход и кэш ещё не создан, создаём пустой
	Проверяем наличие нужного кэша
	Смотрим время изменения файла, чтобы узнать актуален ли кэш
		Если требуемый кэш отсутствует или не актуален, читаем файл заного
			Обновляем кэш
		Если кэш не требует обновления, возвращаем содержимое кэша
	Возвращаем содержимое файла
*/
func (fio *FileIO) getFileContain(name string) ([]byte, error) {
	//if file hasn't been cached yet, read and cache it.
	if fio.FileCached == nil {
		fio.FileCached = make(map[string]*FileCached)
	}
	_, isCacheExist := fio.FileCached[name]
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
	}()

	//if file has been updated after cache, read it again
	fileInfo, _ := file.Stat()
	timeModify := fileInfo.ModTime()	//Timestamp хранит время последнего изменения данных
	if !isCacheExist || !fio.FileCached[name].verifyTimeStamp(&timeModify) {
		fio.FileCached[name] = &FileCached{}

		content, err := ioutil.ReadFile(name)
		if err != nil {
			return nil, err
		}

		/*buf := make([]byte, MaxBaseSize)
		_, err = file.Read(buf)
		if err != nil {
			return nil, err
		}

		buf = bytes.Trim(buf, "\x00")*/
		fio.FileCached[name].cacheFile(content, timeModify)
		return content, nil
	}

	return fio.FileCached[name].Cache, nil
}

func (fio *FileIO) Connect() (err error) {
	_, err = fio.getFileContain(fileNameStr)
	if err != nil {
		return err
	}
	return nil
}

// fetch
func (fio *FileIO) LoadBook(book *Book, name string) (err error) {
	fileContent, err := fio.getFileContain(fileNameStr)
	if err != nil {
		return err
	}
	var bookList []BookJSON
	err = json.Unmarshal(fileContent, &bookList)
	if err != nil {
		return err
	}
	for _, b := range bookList {
		if b.Name == name {
			loadTime, err := time.Parse(dateRFClex, b.AddedAt)
			if err != nil {
				return err
			}
			bk := Book{b.ID, b.Name, b.Author, b.Publisher, b.Kind, b.Size, b.HoldingShelf, loadTime}
			*book = bk
			return nil
		}
	}
	*book = Book{ID:-1, Name:"empty", Author:"empty", Publisher:"empty", Kind:"empty", Size:"empty", HoldingShelf:"empty"}
	return err
}

//TODO: check situation when json file exists but empty
func (fio *FileIO) LoadBooks(books *Books) error {
	fileContent, err := fio.getFileContain(fileNameStr)

	var loadedByJSON []BookJSON
	err = json.Unmarshal(fileContent, &loadedByJSON)
	if err != nil {
		return err
	}

	var book Book
	for _, b := range loadedByJSON {
		if book, err = b.convertFromBookJSON(); err != nil {
			return err
		}
		*books = append(*books, book)
	}
	return nil
}

//Добавление книги в файл
func (fio *FileIO) UploadBook(book Book) error {
	var summary Books
	if err := fio.LoadBooks(&summary); err != nil {
		return err
	}
	summary.AddBook(book)

	var convertedList []BookJSON
	var currentConvert BookJSON
	for i := range summary {
		currentConvert.convertFromBook(summary[i])
		currentConvert.ID = i
		convertedList = append(convertedList, currentConvert)
	}

	content, err := json.MarshalIndent(convertedList, "", "    ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fileNameStr, content, 0666); err != nil {
		return err
	}
	return nil
}

//Запись списка книг в файл. Поступающий список должен быть отсортирован.
func (fio *FileIO) UploadBooks(new Books) error {
	var summaryList Books
	if err := fio.LoadBooks(&summaryList); err != nil {
		return err
	}
	for _, currentNewBook := range new {
		summaryList = append(summaryList, currentNewBook)
	}

	var convertedList []BookJSON
	var currentConvert BookJSON
	for i := range summaryList {
		currentConvert.convertFromBook(summaryList[i])
		currentConvert.ID = i
		convertedList = append(convertedList, currentConvert)
	}

	content, err := json.MarshalIndent(convertedList, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileNameStr, content, 0666)
	if err != nil {
		return err
	}
	return nil
}