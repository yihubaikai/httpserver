package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/goredis"
	"github.com/yihubaikai/gopublic"
	"io"
	"net/http"
	"os"
	//"net/url"
)

//-------------------------------全局变量------------------------------------
var client goredis.Client

const (
	TYPE_json = 1
	TYPE_text = 2
	MaxLength = 1024
)

var (
	MaxWorker = os.Getenv("MAX_WORKERS")
	MaxQueue  = os.Getenv("MAX_QUEUE")
)

type Payload struct {
	Base string `json:"base"`
}
type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// A buffered channel that we can send work requests on.
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

//********************************实现函数************************************

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
	fmt.Println(val)
}

//文档跳转值
func DocumentRedirect(res http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(res, req, url, http.StatusFound)
}

//********************************路由处理************************************
//接受客户端返回的任意数据
func GetNick(res http.ResponseWriter, req *http.Request) {

}

func Home(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.RemoteAddr, req.Host, req.RequestURI)
	req.ParseForm()
	DocumentWrite(res, "home", TYPE_text)
}

//路由
func Routes() {
	http.HandleFunc("/getnick", payloadHandler) //获取二维码
	//http.HandleFunc("/setnick", SetNick) //获取二维码 分别对应大小写
	http.HandleFunc("/", Home) //支付核销函数
}

//--------------------------------------------------------------------------
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				/*if err := job.Payload.UploadToS3(); err != nil {
					fmt.Errorf("Error uploading to S3: %s", err.Error())
				}*/
				fmt.Println("上传数据中...", job)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the body into a string for json decoding
	var content = &PayloadCollection{}
	err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Go through each payload and queue items individually to be posted to S3
	for _, payload := range content.Payloads {

		// let's create a job with the payload
		work := Job{Payload: payload}

		// Push the work onto the queue.
		JobQueue <- work
	}

	w.WriteHeader(http.StatusOK)
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < 10; i++ {
		//worker := NewWorker(d.pool)
		//worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

//初始化   主函数
func main() {

	fmt.Println("初始化Redis数据库...")
	client.Addr = "127.0.0.1:6379"
	client.Db = 0

	fmt.Println("初始化路由...")
	Routes()

	hPub.CreateDateDir("data")

	dispatcher := NewDispatcher(10)
	dispatcher.Run()

	s := hPub.Gettime()
	port := "80"
	fmt.Println(s, "开始监听:", port)
	http.ListenAndServe(":"+port, nil)
	fmt.Println(hPub.Gettime(), "监听失败程序退出")

}
