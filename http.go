package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	TYPE_json = 1
	TYPE_text = 2
)

const tpl = `<html>
<head>
    <title>上传文件</title>
</head>
<body>

<form enctype="multipart/form-data" action="/upload" method="post">

<table width="70%" border="1" cellspacing="0" cellpadding="0">
  <tr>
    <td>  <input type="file" name="uploadfile" /></td>
    </tr>
  <tr>
    <td>
    <input type="hidden" name="token" value=""/>
    <input type="submit" value="upload" />
    </td>
  </tr>
  <tr>
    <td><a href=files>文件访目录</a></td>
  </tr>
</table>


  
  


</form>


</body>
</html`

/*获取当前时间*/
func gettime() string {
	Year := time.Now().Year()     //年[:3]
	Month := time.Now().Month()   //月
	Day := time.Now().Day()       //日
	Hour := time.Now().Hour()     //小时
	Minute := time.Now().Minute() //分钟
	Second := time.Now().Second() //秒
	//Nanosecond:=time.Now().Nanosecond()//纳秒
	var timestr string
	timestr = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", Year, Month, Day, Hour, Minute, Second)
	return timestr
}

/*获取系统当前时间戳*/
func gettimecuo() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UnixNano(), 10)
	timestamp = timestamp[0:13]
	//fmt.Println(timestamp)
	//fmt.Println(t.Unix())
	return timestamp
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	fmt.Println(gettime(), "method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		//token := fmt.Sprintf("%x", h.Sum(nil))

		//t, _ := template.ParseFiles("upload.gtpl")
		//t.Execute(tpl, token)
		DocumentWrite(w, tpl, TYPE_text)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		//filepath := "./upload/" + handler.Filename

		ext := strings.Split(handler.Filename, ".")
		filepath := gettimecuo() + "." + ext[1]
		tab := `<div align="center">
			<table width="70%" border="1" cellspacing="0" cellpadding="0">
			  <tr>
			    <td>文件名</td>
			    <td>`+"<a href=./files/"+filepath+">"+filepath+"</a>"+` </td>
			  </tr>
			  <tr>
			    <td>路径</td>
			    <td><input type=text value=`+"./files/"+filepath+` style="width:100%;"> </td>
			  </tr>
			</table>
			</div>`



		//DocumentWrite(w, "<a href=./files/"+filepath+">文件"+filepath+"已经上传，右击复制链接地址即可</a>", TYPE_text)
		DocumentWrite(w, tab, TYPE_text)
		fmt.Println(filepath)
		f, err := os.OpenFile("./upload/"+filepath, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在upload目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

/*获取当前路径
"path/filepath"
"strings" //需要引入2个库
*/
func getCurrentDir(file string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	ret := strings.Replace(dir, "\\", "/", -1)
	ret += "/" + file
	return ret
}

/* 判断文件是否存在  存在返回 true 不存在返回false*/
func File_Exists(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/*保存文件（优化版）*/
func SaveLog(m_FilePath string, val string) {
	var dir, filename string
	filename = filepath.Base(m_FilePath)
	if len(m_FilePath) > 1 && string([]byte(m_FilePath)[1:2]) == ":" {
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(m_FilePath, filename)
		//print("abspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	} else {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		dir = dir + "/" + m_FilePath
		filename = filepath.Base(m_FilePath)
		dir = strings.TrimSuffix(dir, filename)
		//print("noptabspath:filename:" + filename + "\n" + "dir:" + dir + "\n")
	}

	p := dir + "/" + filename
	p = strings.Replace(p, "\\", "/", -1)
	p = strings.Replace(p, "//", "/", -1)
	//print("fullpath" + p + "\n")
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
	fl, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE, 0644)
	defer fl.Close()

	if err != nil {
		fmt.Println("SaveLog:error")
	} else {
		io.WriteString(fl, val)
	}
}

//文档返回值写出
func DocumentWrite(res http.ResponseWriter, val string, mtype int) {
	//写出返回格式
	if mtype == TYPE_json {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
	} else if mtype == TYPE_text {
		res.Header().Set("Content-Type", "text/html;charset=utf-8")
	} else {
		res.Header().Set("Content-Type", "text/html;charset=utf-8")
	}

	//写出网页响应码
	res.WriteHeader(200)
	//写出结果
	res.Write([]byte(val))
	//服务控制台输出
	//fmt.Println(val)
}

//文档跳转值
func DocumentRedirect(res http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(res, req, url, http.StatusFound)
}
func homepage(res http.ResponseWriter, req *http.Request) { //HOME
	DocumentWrite(res, tpl, TYPE_text)
}

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

 	
	wd = wd + "/upload"
	os.Mkdir(wd, os.ModePerm)  
	fmt.Println(wd)
	fs := http.FileServer(http.Dir(wd))

	mux := http.NewServeMux()
	mux.Handle("/files/", http.StripPrefix("/files", fs))
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/", homepage)

	//设置访问的路由
	fmt.Println(gettime(), "服务器80开始服务。。。")
	err = http.ListenAndServe(":80", mux) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
