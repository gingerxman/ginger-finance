Feature: 初始化系统数据

@full_init
Scenario: 初始化系统数据
	Given ginger登录系统
	When ginger创建公司
		"""
		[{
			"name": "MIX",
			"username": "jobs"
		}, {
			"name": "BabyFace",
			"username": "bill"
		}, {
			"name": "Mocha",
			"username": "tom"
		}]
		"""

