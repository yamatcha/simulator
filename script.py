import os
import subprocess
import requests
import datetime
import math

cmd=['go', 'run', 'main.go']
option=['1']
args=[['1'],['0.1'],['0.01'],['0.001'],['0.0001'],['0.00001'],['0.000001'],['0']]

nowTime=datetime.datetime.now().isoformat()

# for i in range(0,1000):
#     if i==0 or (i%10==0):
#         print([str(i)])
#         args.append(option+[str(i)])


print(args)

ars = option+args

for i in range(len(args)):
	path='./test'+str(i+1)+'.dat'
#	path ='./bufChangeideal.dat'
	with open(path,mode='w') as f:
		print(cmd+option+args[i])
		r=subprocess.run(cmd+option+args[i], stdout=f)
		# print(str(args[i][2]))
		# r=subprocess.check_output(cmd+option+args[i]).decode('utf-8')
		# f.write(str(args[i][0])+' '+r)


url = "https://notify-api.line.me/api/notify"
token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
headers = {"Authorization" : "Bearer "+ token}
payload = {"message" :  nowTime}
# files = {"imageFile": open("test0.dat")}

r = requests.post(url ,headers = headers ,params=payload)

