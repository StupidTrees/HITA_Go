package service

import (
	"encoding/json"
	"fmt"
	repo "hita/repository"
	"hita/utils/api"
	"time"
)

type SyncReq struct {
	Uid      string `form:"uid" json:"uid"`
	LatestId int64  `form:"latestId" json:"latestId"`
}
type RespSync struct {
	Action   string `json:"action"`   //push:客户端上传，pull：客户端拉取
	LatestId int64  `json:"latestId"` //服务端最新的id
	History  string `json:"history"`
	Data     string `json:"data"`
}

type PushReq struct {
	Uid     int64  `form:"uid" json:"uid"`
	History string `form:"history" json:"history"`
	Data    string `form:"data" json:"data"`
}

func (req *SyncReq) Sync() (result RespSync, code int, error error) {
	latestIdServer, err := repo.GetLatestId(req.Uid)
	if err != nil || latestIdServer < req.LatestId { //不存在，或本地超过服务器
		println("push required")
		result.Action = "PUSH"
		result.LatestId = latestIdServer
	} else if latestIdServer == req.LatestId {
		println("none required")
		result.Action = "NONE"
	} else {
		println("pull required")
		result.Action = "PULL"
		result.LatestId = latestIdServer
		histories := repo.GetHistoriesAfter(req.Uid, req.LatestId)
		//fmt.Print(histories)
		requiredIds := map[string]map[string]string{}
		for _, h := range histories {
			if requiredIds[h.Table] == nil {
				requiredIds[h.Table] = map[string]string{}
			}
			if h.Action == "REQUIRE" {
				for _, id := range h.GetIds() {
					requiredIds[h.Table][id] = id
				}
			} else {
				for _, id := range h.GetIds() {
					delete(requiredIds[h.Table], id)
				}
			}
		}

		data := getDataForIds(req.Uid, requiredIds)
		hj, _ := json.Marshal(histories)
		dj, _ := json.Marshal(data)
		result.History = string(hj)
		result.Data = string(dj)
		fmt.Printf("result:%v\n", result)
	}
	code = api.CodeSuccess
	return
}

func (req *PushReq) Push(uid int64, historyList []repo.History, dataMap map[string][]interface{}) (code int, error error) {
	repo.SaveHistories(historyList)
	for _, history := range historyList {
		fmt.Printf("history:%v,%v,%v\n", history.Id, history.Table, history.Action)
		switch history.Action {
		case "REQUIRE":
			{
				switch history.Table {
				case "timetable":
					{
						tts := findDataForIds(history, dataMap, "id")
						var list []repo.Timetable
						for _, ts := range tts {
							tt := repo.Timetable{}
							js, _ := json.Marshal(ts)
							_ = json.Unmarshal(js, &tt)
							//
							tsJs := ts.(map[string]interface{})
							//tt.Id = tsJs["id"].(string)
							tt.UserId = uid
							//tt.Code = tsJs["code"].(string)
							//tt.Name = tsJs["name"].(string)
							//tt.StartTime = repo.MTime(time.Unix(int64(tsJs["startTime"].(float64))/1000, 0))
							//tt.EndTime = repo.MTime(time.Unix(int64(tsJs["endTime"].(float64))/1000, 0))
							st, _ := json.Marshal(tsJs["scheduleStructure"])
							tt.ScheduleStructure = repo.SString(st)
							list = append(list, tt)
						}
						repo.AddTimetables(list)
					}
				case "event":
					{
						tts := findDataForIds(history, dataMap, "id")
						var list []repo.Event
						for _, ts := range tts {
							tt := repo.Event{}
							js, _ := json.Marshal(ts)
							_ = json.Unmarshal(js, &tt)
							tsJs := ts.(map[string]interface{})
							tt.Id = tsJs["id"].(string)
							tt.SubjectId = tsJs["subjectId"].(string)
							tt.TimetableId = tsJs["timetableId"].(string)
							tt.Type = tsJs["type"].(string)
							tt.FromNumber = int8(tsJs["fromNumber"].(float64))
							tt.LastNumber = int8(tsJs["lastNumber"].(float64))
							tt.UserId = uid
							//tt.Teacher = tsJs["teacher"].(string)
							//tt.Place = tsJs["place"].(string)
							tt.Name = tsJs["name"].(string)
							tt.From = repo.MTime(time.Unix(int64(tsJs["from"].(float64))/1000, 0))
							tt.To = repo.MTime(time.Unix(int64(tsJs["to"].(float64))/1000, 0))
							list = append(list, tt)
						}
						repo.AddEvents(list)
					}
				case "subject":
					{
						tts := findDataForIds(history, dataMap, "id")
						var list []repo.TermSubject
						for _, ts := range tts {
							tt := repo.TermSubject{}
							js, _ := json.Marshal(ts)
							_ = json.Unmarshal(js, &tt)
							tt.UserId = uid
							//tt.TimetableId = tsJs["timetableId"].(string)
							//tt.Id = tsJs["id"].(string)
							//tt.Type = tsJs["type"].(string)
							//tt.Name = tsJs["name"].(string)
							//tt.Field = tsJs["field"].(string)
							//tt.School = tsJs["school"].(string)
							//tt.Code = tsJs["code"].(string)
							//tt.Key = tsJs["key"].(string)
							//tt.Color = tsJs["color"].(int32)
							//tt.Credit = tsJs["place"].(float32)
							//tt.CountInSPA = tsJs["countInSPA"].(bool)
							list = append(list, tt)
						}
						repo.AddSubjects(list)
					}
				}
			}
		case "REMOVE":
			{
				switch history.Table {
				case "timetable":
					{
						repo.RemoveTimetablesInIds(history.GetIds())
					}
				case "event":
					{
						repo.RemoveEventsInIds(history.GetIds())
					}
				case "subject":
					{
						repo.RemoveSubjectInIds(history.GetIds())
					}

				}
			}
		}
	}
	code = api.CodeSuccess
	return
}

func findDataForIds(history repo.History, dataMap map[string][]interface{}, keyName string) []interface{} {
	dataList := dataMap[history.Table]
	var res []interface{}
	for _, id := range history.GetIds() {
		for _, data := range dataList {
			if data.(map[string]interface{})[keyName].(string) == id {
				res = append(res, data)
			}
		}
	}
	return res
}

func getDataForIds(uid string, requireIds map[string]map[string]string) map[string][]interface{} {
	result := map[string][]interface{}{}
	for key := range requireIds {
		if result[key] == nil {
			result[key] = []interface{}{}
		}
		var ids []string
		for k := range requireIds[key] {
			ids = append(ids, k)
		}
		//fmt.Printf("require:%v:%v\n", key, ids)
		switch key {
		case "timetable":
			{
				for _, event := range repo.GetTimetablesInIds(uid, ids) {
					js, _ := json.Marshal(event)
					result[key] = append(result[key], string(js))
				}
			}
		case "subject":
			{
				for _, event := range repo.GetSubjectsInIds(uid, ids) {
					js, _ := json.Marshal(event)
					result[key] = append(result[key], string(js))
				}
			}
		case "event":
			{
				for _, event := range repo.GetEventsInIds(uid, ids) {
					js, _ := json.Marshal(event)
					result[key] = append(result[key], string(js))
				}
			}

		}
	}
	return result
}
