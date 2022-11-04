package controller

import (
	"gateway/common/response"
	"gateway/dao"
	"gateway/dto"
	"github.com/gin-gonic/gin"
)

type APPController struct {
}

// APPList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (admin *APPController) APPList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	info := &dao.App{}
	list, total, err := info.APPList(params)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}

	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
		//appCounter, err := tools.FlowCounterHandler.GetCounter(tools.FlowAppPrefix + item.AppID)
		//if err != nil {
		//	response.FailMsg(err.Error())
		//	c.Abort()
		//	return
		//}
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			//RealQpd:  appCounter.TotalCount,
			//RealQps:  appCounter.QPS,
		})
	}
	output := dto.APPListOutput{
		List:  outputList,
		Total: total,
	}
	response.SuccessMsg(output)
	return
}

// APPDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/app_detail [get]
func (admin *APPController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(search)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	response.SuccessMsg(detail)
	return
}

// APPDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (admin *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(search)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	info.IsDelete = 1
	response.SuccessMsg("")
	return
}

// AppAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (admin *APPController) AppAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//验证app_id是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := search.Find(search); err == nil {
		response.FailMsg("租户ID被占用，请重新输入")
		return
	}

	info := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(); err != nil {
		response.FailMsg(err.Error())
		return
	}
	response.SuccessMsg("")
	return
}

// AppUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [post]
func (admin *APPController) AppUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(search)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}

	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(); err != nil {
		response.FailMsg(err.Error())
		return
	}
	response.SuccessMsg("")
	return
}

// AppStatistics godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept  json
// @Produce  json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.StatisticsOutput} "success"
// @Router /app/app_stat [get]
func (admin *APPController) AppStatistics(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//search := &dao.App{
	//	ID: params.ID,
	//}
	//	detail, err := search.Find( search)
	//if err != nil {
	//	response.FailMsg(err.Error())
	//	return
	//}
	//
	//////今日流量全天小时级访问统计
	var todayStat []int64
	//////counter, err := tools.FlowCounterHandler.GetCounter(tools.FlowAppPrefix + detail.AppID)
	//////if err != nil {
	//////	response.FailMsg(err.Error())
	//////	c.Abort()
	//////	return
	//////}
	////currentTime:= time.Now()
	//for i := 0; i <= time.Now().In(lib.TimeLocation).Hour(); i++ {
	//	dateTime:=time.Date(currentTime.Year(),currentTime.Month(),currentTime.Day(),i,0,0,0,lib.TimeLocation)
	//	hourData,_:=counter.GetHourData(dateTime)
	//	todayStat = append(todayStat, hourData)
	//}
	//
	////昨日流量全天小时级访问统计
	var yesterdayStat []int64
	//yesterTime:= currentTime.Add(-1*time.Duration(time.Hour*24))
	//for i := 0; i <= 23; i++ {
	//	dateTime:=time.Date(yesterTime.Year(),yesterTime.Month(),yesterTime.Day(),i,0,0,0,lib.TimeLocation)
	//	hourData,_:=counter.GetHourData(dateTime)
	//	yesterdayStat = append(yesterdayStat, hourData)
	//}
	stat := dto.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	response.SuccessMsg(stat)
	return
}
