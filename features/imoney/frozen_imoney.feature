Feature: 冻结虚拟资产

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
			"exchange_rate": 1
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
Scenario: 用户成功使用虚拟资产
	Given lucy访问'jobs'的商城
	When lucy充值'1000'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 1000
	}
	"""

	When lucy使用'989'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 11
	}
	"""
	Then lucy能获得已冻结虚拟资产'bitcoin'
	"""
	{
		"frozen_amount": 989
	}
	"""

@ginger-finance @imoney
Scenario: 用户成功使用虚拟资产，使用完
	Given lucy访问'jobs'的商城
	When lucy充值'1000'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 1000
	}
	"""

	When lucy使用'1000'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 0
	}
	"""
	Then lucy能获得已冻结虚拟资产'bitcoin'
	"""
	{
		"frozen_amount": 1000
	}
	"""

@ginger-finance @imoney
Scenario: 用户使用超过余额的虚拟资产，失败
	Given lucy访问'jobs'的商城
	When lucy充值'1000'个'bitcoin'
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 1000
	}
	"""

	When lucy使用'1001'个'bitcoin'
	"""
	{
		"error": "frozen_record:not_enough_balance"
	}
	"""
	Then lucy能获得虚拟资产'bitcoin'
	"""
	{
		"balance": 1000
	}
	"""
	Then lucy能获得已冻结虚拟资产'bitcoin'
	"""
	{
		"frozen_amount": 0
	}
	"""
