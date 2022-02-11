package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
  功能介绍:
    mux.Handle("/files/", http.StripPrefix("/files", fs)) //文件存放位置
    mux.HandleFunc("/upload", upload)                     //文件上传
    mux.HandleFunc("/upfile", upfile)                     //文件上传
    mux.HandleFunc("/", homepage)                         //主页页面
    mux.HandleFunc("/caiji", Caiji)                       //采集功能
    mux.HandleFunc("/phoneupfile", phoneupfile)           //获取上传地址
    mux.HandleFunc("/picbeautiful", picbeautiful)         //图片计算评分//https://www.youxiz.net/d/file/upload/1613722972012.png
*/

const (
	TYPE_json = 1
	TYPE_text = 2
)

type CAIJIDATA struct {
	State    string `json:"state"`
	Msg      string `json:"msg"`
	Id       string `json:"id"`
	Title    string `json:"title"`
	Ftitle   string `json:"ftitle"`
	Titlepic string `json:"titlepic"`
	Newstime string `json:"newstime"`
	NewsUrl  string `json:"newsurl"`
	Content  string `json:"content"`
	Item     string `json:"item"`
	Pitem    string `json:"pitem"`
	Sendtime string `json:"sendtime"`
}

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
</html>`

const home_tpl = `<DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">

    <title>tik后台管理系统</title>
    <style type="text/css">
        body, div, h3, h4, li, ol {
            margin: 0;
            padding: 0;
        }
 
        body {
            font: 14px/1.5 'Microsoft YaHei','微软雅黑',Helvetica,Sans-serif;
            min-width: 1200px;
            background: #f0f1f3;
        }
 
        :focus {
            outline: 0;
        }
 
        h3, h4, strong {
            font-weight: 700;
        }
 
        a {
            color: #428bca;
            text-decoration: none;
        }
            a:hover {
                text-decoration: underline;
            }
 
        .error-page {
            background: #f0f1f3;
            padding: 80px 0 180px;
        }
 
        .error-page-container {
            position: relative;
            z-index: 1;
        }
 
        .error-page-main {
            position: relative;
            background: #f9f9f9;
            margin: 0 auto;
            width: 617px;
            -ms-box-sizing: border-box;
            -webkit-box-sizing: border-box;
            -moz-box-sizing: border-box;
            box-sizing: border-box;
            padding: 50px 50px 70px;
        }
 
            .error-page-main:before {
                content: '';
                display: block;
                background: url(img/errorPageBorder.png?1427783409637);
                height: 7px;
                position: absolute;
                top: -7px;
                width: 100%;
                left: 0;
            }
 
            .error-page-main h3 {
                font-size: 24px;
                font-weight: 400;
                border-bottom: 1px solid #d0d0d0;
            }
 
                .error-page-main h3 strong {
                    font-size: 54px;
                    font-weight: 400;
                    margin-right: 20px;
                }
 
            .error-page-main h4 {
                font-size: 20px;
                font-weight: 400;
                color: #333;
            }
 
        .error-page-actions {
            font-size: 0;
            z-index: 100;
        }
 
            .error-page-actions div {
                font-size: 14px;
                display: inline-block;
                padding: 30px 0 0 10px;
                width: 50%;
                -ms-box-sizing: border-box;
                -webkit-box-sizing: border-box;
                -moz-box-sizing: border-box;
                box-sizing: border-box;
                color: #838383;
            }
 
            .error-page-actions ol {
                list-style: decimal;
                padding-left: 20px;
            }
 
            .error-page-actions li {
                line-height: 2.5em;
            }
 
            .error-page-actions:before {
                content: '';
                display: block;
                position: absolute;
                z-index: -1;
                bottom: 17px;
                left: 50px;
                width: 200px;
                height: 10px;
                -moz-box-shadow: 4px 5px 31px 11px #999;
                -webkit-box-shadow: 4px 5px 31px 11px #999;
                box-shadow: 4px 5px 31px 11px #999;
                -moz-transform: rotate(-4deg);
                -webkit-transform: rotate(-4deg);
                -ms-transform: rotate(-4deg);
                -o-transform: rotate(-4deg);
                transform: rotate(-4deg);
            }
 
            .error-page-actions:after {
                content: '';
                display: block;
                position: absolute;
                z-index: -1;
                bottom: 17px;
                right: 50px;
                width: 200px;
                height: 10px;
                -moz-box-shadow: 4px 5px 31px 11px #999;
                -webkit-box-shadow: 4px 5px 31px 11px #999;
                box-shadow: 4px 5px 31px 11px #999;
                -moz-transform: rotate(4deg);
                -webkit-transform: rotate(4deg);
                -ms-transform: rotate(4deg);
                -o-transform: rotate(4deg);
                transform: rotate(4deg);
            }
    </style>
</head>
<body>
    <div class="error-page">
        <div class="error-page-container">
            <div class="error-page-main">
                <h3>
                    <strong><a href=http://www.youxiz.net/>游戏猪</a></strong>WwW.YouXiZ.net
                </h3>
                <div class="error-page-actions">
                    <div>
                        <h4>原因：</h4>
                        <ul>
                            <li>不提供对外访问</li>
                            <li>不提供对外访问</li>
                            <li>不提供对外访问</li>
                        </ul>
                    </div>
                    <div>
                        <h4>游戏猪</h4>
                        <ul>
                            <li><a href=#>点我</a></li>
                            <li><a href=#>点我</a></li>
                            <li><a href=#>点我</a></li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>

`
const manage_tpl = `<DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <title>tik后台管理系统</title>
    <link href="https://cdn.staticfile.org/twitter-bootstrap/3.0.1/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.js"></script>
    <style type="text/css">
        body, div, h3, h4, li, ol {
            margin: 0;
            padding: 0;
        }
 
        body {
            font: 14px/1.5 'Microsoft YaHei','微软雅黑',Helvetica,Sans-serif;
            min-width: 1200px;
            background: #f0f1f3;
        }
 
        :focus {
            outline: 0;
        }
 
        h3, h4, strong {
            font-weight: 700;
        }
 
        a {
            color: #428bca;
            text-decoration: none;
        }
            a:hover {
                text-decoration: underline;
            }
 
        .error-page {
            background: #f0f1f3;
            padding: 80px 0 180px;
        }
 
        .error-page-container {
            position: relative;
            z-index: 1;
        }
 
        .error-page-main {
            position: relative;
            background: #f9f9f9;
            margin: 0 auto;
            width: 617px;
            -ms-box-sizing: border-box;
            -webkit-box-sizing: border-box;
            -moz-box-sizing: border-box;
            box-sizing: border-box;
            padding: 50px 50px 70px;
        }
 
            .error-page-main:before {
                content: '';
                display: block;
                background: url(img/errorPageBorder.png?1427783409637);
                height: 7px;
                position: absolute;
                top: -7px;
                width: 100%;
                left: 0;
            }
 
            .error-page-main h3 {
                font-size: 24px;
                font-weight: 400;
                border-bottom: 1px solid #d0d0d0;
            }
 
                .error-page-main h3 strong {
                    font-size: 54px;
                    font-weight: 400;
                    margin-right: 20px;
                }
 
            .error-page-main h4 {
                font-size: 20px;
                font-weight: 400;
                color: #333;
            }
 
        .error-page-actions {
            font-size: 0;
            z-index: 100;
        }
 
            .error-page-actions div {
                font-size: 14px;
                display: inline-block;
                padding: 30px 0 0 10px;
                width: 50%;
                -ms-box-sizing: border-box;
                -webkit-box-sizing: border-box;
                -moz-box-sizing: border-box;
                box-sizing: border-box;
                color: #838383;
            }
 
            .error-page-actions ol {
                list-style: decimal;
                padding-left: 20px;
            }
 
            .error-page-actions li {
                line-height: 2.5em;
            }
 
            .error-page-actions:before {
                content: '';
                display: block;
                position: absolute;
                z-index: -1;
                bottom: 17px;
                left: 50px;
                width: 200px;
                height: 10px;
                -moz-box-shadow: 4px 5px 31px 11px #999;
                -webkit-box-shadow: 4px 5px 31px 11px #999;
                box-shadow: 4px 5px 31px 11px #999;
                -moz-transform: rotate(-4deg);
                -webkit-transform: rotate(-4deg);
                -ms-transform: rotate(-4deg);
                -o-transform: rotate(-4deg);
                transform: rotate(-4deg);
            }
 
            .error-page-actions:after {
                content: '';
                display: block;
                position: absolute;
                z-index: -1;
                bottom: 17px;
                right: 50px;
                width: 200px;
                height: 10px;
                -moz-box-shadow: 4px 5px 31px 11px #999;
                -webkit-box-shadow: 4px 5px 31px 11px #999;
                box-shadow: 4px 5px 31px 11px #999;
                -moz-transform: rotate(4deg);
                -webkit-transform: rotate(4deg);
                -ms-transform: rotate(4deg);
                -o-transform: rotate(4deg);
                transform: rotate(4deg);
            }
    </style>
</head>
<body>
    <div class="error-page">
        <div class="error-page-container">
            <div class="error-page-main">
               

               <table class="table">
                <thead>
                    <tr>
                        <th>
                            编号
                        </th>
                        <th>
                            IMEI
                        </th>
                        <th>
                            到期时间
                        </th>
                        <th>
                            状态
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <!--tr class="success">
                        <td>1</td>
                        <td><a href=# onClick="ViewData('23423863496729386723896');">23423863496729386723896</a></td>
                        <td>01/04/2012</td>
                        <td>Default</td>
                    </tr --!>
                    {{}}
                </tbody>
            </table>

            <div class="col-md-12 column">
                <form class="form-horizontal" role="form">
                <div class="form-group">
                         <label for="imei" class="col-sm-2 control-label">I M E I</label>
                        <div class="col-sm-10">
                            <input type="email" class="form-control" id="imei" />
                            <input type="hidden" class="form-control" id="id" />
                        </div>
                    </div>
                    <div class="form-group">
                         <label for="endtime" class="col-sm-2 control-label">EndTime</label>
                        <div class="col-sm-10">
                            <input type="email" class="form-control" id="endtime"  value='2021-03-19 14:00:00' />
                        </div>
                    </div>
                    <div class="form-group">
                         <label for="pass" class="col-sm-2 control-label">Password</label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" id="pass"/>
                        </div>
                    </div>
                
                    <div class="form-group">
                        <div class="col-sm-offset-2 col-sm-10">
                             <button type="button" class="btn btn-default" onClick="rachange();">充值</button>
                        </div>
                    </div>
                </form>
            </div>

            </div>
        </div>
    </div>
</body>
<script type="text/javascript">
   function ViewData(id,imei){
      $("#imei").val(imei);
      $("#id").val(id);
    }

    function rachange(){
        var postdata = {};
        postdata["id"]      = $("#id").val();
        postdata["imei"]    = $("#imei").val();
        postdata["endtime"] = $("#endtime").val();
        postdata["pass"]    = $("#pass").val();
        if( postdata["pass"].length<5){
            alert("密码错误");
            return;
        } 
        if( postdata["endtime"].length<5){
            alert("请设置充值时间");
            return;
        } 
         if( postdata["imei"].length<5){
            alert("imei错误");
            return;
        } 

        $.post("/ppadmin", postdata, function(obj){
              if(obj.state=="1"){
                alert(obj.msg);
            }else{
                alert(obj.msg);
                window.location.reload();
            }

              
        },"json");
    }

</script>
</html>



`

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

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

/**************************************************************************/

//根据登陆FLAG获取用户用户QQ
func Caiji_get(item string) (info map[string]string) {
	u := make(map[string]string)
	u["state"] = "1"
	//查询数据库
	tm := hPub.Gettime()
	sqlcmd := "select * from caiji where bsend=0 and item='" + item + "' and (unix_timestamp('" + tm + "') - unix_timestamp(sendtime) )>0  order by sendtime limit 1;"
	if item == "all" {
		sqlcmd = "select * from caiji  where bsend=0 and  (unix_timestamp('" + tm + "') - unix_timestamp(sendtime) )>0  group by item order by rand() limit 1;"
	}
	beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		u["msg"] = ("查询登陆标识失败，请联系管理员:" + item)
		return u
	}
	if len(rs) == 0 {
		u["msg"] = ("查无记录")
		return u
	}

	//获取数据
	for _, row := range rs {
		u["content"] = fmt.Sprintf("%s", row["content"])
		if len(u["content"]) > 0 {
			u["state"] = "0"
			u["id"] = fmt.Sprintf("%s", row["Id"])
			u["title"] = fmt.Sprintf("%s", row["title"])
			u["ftitle"] = fmt.Sprintf("%s", row["ftitle"])
			u["titlepic"] = fmt.Sprintf("%s", row["titlepic"])
			u["newstime"] = fmt.Sprintf("%s", row["newstime"])
			u["newsurl"] = fmt.Sprintf("%s", row["newsurl"])
			u["item"] = fmt.Sprintf("%s", row["item"])
			u["sendtime"] = fmt.Sprintf("%s", row["sendtime"])
			u["pitem"] = fmt.Sprintf("%s", row["pitem"])
			u["msg"] = "succ:" + u["id"]

			//执行一次更新数据库操作
			sqlcmd3 := "update caiji set bsend='1' where id='" + u["id"] + "';"
			Dosql(sqlcmd3)
			break
		}

	}
	beego.Debug(u)
	return u

}

func Caiji_set(newurl string, sqlcmdx string) int {
	sqlcmd := "select id from caiji where newsurl='" + newurl + "'"
	beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		beego.Debug(err)
		return 1
	}

	beego.Debug("Line Num:", len(rs), rs)
	if len(rs) != 0 {
		return 2
	}

	//quntime := hPub.Gettime()
	//sqlcmd2 := "insert into user_qun(qunid,userqq,qunname, qunurl,qunlink,qunstate,quntime) values('" + QunId + "','" + UserQQ + "','" + QunName + "','" + QunUrl + "','" + QunLink + "','" + QunFlag + "','" + quntime + "');"
	beego.Debug(sqlcmdx)
	Dosql(sqlcmdx)
	return 0
}

//添加用户
func Check_User(imei string) int {
	if len(imei) != 15 {
		return 1
	}
	tm := hPub.Gettime()
	//sqlcmd := "select * from phoneuser where imei='" + imei + "' and (unix_timestamp('" + tm + "') - unix_timestamp(endtime) )<0"
	sqlcmd := "select * from phoneuser where imei='" + imei + "'"
	//beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		beego.Debug(err)
		return 1
	}

	beego.Debug("Line Num:", len(rs))
	if len(rs) == 0 {
		sqlcmdx := "insert into phoneuser(imei) values('" + imei + "');"
		beego.Debug(sqlcmdx)
		Dosql(sqlcmdx)
		return 2
	}

	etime := tm
	for _, row := range rs {
		if row["endtime"] == nil {
			continue
		}
		etime = fmt.Sprintf("%s", row["endtime"])
		Dosql(etime)
		break
	}

	beego.Debug(imei, "st", tm, "et", etime)
	if etime == tm {
		beego.Debug("未设置充值时间:", imei, tm, etime)
		return 9
	}

	st := hPub.StrToInt(hPub.Gettime_t2c(tm))
	et := hPub.StrToInt(hPub.Gettime_t2c(etime))
	if st-et > 0 {
		beego.Debug("时间已经过期:", imei, tm, etime)
		return 8
	}
	return 0

}

func Add_time(imei string, add_day string) int {
	if len(imei) != 15 {
		return 1
	}
	sqlcmd := `UPDATE phoneuser c set c.endtime = DATE_ADD(c.endtime, INTERVAL ` + add_day + ` DAY) where imei='` + imei + `';`
	beego.Debug(sqlcmd)
	Dosql(sqlcmd)
	return 0
}

func Racharge(imei string, endtime string) {
	sqlcmd := `UPDATE phoneuser set endtime = '` + endtime + `' where imei='` + imei + `';`
	beego.Debug(sqlcmd)
	Dosql(sqlcmd)
}

/**************************************************************************/

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	fmt.Println(gettime(), "method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
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
		ex := strings.ToLower(ext[1])
		if ex == "php" {
			ext[1] = "txt"
		}
		filepath := gettimecuo() + "." + ext[1]
		tab := `<div align="center">
            <table width="70%" border="1" cellspacing="0" cellpadding="0">
              <tr>
                <td>文件名</td>
                <td>` + "<a href=./files/" + filepath + ">" + filepath + "</a>" + ` </td>
              </tr>
              <tr>
                <td>路径</td>
                <td><input type=text value=` + "./files/" + filepath + ` style="width:100%;"> </td>
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

func upfile(w http.ResponseWriter, req *http.Request) {

	surl := req.FormValue("url")
	fmt.Println(gettime(), req.Method, surl)

	if len(surl) == 0 {
		DocumentWrite(w, "获取的URL:"+surl, TYPE_text)
	} else {
		i, r := Down_URL(surl, "youxiz.net") //down_resource
		if i == 0 {
			DocumentWrite(w, `{"state":"0", "msg":"本地化URL成功", "url":"`+r+`"}`, TYPE_json)
		} else {
			DocumentWrite(w, `{"state":"1", "msg":"本地化URL失败"}`, TYPE_json)
		}
	}

}

func Phone_Manage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		a := manage_tpl
		l := QueryPhoneImei()
		a = strings.Replace(a, "{{}}", l, -1)
		DocumentWrite(w, a, TYPE_text)
		return
	} else {

		id := r.FormValue("id")
		imei := r.FormValue("imei")
		pass := r.FormValue("pass")
		endtime := r.FormValue("endtime")
		beego.Debug(id, imei, pass, endtime)
		if pass != "Abcf8765D4" {
			DocumentWrite(w, `{"state":"1", "msg":"密码错误"}`, TYPE_json)
			return
		}
		if len(imei) != 15 {
			DocumentWrite(w, `{"state":"1", "msg":"imei错误"}`, TYPE_json)
			return
		}
		if len(endtime) < 11 {
			DocumentWrite(w, `{"state":"1", "msg":"时间错误"}`, TYPE_json)
			return
		}

		Racharge(imei, endtime)

		DocumentWrite(w, `{"state":"0", "msg":"充值成功，感谢使用"}`, TYPE_json)
	}
}

func phoneupfile(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("method:", r.Method) //获取请求的方法
	fmt.Println(gettime(), "method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		DocumentWrite(w, tpl, TYPE_text)
		return
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("err:", err)
			DocumentWrite(w, "err", TYPE_text)
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		//filepath := "./upload/" + handler.Filename
		fmt.Println("filename", handler.Filename)
		ext := strings.Split(handler.Filename, ".")
		ex := strings.ToLower(ext[1])
		filepath := ""
		if ex == "php" {
			ext[1] = "txt"
			filepath = hPub.Getday() + "." + ext[1]
		} else {
			filepath = handler.Filename //这样修改有强行替换的作用
		}

		beego.Debug("Filename", ex, ext)
		//如果文件长度不对， 或者文件没有扩展名
		if len(ext) < 2 || len(ext[0]) < 2 {
			beego.Debug("err1", len(ext), len(ext[0]))
			DocumentWrite(w, "err", TYPE_text)
			return
		}

		//这里直接查询该IMEI是否授权
		if len(string(ext[0])) != 15 {
			beego.Debug("err2", len(ext[0]), ext[0])
			DocumentWrite(w, "err", TYPE_text)
			return
		}
		ir := Check_User(ext[0])
		if ir == 9 {
			beego.Debug("未充值", ir)
			DocumentWrite(w, "err9", TYPE_text)
			return
		}
		if ir == 8 {
			beego.Debug("已过期", ir)
			DocumentWrite(w, "err8", TYPE_text)
			return
		}

		if ir != 0 {
			beego.Debug("err9", ir)
			DocumentWrite(w, "err", TYPE_text)
			return
		}

		/*tab := `<div align="center">
		  <table width="70%" border="1" cellspacing="0" cellpadding="0">
		    <tr>
		      <td>文件名</td>
		      <td>` + "<a href=./files/" + filepath + ">" + filepath + "</a>" + ` </td>
		    </tr>
		    <tr>
		      <td>路径</td>
		      <td><input type=text value=` + "./files/" + filepath + ` style="width:100%;"> </td>
		    </tr>
		  </table>
		  </div>`
		  DocumentWrite(w, tab, TYPE_text)
		*/

		DocumentWrite(w, filepath, TYPE_text)
		beego.Debug("succ", filepath)
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

//图片评分系统 --start

type TOKENDATA struct {
	Refresh_token  string `json:"refresh_token"`
	Expires_in     int64  `json:"expires_in"`
	Session_key    string `json:"session_key"`
	Access_token   string `json:"access_token"`
	Scope          string `json:"scope"`
	Session_secret string `json:"session_secret"`
}
type PicResultData struct {
	State       string `json:"state"`
	Msg         string `json:"msg"`
	Face_type   string `json:"face_type"`
	Gender_type string `json:"gender_type"`
	Age         string `json:"age"`
	Beauty      string `json:"beauty"`
}

var baidu_token string = "" //传说这个有效期是30天
func Get_Client_ID() string {
	if len(baidu_token) == 0 {
		url := "https://aip.baidubce.com/oauth/2.0/token"
		s := make(map[string]string)
		s["grant_type"] = "client_credentials"
		s["client_id"] = "WtFI2wN6GNTn6DLqN7G8KfY5"             //应用的API Key
		s["client_secret"] = "q8LaFuGiVNWrazxXw4IccYwe68CUGxLM" //，应用的Secret Key；
		r := hNet.Httppostz(url, s)                             //b, err := json.Marshal(d)
		//fmt.Println(r)
		var t TOKENDATA
		json.Unmarshal([]byte(r), &t)
		/*fmt.Println(t.Refresh_token)
		  fmt.Println(t.Expires_in)
		  fmt.Println(t.Session_key)
		  fmt.Println(t.Access_token)
		  fmt.Println(t.Scope)
		  fmt.Println(t.Session_secret)
		  fmt.Println("err:", err)*/
		baidu_token = t.Access_token
		return t.Access_token
	} else {
		return baidu_token
	}
}

/*
{
    "error_code":0,
    "error_msg":"SUCCESS",
    "log_id":2018955201250,
    "timestamp":1613115638,
    "cached":0,
    "result":{
        "face_num":1,
        "face_list":[
            {
                "face_token":"f3d191eb530a2303283d78bb9390a643",
                "location":{
                    "left":420.15,
                    "top":362,
                    "width":82,
                    "height":75,
                    "rotation":37
                },
                "face_probability":1,
                "angle":{
                    "yaw":10.12,
                    "pitch":19.52,
                    "roll":31.97
                },
                "face_type":{
                    "type":"human",
                    "probability":0.69
                },
                "gender":{
                    "type":"female",
                    "probability":1
                },
                "age":22,
                "beauty":85.87
            }
        ]
    }
}*/

func Check_Image(access_token string, picurl string) string {
	url := "https://aip.baidubce.com/rest/2.0/face/v3/detect?access_token=" + access_token
	s := make(map[string]string)
	s["image"] = picurl
	s["image_type"] = "URL"
	s["face_field"] = "facetype,gender,age,beauty"
	s["max_face_num"] = "2"
	r := hNet.HttppostJson(url, s)
	fmt.Println("r:", r)

	var rt PicResultData

	find := strings.Index(r, "\"error_code\":")
	find2 := strings.Index(r[find:], ",")
	fmt.Println("error_code:", find, find2)
	if find > 0 && find2 > 0 {
		age := r[find+13 : find+find2]
		fmt.Println(age)
		rt.State = age
		if age != "0" {
			rt.State = "1"
			rt.Msg = "error"
			//b, _ := json.Marshal(rt)
			return ",,,," //string(b)
		}
	}

	//判断是否是人类
	if find := strings.Contains(r, "human"); find {
		rt.Face_type = "human"
	} else {
		rt.Face_type = "other"
	}

	//判断性别
	if find := strings.Contains(r, "female"); find {
		rt.Gender_type = "female" //女
	}
	if find := strings.Contains(r, "male"); find {
		rt.Gender_type = "male" //男
	}

	//获取年龄
	find = strings.Index(r, "\"age\":")
	find2 = strings.Index(r[find:], ",")
	fmt.Println("age:", find, find2)
	if find > 0 && find2 > 0 {
		age := r[find+6 : find+find2]
		fmt.Println(age)
		rt.Age = age
	}

	//获取颜值分
	find = strings.Index(r, "\"beauty\":")
	find2 = strings.Index(r[find:], "}")
	if find2 == -1 {
		find2 = strings.Index(r[find:], " ")
	}
	fmt.Println("Beauty:", find, find2)
	if find > 0 && find2 > 0 {
		Beauty := r[find+9 : find+find2]
		fmt.Println(Beauty)
		rt.Beauty = Beauty
	}

	find = strings.Index(r, "\"error_msg\":")
	find2 = strings.Index(r[find:], "\"")
	fmt.Println("error_msg:", find, find2)
	if find > 0 && find2 > 0 {
		age := r[find+13 : find+find2]
		fmt.Println(age)
		rt.Msg = age
	}
	//b, _ := json.Marshal(rt)
	a := rt.Beauty + "," + rt.Face_type + "," + rt.Gender_type + "," + rt.Age + ","
	fmt.Println(a)
	return a //string(b)

}
func picbeautiful(w http.ResponseWriter, req *http.Request) {

	surl := req.FormValue("url")
	fmt.Println(gettime(), req.Method, surl)

	if len(surl) == 0 {
		DocumentWrite(w, "获取的URL:"+surl, TYPE_text)
	} else {
		r := Check_Image(Get_Client_ID(), surl)
		DocumentWrite(w, r, TYPE_text)
	}

}

//图片评分系统 -- end

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
	DocumentWrite(res, home_tpl, TYPE_text)
}

func DownFile(imgPath string, localPath string) (ret int) {
	fileName := path.Base(localPath)
	res, err := http.Get(imgPath)
	if nil != err {
		fmt.Println("A error occurred!")
		return 1
	}

	defer res.Body.Close()

	//获得get请求响应的reader对象
	//reader := bufio.NewReaderSize(res.Body, 32 * 1024)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
		return 1
	}

	writer := bufio.NewWriter(file)
	written, copy_err := io.Copy(writer, res.Body)
	if copy_err != nil {
		fmt.Println("????????????????", copy_err)
		return 1
	}
	file.Close()
	fmt.Printf("Total length: %d", written)
	return 0
}

func bytesToSize(length int) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}
func Down_URL(durl string, _infilename string) (ret int, outPath string) {
	uri, err := url.ParseRequestURI(durl)
	if err != nil {
		fmt.Println("网址错误", durl)
		return 1, ""
	}

	filename := gettimecuo() + "-" + path.Base(uri.Path)

	fmt.Println(filename)

	sysType := runtime.GOOS
	if sysType == "linux" {
		// LINUX系统
		beego.Debug("LINUX系统")
		cmd := exec.Command("/bin/bash", "-c", "curl -k "+durl+" -o upload/"+filename)
		buf, err := cmd.Output()
		if err != nil {
			beego.Debug(err.Error())
		}
		beego.Debug(string(buf))
	}
	if sysType == "windows" {
		// windows系统
		beego.Debug(" windows系统")
		cmd := exec.Command("cmd.exe", "/c", "curl -k "+durl+" -o upload/"+filename)
		buf, err := cmd.Output()
		if err != nil {
			beego.Debug(err.Error())
		}
		beego.Debug(string(buf))
	}

	return 0, filename
}

//下载资源
func down_resource(url string, destDir string) (ret int, outPath string) {

	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		//下载图片资源
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return 1, ""
		}

		//判断文件是否存在
		ext := url[len(url)-4:]
		outPath := "./upload/" + destDir + ext
		outPath2 := destDir + ext
		if File_Exists(outPath) {
			c := gettimecuo()
			outPath = "./upload/" + destDir + "-" + c + ext
			outPath2 = destDir + "-" + c + ext
		}
		fmt.Println(outPath2)

		//只处理jpg图片
		out, create_err := os.Create(outPath)
		if create_err != nil {
			fmt.Println(create_err)
			return 1, ""
		}
		defer out.Close()
		defer resp.Body.Close()
		_, copy_err := io.Copy(out, resp.Body)
		if copy_err != nil {
			fmt.Println("????????????????", copy_err)
			return 1, ""
		}
		out.Close()
		return 0, outPath2
	} else {
		return 1, ""
	}

}

//初始化连接
/*func M_init() {
    //beego.Debug("m_init")
    baseaddress := "127.0.0.1" // Get_config("baseaddress")
    baseport := "3306"         //Get_config("baseport")
    basename := "quntool"      //Get_config("sqlhost")
    sqluser := "quntool"       //Get_config("sqluser") //用户名
    sqlpass := "quntool!@#"    //Get_config("sqlpass") //密码
    hPub.Gettime()
    orm.RegisterDataBase("default", "mysql", sqluser+":"+sqlpass+"@tcp("+baseaddress+":"+baseport+")/"+basename+"?charset=utf8", 200, 200)
    //orm.RegisterModel(new(User))
}*/

//初始化连接
func M_init() {
	//beego.Debug("m_init")
	baseaddress := "127.0.0.1"    // Get_config("baseaddress")
	baseport := "3306"            //Get_config("baseport")
	basename := "caiji"           //Get_config("sqlhost")
	sqluser := "av"               //Get_config("sqluser") //用户名
	sqlpass := "mpnfpc3FcZycx5nt" //Get_config("sqlpass") //密码
	hPub.Gettime()
	orm.RegisterDataBase("default", "mysql", sqluser+":"+sqlpass+"@tcp("+baseaddress+":"+baseport+")/"+basename+"?charset=utf8", 200, 200)
	//orm.RegisterModel(new(User))
}

//切换数据库
func M_Using(dbname string) {
	o := orm.NewOrm()
	o.Using(dbname)
}

//获取有返回值的sql语句，比如select, shwo database
func Getrs(sqlcmd string) (err error, rs []orm.Params) {
	o := orm.NewOrm()
	_, err = o.Raw(sqlcmd).Values(&rs)
	return err, rs
}

//执行sql语句，比如： insert into， update ，deltete， create 等
func Dosql(sqlcmd string) {
	o := orm.NewOrm()
	o.Raw(sqlcmd).Exec()
}

//检查数据库中的表是否存在
//检查数据库中的表是否存在
func Check_Tab_Exists(tabname string) bool {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		err := recover()
		if err != nil {
			beego.Debug(tabname, "查询异常，表", tabname, "不存在")
		}
	}()

	sqlcmd := "select * from `" + tabname + "` limit 0, 1"
	err1, rs := Getrs(sqlcmd)
	if err1 != nil {
		beego.Debug("查询失败 视为表", tabname, "不存在")
		return false //查询失败 视为表不存在
	}
	if len(rs) > 0 {
		beego.Debug("表", tabname, "存在")
		return true
	}
	beego.Debug(tabname, "表补存在")
	return false
}

//显示当前连接中所有的库取值字段：Database
func ShowDatabases() (err error, rs []orm.Params) {
	//var rs []orm.Params
	err, rs = Getrs("show databases")
	beego.Debug("ShowDatabases 取值字段：Database", rs)
	/*if err != nil {
	      beego.Debug(err)
	  } else {
	      for num, row := range rs {
	          beego.Debug(num, row["Database"])
	      }
	  }*/
	return err, rs
}

//显示当前库下所有的数据表 取值字段：Tables_in_square
func ShowTables(dbname string) (err error, rs []orm.Params) {
	M_Using(dbname)
	sqlcmd := "show tables;"
	err, rs = Getrs(sqlcmd)
	beego.Debug("ShowTables 取值字段：Tables_in_square", rs)
	return err, rs
}

//显示当前表下所有的字段  取值字段：Field
//调用这个函数之前必须调用showtables防止未切换数据库
func ShowField(tabname string) (err error, rs []orm.Params) {
	sqlcmd := "desc " + tabname + ";"
	err, rs = Getrs(sqlcmd)
	beego.Debug("ShowField 取值字段：Field", rs)
	return err, rs
}

func QueryURL(newurl string) int {
	sqlcmd := "select id from caiji where newsurl='" + newurl + "'"
	beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		beego.Debug(err)
		return 1
	}

	beego.Debug("Line Num:", len(rs), rs)
	if len(rs) == 0 {
		return 0
	} else {
		return 2
	}
}

func QueryPhoneImei() string {
	sqlcmd := "select id,imei,endtime from phoneuser"
	//beego.Debug(sqlcmd)
	err, rs := Getrs(sqlcmd)
	if err != nil {
		beego.Debug(err)
		return ""
	}

	//beego.Debug("List Num:", len(rs))
	if len(rs) == 0 {
		return ""
	} else {
		ret := ""
		iCount := 0
		for _, row := range rs {
			iCount = iCount + 1
			sCount := fmt.Sprintf("%d", iCount)
			id := fmt.Sprintf("%s", row["id"])
			imei := fmt.Sprintf("%s", row["imei"])
			endtime := fmt.Sprintf("%s", row["endtime"])
			if row["endtime"] == nil {
				endtime = ""
			}
			class_txt := ""
			if iCount%2 == 0 {
				class_txt = ` class="success"`
			}
			ret = ret + `<tr` + class_txt + `>
            <td>` + sCount + `</td>
            <td><a href=# onClick="ViewData('` + id + `','` + imei + `');">` + imei + `</a></td>
            <td>` + endtime + `</td>
            <td>正常</td>
            </tr>`
		}
		return ret
	}
}

func GetData(res http.ResponseWriter, req *http.Request) string {

	var d CAIJIDATA
	d.State = "1"
	d.Msg = "错误"
	item := req.FormValue("item")
	if len(item) > 0 {
		r := Caiji_get(item)
		d.Title = r["title"]
		d.Ftitle = r["ftitle"]
		d.Titlepic = r["titlepic"]
		d.Newstime = r["newstime"]
		d.NewsUrl = r["newsurl"]
		d.Content = r["content"]
		d.Id = r["id"]
		d.Msg = r["msg"]
		d.Item = r["item"]
		d.Pitem = r["pitem"]
		d.Sendtime = r["sendtime"]
		d.State = r["state"]

		b, err := json.Marshal(d)
		if err == nil {
			//c.Ctx.ResponseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
			beego.Debug("获取成功!", string(b))
			//c.Ctx.WriteString(string(b))
			return string(b)
		} else {
			//c.Ctx.ResponseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")
			//c.Ctx.WriteString(r["msg"])
			return r["msg"]
		}
	} else {
		ret := `<html xmlns="http://www.w3.org/1999/xhtml">
                <head>
                <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
                <title>caiji</title>
                </head>

                <body>

                <div align="center" style="width:70%; height:300px;">
                <form id="form1" name="form1" method="post" action="caiji"  >
                  <p>
                  <textarea name="content" cols="" rows="" style="width:70%; height:300px;"></textarea>
                  </p>
                  <p>
                    <input name="subst" type="submit" />
                    </p>
                </form>
                <form id="form1" name="form1" method="post" action="caiji"  >
                  <p>
                  <textarea name="queryurl" cols="" rows="" style="width:70%; height:300px;"></textarea>
                  </p>
                  <p>
                    <input name="subst" type="submit" />
                    </p>
                </form>
                </div>
                <p>&nbsp;</p>
                </body>
                </html>`
		//c.Ctx.WriteString(ret)
		return ret
	}
	return ""
}

func AddData(res http.ResponseWriter, req *http.Request) string {
	//1.设置返回格式
	//c.Ctx.ResponseWriter.Header().Add("Content-Type", "application/json;charset=UTF-8")

	content := req.FormValue("content")
	beego.Debug(content)

	queryurl := req.FormValue("queryurl")
	beego.Debug(queryurl)

	if len(queryurl) > 0 && len(content) == 0 {
		tmpr := QueryURL(queryurl)
		if tmpr == 0 {
			return (`{"state":"0","msg":"可以添加"}`)

		} else {
			return (`{"state":"2","msg":"URL已经采集"}`)
		}
	}

	if len(content) == 0 {
		return (`{"state":"1","msg":"参数错误:无数据"}`)
	}

	var s CAIJIDATA
	err := json.Unmarshal([]byte(content), &s)
	if err != nil {
		return (`{"state":"1","msg":"参数错误:无法解析"}`)
	}

	//2.先进行数据判断
	if len(s.Title) < 5 || len(s.Content) < 6 || len(s.NewsUrl) < 10 {
		return (`{"state":"1","msg":"参数错误:传入的数据标题，内容，和网址不能为空"}`)
	}
	sendtime := s.Sendtime
	if len(sendtime) == 0 {
		sendtime = hPub.Gettime()
	}

	//sqlcmd := "insert into caiji(title,ftitle,titlepic,newstime,content,item,pitem,newsurl,sendtime) values('" + s.Title + "','" + s.Ftitle + "','" + s.Titlepic + "','" + s.Newstime + "','" + s.Content + "','" + s.Item + "','" + s.Pitem + "','" + s.NewsUrl + "','" + sendtime + "');"
	sqlcmd := "insert into caiji(title,ftitle,titlepic,newstime,content,item,pitem,newsurl,sendtime,bsend) values('" + s.Title + "','" + s.Ftitle + "','" + s.Titlepic + "','" + s.Newstime + "','" + s.Content + "','" + s.Item + "','" + s.Pitem + "','" + s.NewsUrl + "','" + sendtime + "','0');"
	beego.Debug(sqlcmd)
	r := Caiji_set(s.NewsUrl, sqlcmd)

	//设置COOKIE
	//c.Ctx.SetCookie("flag", sign, time.Second*60*60*24*7) //一周登陆一次

	//更新跳转
	if r == 0 {
		return (`{"state":"0", "msg":"添加成功"}`)
	} else if r == 2 {
		return (`{"state":"2", "msg":"已经存在"}`)
	} else {
		return (`{"state":"1", "msg":"添加失败"}`)
	}
}

func Caiji(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	if method == "GET" {
		r := GetData(res, req)
		DocumentWrite(res, r, TYPE_text)
		return
	} else {

		qurl := req.FormValue("queryurl")
		if len(qurl) > 0 {
			tmpr := QueryURL(qurl)
			beego.Debug(tmpr)
			if tmpr == 0 {
				DocumentWrite(res, `{"state":"0","msg":"可以添加"}`, TYPE_json)
			} else if tmpr == 2 {
				DocumentWrite(res, `{"state":"2","msg":"URL已经采集"}`, TYPE_json)
			} else {
				DocumentWrite(res, `{"state":"1","msg":"参数错误:无数据"}`, TYPE_json)
			}
			return
		}

		content := req.FormValue("content")
		if len(content) > 0 {
			s := AddData(res, req)
			beego.Debug(s)
			DocumentWrite(res, s, TYPE_json)
			return
		}
	}
	DocumentWrite(res, `{"state":"1", "msg":"NoArgv"}`, TYPE_json)
}

func main() {
	//i := DownFile("https://ss1.bdstatic.com/70cFuXSh_Q1YnxGkpoWK1HF6hhy/it/u=1201719684,3757858959&fm=26&gp=0.jpg", "0.jpg")
	//i := down_resource("https://ss1.bdstatic.com/70cFuXSh_Q1YnxGkpoWK1HF6hhy/it/u=1201719684,3757858959&fm=26&gp=0.jpg", "1.jpg")
	// i,s := down_resource("https://www.80host.com/cloud.html", "1.html")
	//fmt.Println("下载结果:", i, s)
	M_init()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	port := "9566"
	args := os.Args
	if len(args) > 1 {
		port = args[1]
	}

	wd = wd + "/upload"
	os.Mkdir(wd, os.ModePerm)
	fmt.Println(wd)
	fs := http.FileServer(http.Dir(wd))

	mux := http.NewServeMux()
	mux.Handle("/files/", http.StripPrefix("/files", fs))
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/upfile", upfile)
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/caiji", Caiji)

	mux.HandleFunc("/phoneupfile", phoneupfile)   //获取上传地址
	mux.HandleFunc("/picbeautiful", picbeautiful) //计算评分//https://www.youxiz.net/d/file/upload/1613722972012.png
	mux.HandleFunc("/ppadmin", Phone_Manage)

	//设置访问的路由
	fmt.Println(gettime(), "服务器"+port+"开始服务。。。")
	err = http.ListenAndServe(":"+port, mux) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main1() {
	M_init()
	r := Add_time("111111111100000", "1")
	beego.Debug(r)
	r = Check_User("111111111100000")
	beego.Debug(r)
}
