package reporter

import (
	"errors"
	"fmt"
	"path/filepath"
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

func (r *Report) openTemplate() error {
	f, e := excelize.OpenFile(r.templatePath)

	if e != nil {
		log.DefaultLogger.Error("openTemplate: ", e.Error())
		return e
	}

	r.file = f

	return nil
}

func (r *Report) placeholderCellRef(sheetName string, placeholder string) ([]string, error) {
	refs := r.file.SearchSheet(sheetName, placeholder)

	if len(refs) == 0 {
		err := errors.New("Could not find cell reference for: " + sheetName + " - " + placeholder)
		log.DefaultLogger.Error("placeholderCellRef", err.Error())
		return nil, err
	}

	return refs, nil
}

func (r *Report) placeholderRowRef(sheetName string, placeholder string) (int, error) {
	refs := r.file.SearchSheet(sheetName, placeholder)

	if len(refs) == 0 {
		err := errors.New("Could not find cell reference for: " + sheetName + " - " + placeholder)
		log.DefaultLogger.Error("placeholderRowRef", err.Error())
		return 0, err
	}

	split := CELL_REF_REG.FindAllString(refs[0], 2)
	idx, _ := strconv.Atoi(split[1])

	return idx, nil
}

func (r *Report) writeTitle(sheetName string) error {
	refs, err := r.placeholderCellRef(sheetName, "{{title}}")

	if err != nil {
		return err
	}

	r.file.SetCellValue(sheetName, refs[0], sheetName)

	return nil
}

func (r *Report) writeDate(sheetName string) error {
	refs, err := r.placeholderCellRef(sheetName, "{{date}}")

	if err != nil {
		log.DefaultLogger.Error("writeDate: ", err.Error())
		return err
	}

	r.file.SetCellValue(sheetName, refs[0], time.Now().Format("Mon Jan 2 15:04:05"))

	return nil
}

func (r *Report) writeCell(sheetName string, cellRef string, value interface{}) {
	log.DefaultLogger.Info("SheetName: "+sheetName+" -- CellRef: "+cellRef, value)
	r.file.SetCellValue(sheetName, cellRef, value)
}

func (r *Report) SetSheets(panels []api.TablePanel) {
	r.sheets = panels
}

func (r *Report) writeHeaders(sheetName string, columns []api.Column) error {
	idx, err := r.placeholderRowRef(sheetName, "{{headers}}")

	if err != nil {
		log.DefaultLogger.Error("writeHeaders: ", err.Error())
		return err
	}

	log.DefaultLogger.Info(strconv.Itoa(idx))
	for i, column := range columns {
		log.DefaultLogger.Info(column.Text)
		cellRef := r.createCellRef(i, idx)
		log.DefaultLogger.Info(cellRef)
		r.writeCell(sheetName, cellRef, column.Text)
	}

	return nil
}
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (r *Report) createDuplicateTableRows(sheetName string, numberToDuplicate int) (int, error) {
	idx, err := r.placeholderRowRef(sheetName, "{{rows}}")

	duplicateNumber := Min(numberToDuplicate, 500)

	log.DefaultLogger.Info(strconv.Itoa(numberToDuplicate))
	if err != nil {
		log.DefaultLogger.Error("createDuplicateTableRows: ", err.Error())
		return 0, err
	}

	i := 0
	for i < duplicateNumber {
		log.DefaultLogger.Info(strconv.Itoa(i))
		r.file.DuplicateRow(sheetName, idx)
		i += 1
	}

	return idx, nil
}

func (r *Report) createCellRef(columnNumber int, rowNumber int) string {
	return intToCol(columnNumber) + strconv.Itoa(rowNumber)
}

func (r *Report) writeRows(sheetName string, rows [][]interface{}) error {
	idx, err := r.createDuplicateTableRows(sheetName, len(rows))

	if err != nil {
		log.DefaultLogger.Error("writeRows: ", err.Error())
		return err
	}

	log.DefaultLogger.Info("Write Rows: " + "Writing Rows")
	for _, row := range rows {
		for j, value := range row {
			cellRef := r.createCellRef(j, idx)
			log.DefaultLogger.Info("Writing Cell: ", cellRef, value)
			r.writeCell(sheetName, cellRef, value)
		}
		idx += 1
	}

	return nil
}

func (r *Report) Write(auth auth.AuthConfig) error {
	if r.file == nil {
		r.openTemplate()
	}

	for _, s := range r.sheets {
		log.DefaultLogger.Info("Writing Sheet: " + s.Title)

		sIdx := r.file.NewSheet(s.Title)
		err := r.file.CopySheet(1, sIdx)

		if err != nil {
			log.DefaultLogger.Error("Write: ", err.Error())
			return err
		}

		log.DefaultLogger.Info("Getting Data")

		s.GetData(auth)

		log.DefaultLogger.Info("Writing Title: " + s.Title)

		err = r.writeTitle(s.Title)

		if err != nil {
			log.DefaultLogger.Error("Write: ", err.Error())
			return nil
		}

		log.DefaultLogger.Info("Writing Date:", s.Title)

		err = r.writeDate(s.Title)

		if err != nil {
			log.DefaultLogger.Error("Write: ", err.Error())
			return nil
		}

		log.DefaultLogger.Info("Writing Headers")

		err = r.writeHeaders(s.Title, s.Columns)

		if err != nil {
			log.DefaultLogger.Error("Write: ", err.Error())
			return nil
		}

		log.DefaultLogger.Info("Writing Rows")

		err = r.writeRows(s.Title, s.Rows)

		if err != nil {
			log.DefaultLogger.Error("Write: ", err.Error())
			return nil
		}

	}

	log.DefaultLogger.Info("Deleting Sheet: Sheet1")
	r.file.DeleteSheet("Sheet1")

	log.DefaultLogger.Info("Saving Report: " + r.id)

	savePath := filepath.Join("data", r.id+".xlsx")
	if err := r.file.SaveAs(savePath); err != nil {
		log.DefaultLogger.Error("Write: ", err.Error())
	}

	return nil
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
