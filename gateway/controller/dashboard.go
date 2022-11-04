package controller

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/dto"
	"github.com/gin-gonic/gin"
	"time"
)

type DashboardController struct{}

func (service *DashboardController) PanelGroupData(c *gin.Context) {

	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(&dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	app := &dao.App{}
	_, appNum, err := app.APPList(&dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	//counter, err := tools.FlowCounterHandler.GetCounter(tools.FlowTotal)
	//if err != nil {
	//	response.FailMsg(err.Error() 2003, err)
	//	return
	//}
	out := &dto.PanelGroupDataOutput{
		ServiceNum: serviceNum,
		AppNum:     appNum,
		//TodayRequestNum: counter.TotalCount,
		//CurrentQPS:      counter.QPS,
	}
	response.SuccessMsg(out)
}

func (service *DashboardController) ServiceStat(c *gin.Context) {

	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType()
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	legend := []string{}
	//for index, item := range list {
	////	name, ok := tools.LoadTypeMap[item.LoadType]
	//	if !ok {
	//		response.FailMsg("load_type not found")
	//		return
	//	}
	//	list[index].Name = name
	//	legend = append(legend, name)
	//}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	response.SuccessMsg(out)
}

// FlowStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (service *DashboardController) FlowStat(c *gin.Context) {
	//counter, err := tools.FlowCounterHandler.GetCounter(tools.FlowTotal)
	//if err != nil {
	//	response.FailMsg(err.Error() 2001, err)
	//	return
	//}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		//dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		//	hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, 0)
	}

	yesterdayList := []int64{}
	//yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		//	dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		//hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, 0)
	}
	response.SuccessMsg(&dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}
