Feature: 管理虚拟资产

@ginger-finance @imoney
Scenario: 1、系统管理员可以添加虚拟资产
	Given ginger登录系统
	Given ginger配置虚拟资产
	"""
	[{
		"code": "rmb",
		"exchange_rate": 1,
		"is_debtable": true,
		"is_payable": false
	}, {
		"code": "cash",
		"display_name": "现金",
		"exchange_rate": 1
	}, {
		"code": "bitcoin",
		"display_name": "比特币",
		"exchange_rate": 0.9
	}]
	"""
	Then ginger能获取虚拟资产列表
	"""
	[{
		"code": "rmb",
		"display_name": "rmb",
		"exchange_rate": 1,
		"is_debtable": true,
		"is_payable": false
	}, {
		"code": "cash",
		"display_name": "现金",
		"exchange_rate": 1,
		"is_debtable": false,
		"is_payable": false
	}, {
		"code": "bitcoin",
		"display_name": "比特币",
		"exchange_rate": 0.9,
		"is_debtable": false,
		"is_payable": false
	}]
	"""

@ginger-finance @imoney
Scenario: 2、系统管理员可以删除虚拟资产
	Given ginger登录系统
	Given ginger配置虚拟资产
	"""
	[{
		"code": "rmb",
		"exchange_rate": 1,
		"is_debtable": true,
		"is_payable": false
	}, {
		"code": "cash",
		"display_name": "现金",
		"exchange_rate": 1
	}, {
		"code": "bitcoin",
		"display_name": "比特币",
		"exchange_rate": 0.9
	}]
	"""
	When ginger删除虚拟资产"bitcoin"
	Then ginger能获取虚拟资产列表
	"""
	[{
		"code": "rmb",
		"display_name": "rmb",
		"exchange_rate": 1,
		"is_debtable": true,
		"is_payable": false
	}, {
		"code": "cash",
		"display_name": "现金",
		"exchange_rate": 1,
		"is_debtable": false,
		"is_payable": false
	}]
	"""