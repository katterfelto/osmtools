package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

//go:embed templates
var templates embed.FS

func getHome(w http.ResponseWriter, r *http.Request) {
	log := slog.Default()
	log.Info("getHome called")
	templ, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	templ.Execute(w, nil)
}

func postCoop30DayToOsmHandler(w http.ResponseWriter, r *http.Request) {
	log := slog.Default()

	log.Info("Received POST request for CO-OP 30 day to OSM conversion")

	// Handle file upload and processing here
	// For now, just log the file name
	file, header, err := receiveInputFile(log, r)
	if err != nil {
		log.Error("Error retrieving file from form", "error", err)
		http.Error(w, "Error processing file upload", http.StatusInternalServerError)
		return
	}

	// Process the file according to the specified parameters
	if err := processCoopBankStatement(log, file.Name(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, []string{"A", "C", "D", "F", "G", "H", "I", "L", "M", "O", "P", "S", "T", "U", "V"}, map[string]float64{"A": 18, "B": 40, "C": 40, "D": 24, "E": 15, "F": 15, "G": 15}, "A1"); err != nil {
		log.Error("Error processing spreadsheet", "error", err)
		http.Error(w, "Error processing spreadsheet", http.StatusInternalServerError)
		return
	}

	// Stream the processed file back to the client
	returnProcessedFile(log, w, r, fmt.Sprintf("processed-%s", header.Filename), file.Name())
}

func postCoopCustomToOsmHandler(w http.ResponseWriter, r *http.Request) {
	log := slog.Default()

	log.Info("Received POST request for CO-OP Custom to OSM conversion")

	// Handle file upload and processing here
	// For now, just log the file name
	file, header, err := receiveInputFile(log, r)
	if err != nil {
		log.Error("Error retrieving file from form", "error", err)
		http.Error(w, "Error processing file upload", http.StatusInternalServerError)
		return
	}

	// Process the file according to the specified parameters
	if err := processCoopBankStatement(log, file.Name(), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []string{"A", "B", "C", "E", "F", "H", "I", "J", "K", "N", "O", "S", "T", "U", "V", "W"}, map[string]float64{"A": 18, "B": 40, "C": 40, "D": 24, "E": 15, "F": 15, "G": 15}, "A1"); err != nil {
		log.Error("Error processing spreadsheet", "error", err)
		http.Error(w, "Error processing spreadsheet", http.StatusInternalServerError)
		return
	}

	// Stream the processed file back to the client
	returnProcessedFile(log, w, r, fmt.Sprintf("processed-%s", header.Filename), file.Name())
}

func receiveInputFile(log *slog.Logger, r *http.Request) (*os.File, *multipart.FileHeader, error) {
	// Handle file upload and processing here
	// For now, just log the file name
	file, header, err := r.FormFile("inputFile")
	if err != nil {
		log.Error("Error retrieving file from form", "error", err)
		return nil, nil, err
	}
	defer file.Close()

	// save the file to a temporary location
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		log.Error("Error creating temporary file", "error", err)
		return nil, nil, err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error("Error reading file", "error", err)
		return nil, nil, err
	}

	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Error("Error writing to temp file", "error", err)
		return nil, nil, err
	}

	log.Info("Received file, ready to process", "filename", header.Filename, "tempfilename", tempFile.Name())
	return tempFile, header, nil
}

func processCoopBankStatement(log *slog.Logger, filename string, removeRows []int, removeCols []string, setColWidths map[string]float64, deletePictureCell string) error {
	// Open the spreadsheet by the given path.
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Error("Error opening spreadsheet", "error", err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Error("Error closing spreadsheet", "error", err)
		}
	}()

	// Get the default sheet name.
	sheetName := f.GetSheetName(0)

	// Remove the rows from the default sheet
	for i := len(removeRows) - 1; i >= 0; i-- {
		if err := f.RemoveRow(sheetName, removeRows[i]); err != nil {
			log.Error("Error removing row", "error", err)
			return err
		}
	}

	// Remove the columns from the default sheet
	for i := len(removeCols) - 1; i >= 0; i-- {
		if err := f.RemoveCol(sheetName, removeCols[i]); err != nil {
			log.Error("Error removing column", "error", err)
			return err
		}
	}

	// Set the column widths
	for key, value := range setColWidths {
		f.SetColWidth(sheetName, key, key, value)
	}

	if err := f.DeletePicture(sheetName, deletePictureCell); err != nil {
		log.Error("Error deleting picture", "error", err)
		return err
	}

	// Save the spreadsheet by the given path.
	if err := f.Save(); err != nil {
		log.Error("Error saving spreadsheet", "error", err)
		return err
	}

	log.Info("Processing completed successfully", "tempfilename", filename)
	return nil
}

func returnProcessedFile(log *slog.Logger, w http.ResponseWriter, r *http.Request, filename string, tempFilename string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, tempFilename)

	log.Info("Processed file sent", "filename", filename)

	// clean up the temporary file after a short delay to ensure it has been served
	go func() {
		time.Sleep(10 * time.Second)
		if err := os.Remove(tempFilename); err != nil {
			log.Error("Error removing temporary file", "error", err)
		}
	}()
}
