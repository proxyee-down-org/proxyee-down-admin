package handles

import (
	"common"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type payload struct {
	Commits []struct {
		Added    []string
		Modified []string
	}
}

func checkSignature(payloadBts []byte, signature string) bool {
	if len(signature) == 0 {
		return false
	}
	mac := hmac.New(sha1.New, []byte(common.Config.Extension.Secret))
	mac.Write(payloadBts)
	return hex.EncodeToString(mac.Sum(nil)) == signature[5:]
}

const MANIFEST = "manifest.json"

func HandelPayload(payload *payload) {
	//更新插件仓库
	cmd := exec.Command("/bin/sh", "-c", "cd "+common.Config.Extension.Repo+";git pull")
	_, err := cmd.Output()
	if err == nil {
		//扫描入库
		var rootDir = common.Config.Extension.Repo
		files, _ := ioutil.ReadDir(rootDir)
		//遍历所有插件目录，检测版本号是否有变动
		for _, f := range files {
			if f.IsDir() && f.Name()[0:1] != "." {
				manifestBts, _ := ioutil.ReadFile(filepath.Join(rootDir, f.Name(), MANIFEST))
				//读取manifest.json
				var manifestData models.Extension
				//取插件目录下所有文件
				var otherFileNames = ""
				var extensionDir = filepath.Join(rootDir, f.Name())
				filepath.Walk(extensionDir, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						fileName := strings.Replace(path[len(extensionDir):], "\\", "/", -1)
						if fileName[1:2] != "." {
							if len(otherFileNames) != 0 {
								otherFileNames += ","
							}
							otherFileNames += fileName
						}
					}
					return nil
				})
				if json.Unmarshal(manifestBts, &manifestData) == nil {
					manifestData.Files = otherFileNames
					//取当前数据库里的版本号
					var path = "/" + f.Name()
					extensionData, err := models.SelectExtensionByPath(path)
					if err == nil {
						if extensionData != nil {
							//需要更新
							if manifestData.Version > extensionData.Version {
								manifestData.Id = extensionData.Id
								manifestData.UpdateTime = time.Now()
								manifestData.Update()
							}
						} else {
							//第一次入库
							manifestData.Path = path
							manifestData.CreateTime = time.Now()
							manifestData.Insert()
						}
					}
				}
			}
		}
	}
}

func WebHook(w http.ResponseWriter, r *http.Request) {
	payloadBts, _ := ioutil.ReadAll(r.Body)
	signature := r.Header.Get("X-Hub-Signature")
	//验证WebHook合法性
	if checkSignature(payloadBts, signature) {
		//读取响应
		var payload payload
		json.Unmarshal(payloadBts, &payload)
		HandelPayload(&payload)
		io.WriteString(w, "ok")
	} else {
		io.WriteString(w, "fail")
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	r.ParseForm()
	pageSize, err := strconv.Atoi(r.Form.Get("pageSize"))
	if err != nil {
		pageSize = 1
	}
	pdVersion := r.Form.Get("version")
	pdVersionFloat64 := 3.12
	if len(pdVersion) != 0 {
		pdVersionFloat64, err = strconv.ParseFloat(pdVersion, 64)
	}
	page, err := models.SelectExtensionByKeyword(r.Form.Get("keyword"), pdVersionFloat64, pageSize, 10)
	if err == nil {
		bts, err := json.Marshal(page)
		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(bts)
			return
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	w.WriteHeader(500)
	io.WriteString(w, "params error")
}

func CheckExtensionUpdate(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	if strings.EqualFold(r.Method, "options") {
		return
	}
	payloadBts, _ := ioutil.ReadAll(r.Body)
	if len(payloadBts) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, "params error")
		return
	}
	var extensionChecks []models.ExtensionCheck
	err := json.Unmarshal(payloadBts, &extensionChecks)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "params error")
		return
	}
	extensions, err := models.CheckExtensionUpdate(extensionChecks)
	if err == nil {
		bts, err := json.Marshal(extensions)
		if err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(bts)
			return
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	w.WriteHeader(500)
	io.WriteString(w, "params error")
}

func Down(w http.ResponseWriter, r *http.Request) {
	models.SetCORS(w)
	r.ParseForm()
	extId := r.Form.Get("ext_id")
	version := r.Form.Get("version")
	pdVersion := r.Form.Get("pd_version")
	if len(extId) > 0 {
		extIdInt64, err := strconv.ParseInt(extId, 10, 64)
		versionFloat64, err := strconv.ParseFloat(version, 64)
		pdVersionFloat64 := 3.12
		if len(pdVersion) != 0 {
			pdVersionFloat64, err = strconv.ParseFloat(pdVersion, 64)
		}
		if err == nil {
			extensionDown := models.ExtensionDown{
				ExtId:      extIdInt64,
				Version:    versionFloat64,
				PdVersion:  pdVersionFloat64,
				Ip:         models.GetIp(r),
				Ua:         r.UserAgent(),
				CreateTime: time.Now(),
			}
			extensionDown.Insert()
			w.WriteHeader(200)
			io.WriteString(w, "ok")
			return
		}
	}
	w.WriteHeader(400)
	io.WriteString(w, "params error")
}
