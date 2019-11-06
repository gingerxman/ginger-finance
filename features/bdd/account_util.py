# -*- coding: utf-8 -*-
import json
import time
import logging
import requests
from datetime import datetime, timedelta
from bdd.client import RestClient
from bdd import util as bdd_util

class Corp(object):
    def __init__(self, id):
        self.id = id

    def join_platform(self, platform_name):
		client = RestClient()
		
		data = {
			'platform_name': platform_name,
			'corp_id': self.id
		}
		response = client.put("account.platform_member_corp", data)
		assert response.is_success

def __create_corp(username, display_name, corp_type):
	client = RestClient()

	data = {
		'username': username,
		'display_name': display_name or username,
		'password': 'test',
		'type': corp_type
	}
	response = client.put("account.corp", data)
	assert response.is_success

	return Corp(response.data['id'])

def create_general_corp(username, display_name=None):
    return __create_corp(username, display_name, 'general')

def create_platform_corp(username, display_name=None):
    return __create_corp(username, display_name, 'platform')

def create_supplier_corp(username, display_name=None):
    return __create_corp(username, display_name, 'supplier')
