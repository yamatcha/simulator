import os
import subprocess
import requests
import datetime

cmd=['go', 'run', 'main.go']
option=[ '1','0.01']
args=[option]

nowTime=datetime.datetime.now().isoformat()

# for i in range(1,1000):
#     if (i<=10 and i%5==0) or  (10< i and i<=100 and i%10==0) or (100 <i and i<=1000 and i%100==0):
#         print([str(i)])
#         args.append(option+[str(i)])

print(args)


for i in range(len(args)):
    path='./bufSize.dat'
    with open(path,mode='w') as f:
        print(cmd+args[i])
        r=subprocess.run(cmd+args[i],stdout=f)

url = "https://notify-api.line.me/api/notify"
token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
headers = {"Authorization" : "Bearer "+ token}
payload = {"message" :  nowTime}
# files = {"imageFile": open("test0.dat")}

r = requests.post(url ,headers = headers ,params=payload)

