# -*- coding: utf-8 -*-
import json

# from features.util import bdd_util
from features.bdd import util as bdd_util
from behave import *

def create_corp(context, corp_name, username):
	"""
	创建门店
	"""
	response = context.client.put('ginger-account:corp.corp', {
		"corp_name": corp_name,
		"username": username,
		"password": "test"
	})
	bdd_util.assert_api_call_success(response)

@When(u"{user}创建公司")
def step_impl(context, user):
	input_datas = json.loads(context.text)

	for index, input_data in enumerate(input_datas):
		create_corp(context, input_data['name'], input_data['username'])

	# for input_data in input_datas:
	# 	context.execute_steps(u"Given %s登录系统" % input_data['username'])
	# 	context.execute_steps(u"When %s初始化商城" % input_data['username'])
	#
	# context.execute_steps(u"Given %s登录系统" % user)
	# context.execute_steps(u"When %s初始化商城" % user)
