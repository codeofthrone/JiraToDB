package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/GO-SQL-Driver/MySQL"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type JiraObject struct {
	Expend     string       `json:"expend"`
	StartAt    int          `json:"startAt"`
	MaxResults int          `json:"maxResults"`
	Total      int          `json:"total"`
	Issues     []JiraIssues `json:"issues"`
}

type JiraIssues struct {
	Expend string     `json:"expend"`
	ID     string     `json:"id"`
	Self   string     `json:"self"`
	Key    string     `json:"key"`
	Fields JiraFields `json:"fields"`
}

type JiraIssuesChangeLog struct {
	Expend    string                  `json:"expend"`
	ID        string                  `json:"id"`
	Self      string                  `json:"self"`
	Key       string                  `json:"key"`
	Fields    JiraFields              `json:"fields"`
	ChangeLog JiraIssuesChangeLogData `json:"changelog"`
}

type JiraFields struct {
	Summary                       string                  `json:"summary"`
	Progress                      interface{}             `json:"progress"`
	IssueType                     JiraFieldsIssueType     `json:"issuetype"`
	Votes                         interface{}             `json:"votes"`
	Resolution                    JiraFieldsResolution    `json:"resolution"`
	FixVersion                    []JiraFieldsFixVersions `json:"fixVersions"`
	ResoluationDate               string                  `json:"resolutiondate"`
	TimeSpent                     int                     `json:"timespent"`
	Reporter                      JiraFieldsReporter      `json:"reporter"`
	AggregateTimeOriginalEstimate int                     `json:"aggregatetimeoriginalestimate"`
	Updated                       string                  `json:"updated"`
	Created                       string                  `json:"created"`
	Description                   string                  `json:"description"`
	Priority                      JiraFieldsPriority      `json:"priority"`
	DueDate                       string                  `json:"duedate"`
	Status                        JiraFieldsStatus        `json:"status"`
	Labels                        []string                `json:"labels"`
	Assignee                      JiraFieldsReporter      `json:"assignee"`
	Project                       JiraFieldsProject       `json:"project"`
	Version                       []JiraFieldsVersions    `json:"versions"`
	Components                    []JiraFieldsResolution  `json:"components"`
}

type JiraFieldsVersions struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Archived    bool   `json:"archived"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate"`
}

type JiraFieldsProject struct {
	Self         string      `json:"self"`
	ID           string      `json:"id"`
	Key          string      `json:"key"`
	Name         string      `json:"name"`
	EmailAddress string      `json:"emailAddress"`
	AvatarUrls   interface{} `json:"avatarUrls"`
}

type JiraFieldsStatus struct {
	Self        string `json:"self"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

type JiraFieldsPriority struct {
	Self    string `json:"self"`
	IconURL string `json:"iconUrl"`
	Name    string `json:"name"`
	ID      string `json:"id"`
}

type JiraFieldsReporter struct {
	Self         string      `json:"self"`
	Name         string      `json:"name"`
	EmailAddress string      `json:"emailAddress"`
	AvatarUrls   interface{} `json:"avatarUrls"`
	DisplayName  string      `json:"displayName"`
	Active       bool        `json:"active"`
}

type JiraFieldsResolution struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type JiraFieldsIssueType struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Description string `json:"description"`
	iconUrl     string `json:"iconUrl"`
	Name        string `json:"name"`
	SubTask     bool   `json:"subtask"`
}

type JiraFieldsFixVersions struct {
	Self           string `json:"self"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Archived       bool   `json:"archived"`
	Released       bool   `json:"released"`
	ResolutionDate string `json:"resolutiondate"`
	TimeSpent      string `json:"timespent"`
}

type JiraIssuesChangeLogData struct {
	StartAt    int                                `json:"startAt"`
	MaxResults int                                `json:"maxResult"`
	Total      int                                `json:"total"`
	Histories  []JiraIssuesChangeLogDataHistories `json:"histories"`
}

type JiraIssuesChangeLogDataHistories struct {
	ID      string                                  `json:"id"`
	Author  JiraIssuesChangeLogDataHistoriesAuthor  `json:"author"`
	Created string                                  `json:"created"`
	Items   []JiraIssuesChangeLogDataHistoriesItems `json:"items"`
}

type JiraIssuesChangeLogDataHistoriesAuthor struct {
	Self         string      `json:"self"`
	Name         string      `json:"name"`
	EmailAddress string      `json:"emailAddress"`
	AvatarUrls   interface{} `json:"avatarUrls"`
	DisplayName  string      `json:"displayName"`
	Active       bool        `json:"active"`
}

type JiraIssuesChangeLogDataHistoriesItems struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

type ErrorMessageObject struct {
	ErrorMessages []string    `json:"errorMessages"`
	Errors        interface{} `json:"errors"`
}

func SQLcheckItem(MysqlInfo, sqlstr string) int {
	db, err := sql.Open("mysql", MysqlInfo)
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
	}
	var CounterType string
	db.QueryRow(sqlstr).Scan(&CounterType)
	db.Close()
	if CounterType == "" {
		return 0
	} else {
		returnstr, _ := strconv.Atoi(CounterType)
		return returnstr
	}
}

func checkItem(MysqlInfo, mysqlTable string, Colname string, Colvalue string) int {
	db, err := sql.Open("mysql", MysqlInfo)
	if err != nil {
		panic(err.Error())
		log.Fatal(err.Error())
	}
	count_query_string := " SELECT Sn FROM `" + mysqlTable + "` WHERE `" + Colname + "`='" + Colvalue + "'  "
	var CounterType string
	db.QueryRow(count_query_string).Scan(&CounterType)
	db.Close()
	if CounterType == "" {
		return 0
	} else {
		returnstr, _ := strconv.Atoi(CounterType)
		return returnstr
	}
}

func SQLInsertStr(MysqlInfo string, SQLstr string) (string, string, int) {
	db, err := sql.Open("mysql", MysqlInfo)
	var (
		insertcols string
		insertstr  string
		colname    []string
	)
	// log.Println(SQLstr)
	cols, err := db.Query(SQLstr)
	// cols, err := db.Query("SELECT * FROM `Issues` order by SN DESC")
	if err != nil {
		// panic(err.Error())
		log.Fatal(err.Error())
	}
	colname, err = cols.Columns()
	cols.Close()
	if err != nil {
		// panic(err.Error())
		log.Fatal(err.Error())
	}
	for _, value := range colname {
		if value == "Sn" {
			insertcols = "`" + value + "`"
			insertstr = "?"
		} else {
			insertcols = insertcols + ",`" + value + "`"
			insertstr = insertstr + ", ?"
		}
	}
	db.Close()
	return insertcols, insertstr, len(colname)
}

func InsertTable(MysqlInfo string, TableName string, insertcols string, insertstr string, insertdata []interface{}) int {
	db, err := sql.Open("mysql", MysqlInfo)
	sqlStr := "INSERT INTO " + TableName + "(" + insertcols + ") VALUES (" + insertstr + ")"
	// log.Println(sqlStr)
	// log.Println(insertdata)
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Println(TableName, insertcols, insertstr, insertdata)
		log.Fatal(err.Error())
	}
	// stmt.Exec(is...)
	res, err := stmt.Exec(insertdata...)
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
	int64id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}
	idstring := strconv.FormatInt(int64id, 10)
	id, err := strconv.Atoi(idstring)
	if err != nil {
		log.Fatal(err.Error())
	}
	// log.Println("last Insert ID :", id)
	db.Close()
	return id
}

func UpdateTable(MysqlInfo string, sqlstr string) int {
	db, err := sql.Open("mysql", MysqlInfo)
	stmt, err := db.Prepare(sqlstr)
	if err != nil {
		log.Println(sqlstr)
		log.Fatal(err.Error())
		// panic(err.Error())
	}
	// stmt.Exec(is...)
	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}
	int64id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}
	idstring := strconv.FormatInt(int64id, 10)
	id, err := strconv.Atoi(idstring)
	if err != nil {
		log.Fatal(err.Error())
	}
	// log.Println("last Insert ID :", id)
	db.Close()
	return id
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	file, err := os.Open("ProjectList.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	jiraweb := flag.String("jiraweb", "", "Jira Web Site Address, example : -jiraweb=http://jira.sw.studio.htc.com")
	// projectname := flag.String("projectname", "", "Project Name, example : -projectname=TYGH")
	username := flag.String("username", "", "User Name")
	password := flag.String("password", "", "Password")
	mysqldb := flag.String("db", "", "mysql database")
	dbID := flag.String("ID", "", "mysql id")
	dbpw := flag.String("PW", "", "mysql pw")

	for scanner.Scan() {
		fmt.Println(scanner.Text())

		projectname := scanner.Text()
		flag.Parse()
		var (
			insertcols   string
			insertstr    string
			colsize      int
			FixversionID int
			ComponentID  int
			VersionID    int
			LabelID      int
		)
		MysqlInfo := *dbpw + ":" + *dbID + "@/" + *mysqldb
		IssuesCol := "SELECT * FROM `Issues` order by SN DESC"
		insertcols, insertstr, colsize = SQLInsertStr(MysqlInfo, IssuesCol)

		if flag.NFlag() >= 3 {
			client := &http.Client{}
			getProjectIssueList := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=-1", *jiraweb, projectname)
			req, err := http.NewRequest("GET", getProjectIssueList, nil)
			req.SetBasicAuth(*username, *password)
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)

			if err != nil {
				log.Fatal(err)
			}

			// Clone http response
			var bodyBytes []byte
			if resp.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(resp.Body)
			}
			resp_Error := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			var jiraObject JiraObject
			errorMessage := json.NewDecoder(resp.Body).Decode(&jiraObject)
			if errorMessage != nil {
				log.Fatal(errorMessage)
			}

			// If value of Total it seems Project name is error.
			if jiraObject.Total == 0 {
				var errorMessageObject ErrorMessageObject
				json.NewDecoder(resp_Error).Decode(&errorMessageObject)
				for x := range errorMessageObject.ErrorMessages {
					fmt.Println(errorMessageObject.ErrorMessages[x])
				}
				log.Fatal("Issue of Project is empty")
			}

			fmt.Println("Total : ", jiraObject.Total)
			fmt.Println("MaxResult : ", jiraObject.MaxResults)

			for x := range jiraObject.Issues {
				Key := jiraObject.Issues[x].Key
				JIRAID := jiraObject.Issues[x].ID
				fmt.Println(Key)
				// Key = "TYGH-110"
				changeLogCommand := fmt.Sprintf("%s/rest/api/2/issue/%s?expand=changelog", *jiraweb, Key)

				reqChangeLog, err := http.NewRequest("GET", changeLogCommand, nil)
				reqChangeLog.SetBasicAuth(*username, *password)
				reqChangeLog.Header.Add("Content-Type", "application/json")
				respChangeLog, err := client.Do(reqChangeLog)

				if err != nil {
					log.Fatal(err)
				}

				var jiraChangeLog JiraIssuesChangeLog
				json.NewDecoder(respChangeLog.Body).Decode(&jiraChangeLog)

				// update label , fixversion , version , component
				Fixversion := jiraChangeLog.Fields.FixVersion
				Labels := jiraChangeLog.Fields.Labels
				Version := jiraChangeLog.Fields.Version
				Components := jiraChangeLog.Fields.Components

				FixversionID = 0
				VersionID = 0
				LabelID = 0
				ComponentID = 0

				Labelsinsertcols, Labelsinsertstr, LabelsColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM  Labels")
				Labelsinsertdata := make([]interface{}, LabelsColSize)
				Labelinsertcols, Labelinsertstr, LabelColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM  Label")
				Labelinsertdata := make([]interface{}, LabelColSize)
				// remove flag
				sqlstr := " update Label set `Enable`='0' where `Id`='" + JIRAID + "' "
				UpdateTable(MysqlInfo, sqlstr)
				// update current flag
				for i := range Labels {
					if checkItem(MysqlInfo, "Labels", "Name", Labels[i]) == 0 {
						Labelsinsertdata[0] = ""
						Labelsinsertdata[1] = Labels[i]
						InsertTable(MysqlInfo, "Labels", Labelsinsertcols, Labelsinsertstr, Labelsinsertdata)
						LabelID = SQLcheckItem(MysqlInfo, "select sn from Labels where `Name`='"+Labels[i]+"'")
						Labelinsertdata[0] = ""
						Labelinsertdata[1] = JIRAID
						Labelinsertdata[2] = LabelID
						Labelinsertdata[3] = 1
						LabelID = InsertTable(MysqlInfo, "Label", Labelinsertcols, Labelinsertstr, Labelinsertdata)
					} else {
						LabelID = SQLcheckItem(MysqlInfo, "select sn from Labels where `Name`='"+Labels[i]+"'")
						LID := strconv.Itoa(LabelID)
						// update label
						if SQLcheckItem(MysqlInfo, "select sn from Label where `Id`='"+JIRAID+"' AND `Data`='"+LID+"'") == 0 {
							Labelinsertdata[0] = ""
							Labelinsertdata[1] = JIRAID
							Labelinsertdata[2] = LabelID
							Labelinsertdata[3] = 1
							LabelID = InsertTable(MysqlInfo, "Label", Labelinsertcols, Labelinsertstr, Labelinsertdata)
						} else {
							sqlstr := " update Label set `Enable`='1' where `Id`='" + JIRAID + "' "
							UpdateTable(MysqlInfo, sqlstr)
						}
					}
				}

				Fixversionsinsertcols, Fixversionsinsertstr, FixversionsColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Fixversions")
				Fixversionsinsertdata := make([]interface{}, FixversionsColSize)
				Fixversioninsertcols, Fixversioninsertstr, FixversionColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Fixversion")
				Fixversioninsertdata := make([]interface{}, FixversionColSize)
				// remove flag
				sqlstr = " update Fixversion set `Enable`='0' where `Id`='" + JIRAID + "' "
				UpdateTable(MysqlInfo, sqlstr)
				// update current flag
				for i, _ := range Fixversion {
					Fixversionsinsertdata[0] = ""
					Fixversionsinsertdata[1] = Fixversion[i].ID
					Fixversionsinsertdata[2] = Fixversion[i].Name
					Fixversionsinsertdata[3] = Fixversion[i].ResolutionDate
					Fixversionsinsertdata[4] = Fixversion[i].Released
					Fixversionsinsertdata[5] = Fixversion[i].Self
					Fixversionsinsertdata[6] = Fixversion[i].TimeSpent
					Fixversionsinsertdata[7] = Fixversion[i].Archived
					Fixversionsinsertdata[8] = projectname
					if checkItem(MysqlInfo, "Fixversions", "Name", Fixversion[i].Name) == 0 {
						InsertTable(MysqlInfo, "Fixversions", Fixversionsinsertcols, Fixversionsinsertstr, Fixversionsinsertdata)
						FixversionID = SQLcheckItem(MysqlInfo, "select sn from Fixversions where `Id`='"+Fixversion[i].ID+"'")
						Fixversioninsertdata[0] = ""
						Fixversioninsertdata[1] = JIRAID
						Fixversioninsertdata[2] = FixversionID
						Fixversioninsertdata[3] = 1
						FixversionID = InsertTable(MysqlInfo, "Fixversion", Fixversioninsertcols, Fixversioninsertstr, Fixversioninsertdata)
					} else {
						FixversionID = SQLcheckItem(MysqlInfo, "select sn from Fixversions where `Id`='"+Fixversion[i].ID+"'")
						FID := strconv.Itoa(FixversionID)
						// update label
						if SQLcheckItem(MysqlInfo, "select sn from Fixversion where `Id`='"+JIRAID+"' AND `Data`='"+FID+"'") == 0 {
							// fixs y fix n
							Fixversioninsertdata[0] = ""
							Fixversioninsertdata[1] = JIRAID
							Fixversioninsertdata[2] = FID
							Fixversioninsertdata[3] = 1
							InsertTable(MysqlInfo, "Fixversion", Fixversioninsertcols, Fixversioninsertstr, Fixversioninsertdata)
						} else {
							// fixs y fix y
							// update fixversions data

							released := "1"
							if Fixversionsinsertdata[4] == true {
								released = "1"
							} else {
								released = "0"
							}
							Archievd := "1"
							if Fixversionsinsertdata[7] == true {
								Archievd = "1"
							} else {
								Archievd = "0"
							}

							sqlfixversion := "update Fixversions set  " +
								"`ReleaseDate`='" + Fixversionsinsertdata[3].(string) + "' ," +
								"`Released`='" + released + "' ," +
								"`TimeSpent`='" + Fixversionsinsertdata[6].(string) + "' ," +
								"`Archived`='" + Archievd + "' ," +
								"`Name`='" + Fixversionsinsertdata[2].(string) + "' ," +
								"`Project`='" + Fixversionsinsertdata[8].(string) + "' " +
								" where `Sn`='" + FID + "' "
							// log.Println(sqlfixversion)
							UpdateTable(MysqlInfo, sqlfixversion)
							// update fixversion data
							sqlstr := " update Fixversion set `Enable`='1' where `Id`='" + JIRAID + "' AND `Data`='" + FID + "' "
							UpdateTable(MysqlInfo, sqlstr)
						}
					}
				}

				Versionsinsertcols, Versionsinsertstr, VersionsColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Versions")
				Versionsinsertdata := make([]interface{}, VersionsColSize)
				Versioninsertcols, Versioninsertstr, VersionColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Version")
				Versioninsertdata := make([]interface{}, VersionColSize)
				// remove flag
				sqlstr = " update Version set `Enable`='0' where `Id`='" + JIRAID + "' "
				UpdateTable(MysqlInfo, sqlstr)
				// update current flag
				for i, _ := range Version {
					Versionsinsertdata[0] = ""
					Versionsinsertdata[1] = Version[i].ID
					Versionsinsertdata[2] = Version[i].Name
					Versionsinsertdata[3] = Version[i].Archived
					Versionsinsertdata[4] = Version[i].Released
					Versionsinsertdata[5] = Version[i].ReleaseDate
					Versionsinsertdata[6] = Version[i].Self
					Versionsinsertdata[7] = projectname
					if checkItem(MysqlInfo, "Versions", "Name", Version[i].Name) == 0 {
						InsertTable(MysqlInfo, "Versions", Versionsinsertcols, Versionsinsertstr, Versionsinsertdata)
						VersionID = SQLcheckItem(MysqlInfo, "select sn from Versions where `Id`='"+Version[i].ID+"'")
						Versioninsertdata[0] = ""
						Versioninsertdata[1] = JIRAID
						Versioninsertdata[2] = VersionID
						Versioninsertdata[3] = 1
						VersionID = InsertTable(MysqlInfo, "Version", Versioninsertcols, Versioninsertstr, Versioninsertdata)
					} else {
						VersionID = SQLcheckItem(MysqlInfo, "select sn from Versions where `Id`='"+Version[i].ID+"'")
						VID := strconv.Itoa(VersionID)
						// update label
						if SQLcheckItem(MysqlInfo, "select sn from Version where `Id`='"+JIRAID+"' AND `Data`='"+VID+"'") == 0 {
							// fixs y fix n
							Versioninsertdata[0] = ""
							Versioninsertdata[1] = JIRAID
							Versioninsertdata[2] = VID
							Versioninsertdata[3] = 1
							InsertTable(MysqlInfo, "Version", Versioninsertcols, Versioninsertstr, Versioninsertdata)
						} else {
							// fixs y fix y
							// update fixversions data

							released := "1"
							if Versionsinsertdata[3] == true {
								released = "1"
							} else {
								released = "0"
							}
							Archievd := "1"
							if Versionsinsertdata[4] == true {
								Archievd = "1"
							} else {
								Archievd = "0"
							}

							sqlsstr := "update Versions set  " +
								"`ReleaseDate`='" + Versionsinsertdata[5].(string) + "' ," +
								"`Released`='" + released + "' ," +
								"`Archived`='" + Archievd + "' ," +
								"`Name`='" + Versionsinsertdata[2].(string) + "' ," +
								"`Project`='" + Versionsinsertdata[7].(string) + "' " +
								" where `Sn`='" + VID + "' "
							// log.Println(sqlsstr)
							UpdateTable(MysqlInfo, sqlsstr)
							// update fixversion data
							sqlstr := " update Version set `Enable`='1' where `Id`='" + JIRAID + "' AND `Data`='" + VID + "' "
							UpdateTable(MysqlInfo, sqlstr)
						}
					}
				}

				Componentsinsertcols, Componentsinsertstr, ComponentsColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Components")
				Componentsinsertdata := make([]interface{}, ComponentsColSize)
				Componentinsertcols, Componentinsertstr, ComponentColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Component")
				Componentinsertdata := make([]interface{}, ComponentColSize)
				// remove flag
				sqlstr = " update Component set `Enable`='0' where `Id`='" + JIRAID + "' "
				UpdateTable(MysqlInfo, sqlstr)
				// update current flag
				for i, _ := range Components {
					Componentsinsertdata[0] = ""
					Componentsinsertdata[1] = Components[i].ID
					Componentsinsertdata[2] = Components[i].Name
					Componentsinsertdata[3] = Components[i].Description
					Componentsinsertdata[4] = Components[i].Self
					if checkItem(MysqlInfo, "Components", "Name", Components[i].Name) == 0 {
						InsertTable(MysqlInfo, "Components", Componentsinsertcols, Componentsinsertstr, Componentsinsertdata)
						ComponentID = SQLcheckItem(MysqlInfo, "select sn from Components where `Id`='"+Components[i].ID+"'")
						Componentinsertdata[0] = ""
						Componentinsertdata[1] = JIRAID
						Componentinsertdata[2] = ComponentID
						Componentinsertdata[3] = 1
						ComponentID = InsertTable(MysqlInfo, "Component", Componentinsertcols, Componentinsertstr, Componentinsertdata)
					} else {
						ComponentID = SQLcheckItem(MysqlInfo, "select sn from Components where `Id`='"+Components[i].ID+"'")
						VID := strconv.Itoa(ComponentID)
						// update label
						if SQLcheckItem(MysqlInfo, "select sn from Component where `Id`='"+JIRAID+"' AND `Data`='"+VID+"'") == 0 {
							// fixs y fix n
							Componentinsertdata[0] = ""
							Componentinsertdata[1] = JIRAID
							Componentinsertdata[2] = VID
							Componentinsertdata[3] = 1
							InsertTable(MysqlInfo, "Component", Componentinsertcols, Componentinsertstr, Componentinsertdata)
						} else {
							// fixs y fix y
							// update fixversions data
							sqlsstr := "update Components set  " +
								"`Description`='" + Componentsinsertdata[3].(string) + "' ," +
								"`Name`='" + Componentsinsertdata[2].(string) + "' " +
								" where `Sn`='" + VID + "' "
							// log.Println(sqlsstr)
							UpdateTable(MysqlInfo, sqlsstr)
							// update fixversion data
							sqlstr := " update Component set `Enable`='1' where `Id`='" + JIRAID + "' AND `Data`='" + VID + "' "
							UpdateTable(MysqlInfo, sqlstr)
						}
					}
				}

				insertdata := make([]interface{}, colsize)
				insertdata[0] = ""
				insertdata[1] = JIRAID
				insertdata[2] = Key
				insertdata[3] = jiraChangeLog.Fields.Summary
				insertdata[4] = jiraChangeLog.Fields.IssueType.Name
				insertdata[5] = jiraChangeLog.Fields.Resolution.Name
				insertdata[6] = FixversionID
				insertdata[7] = jiraChangeLog.Fields.Reporter.Name
				insertdata[8] = jiraChangeLog.Fields.Description
				insertdata[9] = jiraChangeLog.Fields.Priority.Name
				insertdata[10] = jiraChangeLog.Fields.Status.Name
				insertdata[11] = LabelID
				insertdata[12] = jiraChangeLog.Fields.Assignee.Name
				insertdata[13] = VersionID
				insertdata[14] = ComponentID
				insertdata[15] = jiraChangeLog.Fields.Created
				insertdata[16] = jiraChangeLog.Fields.Updated

				// Check JIRA ID
				if checkItem(MysqlInfo, "Issues", "Id", JIRAID) == 0 {

					if flag.NFlag() >= 6 {
						// log.Println("start Issues")
						// log.Println(insertdata)
						InsertTable(MysqlInfo, "Issues", insertcols, insertstr, insertdata)
					}

				} else {
					if flag.NFlag() >= 6 {
						// log.Println("update Issues")
						var (
							updatecols    string
							updatecolname []string
							insertvals    []interface{}
						)
						db, err := sql.Open("mysql", MysqlInfo)
						cols, err := db.Query("SELECT * FROM `Issues` WHERE Id ='" + JIRAID + "'")
						if err != nil {
							// panic(err.Error())
							log.Fatal(err.Error())
						}
						updatecolname, err = cols.Columns()
						cols.Close()
						if err != nil {
							// panic(err.Error())
							log.Fatal(err.Error())
						}
						updatecounter := 0
						for index, value := range updatecolname {
							if value != "Sn" && value != "Id" && value != "Key" && value != "FixVersions" && value != "Version" && value != "Component" && value != "Labels" {
								Inputvalue, err := db.Query("SELECT " + value + " FROM `Issues` WHERE `Id`='" + JIRAID + "'")
								if err != nil {
									log.Fatal(err.Error())
									// panic(err.Error())
								}
								var inputtmp string
								Inputvalue.Next()
								Inputvalue.Scan(&inputtmp)
								if inputtmp != insertdata[index] {
									// log.Println(value, inputtmp, insertdata[index])
									if updatecounter == 0 {
										updatecols = updatecols + value + "=? "
										insertvals = append(insertvals, insertdata[index].(string))
									} else {
										updatecols = updatecols + ", " + value + "=? "
										insertvals = append(insertvals, insertdata[index].(string))
									}
									updatecounter++
								}
								Inputvalue.Close()
							}

						}
						sqlStr := "update Issues set " + updatecols + " where Id='" + JIRAID + "'"
						// log.Println(sqlStr)
						// log.Println(insertvals)
						stmt, err := db.Prepare(sqlStr)
						if err != nil {
							log.Fatal(err.Error())
							// panic(err.Error())
						}
						_, err = stmt.Exec(insertvals...)
						// res, err := stmt.Exec(insertvals...)
						if err != nil {
							log.Fatal(err.Error())
						}
						// log.Println(res.LastInsertId())
						err = stmt.Close()
						if err != nil {
							log.Fatal(err.Error())
						}

						db.Close()

					}
				}

				for y := range jiraChangeLog.ChangeLog.Histories {
					ChangeLogID := jiraChangeLog.ChangeLog.Histories[y].ID
					// ChangeLogCreated := jiraChangeLog.ChangeLog.Histories[y].Created
					// User := jiraChangeLog.ChangeLog.Histories[y].Author.Name
					// Field := jiraChangeLog.ChangeLog.Histories[y].Items[0].Field
					// FieldType := jiraChangeLog.ChangeLog.Histories[y].Items[0].FieldType
					// From := jiraChangeLog.ChangeLog.Histories[y].Items[0].From
					// Fromstring := jiraChangeLog.ChangeLog.Histories[y].Items[0].FromString
					// To := jiraChangeLog.ChangeLog.Histories[y].Items[0].To
					// Tostring := jiraChangeLog.ChangeLog.Histories[y].Items[0].ToString
					// log.Println(jiraChangeLog.ChangeLog.Histories[y].Items)
					// fmt.Printf("   ==>  ID -> %s , Created -> %s  \n", ChangeLogID, ChangeLogCreated)
					// fmt.Printf("   ==>  JID -> %s , User -> %s  \n", User, User)
					// fmt.Printf("   ==>  field -> %s , Fieldtype -> %s  \n", Field, FieldType)
					// fmt.Printf("   ==>  form -> %s , formstring -> %s  \n", From, Fromstring)
					// fmt.Printf("   ==>  to -> %s , Tostring -> %s  \n", To, Tostring)

					// update all item
					Historyscols, Historysstr, HistorysColSize := SQLInsertStr(MysqlInfo, "SELECT * FROM Historys")
					Historysdata := make([]interface{}, HistorysColSize)

					if SQLcheckItem(MysqlInfo, "select sn from Historys where `Id`='"+ChangeLogID+"'") == 0 {
						Historysdata[0] = ""
						Historysdata[1] = jiraChangeLog.ChangeLog.Histories[y].ID
						Historysdata[2] = JIRAID
						Historysdata[3] = jiraChangeLog.ChangeLog.Histories[y].Created
						Historysdata[4] = jiraChangeLog.ChangeLog.Histories[y].Author.Name
						Historysdata[5] = jiraChangeLog.ChangeLog.Histories[y].Items[0].Field
						Historysdata[6] = jiraChangeLog.ChangeLog.Histories[y].Items[0].FieldType
						Historysdata[7] = jiraChangeLog.ChangeLog.Histories[y].Items[0].From
						Historysdata[8] = jiraChangeLog.ChangeLog.Histories[y].Items[0].FromString
						Historysdata[9] = jiraChangeLog.ChangeLog.Histories[y].Items[0].To
						Historysdata[10] = jiraChangeLog.ChangeLog.Histories[y].Items[0].ToString
						InsertTable(MysqlInfo, "Historys", Historyscols, Historysstr, Historysdata)
					} else {
						tmpHID := SQLcheckItem(MysqlInfo, "select  from Historys where `JId`='"+JIRAID+"'")
						HID := strconv.Itoa(tmpHID)
						sqlsstr := "update Historys set  " +
							"`Id`='" + jiraChangeLog.ChangeLog.Histories[y].ID + "' ," +
							"`Jid`='" + JIRAID + "' ," +
							"`Created`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Created) + "' ," +
							"`User`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Author.Name) + "' ," +
							"`Field`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].Field) + "' ," +
							"`FieldType`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].FieldType) + "' ," +
							"`ColumnFrom`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].From) + "' ," +
							"`Fromstring`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].FromString) + "' ," +
							"`ColumnTo`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].To) + "' ," +
							"`Tostring`='" + template.HTMLEscapeString(jiraChangeLog.ChangeLog.Histories[y].Items[0].ToString) + "' " +
							" where `Sn`='" + HID + "' "
						// log.Println(sqlsstr)
						UpdateTable(MysqlInfo, sqlsstr)
					}
				}
			}
		} else {
			CommandName := strings.Split(os.Args[0], "/")
			fmt.Printf("Please chekc help file by \" %s -h \" \n", CommandName[len(CommandName)-1])
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
