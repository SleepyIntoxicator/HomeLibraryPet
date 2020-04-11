package CLI

import (
	"bufio"
	"os"
	"sync"
)

/*
TODO:
	CLI - модуль для роутинга команд, поступающих из консоли.
	--------------------------------------------------------------------------------
	Примеры 'API':
		get													:	apiRoutingGET
		get		all											:	getAllBooks
		get		range				{first} {last}			:	getBooksInRange
		get		range	id			{first} {last}			:	getBooksInRangeByID
		get		count										:	getCountBooks
		get		books										:	getAllBooks
		get		books	by			{param}	{value...}		:	getBooksByParam
		get		books	sorting	by	{param}					:	getBooksSortingByParam
		get		book	by			{param}	{value...}		:	getBookByParam
		add		book	full								:	addBookFull
		add		book	short								:	addBookShort
		add		book	?		?	?		?				:
		change	books	by			{param}	{value...}		:	changeBooksWhereParam
		change	book	by		id	{value}					:	changeBookByID
		delete	book	by		id	{value}					:	deleteBookByID
		sort	books	by			{param}					:	sortBooksByParam
		create	backup	in					{valuePath}		:	createBackup
		start	list										:	changeState		:=	ExceptionMode
			get books						|	founded 253
			get books by {param} {value}	|	founded 101
			sort books by name	{flags}		|	sorted list sorted 101
			get books by {param} {value}	|	founded 14
			change book by id {value}		|	succ changed
			change books by {param} {value}	|	changed 14
		end		list
		exit
	--------------------------------------------------------------------------------
	Устройство модуля:
		Для начала регистрируются пути
		Представляются в виде дерева
		get
			all
			range
				_
				id
			count
			books
				by
	Пример использования
	cli.HandleFunc("get", apiRoutingGET)
	cli.HandleFunc("get all", getAllBooks)
	...
	cli.HandleFunc("get books by {paramName} {value}, getBooksByParam)
	_
	func getBooksByParam(r string) {

	}
 */

func getBooksByParam(r string) {

}

type Handler func(r string)

type ServeMux struct {
	mu		sync.RWMutex
	m		map[string]MuxEntry
	es		[]MuxEntry
}

type MuxEntry struct {
	h		Handler
	pattern	string
}

func NewServeMux() *ServeMux {return new(ServeMux)}

func main() {
	HandleFunc("abc", getBooksByParam)

	//Listen()
}

func HandleFunc(pattern string, handler Handler) {
}

func Listen(w *os.File, r *os.File) {
	inputScanner := bufio.NewScanner(os.Stdin)
	inputScanner.Scan()
	if err := inputScanner.Err(); err != nil {
		panic(err)
	}

}