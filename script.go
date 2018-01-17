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
     "encoding/json"
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
  // filename := excelName[0:len(excelName)-5]
  // htmlFilename := []string{ filename,".html" }
  // htmlFilenameString := strings.Join(htmlFilename, "")
  HEADER_INDEX := 2
  ROW_INDEX_START := 3
  COLUMN_INDEX_DATA = :2

  xlFile, err := xlsx.OpenFile(excelName)
  if err != nil {
      fmt.Printf("Error: %s", err)
  }

  sheet := xlFile.Sheets[0]
  sites := sheet.Rows[HEADER_INDEX].Cells

  // os.MkdirAll(filename, os.ModePerm)
  templateName := ""
  templateContent := make(map[string]string)

  for _, row := range sheet.Rows[ROW_INDEX_START:] {
    //this is the template file name we are going to be looking for
    //open the html template one at a time and create the outputs

    if(len(row.Cells[0].String()) == 0 && len(row.Cells[1].String()) == 0) {
      continue
    }
    // fmt.Println(row)
    if(len(row.Cells[0].String()) > 0 && len(row.Cells[1].String()) == 0) {
      //write output if any exists then get the next template
      if( len(templateName) > 0 ) {
        os.MkdirAll(templateName, os.ModePerm)

        for _, cell := range sites[COLUMN_INDEX_DATA:] {
          htmlFilename := []string{ templateName,"/",cell.String(),".html" }
          htmlPath := strings.Join(htmlFilename, "")

          err = ioutil.WriteFile(htmlPath, []byte( templateContent[cell.String()] ), 0644)
          if err != nil {
            panic(err)
          }
        }
      }

      //get all the html templates setup
      templateName = row.Cells[0].String()
      templateContent = make(map[string]string)

      htmlTemplateFile := fmt.Sprintf("%s.html", templateName)
      htmlTemplate, err := ioutil.ReadFile(htmlTemplateFile) // just pass the file name
      if err != nil {
          fmt.Printf("Error %s", err)
          continue
      }

      //get all the html templates setup
      for _, cell := range sites[COLUMN_INDEX_DATA:] {
        // fmt.Println(cell.String())
        templateContent[cell.String()] = string(htmlTemplate)
      }
    } else {

      if(len(templateContent) == 0) {
        continue
      }

      elementKey := row.Cells[1]

      for index, cell := range row.Cells[COLUMN_INDEX_DATA:] {
        if(len(elementKey.String()) == 0 || len(cell.String()) == 0) {
          continue
        }

        templateKey := []string{ "[[",elementKey.String(),"]]" }
        templateString := strings.Join(templateKey, "")
        siteIndex := index + COLUMN_INDEX_DATA

        templateContent[sites[siteIndex].String()] = strings.Replace(templateContent[sites[siteIndex].String()], templateString, cell.String(), -1)

        // fmt.Println(templateContent)
      }
    }

    // htmlString := string(htmlTemplate)
    // lastCellIndex := len(row.Cells)-1
    //
    // htmlFilename := []string{ filename,"/",row.Cells[lastCellIndex].String(),".html" }
    // htmlFilenameString := strings.Join(htmlFilename, "")
    //
    // for index, cell := range row.Cells {
    //     templateKey := []string{ "[[[",keys[index].String(),"]]]" }
    //     templateString := strings.Join(templateKey, "")
    //
    //     htmlString = strings.Replace(htmlString, templateString , cell.String(), -1)
    // }

    // err = ioutil.WriteFile(htmlFilenameString, []byte(htmlString), 0644)
    // if err != nil {
    //   panic(err)
    // }
  }
}

type Schedule struct {
  Time string `json:"time"`
}

type Email struct {
    AppId string `json:"app_id"`
    Subject string `json:"subject"`
    From string `json:"from"`
    ReplyTo string `json:"reply_to"`
    Body string `json:"body"`
  }

type Messages struct {
  Email Email `json:"email"`
}

type Message struct {
  AppId string `json:"app_group_id"`
  SegmentId string `json:"segment_id"`
  Broadcast bool `json:"broadcast"`
  Schedule Schedule `json:"schedule"`
  Messages Messages `json:"messages"`
}

func uploadAppboy() {
    url := "https://rest.iad-01.braze.com/messages/schedule/create"
    fmt.Println("URL:>", url)

    timeToSend := "2019-01-01T09:25:25Z"

    jsonObject := Message{"5609dd9b-11b9-443d-a592-4c8e094677c3", "8dab634a-d644-445d-b8bf-0abceb630b7c", true, Schedule{ timeToSend }, Messages{Email{"e7502059-597c-4d14-bf39-eebfa035a24c", "subject test", "test@test.com", "test@test.com", "<b>htmlshouldbehere</b>"}}}
    stringObject, err := json.Marshal(jsonObject)

    fmt.Println("parameters", string(stringObject), err)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(stringObject))
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
    processSpreadsheet("template.xlsx")
    // uploadAppboy()
    // //get all spreadhseets
    // var spreadsheets []string = getSpreadsheets()
    //
    // //process all spreadsheets
    // for _, spreadsheet := range spreadsheets {
    //   processSpreadsheet(spreadsheet)
    // }

    // retrieveCampaignData()
    //upload to appboy
    // uploadAppboy()
}
