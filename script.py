import os
import subprocess
import requests
import datetime

cmd=['go', 'run', 'main.go']
option=[ '2','0.01']
args=[option+['10']]

nowTime=datetime.datetime.now().isoformat()

# for i in range(0,500):
#     if i==0 or (i%10==0):
#         print([str(i)])
#         args.append(option+[str(i)])

print(args)


for i in range(len(args)):
    path='./test'+str(i+2)+'.dat'
    with open(path,mode='w') as f:
        print(cmd+args[i])
        r=subprocess.run(cmd+args[i],stdout=f)
        # print(str(args[i][2]))
        # r=subprocess.check_output(cmd+args[i]).decode('utf-8')
        # f.write(str(args[i][2])+' '+r)


url = "https://notify-api.line.me/api/notify"
token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
headers = {"Authorization" : "Bearer "+ token}
payload = {"message" :  nowTime}
# files = {"imageFile": open("test0.dat")}

r = requests.post(url ,headers = headers ,params=payload)

