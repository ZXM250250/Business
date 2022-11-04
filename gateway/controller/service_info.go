package controller

import (
	"errors"
	"fmt"
	"gateway/common/response"
	"gateway/dao"
	"gateway/dto"
	"gateway/tools"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ServiceController struct{}

// ServiceList 获取服务列表的接口
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//从db中分页读取基本信息
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(params)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	//格式化输出信息
	var outList []dto.ServiceListItemOutput
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(&listItem)
		if err != nil {
			response.SuccessMsg(err)
			return
		}
		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if serviceDetail.Info.LoadType == tools.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == tools.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == tools.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == tools.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == tools.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == tools.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == tools.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == tools.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		//	counter, err := config.FlowCounterHandler.GetCounter(config.FlowServicePrefix + listItem.ServiceName)
		if err != nil {
			response.SuccessMsg(err)
			return
		}

		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	response.SuccessMsg(out)

}
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err := serviceInfo.Find(serviceInfo)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(serviceInfo)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	response.SuccessMsg(serviceDetail)

}

func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	//读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err := serviceInfo.Find(serviceInfo)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(); err != nil {
		response.FailMsg(err.Error())
		return
	}

}

// ServiceStat 服务统计
func (service *ServiceController) ServiceStat(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}
	//	读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(serviceInfo)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}

	counter, err := tools.FlowCounterHandler.GetCounter(tools.FlowServicePrefix + serviceDetail.Info.ServiceName)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}
	//
	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	response.SuccessMsg(&dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})

}

func (service *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg(errors.New("IP列表与权重列表数量不一致").Error())
		return
	}

	tx := dao.GetDB().Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err := serviceInfo.Find(serviceInfo); err == nil {
		tx.Rollback()
		response.FailMsg("服务已存在")
		return
	}

	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Find(httpUrl); err == nil {
		tx.Rollback()
		response.FailMsg("服务接入前缀或域名已存在")
		return
	}

	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	//serviceModel.ID
	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	loadbalance := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadbalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()

}

func (service *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg("IP列表与权重列表数量不一致")
		return
	}

	tx := dao.GetDB().Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err := serviceInfo.Find(serviceInfo)
	if err != nil {
		tx.Rollback()
		response.FailMsg("服务不存在")
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(serviceInfo)
	if err != nil {
		tx.Rollback()
		response.FailMsg("服务不存在")
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := httpRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = params.RoundType
	loadbalance.IpList = params.IpList
	loadbalance.WeightList = params.WeightList
	loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadbalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()
	response.SuccessMsg("")

}

func (admin *ServiceController) ServiceAddTcp(c *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(infoSearch); err == nil {
		response.FailMsg("服务名被占用，请重新输入")
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(tcpRuleSearch); err == nil {
		response.FailMsg("服务端口被占用，请重新输入")
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(grpcRuleSearch); err == nil {
		response.FailMsg("服务端口被占用，请重新输入")
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg("ip列表与权重设置不匹配")
		return
	}

	tx := dao.GetDB().Begin()
	info := &dao.ServiceInfo{
		//	LoadType:    tools.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	httpRule := &dao.TcpRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := httpRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()
	response.SuccessMsg("")
	return
}

func (admin *ServiceController) ServiceUpdateTcp(c *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg("ip列表与权重设置不匹配")
		return
	}

	tx := dao.GetDB().Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(service)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	tcpRule := &dao.TcpRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()
	response.SuccessMsg("")
	return
}

func (admin *ServiceController) ServiceAddGrpc(c *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(infoSearch); err == nil {
		response.FailMsg("服务名被占用，请重新输入")
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(tcpRuleSearch); err == nil {
		response.FailMsg("服务端口被占用，请重新输入")
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(grpcRuleSearch); err == nil {
		response.FailMsg("服务端口被占用，请重新输入")
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg("ip列表与权重设置不匹配")
		return
	}

	tx := dao.GetDB().Begin()
	info := &dao.ServiceInfo{
		//	LoadType:    tools.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := grpcRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()
	response.SuccessMsg("")
	return
}

func (admin *ServiceController) ServiceUpdateGrpc(c *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := c.ShouldBind(params); err != nil {
		response.FailMsg(err.Error())
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		response.FailMsg("ip列表与权重设置不匹配")
		return
	}

	tx := dao.GetDB().Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(service)
	if err != nil {
		response.FailMsg(err.Error())
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	grpcRule := &dao.GrpcRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	//grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := grpcRule.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(); err != nil {
		tx.Rollback()
		response.FailMsg(err.Error())
		return
	}
	tx.Commit()
	response.SuccessMsg("")
	return
}
