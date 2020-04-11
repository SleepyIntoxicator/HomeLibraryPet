package main

import (
	"bufio"
	"fmt"
	lib "goLibryary/Library"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
	Приложение для удобного и быстрого добавления информации о расположении книг в полках.
	Информация о книге передаётся через консоль.
	Доступны операции сортировки, добавления, изменения и обновления информации о книге.
	Используемые технологии:
	GO + Postgres: sqlx,
	Умения для собеседований:
	Создание интерфейса командной строки
	CRUD с базой данных
	Сортировки и структуры данных

	Данные: Книга
	Книга {
		Номер
		Название
		Издательсво
		Тип
		Размер
		Полка №
	}
*/

/*
Виды транспортных документов
На ЖД и авто транспорте - накладная 
На морском - коносамент, чартер, накладная 
На воздушном трнасопрте - чартер,нкаладная 
При смешанных перевозках - смешанный коносамент 

Группы перевозочных документов
- Орагнизационно-распорядительные - инициализирующие перевозку (заявка, заявлени, распоряжение), от них зависит маршрут перевозки.
- Перевозочные - в этих сведениях содержится информация о грузе, об операциях и на смежные документы и распоряжения
- Сопроводительные - используется для дополнительных операций с грузом и о факте их выполнения ( ветеринарный, таможенный и санитарный контроль )
- Акты - акты бывают: общие, о техническом состоянии вагона, о вскрытии вагона, повреждение, коммерческие акты.
- Служебные - книги регистрации извещений, в которых грузополучатель извещается о подаче вагона под погрузку и прибытии груза.
Классификация криптографических протоколов
	- протоколы шифрования\дешифрования
		- Симметричный\асиметричный алгоритм шифрования. Алгоритм выполняется при передаче отправителем сообщения,
			в результате чего, оно превращается в закрытую форму. Т. о. обеспечивается св-во конфиденциальности
		- В целях сохранения целостности обычно совмещаются с алгоритмом вычисления эмитозащитной вставки (ключ шифрования)
		- При использовании асиметричных алг. Ц-сть обеспечивается путём вычисления ЭЦП. Т.о. обеспечивается С-во безотказности и аутентичности сообщения.
	-
		- В основе лежит вычисления ЭЦП с помощью секретного ключа отправителя и процерки ЭЦП с помощью открытого ключа получателя.
			. Открыты ключь берётся из открытого справочника, защищённого от модификации. В случае положительного результата сообзение архивируется.
			. включая ЭЦП и открытый ключь. Операция архиаированич иоэнт не выполняться, если ЭЦП используется только для обеспеения целостности и аутентичности (но не безотказности).
			.В д.с. ЭЦП может быть уничтожено.
	-
		- В основе данной группы протоколов лежит алгоритм проверки того факта, что идентифицируемый объект (пользователь, устр-во или процесс)
		. предъявивший некоторое имя хнает секретную информацию, известную только заявленному объекту. С идентификатором обычно связывают права и полномочия в системе, записынные в защищённой БД.
		. Протокол идентификации может быть расширен до протокола аутентификации, в котором осуществялется проверка провомощности заказываемой услуги.
		. Если в протоколе идентификации используется ЭЦП, то роль секретной информации выполняет секретный ключь ЭЦП.
		. Проверка выполнется в кач-ве открытого ключа. Знание открытого ключа не позволяет определить секретный ключь.
		. Но помогает понять, что он известен отправитель ЭЦП.
	-
		- Данные протоколы совмещают протокоолы аутентификации пользователей с протоколами генерации и распределения ключей. П имеет 2-3 участников. 3-м участником
		. является центр генерации и распределения ключей. Работа протоколов состоит из 3 этапов: генерация, ренистрация и коммуникация.
		. На этапе ЦРК генерирует .. ключи системы: откр и закр. На этапе регистрации: ЦГРК идентифицирует пользователей по документу и генерирует
		. идентифицирующую инфу (в т.ч. маркер безопасности, содержащий системные константы, открытый ключь ЦГРК)
		. На этапе .. генерирутеся .. аутентифицированый ключевой обмен, который заканчивается генерацией .. .
 */

func main() {
	defer recoverAll()
	CLIMain()
}

func recoverAll() {
	if err := recover(); err != nil {
		log.Println("recovered panic:", err)
		_, _ = fmt.Scan()
	}
}

func CLIMain() {
	input := bufio.NewScanner(os.Stdin)
	var strIn string
	var args []string
	var countArgs int

	var db lib.Database
	db = &lib.FileIO{}
	if err := db.Connect(); err != nil {
		log.Fatal("cannot connect to db in main. err:", err)
	}

	for {
		fmt.Print("cmd: ")
		input.Scan()
		if err := input.Err(); err != nil {
			log.Fatal("err with scanning cmd input", err)
		}
		strIn = input.Text()

		args = strings.Split(strIn, " ")
		countArgs = len(args)

		switch args[0] {
		case "exit":
			break
		case "get", "получить":
			if countArgs == 1 {
				fmt.Println("get [-all] [-range [id] --first --last] [-count]")
				fmt.Println("получить [все] [диапазон [ид] (первый) (последний)] [количество]")
				break

			} else
			if args[1] == "all" || args[1] == "все" {
				if countArgs < 2 {
					fmt.Println("arg [2] missed: ")
					break
				}
				var allBooks lib.Books
				if err := db.LoadBooks(&allBooks); err != nil {
					log.Fatal("cannot load books list. err:", err)
				}

				fmt.Println(allBooks.GetBooksTableWithTitle("full books list"))

			} else
			if args[1] == "range" || args[1] == "диапазон" {
				if countArgs < 4 {
					fmt.Println("arg [2-4, 5] missed")
					break
				}

				if countArgs == 4 {
					leftIndex, err := strconv.Atoi(args[2])
					if err != nil {
						log.Println("error in converting from args. err:", err)
						break
					}
					rightIndex, err := strconv.Atoi(args[3])
					if err != nil {
						log.Println("error in converting from args. err:", err)
						break
					}

					var allBooks lib.Books
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}
					if leftIndex < 0 {
						log.Println("Err: left index cannot be less than 0")
						break
					}
					if leftIndex > len(allBooks) {
						log.Println("Err: left index greater than the length of its list. Max index is", len(allBooks))
						break
					}
					if leftIndex < len(allBooks) && rightIndex+1 > len(allBooks) {
						rightIndex = len(allBooks)
					}
					rangeBooks := allBooks[leftIndex-1 : rightIndex+1]
					fmt.Println(rangeBooks.GetBooksTableWithTitle("Books from " + args[2] + " to " + args[3]))

				}
				if countArgs == 5 && args[2] == "id" {
					var allBooks lib.Books
					var err error
					//Parsing ranges
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}
					var leftIndex, rightIndex int
					if leftIndex, err = strconv.Atoi(args[3]); err != nil {
						log.Println("error in converting from args. err:", err)
						break
					}
					if rightIndex, err = strconv.Atoi(args[4]); err != nil {
						log.Println("error in converting from args. err:", err)
						break
					}

					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}

					countBooks := len(allBooks)
					if leftIndex < 0 {
						log.Println("Err: left index cannot be less than 0")
						break
					}
					if leftIndex > countBooks {
						log.Println("Err: left index greater than the length of its list. Max index is", len(allBooks))
						break
					}
					if rightIndex > countBooks {
						rightIndex = countBooks
					}

					founded, err := allBooks.FindBooksWithParams("id", [2]int{leftIndex, rightIndex}, "-range")
					if err != nil {
						log.Println("Cannot find the books with такими parameters. err:", err)
					}
					fmt.Println(founded.GetBooksTableWithTitle("Books from id " + args[3] + " to " + args[4]))
				}
			} else
			if args[1] == "count" || args[1] == "количество" {
				var allBooks lib.Books
				if err := db.LoadBooks(&allBooks); err != nil {
					log.Fatal("Cannot load books list. err: ", err)
				}
				fmt.Println("The list contains " + strconv.Itoa(len(allBooks)) + " book.")
			} else
			if args[1] == "books" || args[1] == "книги" {
				if countArgs == 2 {
					var allBooks lib.Books
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}

					fmt.Println(allBooks.GetBooksTableWithTitle("full books list"))
					break
				} else if countArgs < 4 {
					fmt.Println("arg [2-4] missed: ")
					break
				} else
				if countArgs >= 5 && args[2] == "sorted" && args[3] == "by" {
					sortingParam := args[4]
					sortingFlag := ""
					if countArgs == 6 {
						sortingFlag = args[5]
					}
					var allBooks lib.Books
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}
					err := allBooks.SortBooksWithParams(sortingParam, sortingFlag)
					if err != nil {
						log.Print(err)
					}
					fmt.Print(allBooks.GetBooksTableWithTitle("Sorted list by " + sortingParam))
				} else
				if args[2] == "by" || args[2] == "по" {
					paramName := args[3]
					var paramValue interface{}
					var err error

					/*
					 * Если параметров ровно 4, значит пользователь выбрал полный вариант ввода
					 * Если параметров больше или равно 5, значит пользователь выбрал короткий вариант ввода
					 */
					if countArgs == 4 {
						fmt.Print("\t", paramName, ": ")
						input.Scan()
						if err := input.Err(); err != nil {
							log.Fatal("err with scanning cmd input on param value", err)
						}
						paramValue = input.Text()
					} else
					if countArgs >= 5 {
						if args[3] == "id" || args[3] == "ид" {
							paramValue, err = strconv.Atoi(args[4])
							if err != nil {
								log.Println("id parameter must be an int. err:", err)
								break
							}
							paramName = "id"
						} else {
							startQuotes := IndexRune(strIn, "\"", strings.Index)
							endQuotes := IndexRune(strIn, "\"", strings.LastIndex)
							if startQuotes < 0 || (startQuotes == endQuotes) {
								valueIndex := IndexRune(strIn, paramName, strings.Index)
								valueIndex += utf8.RuneCountInString(paramName) + 1
								paramValue = string([]rune(strIn)[valueIndex:])
							} else {
								paramValue = string([]rune(strIn)[startQuotes+1 : endQuotes])
							}
						}
					}

					var allBooks lib.Books
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}
					founded, err := allBooks.FindBooksWithParams(paramName, paramValue, "") //TODO: think about flags
					if err != nil {
						log.Fatal("Cannot find the books by this param. err:", err)
					}
					switch paramValue.(type) {
					case int:
						fmt.Println(founded.GetBooksTableWithTitle("Founded books by " + paramName + " " + strconv.Itoa(paramValue.(int))))
					case string:
						fmt.Println(founded.GetBooksTableWithTitle("Founded books by " + paramName + " " + paramValue.(string)))
					default:
						fmt.Println(founded.GetBooksTableWithTitle("Founded books by " + paramName))
					}
				}
			} else
			if args[1] == "book" {
				if countArgs < 3 {
					fmt.Println("arg [3] missed: ")
					break
				}

				if args[2] == "by" {
					paramName := args[3]
					fmt.Print("\t", paramName, ":")
					input.Scan()
					if err := input.Err(); err != nil {
						log.Fatal("err with scanning cmd input on param value", err)
					}

					var paramValue interface{}
					var err error
					if paramName == "id" {
						if paramValue, err = strconv.Atoi(input.Text()); err != nil {
							log.Println("id parameter must be an int. err:", err)
							break
						}
					}

					var allBooks lib.Books
					if err := db.LoadBooks(&allBooks); err != nil {
						log.Fatal("cannot load books list. err:", err)
					}
					founded, err := allBooks.FindBooksWithParams(paramName, paramValue, "") //TODO: flags
					if err != nil {
						log.Fatal("todo text error. err:", err)
					}
					fmt.Println(founded.GetBooksTableWithTitle("Founded books by " + paramName))
				}
			}
		case "add":
			if args[1] == "book" {
				if countArgs > 2 {
					if args[2] == "full" {
						paramNames := []string{
							"name",
							"author",
							"publisher",
							"kind",
							"size",
							"holdingShelf"}
						permissionsEmpty := []bool{
							false, false, false, true, true, true,
						}
						params := map[string]string{}
						//TODO: доделать параметры, с ограниченными\рекомендуемыми значениями
						for i, name := range paramNames {
							for {	//input current param, until getting true args
								fmt.Print("\t", name, ": ")
								input.Scan()
								if err := input.Err(); err != nil {
									log.Println("cannot read param value, when adding new book")
									panic(err)
								}
								paramValue := input.Text()
								//Если параметр пропущен, уточняем
								if paramValue == "" {
									if !permissionsEmpty[i] {
										fmt.Println("\t\tThis parameter cannot be skipped. Try again.")
										continue
									} else {
										fmt.Print("\t\tSkip param (y\\n): ")
										input.Scan()
										answer := input.Text()
										if answer == "y" || answer == "skip" {
											fmt.Printf("\t\tBook %s skipped.\n", name)
											params[name] = ""
											break
										}
										fmt.Println() //Вводим параметр заного
										continue
									}
								} else {
									params[name] = paramValue
									break
								}
							}
						}

						newBook, err := lib.NewBook(params["name"], params["author"], params["publisher"],
							params["kind"], params["size"], params["holdingShelf"])
						if err != nil {
							log.Print("Error when adding new book. err: ", err)
							break
						}
						err = db.UploadBook(*newBook)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			break
		}

		if strIn == "exit" || strIn == "выйти" || strIn == "выход"{
			break
		}
	}
}

func OldCLIMain() {
	var input_old string
	for {
		fmt.Print("cmd: ")
		if _, err := fmt.Scanln(&input_old); err != nil {
			fmt.Print(" arg missed\n")
			log.Println("err:\n" + err.Error())
			continue
		}
		if args := strings.Split(input_old, "\""); args[0] == "exit" {
			break;
		} else if args[0] == "get" {
			fmt.Print("\tbooks title: ")
			var title string
			if _, err := fmt.Scanln(&title); err != nil {
				log.Println(err)
				continue
			}
			if args := strings.Split(title, "\""); args[0] == "" {
				title = "All books from file bs.json"
			} else {
				title = args[0]
			}
			var books lib.Books
			var db lib.Database
			db = &lib.FileIO{}
			if err := db.Connect(); err != nil {
				panic(err)
			}
			if err := db.LoadBooks(&books); err != nil {
				panic(err)
			}
			fmt.Println(books.GetBooksTableWithTitle(title))
		} else if args[0] == "load" {
			var book *lib.Book
			book, err := lib.NewBook("added_name","added_author",
				"added_publisher","added_kind","added_size",
				"added_shelf")
			if err != nil {
				panic(err)
			}

			var db lib.Database
			db = &lib.FileIO{}
			if err := db.Connect(); err != nil {
				panic(err)
			}
			if err := db.UploadBook(*book); err != nil {
				panic(err)
			}
		} else if args[0] == "loads" {
			var book *lib.Book
			book, err := lib.NewBook("added_name","added_author",
				"added_publisher","added_kind","added_size",
				"added_shelf")
			if err != nil {
				panic(err)
			}

			var books lib.Books
			for i := 0; i < 10; i++ {
				book.Name = "added_name_" + strconv.Itoa(i)
				books = append(books, *book)
			}

			var db lib.Database
			db = &lib.FileIO{}
			if err := db.Connect(); err != nil {
				panic(err)
			}
			if err := db.UploadBooks(books); err != nil {
				panic(err)
			}
		}
	}
}

func ScenarioLoadMain() {
		var db lib.Database
		db = &lib.FileIO{}

		err := db.Connect()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Disconnect()

		var book lib.Book
		err = db.LoadBook(&book, "Совершенный код")
		if err != nil {
			panic(err)
		}

		fmt.Printf("The book: [%d]  %s  of  %s\n", book.ID, book.Name, book.Author)
		err = db.LoadBook(&book, "Современные операционные системы 4-е издание")
		if err != nil {
			panic(err)
		}

		fmt.Printf("The book: [%d]  %s  of  %s\n", book.ID, book.Name, book.Author)
		//----------------------------------------------------------------------------------------
		//-----------------------------[FindBooksWithParams]---------------------------------------
		fmt.Print("\n\n\n")
		fmt.Println(strings.Repeat("-", 179))
		fmt.Println("---------------------------------[Find books with params - name=Объектно flags=\"\"]----------------------------------------------------------------------------------------------------------------")

		var library lib.Books
		err = db.LoadBooks(&library)
		if err != nil {
			log.Fatal(err)
		}
		var libs *lib.Books
		libs, err = library.FindBooksWithParams("name", "Объектно", "")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(libs.GetBooksTableWithTitle("Find books with params - name=Объектно flags=\"\""))
		//----------------------------------------------------------------------------------------
		//-----------------------------[FindBooksWithParams]---------------------------------------
		fmt.Print("\n\n\n")
		fmt.Println(strings.Repeat("-", 179))
		fmt.Println("---------------------------------[Find books with params - date=now flags=\"-n\"]----------------------------------------------------------------------------------------------------------------")

		libs, err = library.FindBooksWithParams("date", "02.03.2020", "-a")
		if err != nil {
			panic(err)
		}
		fmt.Print(libs.GetBooksTable())
		//----------------------------------------------------------------------------------------
		//-----------------------------[SortBooksWithParams]---------------------------------------
		fmt.Print("\n\n\n")
		fmt.Println(strings.Repeat("-", 179))
		fmt.Println("---------------------------------[Sort books with params - date flag=\"\"]----------------------------------------------------------------------------------------------------------------")

		err = libs.SortBooksWithParams("date", "")
		if err != nil {
			panic(err)
		}
		fmt.Print(libs.GetBooksTable())
		//----------------------------------------------------------------------------------------
		//----------------------------------------------------------------------------------------
		fmt.Print("\n\n\n")
		fmt.Println(strings.Repeat("-", 179))

		fmt.Print(library.GetBooksTable())
}


func IndexRune(str, substr string, indexFunc func(s, ss string) int) int {
	byteIndex := indexFunc(str, substr)
	if byteIndex == -1 {
		return -1
	}
	return utf8.RuneCountInString(str[:byteIndex])
}