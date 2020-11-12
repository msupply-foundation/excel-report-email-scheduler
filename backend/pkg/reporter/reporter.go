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
	f, err := excelize.OpenFile(r.templatePath)
	if err != nil {
		log.DefaultLogger.Error("Could not open template: ", err.Error())
		return err
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
	idx, err := strconv.Atoi(split[1])
	if err != nil {
		log.DefaultLogger.Error("placeholderRowRef: Atoi: " + err.Error())
		return 0, err
	}

	return idx, nil
}

func (r *Report) writeTitle(sheetName string) error {
	refs, err := r.placeholderCellRef(sheetName, "{{title}}")
	if err != nil {
		log.DefaultLogger.Error("writeTitle: placeholderCellRef: " + err.Error())
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

	for i, column := range columns {
		cellRef := r.createCellRef(i, idx)
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
	if err != nil {
		log.DefaultLogger.Error("createDuplicateTableRows: " + err.Error())
		return 0, err
	}

	duplicateNumber := Min(numberToDuplicate, 500)
	i := 0
	for i < duplicateNumber {
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
		log.DefaultLogger.Error("writeRows: createDuplicateTableRows: " + err.Error())
		return err
	}

	for _, row := range rows {
		for j, value := range row {
			cellRef := r.createCellRef(j, idx)
			r.writeCell(sheetName, cellRef, value)
		}
		idx += 1
	}
	return nil
}

func (r *Report) Write(auth auth.AuthConfig) error {
	log.DefaultLogger.Info(fmt.Sprintf("Starting to create report %s...", r.id))
	if r.file == nil {
		if err := r.openTemplate(); err != nil {
			return err
		}
	}

	for _, s := range r.sheets {
		log.DefaultLogger.Info(fmt.Sprintf("Creating new sheet %s", s.Title))
		sIdx := r.file.NewSheet(s.Title)

		if err := r.file.CopySheet(1, sIdx); err != nil {
			log.DefaultLogger.Error("Write: copySheet: " + err.Error())
			return err
		}

		if err := s.GetData(auth); err != nil {
			log.DefaultLogger.Error("Write: GetData: " + err.Error())
			return err
		}

		if err := r.writeTitle(s.Title); err != nil {
			log.DefaultLogger.Error("Write: writeTitle: " + err.Error())
			return err
		}

		if err := r.writeDate(s.Title); err != nil {
			log.DefaultLogger.Error("Write: writeDate: " + err.Error())
			return err
		}

		if err := r.writeHeaders(s.Title, s.Columns); err != nil {
			log.DefaultLogger.Error("Write: writeHeaders: " + err.Error())
			return err
		}

		if err := r.writeRows(s.Title, s.Rows); err != nil {
			log.DefaultLogger.Error("Write: writeRows: " + err.Error())
			return err
		}
	}

	r.file.DeleteSheet("Sheet1")

	log.DefaultLogger.Info("Saving report...")

	savePath := filepath.Join("data", r.id+".xlsx")
	if err := r.file.SaveAs(savePath); err != nil {
		log.DefaultLogger.Error("Write: ", err.Error())
	}

	log.DefaultLogger.Info(fmt.Sprintf("Report finished! %s :tada", r.id))

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
