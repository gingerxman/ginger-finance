#!/usr/bin/env python2
#-*- coding: utf-8 -*-

import os,sys,signal,time	   	   	   	   	
from fsevents import Observer
from fsevents import Stream  

SERVER_PID = 0
LAST_REBUILD_TIME = 0

def callback(FileEvent):
	# attributes of FileEvent：mask, cookie and name.
	# mask: 512-delete;256-create;2-changed;...
	is_target_file = False
	if FileEvent.name.endswith(".go") or FileEvent.name.endswith(".json"):
		is_target_file = True

	if not is_target_file:
		return

	is_need_rebuild = False
	file_name = FileEvent.name
	if FileEvent.mask == 256:
		print 'delete %s, rebuild' % file_name
		is_need_rebuild = True
	if FileEvent.mask == 2:
		print 'change %s, rebuild' % file_name
		is_need_rebuild = True

	if is_need_rebuild:
		with open('_pid', 'rb') as f:
			pid = int(f.read().strip())
			print 'kill process ', pid
			try:  
				os.kill(pid, signal.SIGKILL)  
			except OSError:
				pass

def child():
	observer = Observer()
	stream = Stream(callback, ".", file_events=True)
	observer.schedule(stream)
	observer.start()

def app_child():
	os.execvp("./ginger-mall", [""])

class Watcher:  
	""" 
	创建一个做苦工的子进程。然后父进程等待KeyboardInterrupt并杀掉子进程。
	"""  
	def __init__(self):  
		self.child = os.fork()  
		if self.child == 0:  
			child()  
		else:  
			self.watch()
			
			#self.watch()
			#self.app_child = os.fork()
			#if self.app_child == 0:
			#	app_child()
			#else:
				#with open('_pid', 'wb') as f:
				#	f.write(str(self.app_child))
				#self.watch()  
  
	def watch(self):  
		while True:
			import subprocess
			proc = subprocess.Popen("./ginger-finance", stdout=subprocess.PIPE)
			with open('_pid', 'wb') as f:
				f.write(str(proc.pid))

			try:
				while True:
					line = proc.stdout.readline()
					if line != '':
						print line.rstrip()
					else:
						print '\n\n========= Compile EEL ========='
						os.system('go build -o ginger-finance -mod=vendor -v ./main.go')
						break
			except KeyboardInterrupt:
				print 'Ctrl+C received, stop'
				break
		
		self.kill()  
		sys.exit()
		# try:  
		# 	os.wait()  
		# except KeyboardInterrupt:  
		# 	#捕获 Control+C，杀掉子进程 
		# 	print 'KEYBOARDINTERRUPT\n'  
		# 	self.kill()  
		# 	sys.exit()  
		# except:
		# 	print 'haha'
  
	def kill(self):  
		try:  
			os.kill(self.child, signal.SIGKILL)  
		except OSError: pass  

		# try:  
		# 	os.kill(self.app_child, signal.SIGKILL)
		# except OSError: pass

		import time
		time.sleep(1)
  

if __name__ == '__main__':
	exit_status = os.system('go build -o ginger-finance -mod=vendor -v ./main.go')
	if exit_status != 0:
		print '[Error] Compile Fail!!!'
	else:
		Watcher()