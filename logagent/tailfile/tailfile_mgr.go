package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

type tailTaskMgr struct {
	tailTasks        map[string]*tailTask       //所有的tailTask任务
	collectEntryList []common.CollectEntry      //所有的配置项
	confChan         chan []common.CollectEntry //所有等待新配置的通道

}

var (
	ttMgr *tailTaskMgr
)

// Init 初始化日志管理
func Init(allConf []common.CollectEntry) (err error) {

	ttMgr = &tailTaskMgr{
		tailTasks:        make(map[string]*tailTask, 20),
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		task := newTailTask(conf.Path, conf.Topic)
		err = task.init()
		if err = task.init(); err != nil {
			logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTasks[task.path] = task // 把创建的这个tailTask任务登记在册,方便后续管理
		//挂起一个后台的g去收集日志
		go task.run()
	}
	//启动一个后台的g去等待新的配置的变化
	go ttMgr.watch()
	return

}

func (t *tailTaskMgr) watch() {
	for {
		//这个通道收到信息说明 配置发生了变动
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd, conf:%v, start manage tailTask...", newConf)
		for _, conf := range newConf {
			//原来就存在的配置就不需要改动
			if t.isExist(conf) {
				continue
			}
			//原来没有的我要新建一个tailTask任务
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.init()
			if err != nil {
				logrus.Errorf("create tailObj for path:%s failed, err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			t.tailTasks[tt.path] = tt // 把创建的这个tailTask任务登记在册,方便后续管理

			//挂起一个后台的g去收集日志
			go tt.run()
		}
		//原来没有的现在要把tailTask停掉

		for key, task := range t.tailTasks {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				logrus.Infof("the task collect path:%s need to stop.", task.path)
				delete(t.tailTasks, key) //从管理类中删除掉
				task.cancel()
			}
		}

	}
}

//判断tailTaskMap中是否存在该收集项
func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTasks[conf.Path]
	return ok

}

func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
