package parse

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

var Faculties map[string]string
var Groups map[string]string

const url = "https://ruz.spbstu.ru"

func Start() {
	findFaculty()
	findGroup()
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func GetHtml(url string) *http.Response {
	res, err := http.Get(url)
	CheckError(err)

	if res.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res
}

func findFaculty() {
	Faculties = make(map[string]string)
	res := GetHtml(url)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckError(err)
	doc.Find("ul.faculty-list__list>li.faculty-list__item").Each(func(i int, item *goquery.Selection) {
		a := item.Find("a.faculty-list__link")
		title := strings.TrimSpace(a.Text())
		link, _ := a.Attr("href")
		Faculties[title] = url + link
	})
}

func findGroup() {
	Groups = make(map[string]string)

	ch := make(chan [2]string)
	defer close(ch)
	check := make(chan bool)
	defer close(check)

	for faculty := range Faculties {
		go func(_faculty string) {
			res := GetHtml(_faculty)
			doc, err := goquery.NewDocumentFromReader(res.Body)
			CheckError(err)
			doc.Find("ul.groups-list>li.groups-list__item").Each(func(i int, item *goquery.Selection) {
				a := item.Find("a.groups-list__link")
				title := strings.TrimSpace(a.Text())
				link, _ := a.Attr("href")
				ch <- [2]string{title, url + link}
			})
			res.Body.Close()
			check <- true
		}(Faculties[faculty])
	}

	for i := 0; i < len(Faculties); {
		select {
		case group := <-ch:
			Groups[group[0]] = group[1]
		case <-check:
			i++
		}
	}
}

type TimeTable []TimeTableDay

type TimeTableDay struct {
	Date    string
	Lessons []Lesson
}

type Lesson struct {
	Title, LessonType, Teacher, Place string
}

//func GetTimeTable() TimeTable {
//	res := make(TimeTable, 0)
//	doc.Find("ul.schedule>li.schedule__day").Each(func(i int, day *goquery.Selection) {
//		table := TimeTableDay{}
//		table.Date = day.Find("div.schedule__date").Text()
//
//		day.Find("ul.schedule__lessons>li.lesson").Each(func(i int, lesson *goquery.Selection) {
//			lessonStruct := Lesson{}
//			lessonStruct.Title = lesson.Find("div.lesson__subject").Text()
//			lessonStruct.LessonType = lesson.Find("div.lesson__type").Text()
//			lessonStruct.Teacher = lesson.Find("div.lesson__teachers").Text()
//			lessonStruct.Place = lesson.Find("div.lesson__places").Text()
//			table.Lessons = append(table.Lessons, lessonStruct)
//		})
//		res = append(res, table)
//	})
//	return res
//}
