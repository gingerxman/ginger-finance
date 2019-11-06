# -*- coding: utf-8 -*-
import json
import time
import logging
import requests
from datetime import datetime, timedelta


import db_util
from features.bdd.client import RestClient

exec_sql = db_util.exec_sql

tc = None

def convert_to_same_type(a, b):
	def to_same_type(target, other):
		target_type = type(target)
		other_type = type(other)
		if other_type == target_type:
			return True, target, other

		if (target_type == float) or (other_type == float):
			target_type = float
		if (target_type == int) or (target_type == float):
			try:
				other = target_type(other)
				return True, target, other
			except:
				return False, target, other

		return False, target, other

	is_success, new_a, new_b = to_same_type(a, b)
	if is_success:
		return new_a, new_b
	else:
		is_success, new_b, new_a = to_same_type(b, a)
		if is_success:
			return new_a, new_b

	return a, b


###########################################################################
# assert_dict: 验证expected中的数据都出现在了actual中
###########################################################################
def assert_dict(expected, actual):
	global tc
	is_dict_actual = isinstance(actual, dict)
	for key in expected:
		expected_value = expected[key]
		if is_dict_actual:
			actual_value = actual[key]
		else:
			actual_value = getattr(actual, key)

		if isinstance(expected_value, dict):
			assert_dict(expected_value, actual_value)
		elif isinstance(expected_value, list):
			assert_list(expected_value, actual_value, {'key':key})
		else:
			try:
				expected_value, actual_value = convert_to_same_type(expected_value, actual_value)
				tc.assertEquals(expected_value, actual_value)
			except Exception, e:
				items = ['\n<<<<<', 'e: %s' % str(expected), 'a: %s' % str(actual), 'key: %s' % key, e.args[0], '>>>>>\n']
				e.args = ('\n'.join(items),)
				raise e


def assert_list(expected, actual, options=None):
	"""
	验证expected中的数据都出现在了actual中
	"""
	global tc
	try:
		hint = 'list length DO NOT EQUAL: %d != %d' % (len(expected), len(actual))

		if options and 'key' in options:
			hint = '%s - %s' % (options['key'], hint)
		tc.assertEquals(len(expected), len(actual), hint)
	except:
		if options and 'key' in options:
			print '      Outer Compare Dict Key: ', options['key']
		raise

	for i in range(len(expected)):
		expected_obj = expected[i]
		actual_obj = actual[i]
		if isinstance(expected_obj, dict):
			assert_dict(expected_obj, actual_obj)
		else:
			expected_obj, actual_obj = convert_to_same_type(expected_obj, actual_obj)
			tc.assertEquals(expected_obj, actual_obj)



def assert_api_call(response, context):
	if context.text:
		input_data = json.loads(context.text)
		if 'error' in input_data:
			assert_api_call_failed(response, input_data['error'])
			return False
		else:
			assert_api_call_success(response)
			return True
	else:
		assert_api_call_success(response)
		return True


###########################################################################
# assert_api_call_success: 验证api调用成功
###########################################################################
def assert_api_call_success(response):
	if 200 != response.body['code']:
		buf = []
		buf.append('>>>>>>>>>>>>>>> response <<<<<<<<<<<<<<<')
		buf.append(str(response))
		logging.error("API calling failure: %s" % '\n'.join(buf))
	assert 200 == response.body['code'], "code != 200, call api FAILED!!!!"


###########################################################################
# assert_api_call_failed: 验证api调用失败
###########################################################################
def assert_api_call_failed(response, expected_err_code=None):
	if 200 == response.body['code']:
		buf = []
		buf.append('>>>>>>>>>>>>>>> response <<<<<<<<<<<<<<<')
		buf.append(str(response))
		logging.error("API calling not expected: %s" % '\n'.join(buf))
	assert 200 != response.body['code'], "code == 200, call api NOT EXPECTED!!!!"
	if expected_err_code:
		actual_err_code = str(response.body['errCode'])
		assert expected_err_code in actual_err_code, "errCode(%s) != '%s', error code FAILED!!!" % (actual_err_code, expected_err_code)



###########################################################################
# assert_expected_list_in_actual: 验证expected中的数据都出现在了actual中
###########################################################################
def assert_expected_list_in_actual(expected, actual):
	global tc

	for i in range(len(expected)):
		expected_obj = expected[i]
		actual_obj = actual[i]
		if isinstance(expected_obj, dict):
			assert_dict(expected_obj, actual_obj)
		else:
			try:
				tc.assertEquals(expected_obj, actual_obj)
			except Exception, e:
				items = ['\n<<<<<', 'e: %s' % str(expected), 'a: %s' % str(actual), 'key: %s' % key, e.args[0], '>>>>>\n']
				e.args = ('\n'.join(items),)
				raise e


###########################################################################
# print_json: 将对象以json格式输出
###########################################################################
def print_json(obj):
	print json.dumps(obj, indent=True)


def table2dict(context):
	expected = []
	for row in context.table:
		data = {}
		for heading in row.headings:
			if ':' in heading:
				real_heading, value_type = heading.split(':')
			else:
				real_heading = heading
				value_type = None
			value = row[heading]
			if value_type == 'i':
				value = int(value)
			if value_type == 'f':
				value = float(value)
			data[real_heading] = value
		expected.append(data)
	return expected


def get_date(str):
	"""
		将字符串转成datetime对象
		今天 -> 2014-4-18
	"""
	# 处理expected中的参数
	str = str.strip()
	today = datetime.now()
	if u'本月' in str:
		month = today.strftime('%Y-%m-')
		str = str.replace(u'本月', month)
	elif u'上月' in str:
		month = (today - timedelta(31)).strftime('%Y-%m-')
		str = str.replace(u'上月', month)

	is_specify_time = (' ' in str)
	if u'今天' in str:
		delta = 0
	elif u'昨天' in str:
		delta = -1
	elif u'前天' in str:
		delta = -2
	elif u'明天' in str:
		delta = 1
	elif u'后天' in str:
		delta = 2
	elif u'天后' in str:
		if is_specify_time:
			delta = int(str.split(' ')[0][:-2])
		else:
			delta = int(str[:-2])
	elif u'天前' in str:
		if is_specify_time:
			delta = 0 - int(str.split(' ')[0][:-2])
		else:
			delta = 0 - int(str[:-2])
	else:
		is_specify_time = False
		tmp = str.split(' ')
		if len(tmp) == 1:
			strp = "%Y-%m-%d"
		elif len(tmp[1]) == 8:
			strp = "%Y-%m-%d %H:%M:%S"
		elif len(tmp[1]) == 5:
			strp = "%Y-%m-%d %H:%M"
		return datetime.strptime(str, strp)

	if is_specify_time:
		date = (today + timedelta(delta)).strftime('%Y-%m-%d')
		specified_time = str.split(' ')[1]
		if len(specified_time) == 8:
			strp = "%Y-%m-%d %H:%M:%S"
		elif len(specified_time) == 5:
			strp = "%Y-%m-%d %H:%M"
		str = '%s %s' % (date, specified_time)
		return datetime.strptime(str, strp)
	else:
		return today + timedelta(delta)


def get_date_str(str):
	date = get_date(str)
	return date.strftime('%Y-%m-%d')

#
# #获得corp user对应的corp的id
# def get_corp_id_for_corpuser(client, username):
# 	data = {
# 		"username": username,
# 		"password": '55e421ee9bdc9d9f6b6c6518E590b0ee'
# 	}
# 	resp = client.put('ginger-account:login.logined_corp_user', data)
#
# 	return resp.data['cid']
#
# def get_corp_uuid_for_corpuser(client, username):
# 	data = {
# 		"username": username
# 	}
# 	resp = client.get('ginger-account:dev.corp_uuid', data)
#
# 	return resp.data['uuid']
#
# def get_corp_token_for_corpuser(client, username):
# 	corp_id = get_corp_id_for_corpuser(client, username)
# 	data = {
# 		"corp_id": corp_id
# 	}
# 	resp = client.get('ginger-account:dev.corp_token', data)
#
# 	return resp.data['__cs']
#
#
# def get_users_by_names(client, names):
# 	users = []
# 	for name in names:
# 		resp = get_user_by_name(client, name)
# 		users.append(resp.data)
# 	return users
#
#
# def get_user_by_name(client, name):
# 	resp = client.put("ginger-account:login.logined_bdd_user", {
# 		'name': name,
# 	})
# 	assert_api_call_success(resp)
# 	return resp
#
#
# def get_user_by_id(client, id):
# 	resp = client.put("ginger-account:login.dev_logined_user", {
# 		'user_id': id,
# 	})
# 	return resp
