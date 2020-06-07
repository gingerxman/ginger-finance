# -*- coding: utf-8 -*-
import pymysql as mysql

def get_db_name():
	with open('conf/app.conf', 'rb') as f:
		for line in f:
			if "${_DB_NAME||" in line:
				beg = line.find('||')+2
				end = line.find('}', beg)

				return line[beg:end]

def exec_sql(sql, params):
	db_name = get_db_name()
	DB = mysql.connect(host='localhost',user='root',passwd='root',db=db_name,charset='utf8')
	cursor = DB.cursor()
	try:
		cursor.execute(sql, params)
		if sql.startswith("select"):
			columns = [col[0] for col in cursor.description]
			results = cursor.fetchall()
			return [dict(zip(columns, row)) for row in results]
		elif sql.startswith("update"):
			DB.commit()
	except:
		DB.rollback()
		print "sql: {} {}".format(sql, params)
		print '[db_util] exception!!!!!!!!!'
	finally:
		cursor.close()
		DB.close()

if __name__ == '__main__':
	for result in exec_sql("select * from auth_user", None):
		print result