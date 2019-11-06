# -*- coding: utf-8 -*-

import os

import sys
path = os.path.abspath(os.path.join('.', '..'))
sys.path.insert(0, path)
reload(sys)
sys.setdefaultencoding('utf8')

import unittest
from features.bdd import util as bdd_util

def __clear_all_account_data():
	"""
	清空账号数据
	"""
	#User.objects.Where(id__gt=3).delete()


clean_modules = []
def __clear_all_app_data():
	"""
	清空应用数据
	"""
	if len(clean_modules) == 0:
		for clean_file in os.listdir('./features/clean'):
			if clean_file.startswith('__'):
				continue

			if clean_file.startswith('.'):
				#skip .DS_st in mac
				continue

			module_name = 'features.clean.%s' % clean_file[:-3]
			module = __import__(module_name, {}, {}, ['*',])	
			clean_modules.append(module)

	for clean_module in clean_modules:
		clean_module.clean()


def before_all(context):
	__clear_all_account_data()

	#创建test case，使用assert
	context.tc = unittest.TestCase('__init__')
	bdd_util.tc = context.tc


def after_all(context):
	pass


def before_scenario(context, scenario):
	context.scenario = scenario

	context.execute_steps(u"Given 重置服务")


def after_scenario(context, scenario):
	pass

