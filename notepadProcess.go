package main

import (
	"sync"
	"os"
	"sync/atomic"
	"github.com/smtc/glog"
	"github.com/howeyc/fsnotify"
	"github.com/guotie/deferinit"
)

func init() {
	deferinit.AddRoutine(notepadProcess)
}

type counter struct {
	val int32
}

func (c *counter) increment() {
	atomic.AddInt32(&c.val, 1)
}

/**
记事本处理文件
如当前目录不存在则不做任何处理,该方法直接不做任何处理
创建人:邵炜
创建时间:2016年4月12日11:37:41
*/
func notepadProcess(ch chan struct{}, wg *sync.WaitGroup) {
	fi, err := os.Stat(notepadProcessDir)
	if err != nil {
		glog.Error("notepadProcess: file data is error! err: %s \n", err.Error())
		return
	}
	if !fi.IsDir() {
		glog.Error("notepadProcess: message file name :%s is not defind! \n", notepadProcessDir)
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Error("notepadProcess: fsnotify newWatcher is error! err: %s \n", err.Error())
		return
	}
	var  modifyReceived counter
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				glog.Info("messageCenterFtp: fsnotify watcher fileName: %s is change!  ev: %v \n", ev.Name, ev)
				if ev.IsModify() {
					modifyReceived.increment()
					if modifyReceived.val % 2 == 0 {
						go func(filePath string) {
							readFileMobs(ev.Name)
						}(ev.Name)
					}
				}
			case err := <-watcher.Error:
				glog.Error("messageCenterFtp: fsnotify watcher is error! err: %s \n", err.Error())
			}
		}
		done <- true
	}()
	err = watcher.WatchFlags(notepadProcessDir,fsnotify.FSN_MODIFY)
	if err != nil {
		glog.Error("messageCenterFtp watch error. messageCenterDir: %s  err: %s \n", notepadProcessDir, err.Error())
	}

	// Hang so program doesn't exit
	<-ch

	/* ... do stuff ... */
	watcher.Close()
	wg.Done()
}