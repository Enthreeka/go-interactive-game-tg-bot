package excel

import (
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/xuri/excelize/v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"sync"
	"time"
)

const filename = "users.xlsx"

type Excel struct {
	log *logger.Logger
	mu  sync.Mutex
}

func NewExcel(log *logger.Logger) *Excel {
	return &Excel{log: log}
}

func (e *Excel) GenerateUserResultsExcelFile(results []entity.UserResult, contestID int, username string) (string, error) {
	start := time.Now()

	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			e.log.Error("failed to close excel: %v", err)
		}
	}()

	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	headers := map[string]string{
		"A1": "TG Username",
		"B1": "User ID",
		"C1": "Result ID",
		"D1": "Total Points",
	}

	for cell, value := range headers {
		f.SetCellValue(sheetName, cell, value)
	}

	for i, result := range results {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), result.User.TGUsername)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), result.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), result.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), result.TotalPoints)
	}

	filename := fmt.Sprintf("contest_%d_results.xlsx", contestID)
	err := f.SaveAs(filename)
	if err != nil {
		e.log.Error("failed to save file: %s", filename)
		return "", err
	}

	end := time.Since(start)
	e.log.Info("[%s] by [%s] Время генерации файла: %f", filename, username, end.Seconds())
	return filename, nil
}

func (e *Excel) GenerateForUserResultsExcelFile(results []entity.UserResult, contestID int, username string) (string, error) {
	start := time.Now()

	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			e.log.Error("failed to close excel: %v", err)
		}
	}()

	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	headers := map[string]string{
		"A1": "TG Username",
		"B1": "Total Points",
	}

	for cell, value := range headers {
		f.SetCellValue(sheetName, cell, value)
	}

	for i, result := range results {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), result.User.TGUsername)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), result.TotalPoints)
	}

	filename := fmt.Sprintf("contest_%d_results.xlsx", contestID)
	err := f.SaveAs(filename)
	if err != nil {
		e.log.Error("failed to save file: %s", filename)
		return "", err
	}

	end := time.Since(start)
	e.log.Info("[%s] by [%s] Время генерации файла: %f", filename, username, end.Seconds())
	return filename, nil
}

func (e *Excel) GetExcelFile(fileName string) (*[]byte, error) {
	//e.mu.Lock()

	file, err := os.Open(fileName)
	if err != nil {
		e.log.Error("os.Open: failed to open file: %v", err)
		return nil, err
	}
	defer func() {
		//e.mu.Unlock()
		err := file.Close()
		if err != nil {
			e.log.Error("%v", err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		e.log.Error("file.Stat: failed to get file stat: %v", err)
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileID := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: make([]byte, fileSize),
	}

	if _, err = file.Read(fileID.Bytes); err != nil {
		e.log.Error("file.Read: failed to get read file: %v", err)
		return nil, err
	}

	return &fileID.Bytes, nil
}
