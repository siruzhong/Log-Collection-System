package tail

import (
	"github.com/sirupsen/logrus"
	"logs-collection-system/conf"
)

// tailTaskManager tailTask管理者结构体
type tailTaskManager struct {
	tailTaskMap   map[string]*tailTask  // 所有的tailTask集合
	logConfigList []conf.LogConfig      // 所有配置项
	confChan      chan []conf.LogConfig // 配置管道:用来等待新配置
}

// ttManager 全局tailTask管理者
var ttManager *tailTaskManager

// Init 初始化全局ttManager变量
func Init(logConfigList []conf.LogConfig) (err error) {
	// 创建日志收集任务管理者
	ttManager = &tailTaskManager{
		tailTaskMap:   make(map[string]*tailTask, 30),
		logConfigList: logConfigList,
		confChan:      make(chan []conf.LogConfig, 32),
	}
	// 遍历所有配置项,为每个配置项创建一个日志收集任务
	for _, logConfig := range logConfigList {
		tailTask := newTailTask(logConfig)                  // 创建日志收集任务
		ttManager.tailTaskMap[tailTask.filepath] = tailTask // 存储该日志收集任务
		logrus.Infof("create a tail task for file:%s", logConfig.Path)
		// 起一个后台goroutine将日志写入msgChan
		go tailTask.readFileToChan()
	}
	// 起一个goroutine监听获取新配置
	go ttManager.watchConf()
	return
}

// watchConf 循环监听获取新配置
func (t *tailTaskManager) watchConf() {
	for {
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd,conf:%v, start manage tailTask...", newConf)
		logrus.Infof("old conf:%v", t.tailTaskMap)
		// 遍历新配置的每个配置项
		for _, conf := range newConf {
			// 1.原来有的tailTask不变
			_, ok := t.tailTaskMap[conf.Path]
			if ok {
				continue
			}
			// 2.原来没有的tailTask需要去创建
			tailTask := newTailTask(conf)
			t.tailTaskMap[conf.Path] = tailTask
			go tailTask.readFileToChan()
		}
		// 3.原来有的现在没有的tailTask需要删除
		for path, tt := range t.tailTaskMap {
			var flag bool
			for _, conf := range newConf {
				if path == conf.Path {
					flag = true
					continue
				}
			}
			if !flag {
				logrus.Infof("the tail task is going to stop,path:%s", tt.filepath)
				delete(t.tailTaskMap, path) // 删除
				tt.cancel()
			}
		}
	}
}

// SendNewConf 接受新配置
func SendNewConf(newConf []conf.LogConfig) {
	ttManager.confChan <- newConf
}
