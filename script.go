package main

import (
    "fmt"
    "github.com/tealeg/xlsx"
    "io/ioutil"
    "strings"
     "path/filepath"
     "os"
     "net/http"
     "bytes"
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

func uploadAppboy() {
    /*
    {
      "app_id": (required, string) see App Identifier above,
      "subject": (optional, string),
      "from": (required, valid email address in the format "Display Name <email@address.com>"),
      "reply_to": (optional, valid email address in the format "email@address.com" - defaults to your app group's default reply to if not set),
      "body": (required unless email_template_id is given, valid HTML),
      "plaintext_body": (optional, valid plaintext, defaults to autogenerating plaintext from "body" when this is not set),
      "preheader"*: (optional, string) Recommended length 50-100 characters.
      "email_template_id": (optional, string) If provided, we will use the subject/body values from the given email template UNLESS they are specified here, in which case we will override the provided template,
      "message_variation_id": (optional, string) used when providing a campaign_id to specify which message variation this message should be tracked under,
      "extras": (optional, valid Key Value Hash), extra hash - for SendGrid customers, this will be passed to SendGrid as Unique Arguments,
      "headers": (optional, valid Key Value Hash), hash of custom extensions headers. Currently, only supported for SendGrid customers
    }
    */

    url := "http://restapi3.apiary.io/notes"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func retrieveCampaignData() {
  	campaignURL := "https://rest.iad-01.braze.com/campaigns/data_series"
    appGroupId := os.Getenv("CM_BRAZE_GROUP_ID")
    campaignId := ""
    dataUrl := fmt.Sprintf("%s?app_group_id=%s&campaign_id=%s&length=%s", campaignURL, appGroupId, campaignId, "14")

    fmt.Println(dataUrl)
    /*
    app_group_id	Yes	String	App Group API Identifier
    campaign_id	Yes	String	Campaign API Identifier
    length	Yes	Integer	Max number of days before ending_at to include in the returned series - must be between 1 and 100 inclusive
    ending_at	No	DateTime (ISO 8601 string)	Date on which the data series should end - defaults to time of the request
    */

    req, err := http.NewRequest("GET", dataUrl, nil)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
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
