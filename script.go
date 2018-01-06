package main

import (
    "fmt"
    "github.com/tealeg/xlsx"
    "io/ioutil"
    "strings"
     "path/filepath"
     "os"
)
// this is a comment

func getSpreadsheets() []string {

  var spreadsheets []string

  err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
      if filepath.Ext(path) != ".xlsx" {
        return nil
      }
      spreadsheets = append(spreadsheets, path)
      return nil
  })
  if err != nil {
      panic(err)
  }

  return spreadsheets
}

func processSpreadsheet(excelName string) {
  filename := excelName[0:len(excelName)-5]
  htmlFilename := []string{ filename,".html" }
  htmlFilenameString := strings.Join(htmlFilename, "")

  htmlTemplate, err := ioutil.ReadFile(htmlFilenameString) // just pass the file name
  if err != nil {
      fmt.Print(err)
  }

  xlFile, err := xlsx.OpenFile(excelName)
  if err != nil {
      fmt.Printf("error")
  }

  sheet := xlFile.Sheets[0]
  keys := sheet.Rows[0].Cells

  os.MkdirAll(filename, os.ModePerm)

  for _, row := range sheet.Rows[1:] {

    htmlString := string(htmlTemplate)
    lastCellIndex := len(row.Cells)-1

    htmlFilename := []string{ filename,"/",row.Cells[lastCellIndex].String(),".html" }
    htmlFilenameString := strings.Join(htmlFilename, "")

    for index, cell := range row.Cells {
        templateKey := []string{ "[[[",keys[index].String(),"]]]" }
        templateString := strings.Join(templateKey, "")

        htmlString = strings.Replace(htmlString, templateString , cell.String(), -1)
    }

    err = ioutil.WriteFile(htmlFilenameString, []byte(htmlString), 0644)
    if err != nil {
      panic(err)
    }
  }
}

func main() {

    //get all spreadhseets
    var spreadsheets []string = getSpreadsheets()

    //process all spreadsheets
    for _, spreadsheet := range spreadsheets {
      processSpreadsheet(spreadsheet)
    }

    //upload to appboy
}
