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
	"sync"
)

type timeData struct {
	year  string
	month string
	day   string
}

type data struct {
	time timeData
	data string
}

func isNumber(n string) bool {
	_, err := strconv.Atoi(n)
	if err != nil {
		return false
	}
	return true
}

func setAllTimeData(body string) []timeData {
	timeArray := make([]timeData, 0)

	for _, line := range strings.Split(strings.TrimRight(body, "\n"), "\n") {
		if len(line) < 74 {
			continue
		}
		if isNumber(line[80:84]) == false || isNumber(line[85:87]) == false || isNumber(line[88:90]) == false {
			continue
		}

		tmpTime := timeData{year: line[80:84], month: line[85:87], day: line[88:90]}
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

func getDetailDataByData(tmp string) (md5, sha1, sha256, sha string) {
	words := strings.Fields(tmp)
	for i := 0; i < len(words); i++ {
		if i%4 == 0 {
			md5 += words[i] + "\n"
			continue
		}
		if i%4 == 1 {
			sha1 += words[i] + "\n"
			continue
		}
		if i%4 == 2 {
			sha256 += words[i] + "\n"
			continue
		}
		if i%4 == 3 {
			sha += words[i] + "\n"
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

func writeFile(path string, text string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	textWriter := bufio.NewWriter(file)

	_, err = textWriter.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
	textWriter.Flush()
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

func createFileAndWrite(path string, text string) {
	createFile(path)
	writeFile(path, text)
}

func getResult(array []data, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(array); i++ {
		md5, sha1, sha256, sha := getDetailDataByData(array[i].data)
		yearPath := "data/" + array[i].time.year
		monthPath := yearPath + "/" + array[i].time.month
		dayPath := monthPath + "/" + array[i].time.day
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

func getDataByTime(timeArrayData []timeData, tmp *[]data, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(timeArrayData); i++ {
		timeString := timeArrayData[i].year + "-" + timeArrayData[i].month + "-" + timeArrayData[i].day
		url := "https://malshare.com/daily/" + timeString + "/malshare_fileList." + timeString + ".all.txt"
		tmpData := data{timeArrayData[i], fetchData(url)}
		fmt.Println(i)
		*tmp = append(*tmp, tmpData)
	}
}

func main() {
	var wg sync.WaitGroup
	timeStringData := fetchData("https://malshare.com/daily/")
	timeArrayData := setAllTimeData(timeStringData)
	createFolder("data")
	dataArray := make([]data, 0)
	step := len(timeArrayData) / 100
	wg.Add(100)
	for i := 0; i < 99; i++ {
		go getDataByTime(timeArrayData[step*i:step*i+step], &dataArray, &wg)
	}
	go getDataByTime(timeArrayData[step*99:len(timeArrayData)], &dataArray, &wg)

	wg.Wait()
	step = len(dataArray) / 100
	wg.Add(100)
	for i := 0; i < 99; i++ {
		go getResult(dataArray[step*i:step*i+step], &wg)
	}
	go getResult(dataArray[step*99:len(timeArrayData)], &wg)
	wg.Wait()

	fmt.Println("success")
}

//func indexHrefInLine(line, target string) int {
//	return strings.Index(line, target)
//}
