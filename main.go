package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type time struct {
	year  string
	month string
	day   string
}

func isNumber(n string) bool {
	_, err := strconv.Atoi(n)
	if err != nil {
		return false
	}
	return true
}

func setAllTimeData(body string) []time {
	timeArray := make([]time, 0)

	for _, line := range strings.Split(strings.TrimRight(body, "\n"), "\n") {
		if len(line) < 74 {
			continue
		}
		if isNumber(line[80:84]) == false || isNumber(line[85:87]) == false || isNumber(line[88:90]) == false {
			continue
		}

		tmpTime := time{year: line[80:84], month: line[85:87], day: line[88:90]}
		timeArray = append(timeArray, tmpTime)
	}

	return timeArray
}

func fetchData(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}

func setDetailDataInTime(time time) (md5, sha1, sha256, sha []string) {
	timeString := time.year + "-" + time.month + "-" + time.day
	url := "https://malshare.com/daily/" + timeString + "/malshare_fileList." + timeString + ".all.txt"
	words := strings.Fields(fetchData(url))
	for i := 0; i < len(words); i++ {
		if i%4 == 0 {
			md5 = append(md5, words[i])
			continue
		}
		if i%4 == 1 {
			sha1 = append(sha1, words[i])
			continue
		}
		if i%4 == 2 {
			sha256 = append(sha256, words[i])
			continue
		}
		if i%4 == 3 {
			sha = append(sha, words[i])
			continue
		}
	}
	return
}

func readFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatal(err)
	}
}

func writeFile(path string, text []string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	textWriter := bufio.NewWriter(file)

	for i := 0; i < len(text); i++ {
		value := text[i] + "\n"
		_, err = textWriter.WriteString(value)
		if err != nil {
			log.Fatal(err)
		}
		textWriter.Flush()
	}
	fmt.Println("Data written to file successfully...")
}

func createFile(path string) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	f.Close()
}

func createFolder(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return
	}
}

func createFileAndWrite(path string, text []string) {
	createFile(path)
	writeFile(path, text)
}

func main() {
	timeStringData := fetchData("https://malshare.com/daily/")
	timeArrayData := setAllTimeData(timeStringData)
	createFolder("data")
	for i := 0; i < len(timeArrayData); i++ {
		md5, sha1, sha256, sha := setDetailDataInTime(timeArrayData[i])
		yearPath := "data/" + timeArrayData[i].year
		monthPath := yearPath + "/" + timeArrayData[i].month
		dayPath := monthPath + "/" + timeArrayData[i].day
		filePath := dayPath + "/"
		createFolder(yearPath)
		createFolder(monthPath)
		createFolder(dayPath)
		createFileAndWrite(filePath+"md5.txt", md5)
		createFileAndWrite(filePath+"sha1.txt", sha1)
		createFileAndWrite(filePath+"sha256.txt", sha256)
		createFileAndWrite(filePath+"sha.txt", sha)
	}

}

//func indexHrefInLine(line, target string) int {
//	return strings.Index(line, target)
//}
