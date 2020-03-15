package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	dirPath := os.Args[1]
	//проверка на указание папки
	if len(os.Args) == 1 {
		log.Fatal("No folder declared")
	}
	//читаем содержимое папки
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	//проверка есть ли в папке файлы
	if len(files) == 0 {
		log.Fatal("No files in specified directory")
	}
	//Создаем выходной файл
	file, err := os.Create("output.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	tokens, err := makeIndexFile(dirPath, files)
	if err != nil {
		log.Fatal(err)
	}
	//Пишем в выходной файл
	for key, value := range tokens {
		_, err = file.WriteString("'" + key + "'" + " >> {" + value + "},\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func makeIndexFile(dirPath string, files []os.FileInfo) (map[string]string, error) {
	tokens := map[string]string{}
	for i, file := range files {
		currentFile, err := ioutil.ReadFile(dirPath + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		re, err := regexp.Compile("[.,?!;:]")
		if err != nil {
			return nil, err
		}
		// Убираем всю пунктуацию
		premove := re.ReplaceAllString(string(currentFile), "")
		// Разбиваем на стринги по пробелу
		str := strings.Split(premove, " ")
	wordLoop:
		for j := range str {
			if str[j] == "" {
				continue wordLoop
			}
			value, ok := tokens[str[j]]
			//проверяем если есть данное значение в списке, добавляем номер файла. если нет, дописываем в файл значение и номер файла
			//эту часть принципиально не знал как реализовать. взято из репы другого студента. не могу понять что и зачем делает этот код "strconv.Itoa(int(value[len(value)-1])-48)"
			if ok && strconv.Itoa(int(value[len(value)-1])-48) != strconv.Itoa(i+1) {
				tokens[str[j]] = value + "," + strconv.Itoa(i+1)
			} else if !ok {
				tokens[str[j]] = strconv.Itoa(i + 1)
			}
		}
	}
	return tokens, nil
}
