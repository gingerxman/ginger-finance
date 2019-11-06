Feature: 虚拟资产充值

	Background:
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
			"exchange_rate": 1
		}, {
			"code": "bitcoin",
			"exchange_rate": 0.9
		}]
		"""
		When ginger创建公司
		"""
		[{
			"name": "MIX",
			"username": "jobs"
		}, {
			"name": "BabyFace",
			"username": "bill"
		}]
		"""


@ginger-finance @imoney
Scenario: 用户可以充值虚拟资产
	Given lucy访问'jobs'的商城
	When lucy充值'1000'个'cash'
	Then lucy能获得虚拟资产'cash'
	"""
	{
		"balance": 1000
	}
	"""
	When lucy充值'1000'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 1000
	}
	"""
	Given lily访问'jobs'的商城
	Then lily能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 0
	}
	"""
	When lily充值'2000'个'bitcoin'
	Then lily能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 2000
	}
	"""
