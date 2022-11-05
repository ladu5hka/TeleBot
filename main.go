package main

import (
	prs "GolangProjects/Parse"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://ruz.spbstu.ru"
	res := prs.GetHtml(url)
	defer res.Body.Close()

	docFaculty, err := goquery.NewDocumentFromReader(res.Body)
	prs.CheckError(err)

	faculties := prs.FindFaculty(url, docFaculty)
	facUrl := faculties["Институт промышленного менеджмента, экономики и торговли"]

	res = prs.GetHtml(facUrl)
	defer res.Body.Close()

	docGroups, err := goquery.NewDocumentFromReader(res.Body)
	prs.CheckError(err)

	groups := prs.FindGroup(url, docGroups)
	groupUrl := groups["3733801/00401"]

	res = prs.GetHtml(groupUrl)
	defer res.Body.Close()

	docTimeTable, err := goquery.NewDocumentFromReader(res.Body)
	prs.CheckError(err)

	fmt.Println(prs.GetTimeTable(docTimeTable))
}
