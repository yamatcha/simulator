import os
import subprocess
import requests
import datetime
import math

cmd='go run main.go'
csvpaths = ["/csv/chicago.csv"]
modes=['2']
timeWidths =['0.01']
memsizes =['1000','500']
entryargs=['20','10','5']
protocols=['-protocol=TCP -selectedPort=443','-protocol=TCP -selectedPort=80', '-protocol=TCP -selectedPort=443,80','']

nowTime=datetime.datetime.now().isoformat()

l = []
for timeWidth in timeWidths:
	for csvpath in csvpaths:
		for mode in modes:
			for memsize in memsizes:
				for entryarg in entryargs:
					bufSize=str(int(memsize)//int(entryarg))
					for protocol in protocols:
						l.append(cmd+' -csvPath=/home/soju'+csvpath+' -mode='+mode+' -timeWidth='+timeWidth+' -bufSize=' +bufSize+' -entrySize='+entryarg+' '+protocol)

for i in l:
	print(i)

path='/home/soju/researchResult/protocolfilter/test1.txt'
with open(path,mode='a') as f:
	r=''
	for i in l:
		r=subprocess.check_output(i.split()).decode('utf-8')+' '
		print(i)
		print(r.replace("\n","")+'\n')
		f.write(i+' '+r.replace("\n","")+'\n')
	f.flush()
