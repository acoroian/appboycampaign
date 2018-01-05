package main

import (
    "fmt"
    "github.com/tealeg/xlsx"
    "io/ioutil"
    "strings"
)
// this is a comment

func main() {
    htmlTemplate, err := ioutil.ReadFile("generic_promo_template.html") // just pass the file name
    if err != nil {
        fmt.Print(err)
    }

    excelFileName := "html_template.xlsx"
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
        fmt.Printf("error")
    }

    sheet := xlFile.Sheets[0]

    keys := sheet.Rows[0].Cells

    for _, row := range sheet.Rows[1:] {

      htmlString := string(htmlTemplate)
      lastCellIndex := len(row.Cells)-1

      htmlFilename := []string{ row.Cells[lastCellIndex].String(),".html" }
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
