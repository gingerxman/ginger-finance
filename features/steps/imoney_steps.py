# -*- coding:utf-8 -*-

import json
from datetime import datetime

from behave import *

from features.bdd import util as bdd_util
from features.steps import step_util
from features.bdd.client import RestClient

def get_imoney_id_by_code(code):
	objs = bdd_util.exec_sql("select * from imoney_imoney where code = %s", [code])
	return objs[0]['id']

def get_frozen_record_id_by_code(imoney_code, amount):
	sql = "select * from account_frozen_record where imoney_code = %s and amount = %s", [imoney_code, amount]
	objs = bdd_util.exec_sql("select * from account_frozen_record where imoney_code = %s and amount = %s", [imoney_code, amount])
	return objs[0]['id']

IMONEY_STATUS_NAME2STATUS = {
	u'启用': 1,
	u'停用': 0
}

IMONEY_STATUS_STATUS2NAME = {
	1: u'启用',
	0: u'停用'
}


@Given(u"{corp_user}配置虚拟资产")
def step_impl(context, corp_user):
	imoney_configs = json.loads(context.text)

	for imoney_config in imoney_configs:
		code = imoney_config['code']
		exchange_rate = imoney_config['exchange_rate']
		display_name = imoney_config.get('display_name', code)
		is_payable = imoney_config.get('is_payable', True)
		is_debtable = imoney_config.get('is_debtable', False)

		resp = context.client.put("imoney.imoney", {
			"code": code,
			"exchange_rate": exchange_rate,
			"display_name": display_name,
			"is_payable": is_payable,
			"is_debtable": is_debtable
		})
		bdd_util.assert_api_call_success(resp)

@given(u"系统为'{user_name}'转账'{amount}'个'{imoney_code}'")
def step_impl(context, user_name, amount, imoney_code):
	user_id = step_util.get_user_id_by_name(context.client, user_name)
	data = {
		"source_user_id": 0,#client.cur_user_id,
		"dest_user_id": user_id,
		"imoney_code": imoney_code,
		"amount": int(amount) * 100,
		"bid": "bdd"
	}
	resp = context.client.put('imoney.transfer', data)
	bdd_util.assert_api_call_success(resp)

@when(u"{user}充值'{amount}'个'{imoney_code}'")
def step_impl(context, user, amount, imoney_code):
	step = u"Given 系统为'%s'转账'%s'个'%s'" % (user, amount, imoney_code)
	context.execute_steps(step)


@when(u"{user}使用'{amount}'个'{imoney_code}'")
def step_impl(context, user, amount, imoney_code):
	response = context.client.put('imoney.frozen_record', {
		'imoney_code': imoney_code,
		'amount': int(amount) * 100,
		'type': 'consume',
		'remark': ''
	})
	bdd_util.assert_api_call(response, context)

@when(u"{user}取消对'{amount}'个'{imoney_code}'的使用")
def step_impl(context, user, amount, imoney_code):
	amount = int(amount)
	record_id = get_frozen_record_id_by_code(imoney_code, 100*amount)
	response = context.client.delete('imoney.frozen_record', {
		'id': record_id
	})
	bdd_util.assert_api_call(response, context)

@then(u"{user}能获得已冻结虚拟资产'{imoney_code}'")
def step_impl(context, user, imoney_code):
	expected = json.loads(context.text)
	response = context.client.get('imoney.frozen_record', {
		'imoney_code': imoney_code
	})
	actual = response.data
	actual['frozen_amount'] = bdd_util.format_price((actual['frozen_amount']))

	bdd_util.assert_dict(expected, actual)

@then(u"{user}拥有'{amount}'个'{imoney_code}'")
def step_impl(context, user, amount, imoney_code):
	response = context.client.get('imoney.balance.', {
		'imoney_code': imoney_code,
	})
	balance = response.data['valid_balance']
	assert amount * 100 == balance, u'账户余额不对: expect(%s), actual(%s)' % (amount, balance)

@then(u"{user}能获得虚拟资产'{imoney_code}'")
def step_impl(context, user, imoney_code):
	response = context.client.get('imoney.balance', {
		'imoney_code': imoney_code,
	})

	actual = bdd_util.format_price(response.data)
	expected = json.loads(context.text)['balance']

	assert actual == expected, 'actual(%s) != expected(%d)' % (actual, expected)

@when(u"{user}通过虚拟资产'{imoney_code}'的数量为'{count}'的提现申请")
def step_impl(context, user, imoney_code, count):
	response = context.client.get('/imoney/withdraw_records', {
		'imoney_code': imoney_code
	})

	count = int(count)
	target_record = None
	for record in response.data['records']:
		if record['amount'] == count:
			target_record = record
			break

	response = context.client.put('/imoney/finished_withdraws/', {
		'ids': json.dumps([target_record['id']]),

	})
	bdd_util.assert_api_call_success(response)

def __update_withdraw_id(old_id, new_id):
	"""
	修改提现记录的id
	"""
	sql = """
		update imoney_withdraw 
		set id={} where id={};
	""".format(new_id, old_id)
	db_util.SQLService.use('plutus').execute_sql(sql)

	sql = """
		update flow_task 
		set resource_id={} where resource_id={};
	""".format(new_id, old_id)
	db_util.SQLService.use('plutus').execute_sql(sql)

	sql = """
			update flow_task_log 
			set resource_id={} where resource_id={};
		""".format(new_id, old_id)
	db_util.SQLService.use('plutus').execute_sql(sql)


@given(u"系统设置手动完成已同意通过的提现为'{action}'")
def step_impl(context, action):

	context.AUTO_PASS = False if action == u'开启' else True

@when(u"系统通过已同意的提现")
def finish_passed_withdraws(context):
	response = context.client.put('commands.handle_passed_withdraw', {})
	bdd_util.assert_api_call_success(response)

@when(u"{user}申请提现")
def step_impl(context, user):
	data = json.loads(context.text)

	user_role = user_util.get_user_role(user).lower()
	channel = WITHDRAW_STR2TYPE[data['channel']]

	response = context.client.put('imoney.withdraw', {
		'imoney_code': data['imoney'],
		'amount': data['amount'],
		'user_type': user_role,
		'channel': channel
	})
	print response.body
	bdd_util.assert_api_call(response, context)

	if response.body['code'] == 200:
		record_id = response.data['id']
		if data.has_key('id'):
			record_id = data['id']
			__update_withdraw_id(response.data['id'], data['id'])

		# if imoney_models.Withdraw.select().dj_where(status=imoney_models.IMONEY_WITHDRAW_STATUS['PASS'], id=record_id).first():
		# 	if getattr(context, 'AUTO_PASS', True):
		# 		finish_passed_withdraws(context)


def __get_withdraw_status_text(status):
	if status == 'requesting':
		return u'待处理'
	elif status == 'success':
		return u'已完成'
	elif status == 'rejected':
		return u'已驳回'
	elif status == 'waiting':
		return u'暂不处理'
	elif status == 'failed':
		return u'人工处理'
	elif status == 'pass':
		return u'已同意'
	else:
		return u'未知'

WITHDRAW_STR2TYPE = {
	u'微信': 'weixin',
	u'银行卡': 'bank',
}

@then(u"{user}能获得虚拟资产'{imoney_code}'的提现记录列表")
def step_impl(context, user, imoney_code):
	expected = json.loads(context.text)
	response = context.client.get('imoney.withdraw_records', {
		'imoney_code': imoney_code,
	})
	actual = response.data['records']
	print actual
	for record in actual:
		record['status'] = __get_withdraw_status_text(record['status'])

	bdd_util.assert_list(expected, actual)


@then(u"{user}能获得提现记录列表")
def step_impl(context, user):
	expected = json.loads(context.text)
	response = context.client.get('imoney.withdraws', {
		'with_options': json.dumps({'with_flow_logs': True})
	})
	actual = response.data['records']
	print actual
	for record in actual:
		record['status'] = __get_withdraw_status_text(record['status'])
		record['imoney'] = record['imoney_code']
		record['user'] = user_util.get_username_by_corp_user_id(record['user_id'], record['extra_data']['user_type'] == 'corp')
		record['flow_logs'] = map(lambda log: u'{}{}'.format(log['username'], log['remark']), record['flow_logs'])

	bdd_util.assert_list(expected, actual)

@when(u"{user}审核通过提现申请")
def step_impl(context, user):
	datas = json.loads(context.text)
	for data in datas:
		response = context.client.put('imoney.finished_withdraw', {
			'id': data['id'],
			'remark': data.get('remark', '')
		})
		bdd_util.assert_api_call(response, context)

@when(u"{user}暂不处理提现申请")
def step_impl(context, user):
	datas = json.loads(context.text)
	for data in datas:
		response = context.client.put('imoney.waited_withdraws', {
			'ids': json.dumps([data['id']]),
			'remark': data.get('remark', '')
		})
		bdd_util.assert_api_call(response, context)

@when(u"{user}驳回提现申请")
def step_impl(context, user):
	datas = json.loads(context.text)
	for data in datas:
		response = context.client.put('imoney.rejected_withdraws', {
			'ids': json.dumps([data['id']]),
			'remark': data.get('remark', '')
		})
		bdd_util.assert_api_call(response, context)

@when(u"系统发起周期提现")
def step_impl(context):
	response = context.client.put('commands.do_periodic_withdraw', {})
	bdd_util.assert_api_call(response, context)

@when(u"{source_user}向'{dest_user}'转账'{amount}'个'{imoney_code}'")
def step_impl(context, source_user, dest_user, amount, imoney_code):
	response = context.client.put('imoney.transfer', {
		'bid': BidFactory().generate_order_bid(),
		'source_user_id': user_util.get_user_id_by_name(source_user) if source_user not in ('manager', 'xiaocheng') else 0,
		'dest_user_id': user_util.get_user_id_by_name(dest_user) if dest_user not in ('manager', 'xiaocheng') else 0,
		'imoney_code': imoney_code,
		'amount': amount,
	})
	bdd_util.assert_api_call(response, context)

@when(u"系统发起批量交易")
def step_impl(context):
	data = json.loads(context.text)
	for d in data:
		d['source_user_id'] = user_util.get_user_id_by_name(d['source_user']) if d['source_user'] not in ('manager', 'xiaocheng') else 0
		d['dest_user_id'] = user_util.get_user_id_by_name(d['dest_user']) if d['dest_user'] not in ('manager', 'xiaocheng') else 0
		del d['source_user']
		del d['dest_user']
	print data
	response = context.client.put('imoney.transfers', {
		'transfer_data': json.dumps(data)
	})
	bdd_util.assert_api_call(response, context)

@then(u"{user}能获取用户的'{imoney_code}'账户信息")
def step_impl(context, user, imoney_code):
	expected = json.loads(context.text)
	user_ids = []
	for e in expected:
		uid = user_util.get_user_id_by_name(e['user'])
		user_ids.append(uid)
		e['user_id'] = uid
		del e['user']

	resp = context.client.get('imoney.users_balance', {
		'user_ids': json.dumps(user_ids),
		'imoney_code': imoney_code,
		'with_options': json.dumps({
			'with_total_income': True,
		})
	})
	print resp

	user_id2actual_data = {a['user_id']: a for a in resp.data}

	for d in expected:
		bdd_util.assert_dict(d, user_id2actual_data[d['user_id']])


@then(u"{user}能获取'{imoney_code}'账户信息")
def step_impl(context, user, imoney_code):
	resp = context.client.get('imoney.balance', {
		'imoney_code': imoney_code,
		'_v': 2,
		'with_options': json.dumps({
			'with_total_income': True,
		})
	})
	expected = json.loads(context.text)
	actual = resp.data
	bdd_util.assert_dict(expected, actual)

@then(u"{user}能获取'{imoney_code}'余额")
def step_impl(context, user, imoney_code):
	expected = json.loads(context.text)
	resp = context.client.get('imoney.balance', {
		'imoney_code': imoney_code,
		'_v': 2,
	})

	actual = resp.data['valid_balance']
	print expected, actual
	assert expected == actual

@then(u"系统专用账户'{account_code}'的余额为")
def step_impl(context, account_code):
	expected = json.loads(context.text)
	sql = """
		select balance from account_account 
		where code='{}';
	""".format(account_code)
	record = SQLService.use('plutus').execute_sql(sql).fetchone()

	assert(float(expected) == float(record[0]))