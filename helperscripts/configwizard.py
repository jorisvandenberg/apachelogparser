#!/usr/local/bin/python3.11

import os
import configparser
import inquirer
import re
import sys
import platform

ini_data = {}

def write_ini_file(filename):
	global ini_data  # access the global variable
	config = configparser.ConfigParser()
	for section, options in ini_data.items():
		config.add_section(section)
		for option, value in options.items():
			config.set(section, option, value)
	with open(filename, 'w') as f:
		config.write(f)

def get_configfilename():
	questions = [
		inquirer.Path('configfile_dir', message="Where is the config file dir (existing full path)", path_type=inquirer.Path.DIRECTORY,),
		inquirer.Text('configfile_name', message="What is the name of the configfile (ending in .ini)", validate=validate_file_extension_ini),
	]
	answers = inquirer.prompt(questions)
	
	configfile_dir = answers['configfile_dir']
	configfile_name = answers['configfile_name']
	return configfile_dir+configfile_name

def get_default_log_path():
	distrodb = {
		'win': {
			'list': ['Windows'],
			'path': 'C:\\Program Files\\Apache\\Apache2.4\\logs\\'
		},
		'apache2': {
			'list': ['Darwin', 'Ubuntu', 'SUSE', 'OpenSUSE', 'Debian GNU/Linux', 'debian'],
			'path': '/var/log/apache2/'
		},
		'httpd': {
			'list': ['Fedora', 'CentOS', 'RHEL', 'Red Hat Enterprise Linux (RHEL)'],
			'path': '/var/log/apache2/'
		}
	}
	
	os_name = platform.system()

	# Detect the Linux distribution, if applicable
	if os_name == 'Linux':
		try:
			with open('/etc/os-release', 'r') as f:
				distro_info = {}
				for line in f:
					if '=' in line:
						key, value = line.strip().split('=', 1)
						distro_info[key.lower()] = value.strip().strip('"')
				distro = distro_info.get('name') or distro_info.get('id')
		except:
			distro = None
	else:
		distro = os_name
	
	for key, value in distrodb.items():
		if distro in value['list']:
			return value['path']
	
	return ""

	
def fill_input_section():
	osdefault = get_default_log_path()
	global ini_data
	domain = ini_data['general']['mydomain']
	pattern = pattern = r'^' + re.escape(domain) + r'.*'
	questions = [
		inquirer.Path('logfile_dir', message="Where can i find the logfiles?", path_type=inquirer.Path.DIRECTORY, default=osdefault),
		inquirer.Text('logfileregex', message="Regex which logiles i need to parse", default=pattern),
		inquirer.Text('parseregex', message="regex to match log format values (or clf)", default='clf'),
		inquirer.List("fullloadcheck", message="Do i need to parse every line of every logfile? (false = only lines newer than last load)", choices=["true", "false"], default="false"),
	]
	
	answers = inquirer.prompt(questions)
	logfile_dir = answers['logfile_dir']
	logfileregex = answers['logfileregex']
	parseregex = answers['parseregex']
	fullloadcheck = answers['fullloadcheck']
	ini_data['input'] = {
        'logfilepath': logfile_dir,
		'logfileregex': logfileregex,
		'parseregex': parseregex,
		'fullloadcheck': fullloadcheck,
    }

def fill_general_section():
	global ini_data
	questions = [
		inquirer.Path('database_dir', message="Where is the database (existing full path)", path_type=inquirer.Path.DIRECTORY,),
		inquirer.Text('database_file', message="What is the name of the database (ending in .db)", validate=validate_file_extension_db_sqlite3),
		inquirer.Text('timeformat', message="enter a valid timeformat", default="02/Jan/2006:15:04:05 -0700"),
		inquirer.Text('mydomain', message="enter your top level domain (mydomain.com)"),
		inquirer.List("writelog", message="Do i need to write logfiles?", choices=["true", "false"], default="true"),
	]
	answers = inquirer.prompt(questions)
	
	database_dir = answers['database_dir']
	database_file = answers['database_file']
	timeformat = answers['timeformat']
	mydomain = answers['mydomain']
	writelog_output = answers['writelog']

	# check if the file exists
	"""
	if !os.path.isfile(data_dir):
		os.exit(1)
	"""
	ini_data['general'] = {
        'dbpath': database_dir + database_file,
        'timeformat': timeformat,
        'mydomain': mydomain,
		'writelog' : writelog_output,
    }

def validate_file_extension_ini(answers, current):
	if not re.match(r".+\.ini$", current):
		raise inquirer.errors.ValidationError("", reason="I don't like config filename!")
	return True

def validate_file_extension_db_sqlite3(answers, current):
	if not re.match(r".+\.db$", current):
		raise inquirer.errors.ValidationError("", reason="I don't like database filename!")
	return True

def main():
	config_filename = get_configfilename()
	fill_general_section()
	fill_input_section()
	write_ini_file(config_filename)

if __name__ == "__main__":
	main()