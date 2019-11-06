# -*- coding:utf-8 -*-
#pylint: disable=E0602,E0102
#获得corp user对应的corp的id

def get_corp_id_for_corpuser(client, username):
	data = {
		"username": username,
		"password": 'db7c6f3cf1ddda9498dd0148b87038f1'
	}
	resp = client.put('ginger-account:login.logined_corp_user', data)

	return resp.data['cid']

def get_corp_uuid_for_corpuser(client, username):
	data = {
		"username": username
	}
	resp = client.get('ginger-account:dev.corp_uuid', data)

	return resp.data['uuid']

def get_corp_token_for_corpuser(client, username):
	corp_id = get_corp_id_for_corpuser(client, username)
	data = {
		"corp_id": corp_id
	}
	resp = client.get('ginger-account:dev.corp_token', data)

	return resp.data['__cs']


def get_users_by_names(client, names):
	users = []
	for name in names:
		resp = get_user_by_name(client, name)
		users.append(resp.data)
	return users


def get_user_by_name(client, name):
	resp = client.put("ginger-account:login.logined_bdd_user", {
		'username': name,
		"type": "user"
	})
	return resp

def get_user_id_by_name(client, name):
	resp = get_user_by_name(client, name)
	return resp.data['id']

