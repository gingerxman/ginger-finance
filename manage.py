#!/usr/bin/env python

import sys

def load_local_command(command):
    module_name = 'command.%s' % command
    #try:
    print 'load local command: ', module_name
    module = __import__(module_name, {}, {}, ['*',])
    return module
    #except Exception:
    #    return None

def run_command(command):
    command_module = load_local_command(command)

    if not command_module:
        print 'no command named: ', command
    else:
        instance = getattr(command_module, 'Command')()
        try:
            instance.handle(*sys.argv[2:])
        except TypeError, e:
            print '[ERROR]: wrong command arguments, usages:'
            print instance.help
            print 'Exception: {}'.format(e)

if __name__ == '__main__':
    command = sys.argv[1]
    run_command(command)

