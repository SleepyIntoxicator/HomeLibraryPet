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

func IndexRune(str, substr string, indexFunc func(s, ss string) int) int {
	byteIndex := indexFunc(str, substr)
	if byteIndex == -1 {
		return -1
	}
	return utf8.RuneCountInString(str[:byteIndex])
}