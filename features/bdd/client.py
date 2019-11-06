# -*- coding: utf-8 -*-
import json
import requests
import os
import sys

CUR_CONTEXT = None

_SERVICE_INFO = None
def get_service_info():
	global _SERVICE_INFO
	if not _SERVICE_INFO:
		service_name = None
		service_port = None
		with open('./conf/app.conf', 'rb') as f:
			for line in f:
				line = line.strip()
				if line.startswith('SERVICE_NAME'):
					service_name = line.split(' ')[-1]
				if line.startswith('HTTP_PORT'):
					service_port = line.split(' ')[-1]

		_SERVICE_INFO = {
			'name': service_name,
			'port': service_port
		}

	return _SERVICE_INFO

def is_in_dev_machine():
	if os.name == 'nt':
		return True

	if sys.platform == 'darwin':
		return True

	return False

class ApiResponse(object):
	"""
	api call的response
	"""
	def __init__(self, response):
		self.text = response.text.strip()
		self.raw_response = response

	@property
	def json_data(self):
		return json.loads(self.text)['data']

	@property
	def data(self):
		return json.loads(self.text)['data']

	@property
	def body(self):
		return json.loads(self.text)

	@property
	def json(self):
		return json.loads(self.text)

	@property
	def is_success(self):
		"""
		判断该次请求是否成功
		"""
		r = self.raw_response
		if r.status_code != 200:
			assert False, "http status code is %d, http call is FAILED!!!!" % r.status_code
			return False

		if 'html>' in self.text:
			assert False, "NOT a valid json string, call api FAILED!!!!"
			return False

		if self.json['code'] == 200:
			return True
		else:
			print '-*-' * 20
			print self.text
			print '-*-' * 20
			assert 200 == self.json['code'], "json[code] != 200, call api FAILED!!!!"
			return False

	@property
	def is_fail(self):
		if self.json["code"] != 200:
			return True
		else:
			return False

	def __repr__(self):
		return self.text

class RestClient(object):
	"""
	访问rest资源的client
	"""
	def __init__(self):
		self.jwt_token = None
		self.cookies = {}

	def add_cookie(self, key, value):
		self.cookies[key] = value

	def __get_url(self, type, resource):
		service_info = get_service_info()

		service = None
		if ':' in resource:
			service, resource = resource.split(':')

		if not service is None and service == service_info['name']:
			if is_in_dev_machine():
				#本机开发，直接访问127.0.0.1
				service = None

		if not service:
			if not is_in_dev_machine():
				#自动化测试环境，没有携带service的，自动转为当前service
				servivce = service_info['name']

		pos = resource.rfind('.')
		if pos == -1:
			raise RuntimeError('INVALID RESOURCE: %s' % resource)
		app = resource[:pos].replace('.', '/')
		app_resource = resource[pos+1:]

		if service:
			url = 'http://devapi.vxiaocheng.com/%s/%s/%s/' % (service, app, app_resource)
		else:
			url = 'http://127.0.0.1:%s/%s/%s/' % (service_info['port'], app, app_resource)
		if type == 'put':
			url = '%s?_method=put' % url
		elif type == 'delete':
			url = '%s?_method=delete' % url

		print "url: ", url

		return url

	def get(self, resource, data={}, context=None):
		url = self.__get_url('get', resource)

		headers = {}
		if self.jwt_token:
			headers = {
				'AUTHORIZATION': self.jwt_token
			}
		r = requests.get(url, data, headers=headers, cookies=self.cookies)
		return ApiResponse(r)

	def post(self, resource, data={}, context=None):
		url = self.__get_url('post', resource)

		headers = {}
		if self.jwt_token:
			headers = {
				'AUTHORIZATION': self.jwt_token
			}
		r = requests.post(url, data, headers=headers, cookies=self.cookies)
		return ApiResponse(r)

	def put(self, resource, data={}, context=None):
		url = self.__get_url('put', resource)

		headers = {}
		if self.jwt_token:
			headers = {
				'AUTHORIZATION': self.jwt_token
			}
		r = requests.post(url, data, headers, headers=headers, cookies=self.cookies)
		return ApiResponse(r)

	def delete(self, resource, data={}, context=None):
		url = self.__get_url('delete', resource)

		headers = {}
		if self.jwt_token:
			headers = {
				'AUTHORIZATION': self.jwt_token
			}
		r = requests.post(url, data, headers=headers, cookies=self.cookies)
		return ApiResponse(r)

def login(type, user, password=None, **kwargs):
	if not password:
		password = 'db7c6f3cf1ddda9498dd0148b87038f1'

	client = RestClient()
	if type == 'app':
		resp = client.put("ginger-account:login.logined_bdd_user", {
			'username': user,
			'type': 'user'
		})
		assert resp.is_success
	elif type == 'backend':
		data = {
			"username": user,
			"password": password
		}
		resp = client.put('ginger-account:login.logined_corp_user', data)
		assert resp.is_success

	print '-*-' * 20
	print resp.data
	print '-*-' * 20
	client.jwt_token = resp.data['jwt']
	client.cur_user_id = resp.data['id']

	if 'context' in kwargs:
		context = kwargs['context']
		if context:
			context.client = client

	return client