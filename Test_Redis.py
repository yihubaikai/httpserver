import requests,json
import time
import sys,os,io
import codecs
from datetime import datetime
import redis
port = "80"#sys.argv[1]

#添加昵称
def Add_huoshan_Nick(userid, text,_time,flag="0",base="huoshan"):
	data = {
		"userid"	:userid,
		#"usernick"  :usernick,
		"flag"	  :flag,
		"starttime" :_time,	#datetime.now().strftime( "%Y-%m-%d %H:%M:%S" ),
		"base"	  :base
	}
	arr = text.split("|")
	#print(arr)
	if len(arr)>=7:
		data["usernick"] = arr[1]
		data["replay"] = arr[2]
		data["city"] = arr[3]
		data["mynick"] = arr[4]
		data["fromid"] = arr[5]

	r = {}
	try:
		#print(data)
		url = "http://pay.ggpaygg.com:"+ port + "/"
		r = requests.post(url=url,data = data )
		if(r.status_code != 200):
			print("error:" + str(r.status_code))
			return

		#print(r.text)
	except Exception as e:
		return
	
	#print(r)
	return


def printx(text):
	print("*" * 20)
	print(text)
	print("*" * 20)


def SaveLog(filepath, text, isNewfile=False):
	if(isNewfile==True and os.path.isfile(filepath) and os.path.exists(filepath)):
		os.remove(filepath)
	with open(filepath,"a",encoding='utf-8') as f:
	#with open(filepath,"a") as f:
		f.write(text)
		#f.flush()
def Main_HuoShan_Redis():
	red = redis.Redis(host='node1.m8p.net',port=6379,db=1,password='p@ssw0rd')
	#file = open(sFile)
	print("*" * 66)
	printx("开始 ，火山， 连接数据库...")
	
	#path = os.getcwd() + "\\data\\huoshan1.json"
	path = "D:\\code\\4.es_code\\transfer_data\\data\\huoshan1.json"
	print(path)
	file = codecs.open(path,'r','utf-8')
	icount = 0
	a=datetime.now() 
	s = datetime.now().strftime( "%Y-%m-%d %H:%M:%S" )
	print(s,a)
	while 1:
			line = file.readline()
			if not line:
				break
			if(len(line)<2):
				continue
			line = line.replace("\n", "")
			line = line.replace("\r", "")

			#解析json
			obj = json.loads(line)
			#print(obj)
			if("nick" in obj ) and ("class" in obj) and ("flag" in obj) and ("starttime" in obj):
				nick = obj["nick"]
				text = obj["class"]
				flag = obj["flag"]
				time = obj["starttime"]
				red.set(nick, text)
				#Add_huoshan_Nick(nick, text,time,flag,"huoshan")

				print(text)
				icount = icount + 1
				if(icount%1000==0):
					b=datetime.now()
					tstamp = (b-a).seconds
					print(s,datetime.now().strftime( "%Y-%m-%d %H:%M:%S" ),tstamp)
				if(icount>=10000):
					b=datetime.now()
					tstamp = (b-a).seconds
					print(s,datetime.now().strftime( "%Y-%m-%d %H:%M:%S" ),tstamp)
					SaveLog("log.txt", port + "/" +  str(tstamp) + "/"+ str(os.getpid()) + "/"  + s + "/" + datetime.now().strftime( "%Y-%m-%d %H:%M:%S" ) + "\n")
					break
				#print(line)
				#break
	#r = cli.Get_Nick("0","1","huoshan_user","user")
	#print("获取昵称:", r)
	print("完成")

if __name__ == '__main__':
	Main_HuoShan_Redis()




