import os
import subprocess
import requests
import datetime
import math

cmd='go run main.go'
traces = ['0','1','2']
options=['2','3']
times =['0.01','0.001','0.0001']
memsizes =['1000','500']
entryargs=['20','10','5']
tracename=['wide','chicago','sinet']


nowTime=datetime.datetime.now().isoformat()

# for i in range(0,1000):
#     if i==0 or (i%10==0):
#         print([str(i)])
#         args.append(option+[str(i)])

l = []
for time in times:
	for trace in traces:
		for option in options:
			for memsize in memsizes:
				for entryarg in entryargs:
					# print(cmd,traces+' '+option+' '+time+' ' +str(int(memsize)/int(entryarg))+' '+entryarg)
					l.append(cmd+' '+trace+' '+option+' '+time+' ' +str(int(memsize)//int(entryarg))+' '+entryarg)

for i in l:
	print(i)
# print(l)


# for i in range(len(l)):
	# path='./test'+str(i+1)+'.dat'
#	path ='./bufChangeideal.dat'
path='test2.txt'
with open(path,mode='a') as f:
	for time in times:
		for trace in traces:
			f.write(tracename[int(trace)]+'\n')
			# for option in options:
			for memsize in memsizes:
				for entryarg in entryargs:
				# print(cmd,traces+' '+option+' '+time+' ' +str(int(memsize)/int(entryarg))+' '+entryarg)
					# l.append(cmd+' '+trace+' '+option+' '+time+' ' +str(int(memsize)//int(entryarg))+' '+entryarg)
					r=''
					for option in options:
						# print(cmd.split()+[trace,option,time,str(int(memsize)//int(entryarg)),entryarg])
						r+=subprocess.check_output('ssh lemon cd go/src/github.com/yamatcha/simulator | '.split()+cmd.split()+[trace,option,time,str(int(memsize)//int(entryarg)),entryarg]).decode('utf-8')+' '
					print(r.replace("\n","")+'\n')
					f.write(str(int(memsize)//int(entryarg))+' '+entryarg+' '+r.replace("\n","")+'\n')
					f.flush()

		# print(cmd+option+args[i])
		# r=subprocess.run(cmd+option+args[i], stdout=f)
		# print(str(args[i][2]))
		# r=subprocess.check_output(cmd+option+args[i]).decode('utf-8')
		# f.write(str(args[i][0])+' '+r)


url = "https://notify-api.line.me/api/notify"
token = "WcudzQXjoEgLad8EA68AkLe98Tl5mxEjbVhgOjdBIZH"
headers = {"Authorization" : "Bearer "+ token}
payload = {"message" :  nowTime}
# files = {"imageFile": open("test0.dat")}

r = requests.post(url ,headers = headers ,params=payload)

