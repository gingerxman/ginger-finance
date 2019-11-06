# -*- coding: utf-8 -*-

import sys, os, time
import subprocess
import urllib

from django.conf import settings
import requests

from terminal import Terminal

MANUAL_START_SERVER = False

def get_binary_path():
	if 'win' in sys.platform:
		if sys.platform == 'darwin':
			return './%s' % settings.SERVICE_NAME
		else:
			return '%s.exe' % settings.SERVICE_NAME
	else:
		return './%s' % settings.SERVICE_NAME

def is_test_target_running():
	"""
	判断server是否处于运行状态
	:return:
	"""
	"""
	:return: 
	"""
	url = 'http://127.0.0.1:%s/op/health/' % settings.SERVICE_PORT
	r = requests.get(url, timeout=1)
	if r.status_code == 200:
		return True
	else:
		return False

def start_test_target():
	"""
	启动待测试的程序
	"""
	if is_test_target_running():
		return False

	print '[bdd] >>>>>>>> build program <<<<<<<<'
	cmd = "go build %s" % settings.SERVICE_NAME
	print '%s ...' % cmd
	os.system(cmd)
	print '[bdd] >>>>>>>> run program %s <<<<<<<<' % get_binary_path()
	script = [get_binary_path(),]
	terminal = Terminal()
	terminal.execute(None, "robert", script, os.getcwd(), True, None)

	while True:
		print '[bdd] waiting test target...'
		time.sleep(1)
		try:
			url = 'http://127.0.0.1:%s/op/health/' % settings.SERVICE_PORT
			print url
			r = requests.get(url, timeout=1)
			if r.status_code == 200 and 'running' in r.text:
				print '[bdd] test target is ready'
				break
		except requests.exceptions.ConnectTimeout:
			print '[bdd] wait timeout, retry'
		except requests.exceptions.ConnectionError:
			print '[bdd] connection refused, retry'

	global MANUAL_START_SERVER
	MANUAL_START_SERVER = True
	return True


def stop_test_target():
	"""
	启动待测试的程序
	"""
	if MANUAL_START_SERVER:
		url = 'http://127.0.0.1:%s/_server/api/shutdown/' % settings.SERVICE_PORT
		r = requests.get(url, timeout=1)
		print r.text
	else:
		print 'skip stop test target'

