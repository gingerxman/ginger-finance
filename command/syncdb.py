# -*- coding: utf-8 -*-

import os
import json
import shutil

from command_base import BaseCommand

class Command(BaseCommand):
	help = "syncdb"
	args = ''
	
	def handle(self, **options):
		cmd = "go run command/syncdb.go"
		print 'run> ', cmd
		os.system(cmd)
