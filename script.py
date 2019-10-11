import os
import subprocess
import requests
import datetime

cmd=['go', 'run', 'main.go']
args=[[ '1','0.01','1000'],[ '2','0.01','1000']]

nowTime=datetime.datetime.now().isoformat()

for i in range(len(args)):
    path='./test'+str(i)+'.dat'
    with open(path,mode='w') as f:
        subprocess.run(cmd+args[i],stdout=f,shell=True)

url = "https://notify-api.line.me/api/notify"
token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
headers = {"Authorization" : "Bearer "+ token}
payload = {"message" :  nowTime+" simuilation finish"}
# files = {"imageFile": open("test0.dat")}

r = requests.post(url ,headers = headers ,params=payload)

