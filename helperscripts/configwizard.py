#!/usr/bin/python3

import os
import configparser
import inquirer
import re
import sys
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

def fill_general_section():
	global ini_data
	questions = [
		inquirer.Path('database_dir', message="Where is the database (existing full path)", path_type=inquirer.Path.DIRECTORY,),
		inquirer.Text('database_file', message="What is the name of the database (ending in .db)", validate=validate_file_extension_db_sqlite3),
		inquirer.Text('timeformat', message="enter a valid timeformat", default="02/Jan/2006:15:04:05 -0700"),
		inquirer.Text('mydomain', message="enter your top level domain (mydomain.com)"),
		inquirer.List("writelog", message="Do i need to write logfiles to the output directory?", choices=["yes", "no"], default="no"),
	]
	answers = inquirer.prompt(questions)
	
	database_dir = answers['database_dir']
	database_file = answers['database_file']
	timeformat = answers['timeformat']
	mydomain = answers['mydomain']
	writelog = answers['writelog']
	
	if writelog == "yes":
		writelog_output = 'true'
	else:
		writelog_output = 'false'

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
	
	write_ini_file(config_filename)
"""
config = configparser.ConfigParser()

# Add sections to the INI file
config.add_section('Section1')
config.add_section('Section2')

# Add options to the sections
config.set('Section1', 'Option1', 'Value1')
config.set('Section1', 'Option2', 'Value2')
config.set('Section2', 'Option1', 'Value3')
config.set('Section2', 'Option2', 'Value4')

# Write the INI file to disk
with open('example.ini', 'w') as configfile:
    config.write(configfile)
	"""
	
	
	
"""
	questions = [
		inquirer.Text('config_file', message="Enter the path and name for the config file", default='config.ini'),
		inquirer.Path('data_dir', message="Select the path to the data directory", path_type=inquirer.Path.DIRECTORY),
		inquirer.Path('output_dir', message="Select the path to the output directory", path_type=inquirer.Path.DIRECTORY),
		inquirer.Path('database_dir', message="Where is the database (full path)", path_type=inquirer.Path.FILE, validate=validate_file_extension),
		inquirer.Confirm('debug', message="Enable debugging?", default=False),
	]

	# Run the wizard and get the user's answers
	answers = inquirer.prompt(questions)

	# Verify that the directories and files exist
	if not os.path.isdir(answers['data_dir']):
		print("Error: The data directory does not exist.")
		exit()

	if not os.path.isdir(answers['output_dir']):
		print("Error: The output directory does not exist.")
		exit()

	# Write the config file
	config = configparser.ConfigParser()
	config['DEFAULT'] = {
		'data_dir': answers['data_dir'],
		'output_dir': answers['output_dir'],
		'debug': str(answers['debug']),
	}

	with open(answers['config_file'], 'w') as f:
		config.write(f)
"""
if __name__ == "__main__":
	main()