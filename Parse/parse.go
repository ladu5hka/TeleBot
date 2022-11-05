package parse

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"sync"
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
		log.Fatal(err)
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
	groups := sync.Map{}
	wg := sync.WaitGroup{}
	for faculty := range Faculties {
		wg.Add(1)
		go func(_faculty string) {
			res := GetHtml(Faculties[_faculty])
			doc, err := goquery.NewDocumentFromReader(res.Body)
			CheckError(err)
			doc.Find("ul.groups-list>li.groups-list__item").Each(func(i int, item *goquery.Selection) {
				a := item.Find("a.groups-list__link")
				title := strings.TrimSpace(a.Text())
				link, _ := a.Attr("href")
				groups.Store(title, url+link)
			})
			res.Body.Close()
			wg.Done()
		}(faculty)
	}
	wg.Wait()
	groups.Range(func(key any, value any) bool {
		Groups[key.(string)] = value.(string)
		return true
	})
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
