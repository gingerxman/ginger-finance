# -*- coding: utf-8 -*-

import os
import json
import shutil

from command_base import BaseCommand

class Command(BaseCommand):
	help = "syncdb"
	args = ''
	
	def handle(self, **options):
		cmd = "go run command/cmd.go orm sqlall"
		print 'run> ', cmd
		os.system(cmd)
