package tail

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logs-collection-system/conf"
	"logs-collection-system/kafka"
	"time"
)

// tailTask 日志收集任务结构体
type tailTask struct {
	filepath string     // 文件路径
	topic    string     // 写入etcd的主题
	tail     *tail.Tail // 文件追踪
	ctx      context.Context
	cancel   context.CancelFunc
}

// newTailTask 创建一个日志收集任务
func newTailTask(logConfig conf.LogConfig) (tt *tailTask) {
	// 创建tail对象
	tailConfig := tail.Config{
		ReOpen: true, // 重新打开变更的文件
		Follow: true, // 持续等待文件中的新行
		Location: &tail.SeekInfo{
			Offset: 0, // whence表示给offset参数一个定义，表示要从哪个位置开始偏移
			Whence: 2, // 0代表从文件开头开始算起，1代表从当前位置开始算起，2代表从文件末尾算起
		},
		MustExist: false, // 文件可以为空
		Poll:      true,  // 轮询文件更改
	}
	tail, err := tail.TailFile(logConfig.Path, tailConfig)
	if err != nil {
		logrus.Errorf("Tail:create tail for path:%s err=%v", logConfig.Path, err)
		return
	}
	// 创建context和cancel
	ctx, cancel := context.WithCancel(context.Background())
	// 创建tailTask对象
	tt = &tailTask{topic: logConfig.Topic, filepath: logConfig.Path, tail: tail, ctx: ctx, cancel: cancel}
	return
}

// readFileToChan 读取文件中最新内容写入msgChan
func (tt *tailTask) readFileToChan() (err error) {
	// 循环读取文件数据
	for {
		select {
		// 退出goroutine的条件:只要调用了ctx.cancel(),就会收到信号
		case <-tt.ctx.Done():
			logrus.Infof("path:%s is stopped...", tt.filepath)
			return
		// 循环读取文件数据
		case line, ok := <-tt.tail.Lines:
			if !ok {
				logrus.Errorf("tail file close reopen, filename:%s\n", tt.filepath)
				time.Sleep(time.Second) // 读取出错等1s
				continue
			}
			// 读取文件最新行首先写入msgChan管道，然后kafka Client写入kafka，使两步成为异步操作
			msg := &sarama.ProducerMessage{}
			msg.Topic = tt.topic
			msg.Value = sarama.StringEncoder(line.Text)
			// 将消息写入管道
			kafka.SendMesToChan(msg)
		}
	}
	return
}
