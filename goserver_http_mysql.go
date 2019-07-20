package main
 
import (
	"fmt"
	"html/template"
	"log"
	"net/http"
    "encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"container/list"
 
)

//******全局变量**********************
var db *sql.DB;
var sql_list list.List;



func heart(res http.ResponseWriter, req *http.Request) {
	admin := req.FormValue("admin")
	password := req.FormValue("password")
 
	fmt.Println(admin)
	fmt.Println(password)
 
	if admin != "admin" || password != "admin888" {
		res.Write([]byte("Login Fail,Please Try Again!"))
	} else {
		res.Write([]byte("<a href=/succ>登陆成功</a>"))
	}
}

func pcrefresh(res http.ResponseWriter, req *http.Request) {
	admin := req.FormValue("admin")
	password := req.FormValue("password")
 
	fmt.Println(admin)
	fmt.Println(password)
 
	if admin != "admin" || password != "admin888" {
		res.Write([]byte("login faild"))
	} else {
		res.Write([]byte("login succ"))
	}
}

func getqrcode(res http.ResponseWriter, req *http.Request){
	type aa struct {
        Status string
        Msg string
    }

    
    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(200)
    u:=aa{"0","你好"}
 	if result,err:=json.Marshal(&u);err==nil{
		res.Write([]byte( string(result)))
	}
	fmt.Println("输出json")
	//res.Write(string(arr))
    //json.NewEncoder(res).Encode(&u)	//res.Write([]byte("获取二维码")
}

func Routes(){
	http.HandleFunc("/online", 		heart) 		//登陆
	http.HandleFunc("/pcrefresh",  	pcrefresh)  //支付
	http.HandleFunc("/getqrcode",	getqrcode)  //获取二维码

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {//HOME
		t, err := template.ParseFiles("login.html")
		if err != nil {
			log.Println("err")
		}
		t.Execute(res, nil)
	})
 
}


//初始化 自动系统调用
func init() {
	db, _ = sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/test?charset=utf8&allowOldPasswords=1") 
	db.SetMaxOpenConns(2000) 
	db.SetMaxIdleConns(1000) 
	db.Ping()
	sql_list := list.New()
	for i := 0; i < 10; i++ {
		sql_list.PushBack(i)
	}
	for i := sql_list.Front(); i != nil; i = i.Next() {
    	//fmt.Println(i.Value)
	}
}
 

func testmysql(){
	/*db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/test?charset=utf8&allowOldPasswords=1") 
	if err != nil {
		fmt.Println(err) 
		return 
	} */
	//defer db.Close() 
	
	//插入
	stmt, _ := db.Prepare(`INSERT into tabx (name) values (?)`) 
	res, _ := stmt.Exec("tony")
	id, _ := res.LastInsertId()
	fmt.Println(id)
 
 
 	//查询
	rows,_ := db.Query("select * from tabx")
	for rows.Next(){
	    var id int
	    var name string
	    rows.Columns()
	    if err := rows.Scan(&id,&name); err != nil {
	        log.Fatal(err)
	    }
	    fmt.Printf("id:%d name:%s\n", id, name)
	}
	rows.Close()
}

func main() {
	//testmysql()
	Routes()
 	port := "80"
	fmt.Println("start http server at:", port)
	http.ListenAndServe(":"+port, nil)
	db.Close()
}