/**
 * @Author: niaoshuai
 * @Date: 2020/8/9 8:51 上午
 */
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	logger "svn-hook/pkg/log"
)

var (
	repos           string
	txn             string
	jenkinsDataPath string
	logPath         string
	svnExecPath     string
)

type Jenkins struct {
	description string        `json:"description"`
	Data        []JenkinsData `json:"data"`
}

type JenkinsData struct {
	SvnServer     string `json:"svn_server"`
	JenkinsServer string `json:"jenkins_server"`
}

func main() {
	// flag
	flag.StringVar(&repos, "REPOS", "", "REPOS")
	flag.StringVar(&txn, "TXN", "", "TXN")
	flag.StringVar(&jenkinsDataPath, "JENKINS_DATA", "jenkins-data.json", "jenkinsData")
	flag.StringVar(&logPath, "LOG_PATH", "svn-hook.log", "log")
	flag.StringVar(&svnExecPath, "SVNLOOK_PATH", "svnlook", "svnlook path")
	flag.Parse()

	// check value
	if jenkinsDataPath == "" {
		log.Fatal("jenkinsData not null")
	}

	if repos == "" {
		log.Fatal("REPOS not null")
	}

	if txn == "" {
		log.Fatal("TXN not null")
	}

	// 初始化日志
	logger.InitLog(logPath)
	// 读取文件
	jenkins := readJenkinsData(jenkinsDataPath)
	// exec svn log command
	svnLookLog()
	// exec svn dirChange commandLookLog()
	dirChangeLog := svnLookDirChangeLog()

	for _, jenkinsData := range jenkins.Data {
		// check
		if strings.Contains(dirChangeLog, jenkinsData.SvnServer) {
			// 发送 Http 调用
			jenkinsApi(jenkinsData.JenkinsServer)
		}
	}
}

// 调用 JENKINS JS API
func jenkinsApi(url string) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "123456")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bodyBytes))
}

// 读取 jenkins data
func readJenkinsData(jenkinsDataPath string) *Jenkins {
	f, err := os.Open(jenkinsDataPath)
	defer f.Close()
	if err != nil {
		logger.Fatal(err)
	}
	r := io.Reader(f)
	ret := &Jenkins{}
	if err = json.NewDecoder(r).Decode(ret); err != nil {
		logger.Fatal(err)
	}
	return ret
}

// 读取 svn 目录改变
func svnLookDirChangeLog() string {
	svnLookDirChangeCmd := exec.Command(svnExecPath, "dirs-changed", "-r", txn, repos)
	var out2 bytes.Buffer
	svnLookDirChangeCmd.Stdout = &out2
	err := svnLookDirChangeCmd.Start()
	if err != nil {
		logger.Fatal(err)
	}
	err = svnLookDirChangeCmd.Wait()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(out2.String())
	return out2.String()
}

// 读取 svn log
func svnLookLog() string {
	svnLookLogCmd := exec.Command(svnExecPath, "log", repos, "-r", txn)
	var out1 bytes.Buffer
	svnLookLogCmd.Stdout = &out1
	err := svnLookLogCmd.Start()
	if err != nil {
		logger.Fatal(err)
	}
	err = svnLookLogCmd.Wait()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(out1.String())
	return out1.String()
}
