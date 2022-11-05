package parse

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetHtml(url string) *http.Response {
	res, err := http.Get(url)
	CheckError(err)

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res
}

func FindFaculty(url string, doc *goquery.Document) map[string]string {
	faculties := make(map[string]string)
	doc.Find("ul.faculty-list__list>li.faculty-list__item").Each(func(i int, item *goquery.Selection) {
		a := item.Find("a.faculty-list__link")
		title := strings.TrimSpace(a.Text())
		link, _ := a.Attr("href")
		faculties[title] = url + link
	})
	return faculties
}

func FindGroup(url string, doc *goquery.Document) map[string]string {
	groups := make(map[string]string)
	doc.Find("ul.groups-list>li.groups-list__item").Each(func(i int, item *goquery.Selection) {
		a := item.Find("a.groups-list__link")
		title := strings.TrimSpace(a.Text())
		link, _ := a.Attr("href")
		groups[title] = url + link
	})
	return groups
}

type TimeTable []TimeTableDay

type TimeTableDay struct {
	Date    string
	Lessons []Lesson
}

type Lesson struct {
	Title, LessonType, Teacher, Place string
}

func GetTimeTable(doc *goquery.Document) TimeTable {
	res := make(TimeTable, 0)
	doc.Find("ul.schedule>li.schedule__day").Each(func(i int, day *goquery.Selection) {
		table := TimeTableDay{}
		table.Date = day.Find("div.schedule__date").Text()

		day.Find("ul.schedule__lessons>li.lesson").Each(func(i int, lesson *goquery.Selection) {
			lessonStruct := Lesson{}
			lessonStruct.Title = lesson.Find("div.lesson__subject").Text()
			lessonStruct.LessonType = lesson.Find("div.lesson__type").Text()
			lessonStruct.Teacher = lesson.Find("div.lesson__teachers").Text()
			lessonStruct.Place = lesson.Find("div.lesson__places").Text()
			table.Lessons = append(table.Lessons, lessonStruct)
		})
		res = append(res, table)
	})
	return res
}
