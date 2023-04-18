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
			#print("%s %s\n" % (config,value))
			config.set(section, option, value)
	with open(filename, 'w') as f:
		config.write(f)

def get_configfilename():
	questions = [
		inquirer.Path('configfile_dir', message="Where is the config file dir (existing full path)", path_type=inquirer.Path.DIRECTORY, default="/root/"),
		inquirer.Text('configfile_name', message="What is the name of the configfile (ending in .ini)", validate=validate_file_extension_ini, default="test.ini"),
	]
	answers = inquirer.prompt(questions)
	
	configfile_dir = answers['configfile_dir']
	configfile_name = answers['configfile_name']
	return configfile_dir+configfile_name

def get_default_log_path():
	distrodb = {
		'win': {
			'list': ['Windows'],
			'logpath': 'C:\\Program Files\\Apache\\Apache2.4\\logs\\',
			'outputpath': 'C:\\Program Files\\Apache\\Apache2.4\\www\\',
		},
		'apache2': {
			'list': ['Darwin', 'Ubuntu', 'SUSE', 'OpenSUSE', 'Debian GNU/Linux', 'debian'],
			'logpath': '/var/log/apache2/',
			'outputpath': '/var/www/html/',
		},
		'httpd': {
			'list': ['Fedora', 'CentOS', 'RHEL', 'Red Hat Enterprise Linux (RHEL)'],
			'logpath': '/var/log/apache2/',
			'outputpath': '/var/www/html/',
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
			return [value['logpath'], value['outputpath']]
	
	return ["", ""]
	
def fill_input_section():
	osdefault = get_default_log_path()
	global ini_data
	domain = ini_data['general']['mydomain']
	pattern = pattern = r'^' + re.escape(domain) + r'.*'
	questions = [
		inquirer.Path('logfile_dir', message="Where can i find the logfiles?", path_type=inquirer.Path.DIRECTORY, default=osdefault[0]),
		inquirer.Text('logfileregex', message="Regex which logiles i need to parse", default=pattern),
		inquirer.Text('parseregex', message="regex to match log format values (or clf)", default='clf'),
		inquirer.List("fullloadcheck", message="Do i need to parse every line of every logfile? (false = only lines newer than last load)", choices=["true", "false"], default="false"),
		inquirer.Text('order', message="the order in which the regex will find the necessary fields", default='123456789'),
	]
	
	answers = inquirer.prompt(questions)
	logfile_dir = answers['logfile_dir']
	logfileregex = answers['logfileregex']
	parseregex = answers['parseregex']
	fullloadcheck = answers['fullloadcheck']
	order = answers['order']
	ini_data['input'] = {
        'logfilepath': logfile_dir,
		'logfileregex': logfileregex,
		'parseregex': parseregex,
		'fullloadcheck': fullloadcheck,
		'parserfield_ip': order[0],
		'parserfield_datetime': order[1],
		'parserfield_method': order[2],
		'parserfield_request': order[3],
		'parserfield_httpversion': order[4],
		'parserfield_returncode': order[5],
		'parserfield_httpsize': order[6],
		'parserfield_referrer': order[7],
		'parserfield_useragent': order[8],
    }
	
def fill_output_section():
	osdefault = get_default_log_path()
	global ini_data
	questions = [
		inquirer.Path('outputpath', message="Where do i need to write the output to?", path_type=inquirer.Path.DIRECTORY, default=osdefault[1]),
		inquirer.List("emptyoutputpath", message="Do i need to remove all html files from the output before creating new ones?", choices=["true", "false"], default="true"),
		inquirer.Text('number_of_days_detailed', message="how many days to return detailed info for (in tables)", default='31'),
		inquirer.Text('max_number_of_days', message="how many days to return as a fallback", default='124'),
		inquirer.Text('assethost', message="where do i find go-echarts javascript files?", default='https://go-echarts.github.io/go-echarts-assets/assets/'),
		inquirer.List("zipoutput", message="Do i need to create a zipfile with the output?", choices=["true", "false"], default="false"),
		inquirer.Text("zippath", message="If i create a zipfile with the output, can you give me the full path?", default="./output.zip"),
		inquirer.Text("numberofreferrers", message="limit number of refferers in table", default="31"),
	]
	
	answers = inquirer.prompt(questions)
	outputpath = answers['outputpath']
	emptyoutputpath = answers['emptyoutputpath']
	number_of_days_detailed = answers['number_of_days_detailed']
	max_number_of_days = answers['max_number_of_days']
	assethost = answers['assethost']
	zipoutput = answers['zipoutput']
	zippath = answers['zippath']
	numberofreferrers = answers['numberofreferrers']
	
	ini_data['output'] = {
        'outputpath': outputpath,
        'emptyoutputpath': emptyoutputpath,
        'number_of_days_detailed': number_of_days_detailed,
        'max_number_of_days': max_number_of_days,
        'assethost': assethost,
        'zipoutput': zipoutput,
        'zippath': zippath,
        'numberofreferrers': numberofreferrers,
    }

def fill_ignorevisitorips_section():
	global ini_data
	ini_data['ignorevisitorips'] = {}
	answers = {}
	print("add as many ip's you want to ignore while parsing as you want..\n")
	print("when finished, leave the prompts blank!\n")
	while True:
		questions = [
			inquirer.Text('key', message="reference"),
			inquirer.Text('value', message="ip to ignore")
		]
		response = inquirer.prompt(questions)

		if not response['key'] and not response['value']:
			break

		answers[response['key']] = response['value']
		for key, value in answers.items():
			ini_data['ignorevisitorips'][key] = [value][0]

def fill_ignorehostagents_section():
	global ini_data
	ini_data['ignorehostagents'] = {}
	answers = {}
	print("add as many hostagentss you want to ignore while parsing as you want..\n")
	print("globs welcome!!!..\n")
	print("when finished, leave the prompts blank!\n")
	while True:
		questions = [
			inquirer.Text('key', message="reference"),
			inquirer.Text('value', message="hostagent to ignore")
		]
		response = inquirer.prompt(questions)

		if not response['key'] and not response['value']:
			break

		answers[response['key']] = response['value']
		for key, value in answers.items():
			ini_data['ignorehostagents'][key] = [value][0]

def fill_ignorereferrers_section():
	global ini_data
	ini_data['ignorereferrers'] = {}
	answers = {}
	print("add as many referrers you want to ignore while parsing as you want..\n")
	print("globs welcome!!!..\n")
	print("when finished, leave the prompts blank!\n")
	while True:
		questions = [
			inquirer.Text('key', message="reference"),
			inquirer.Text('value', message="referrer to ignore")
		]
		response = inquirer.prompt(questions)

		if not response['key'] and not response['value']:
			break

		answers[response['key']] = response['value']
		for key, value in answers.items():
			ini_data['ignorereferrers'][key] = [value][0]

def fill_ignoredrequests_section():
	global ini_data
	ini_data['ignoredrequests'] = {}
	answers = {}
	print("add as many requests you want to ignore while parsing as you want..\n")
	print("globs welcome!!!..\n")
	print("when finished, leave the prompts blank!\n")
	while True:
		questions = [
			inquirer.Text('key', message="reference"),
			inquirer.Text('value', message="request to ignore")
		]
		response = inquirer.prompt(questions)

		if not response['key'] and not response['value']:
			break

		answers[response['key']] = response['value']
		for key, value in answers.items():
			ini_data['ignoredrequests'][key] = [value][0]

def fill_general_section():
	global ini_data
	questions = [
		inquirer.Path('database_dir', message="Where is the database (existing full path)", path_type=inquirer.Path.DIRECTORY, default="/root/"),
		inquirer.Text('database_file', message="What is the name of the database (ending in .db)", validate=validate_file_extension_db_sqlite3, default="db.db"),
		inquirer.Text('timeformat', message="enter a valid timeformat", default="02/Jan/2006:15:04:05 -0700"),
		inquirer.Text('mydomain', message="enter your top level domain (mydomain.com)", default="merel.mobi"),
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

def fill_allstats_sections():
	global ini_data
	statdb = {
		'conf_stat_raw_PerHour_hits': {
			'enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enabled?'},
			'table_enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable table output?'},
			'table_title': {'type': 'string', 'validate': '', 'default': 'Number of raw hits per hour over th last |number_of_days_detailed| days', 'helpmessage': 'table title'},
			'table_description' : {'type': 'string', 'validate': '', 'default': 'Count of all raw sucessfull hits per hour.', 'helpmessage': 'table description'},
			'table_pagecontent' : {'type': 'string', 'validate': '', 'default': 'We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you\'ll get stats untill 15:34 x days ago).', 'helpmessage': 'table pagecontent'},
			'table_pagefooter' : {'type': 'string', 'validate': '', 'default': 'only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they\'ll never be shown!', 'helpmessage': 'table pagefooter'},
			'table_filename' : {'type': 'string', 'validate': '', 'default': 'raw_PerHour_hits_table.html', 'helpmessage': 'table filename'},
			'table_index_name' : {'type': 'string', 'validate': '', 'default': 'table raw hits per hour last |number_of_days_detailed| days', 'helpmessage': 'table indexfile linktext'},
			'table_index_group' : {'type': 'string', 'validate': '', 'default': 'raw hits', 'helpmessage': 'table indexfile group'},
			'table_index_order' : {'type': 'int', 'validate': '', 'default': '1', 'helpmessage': 'table indexfile order'},
			'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable linegrapht output?'},
			'linegraph_title' : {'type': 'string', 'validate': '', 'default': 'linegraph raw hits per hour last |number_of_days_detailed| days', 'helpmessage': 'linegraph title'},
			'linegraph_description' : {'type': 'string', 'validate': '', 'default': 'Count of all raw succesfull hits per hour.', 'helpmessage': 'linegraph description'},
			'linegraph_filename' : {'type': 'string', 'validate': '', 'default': 'raw_PerHour_hits_linegraph.html', 'helpmessage': 'linegraph filename'},
			'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': 'raw hits', 'helpmessage': 'linegraph indexfile group'},
			'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '15', 'helpmessage': 'linegraph indexfile order'},
			'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': 'false', 'helpmessage': 'enable comparable linegraphs?'},
			'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs title'},
			'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs description'},
			'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs filename'},
			'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage group'},
			'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage order'},
			'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs parameters'},
		},
		'conf_stat_raw_PerDay_hits': {
			'enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enabled?'},
			'table_enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable table output?'},
			'table_title': {'type': 'string', 'validate': '', 'default': 'Number of raw hits per day over th last |number_of_days_detailed| days', 'helpmessage': 'table title'},
			'table_description' : {'type': 'string', 'validate': '', 'default': 'Count of all raw sucessfull hits per day.', 'helpmessage': 'table description'},
			'table_pagecontent' : {'type': 'string', 'validate': '', 'default': 'We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you\'ll get stats untill 15:34 x days ago).', 'helpmessage': 'table pagecontent'},
			'table_pagefooter' : {'type': 'string', 'validate': '', 'default': 'only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they\'ll never be shown!', 'helpmessage': 'table pagefooter'},
			'table_filename' : {'type': 'string', 'validate': '', 'default': 'raw_PerDay_hits_table.html', 'helpmessage': 'table filename'},
			'table_index_name' : {'type': 'string', 'validate': '', 'default': 'table raw hits per day last |number_of_days_detailed| days', 'helpmessage': 'table indexfile linktext'},
			'table_index_group' : {'type': 'string', 'validate': '', 'default': 'raw hits', 'helpmessage': 'table indexfile group'},
			'table_index_order' : {'type': 'int', 'validate': '', 'default': '2', 'helpmessage': 'table indexfile order'},
			'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable linegrapht output?'},
			'linegraph_title' : {'type': 'string', 'validate': '', 'default': 'linegraph raw hits per day last |number_of_days_detailed| days', 'helpmessage': 'linegraph title'},
			'linegraph_description' : {'type': 'string', 'validate': '', 'default': 'Count of all raw succesfull hits per day.', 'helpmessage': 'linegraph description'},
			'linegraph_filename' : {'type': 'string', 'validate': '', 'default': 'raw_PerDay_hits_linegraph.html', 'helpmessage': 'linegraph filename'},
			'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': 'raw hits', 'helpmessage': 'linegraph indexfile group'},
			'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '3', 'helpmessage': 'linegraph indexfile order'},
			'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable comparable linegraphs?'},
			'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': 'linegraph raw hits per day, week by week last 4 weeks&linegraph raw hits per day, month by month last 4 months', 'helpmessage': 'comparable linegraphs title'},
			'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '4 week comparison of all raw succesfull hits.&4 mont comparison of all raw succesfull hits.', 'helpmessage': 'comparable linegraphs description'},
			'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': 'raw_PerDay_hits_4WeeksLinegraph.html&raw_PerDay_hits_4MonthsLinegraph.html', 'helpmessage': 'comparable linegraphs filename'},
			'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': 'raw hits&raw hits', 'helpmessage': 'comparable linegraphs indexpage group'},
			'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '4&5', 'helpmessage': 'comparable linegraphs indexpage order'},
			'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': 'day,4,8,week&day,4,32,month', 'helpmessage': 'comparable linegraphs parameters'},
		},
		'conf_stat_unique_PerHour_hits': {
			'enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enabled?'},
			'table_enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable table output?'},
			'table_title': {'type': 'string', 'validate': '', 'default': 'Number of unique hits per hour over th last |number_of_days_detailed| days', 'helpmessage': 'table title'},
			'table_description' : {'type': 'string', 'validate': '', 'default': 'Count of all unique sucessfull hits per hour.', 'helpmessage': 'table description'},
			'table_pagecontent' : {'type': 'string', 'validate': '', 'default': 'A unique hit = the first time this hour one use generates a hit.|We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you\'ll get stats untill 15:34 x days ago).', 'helpmessage': 'table pagecontent'},
			'table_pagefooter' : {'type': 'string', 'validate': '', 'default': 'only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they\'ll never be shown!', 'helpmessage': 'table pagefooter'},
			'table_filename' : {'type': 'string', 'validate': '', 'default': 'unique_PerHour_hits_table.html', 'helpmessage': 'table filename'},
			'table_index_name' : {'type': 'string', 'validate': '', 'default': 'table unique hits per hour last |number_of_days_detailed| days', 'helpmessage': 'table indexfile linktext'},
			'table_index_group' : {'type': 'string', 'validate': '', 'default': 'unique hits', 'helpmessage': 'table indexfile group'},
			'table_index_order' : {'type': 'int', 'validate': '', 'default': '1', 'helpmessage': 'table indexfile order'},
			'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable linegrapht output?'},
			'linegraph_title' : {'type': 'string', 'validate': '', 'default': 'linegraph unique hits per hour last |number_of_days_detailed| days', 'helpmessage': 'linegraph title'},
			'linegraph_description' : {'type': 'string', 'validate': '', 'default': 'Count of all unique succesfull hits per hour.', 'helpmessage': 'linegraph description'},
			'linegraph_filename' : {'type': 'string', 'validate': '', 'default': 'unique_PerHour_hits_linegraph.html', 'helpmessage': 'linegraph filename'},
			'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': 'unique hits', 'helpmessage': 'linegraph indexfile group'},
			'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '4', 'helpmessage': 'linegraph indexfile order'},
			'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': 'false', 'helpmessage': 'enable comparable linegraphs?'},
			'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs title'},
			'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs description'},
			'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs filename'},
			'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage group'},
			'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage order'},
			'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs parameters'},
		},
		'conf_stat_unique_PerDay_hits': {
			'enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enabled?'},
			'table_enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable table output?'},
			'table_title': {'type': 'string', 'validate': '', 'default': 'Number of unique hits per day over th last |number_of_days_detailed| days', 'helpmessage': 'table title'},
			'table_description' : {'type': 'string', 'validate': '', 'default': 'Count of all unique sucessfull hits per day.', 'helpmessage': 'table description'},
			'table_pagecontent' : {'type': 'string', 'validate': '', 'default': 'Unique hit = the first hit of one user this day|We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you\'ll get stats untill 15:34 x days ago).', 'helpmessage': 'table pagecontent'},
			'table_pagefooter' : {'type': 'string', 'validate': '', 'default': 'only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they\'ll never be shown!', 'helpmessage': 'table pagefooter'},
			'table_filename' : {'type': 'string', 'validate': '', 'default': 'unique_PerDay_hits_table.html', 'helpmessage': 'table filename'},
			'table_index_name' : {'type': 'string', 'validate': '', 'default': 'table unique hits per day last |number_of_days_detailed| days', 'helpmessage': 'table indexfile linktext'},
			'table_index_group' : {'type': 'string', 'validate': '', 'default': 'unique hits', 'helpmessage': 'table indexfile group'},
			'table_index_order' : {'type': 'int', 'validate': '', 'default': '3', 'helpmessage': 'table indexfile order'},
			'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable linegrapht output?'},
			'linegraph_title' : {'type': 'string', 'validate': '', 'default': 'linegraph unique hits per day last |number_of_days_detailed| days', 'helpmessage': 'linegraph title'},
			'linegraph_description' : {'type': 'string', 'validate': '', 'default': 'Count of all unique succesfull hits per day.', 'helpmessage': 'linegraph description'},
			'linegraph_filename' : {'type': 'string', 'validate': '', 'default': 'unique_PerDay_hits_linegraph.html', 'helpmessage': 'linegraph filename'},
			'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': 'unique hits', 'helpmessage': 'linegraph indexfile group'},
			'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '3', 'helpmessage': 'linegraph indexfile order'},
			'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable comparable linegraphs?'},
			'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': 'linegraph unique hits per day, week by week last 4 weeks&linegraph unique hits per day, month by month last 4 months', 'helpmessage': 'comparable linegraphs title'},
			'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '4 week comparison of all unique succesfull hits.&4 month comparison of all unique succesfull hits.', 'helpmessage': 'comparable linegraphs description'},
			'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': 'unique_PerHour_hits_4WeeksLinegraph.html&unique_PerHour_hits_4MonthsLinegraph.html', 'helpmessage': 'comparable linegraphs filename'},
			'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': 'unique hits&unique hits', 'helpmessage': 'comparable linegraphs indexpage group'},
			'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '5&6', 'helpmessage': 'comparable linegraphs indexpage order'},
			'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': 'day,4,8,week&day,4,32,month', 'helpmessage': 'comparable linegraphs parameters'},
		},
		'conf_stat_raw_PerHour_ReferringUrls': {
			'enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enabled?'},
			'table_enabled': {'type': 'bool', 'validate': '', 'default': 'true', 'helpmessage': 'enable table output?'},
			'table_title': {'type': 'string', 'validate': '', 'default': 'Number of raw referrers per hour over the last |number_of_days_detailed| days', 'helpmessage': 'table title'},
			'table_description' : {'type': 'string', 'validate': '', 'default': 'Count of all raw succesfull referrers', 'helpmessage': 'table description'},
			'table_pagecontent' : {'type': 'string', 'validate': '', 'default': 'We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you\'ll get stats untill 15:34 x days ago).', 'helpmessage': 'table pagecontent'},
			'table_pagefooter' : {'type': 'string', 'validate': '', 'default': 'only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they\'ll never be shown!', 'helpmessage': 'table pagefooter'},
			'table_filename' : {'type': 'string', 'validate': '', 'default': 'raw_PerHour_ReferringUrls_table.html', 'helpmessage': 'table filename'},
			'table_index_name' : {'type': 'string', 'validate': '', 'default': 'table raw referrers per hour over last |number_of_days_detailed| days', 'helpmessage': 'table indexfile linktext'},
			'table_index_group' : {'type': 'string', 'validate': '', 'default': 'referrers', 'helpmessage': 'table indexfile group'},
			'table_index_order' : {'type': 'int', 'validate': '', 'default': '2', 'helpmessage': 'table indexfile order'},
			'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': 'false', 'helpmessage': 'enable linegrapht output?'},
			'linegraph_title' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph title'},
			'linegraph_description' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph description'},
			'linegraph_filename' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph filename'},
			'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph indexfile group'},
			'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '', 'helpmessage': 'linegraph indexfile order'},
			'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': 'false', 'helpmessage': 'enable comparable linegraphs?'},
			'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs title'},
			'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs description'},
			'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs filename'},
			'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage group'},
			'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage order'},
			'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs parameters'},
		}
		# 'conf_stat_raw_PerHour_hits': {
			# 'enabled': {'type': 'bool', 'validate': '', 'default': '', 'helpmessage': 'enabled?'},
			# 'table_enabled': {'type': 'bool', 'validate': '', 'default': '', 'helpmessage': 'enable table output?'},
			# 'table_title': {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table title'},
			# 'table_description' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table description'},
			# 'table_pagecontent' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table pagecontent'},
			# 'table_pagefooter' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table pagefooter'},
			# 'table_filename' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table filename'},
			# 'table_index_name' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table indexfile linktext'},
			# 'table_index_group' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'table indexfile group'},
			# 'table_index_order' : {'type': 'int', 'validate': '', 'default': '', 'helpmessage': 'table indexfile order'},
			# 'linegraph_enabled' : {'type': 'bool', 'validate': '', 'default': '', 'helpmessage': 'enable linegrapht output?'},
			# 'linegraph_title' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph title'},
			# 'linegraph_description' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph description'},
			# 'linegraph_filename' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph filename'},
			# 'linegraph_index_group' : {'type': 'string', 'validate': '', 'default': '', 'helpmessage': 'linegraph indexfile group'},
			# 'linegraph_index_order' : {'type': 'int', 'validate': '', 'default': '', 'helpmessage': 'linegraph indexfile order'},
			# 'linegraph_compare_x_days_weeks_months_enabled' : {'type': 'bool', 'validate': '', 'default': '', 'helpmessage': 'enable comparable linegraphs?'},
			# 'linegraph_compare_x_days_weeks_months_title' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs title'},
			# 'linegraph_compare_x_days_weeks_months_description' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs description'},
			# 'linegraph_compare_x_days_weeks_months_filename' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs filename'},
			# 'linegraph_compare_x_days_weeks_months_index_group' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage group'},
			# 'linegraph_compare_x_days_weeks_months_index_order' : {'type': 'list_int', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs indexpage order'},
			# 'linegraph_compare_x_days_weeks_months_parameters' : {'type': 'list_string', 'validate': '', 'default': '', 'helpmessage': 'comparable linegraphs parameters'},
		# }
	}
	for section, sectiondb in statdb.items():
		print(f"section: {section}")
		questions = []
		for key, valuedb in sectiondb.items():
			if valuedb['type'] == 'bool':
				questions.append(inquirer.List(key, message=valuedb['helpmessage'], choices=["true", "false"], default=valuedb['default']))
			elif valuedb['type'] in ['string', 'int', 'list_string', 'list_int']:
				questions.append(inquirer.Text(key, message=valuedb['helpmessage'], default=valuedb['default']))
			else:
				questions.append(inquirer.Text(key, message=valuedb['helpmessage'], default=valuedb['default']))
		#print("%s" % questions)
		answers = inquirer.prompt(questions)
		ini_data[section] = {}
			

def main():
	config_filename = get_configfilename()
	fill_general_section()
	fill_input_section()
	fill_output_section()
	fill_ignorevisitorips_section()
	fill_ignorehostagents_section()
	fill_ignorereferrers_section()
	fill_ignoredrequests_section()
	fill_allstats_sections()
	write_ini_file(config_filename)

if __name__ == "__main__":
	main()