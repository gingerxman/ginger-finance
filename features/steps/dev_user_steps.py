# -*- coding:utf-8 -*-
#pylint: disable=E0602,E0102

import json

from behave import *
from features.bdd import util as bdd_util

# from db import user_models

NAME2INFO = {
	'zhouxun': {
		'display_name': u'周迅',
		'phone': '13811223300',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/zhouxun/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/zhouxun/cover.jpg',
		'gender': 'female',
		'province': u'北京市',
		'city': u'海淀区',
		'slogan': u'算什么男人...',
		'latitude': 40.032981,
		'longitude': 116.320205
	},
	'yangmi': {
		'display_name': u'杨幂',
		'phone': '13811223301',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/yangmi/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/yangmi/cover.jpg',
		'gender': 'female',
		'province': u'北京市',
		'city': u'海淀区',
		'slogan': u'我在这里欢笑 我在这里哭泣 北京北京',
		'latitude': 40.032981,
		'longitude': 116.320205
	},
	'sunli': {
		'display_name': u'孙俪',
		'phone': '13811223302',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/sunli/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/sunli/cover.jpg',
		'gender': 'female',
		'province': u'上海市',
		'city': u'徐汇区',
		'slogan': u'',
		'latitude': 31.194962,
		'longitude': 121.337729
	},
	'gal': {
		'display_name': u'盖尔.加朵',
		'phone': '13811223303',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/gal/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/gal/cover.jpg',
		'gender': 'male',
		'province': u'江苏省',
		'city': u'无锡市',
		'slogan': u'',
		'latitude': 31.525020,
		'longitude': 120.227050
	},
	'tangwei': {
		'display_name': u'汤唯',
		'phone': '13811223304',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/tangwei/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/tangwei/cover.jpg',
		'gender': 'male',
		'province': u'江苏省',
		'city': u'无锡市',
		'slogan': u'',
		'latitude': 31.525020,
		'longitude': 120.227050
	},
	'baby': {
		'display_name': u'AngelaBaby',
		'phone': '13811223305',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/baby/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/baby/cover.jpg',
		'gender': 'female',
		'province': u'上海市',
		'city': u'徐汇区',
		'slogan': u'',
		'latitude': 31.194962,
		'longitude': 121.337729
	},
	'zhaoliyin': {
		'display_name': u'赵丽颖',
		'phone': '13811223306',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/zhaoliyin/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/zhaoliyin/cover.jpg',
		'gender': 'female',
		'province': u'江苏省',
		'city': u'无锡市',
		'slogan': u'',
		'latitude': 31.525020,
		'longitude': 120.227050
	},
	'liuyan': {
		'display_name': u'柳岩',
		'phone': '13811223307',
		'avatar': 'http://resource.vxiaocheng.com/veeno/demo/girls/liuyan/avatar.jpg',
		'cover': 'http://resource.vxiaocheng.com/veeno/demo/girls/liuyan/cover.jpg',
		'gender': 'female',
		'province': u'江苏省',
		'city': u'南京市',
		'slogan': u'每屏泡泡棒棒哒~',
		'latitude': 32.023668,
		'longitude': 118.787248
	}
}

@given(u"{username}注册为艺人-droped")
def step_impl(context, username):
	extra_data = NAME2INFO.get(username, {})
	if context.text:
		extra_data.update(json.loads(context.text))
	client = bdd_util.login('microapp_user', username, context=context, extra_data=extra_data)
	context.client = client

	#创建艺人
	artist_data = {
		"name": username,
		"avatar": extra_data.get('avatar', ''),
		"age": extra_data.get('age', 18),
		"sex": extra_data.get("gender", "female"),
		"country": extra_data.get("country", u"中国"),
		"province": extra_data.get("province", u"江苏省"),
		"city": extra_data.get("city", u"南京市")
	}
	response = context.client.put("coral:artist.artist", artist_data)
	bdd_util.assert_api_call_success(response)

	#更新艺人信息
	artist_data['region'] = u'%s,%s' % (extra_data.get("province", u"江苏省"), extra_data.get("city", u"南京市"))
	update_data = artist_data
	if context.text:
		# from db import tag_models
		# tag_ids = None
		# input_data = json.loads(context.text)
		# if 'tags' in input_data:
		# 	tag_ids = [db_model.id for db_model in tag_models.Tag.select().dj_where(name__in=input_data['tags'])]
		# if tag_ids:
		# 	update_data['tag_ids'] = json.dumps(tag_ids)

		input_data = json.loads(context.text)
		for field in ['weight', 'height', 'bwh']:
			if field in input_data:
				update_data[field] = input_data[field]
	response = context.client.post("coral:artist.artist", update_data)
	bdd_util.assert_api_call_success(response)

	#更新user相关信息
	user_data = {
		"name": username,
		'birthday': '1995-03-08',
		'slogan': extra_data.get('slogan', ''),
		'region': u'%s,%s' % (extra_data.get("province", u"江苏省"), extra_data.get("city", u"南京市"))
	}
	response = context.client.post("skep:account.user", user_data)
	bdd_util.assert_api_call_success(response)

	#更新user location
	location_data = {
		'lat': extra_data.get('latitude', ''),
		'lng': extra_data.get('longitude', ''),
		'region': u'%s,%s' % (extra_data.get("province", u"江苏省"), extra_data.get("city", u"南京市"))
	}
	response = context.client.post("skep:account.user_location", location_data)
	bdd_util.assert_api_call_success(response)

