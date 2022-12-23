package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"helloworld-api/database"
	"helloworld-api/tables"
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

func getDetailDataByData(tmp string) (md5, sha1, sha256, sha []string) {
	words := strings.Fields(tmp)
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
		time := array[i].time.year + "-" + array[i].time.month + "-" + array[i].time.day
		insertDataToDB(time, "md5", md5)
		insertDataToDB(time, "sha1", sha1)
		insertDataToDB(time, "sha256", sha256)
		insertDataToDB(time, "sha", sha)
		fmt.Println("Adding data to database")
	}
}

func insertDataToDB(time string, tmpType string, data []string) {
	ctx := context.Background()
	for i := 0; i < len(data); i++ {
		tmpData := tables.Data{Time: time, Data: data[i], Type: tmpType}
		db.NewInsert().
			Model(&tmpData).
			Exec(ctx)
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

var db *bun.DB

func Crawl() {
	db = database.ConnectDatabase()
	var wg sync.WaitGroup
	timeStringData := fetchData("https://malshare.com/daily/")
	timeArrayData := setAllTimeData(timeStringData)
	createFolder("data")
	dataArray := make([]data, 0)
	step := len(timeArrayData[0:210]) / 20
	wg.Add(20)
	for i := 0; i < 19; i++ {
		go getDataByTime(timeArrayData[step*i:step*i+step], &dataArray, &wg)
	}
	go getDataByTime(timeArrayData[step*19:len(timeArrayData[0:210])], &dataArray, &wg)

	wg.Wait()

	step = len(dataArray) / 20
	wg.Add(20)
	for i := 0; i < 19; i++ {
		go getResult(dataArray[step*i:step*i+step], &wg)
	}
	go getResult(dataArray[step*19:len(timeArrayData[0:210])], &wg)
	wg.Wait()
}

//func indexHrefInLine(line, target string) int {
//	return strings.Index(line, target)
//}
