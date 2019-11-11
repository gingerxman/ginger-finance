Feature: 解冻已冻结的虚拟资产

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

	@ginger-finance @imoney @wip
	Scenario: 解冻用户使用的虚拟资产
		Given lucy访问'jobs'的商城
		When lucy充值'1000'个'bitcoin'
		Then lucy能获得虚拟资产'bitcoin'
		"""
		{
			"balance": 1000
		}
		"""

		# 冻结资产
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

		# 解冻资产
		When lucy取消对'989'个'bitcoin'的使用
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
