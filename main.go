package main

import (
	"flag"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

var (
	configFn    = flag.String("config", "./config.json", "config file path") //配置文件地址
	debugFlag   = flag.Bool("d", false, "debug mode")                        //是否为调试模式
	notepadProcessDir string //监控目录
	interFaceApi string //订购接口请求地址
)

/**
主函数入口
创建人:邵炜
创建时间:2016年2月26日11:22:03
*/
func main() {

	if checkPid() { //判断程序是否启动
		return
	}

	flag.Parse()

	serverRun(*configFn, *debugFlag)

	c := make(chan os.Signal, 1)
	writePid()
	// 信号处理
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 等待信号
	<-c

	serverExit()
	rmPidFile()
	os.Exit(0)
}

/**
服务启动
创建人:邵炜
创建时间:2016年2月26日11:22:16
输入参数: cfn(配置文件地址) debug(是否为调试模式)
*/
func serverRun(cfn string, debug bool) {
	config.ReadCfg(cfn)
	logInit(debug)

	notepadProcessDir=strings.TrimSpace(config.GetStringMust("notepadProcessDir"))
	interFaceApi=strings.TrimSpace(config.GetStringMust("interFaceApi"))

	//初始化所有go文件中的init内方法
	deferinit.InitAll()
	glog.Info("init all module successfully \n")

	//设置多CPU运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	glog.Info("set many cpu successfully \n")

	//启动所有go文件中的init方法
	deferinit.RunRoutines()
	glog.Info("init all run successfully \n")
}

/**
结束进程
创建人:邵炜
创建时间:2016年3月7日14:21:24
*/
func serverExit() {
	// 结束所有go routine
	deferinit.StopRoutines()
	glog.Info("stop routine successfully.\n")

	deferinit.FiniAll()
	glog.Info("fini all modules successfully.\n")

	glog.Close()
}


