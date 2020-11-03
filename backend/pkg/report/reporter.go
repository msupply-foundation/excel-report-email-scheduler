package report

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/simple-datasource-backend/pkg/api"
	"github.com/grafana/simple-datasource-backend/pkg/auth"
)

var CELL_REF_REG = regexp.MustCompile("([A-Za-z]+)|([0-9]+)")

// TODO: Write a function to generate column letter given an int
var COLUMN_LETTERS = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT"}

func intToCol(i int) string {
	i += 1
	result := ""
	for i > 0 {
		rem := i % 26
		if rem == 0 {
			result = "Z" + result
			i = i/26 - 1
		} else {
			result = fmt.Sprintf("%c", 65+rem-1) + result
			i = i / 26
		}
	}
	return result
}

type Reporter struct {
	templatePath string
	reports      map[string]Report
}

type Report struct {
	id           string
	templatePath string
	file         *excelize.File
	sheets       []api.TablePanel
}

func NewReport(id string, templatePath string) *Report {
	return &Report{id: id, templatePath: templatePath}
}

func (r *Report) openTemplate() {
	f, _ := excelize.OpenFile(r.templatePath)
	r.file = f
}

func (r *Report) placeholderCellRef(sheetName string, placeholder string) []string {
	refs := r.file.SearchSheet(sheetName, placeholder)
	return refs
}

func (r *Report) placeholderRowRef(sheetName string, placeholder string) int {
	refs := r.file.SearchSheet(sheetName, placeholder)
	split := CELL_REF_REG.FindAllString(refs[0], 2)
	idx, _ := strconv.Atoi(split[1])

	return idx
}

func (r *Report) writeTitle(sheetName string) {
	refs := r.placeholderCellRef(sheetName, "{{title}}")
	r.file.SetCellValue(sheetName, refs[0], sheetName)
}

func (r *Report) writeDate(sheetName string) {
	refs := r.placeholderCellRef(sheetName, "{{date}}")
	r.file.SetCellValue(sheetName, refs[0], time.Now().Format("Mon Jan 2 15:04:05"))
}

func (r *Report) writeCell(sheetName string, cellRef string, value interface{}) {
	r.file.SetCellValue(sheetName, cellRef, value)
}

func (r *Report) SetSheets(panels []api.TablePanel) {
	r.sheets = panels
}

func (r *Report) writeHeaders(sheetName string, columns []api.Column) {
	idx := r.placeholderRowRef(sheetName, "{{headers}}")
	log.DefaultLogger.Info(strconv.Itoa(idx))
	for i, column := range columns {
		log.DefaultLogger.Info(column.Text)
		cellRef := r.createCellRef(i, idx)
		log.DefaultLogger.Info(cellRef)
		r.writeCell(sheetName, cellRef, column.Text)
	}
}

func (r *Report) createDuplicateTableRows(sheetName string, numberToDuplicate int) int {
	idx := r.placeholderRowRef(sheetName, "{{rows}}")

	i := 0
	for i < numberToDuplicate {
		r.file.DuplicateRow(sheetName, idx)
		i += 1
	}

	return idx
}

func (r *Report) createCellRef(columnNumber int, rowNumber int) string {
	return intToCol(columnNumber) + strconv.Itoa(rowNumber)
}

func (r *Report) writeRows(sheetName string, rows [][]interface{}) {
	idx := r.createDuplicateTableRows(sheetName, len(rows))

	for _, row := range rows {
		for j, value := range row {
			cellRef := r.createCellRef(j, idx)
			r.writeCell(sheetName, cellRef, value)
		}
		idx += 1
	}
}

func (r *Report) Write(auth auth.AuthConfig) {
	if r.file == nil {
		r.openTemplate()
	}

	for _, s := range r.sheets {
		sIdx := r.file.NewSheet(s.Title)
		r.file.CopySheet(1, sIdx)
		s.GetData(auth)
		r.writeTitle(s.Title)
		r.writeDate(s.Title)
		r.writeHeaders(s.Title, s.Columns)
		r.writeRows(s.Title, s.Rows)
	}

	r.file.DeleteSheet("Sheet1")

	if err := r.file.SaveAs("./data/" + r.id + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}

func NewReporter(templatePath string) *Reporter {
	return &Reporter{templatePath: templatePath}
}

func (r *Reporter) createPath(scheduleID string) string {
	return scheduleID + ".xlsx"
}

func (r *Reporter) CreateNewReport(scheduleID string) *Report {
	report := NewReport(scheduleID, r.templatePath)
	return report
}
