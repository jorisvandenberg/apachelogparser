#!/usr/bin/python3
import sqlite3
import argparse
import os


def main():
	parser = argparse.ArgumentParser(description='Description of your script')
	parser.add_argument('--inputdb', type=str, help='Path to the input file', required=True)
	parser.add_argument('--outputdb', type=str, help='Path to the output file', required=True)
	args = parser.parse_args()
	# Check if input file exists
	if not os.path.isfile(args.inputdb):
		print("Error: Input file does not exist!")
		exit()
	# Check if output file does not exist
	if os.path.isfile(args.outputdb):
		print("Error: Output file already exists!")
		exit()
	with sqlite3.connect(args.outputdb) as conn:
		cursor = conn.cursor()
		cursor.execute('''CREATE TABLE `request` (`id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, `request` TEXT NOT NULL )''')
		cursor.execute(''' CREATE TABLE `referrer` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`referrer` TEXT NOT NULL) ''')
		cursor.execute(''' CREATE TABLE `alreadyloaded` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, `hash` TEXT NOT NULL) ''')
		cursor.execute(''' CREATE TABLE `visit` (`id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, `referrer` INTEGER NOT NULL, `request` INTEGER NOT NULL, `visit_timestamp` INTEGER NOT NULL, `user` INTEGER NOT NULL, `statuscode` INTEGER, `httpsize` INTEGER,  FOREIGN KEY(`request`) REFERENCES `request`(`id`), FOREIGN KEY(`referrer`) REFERENCES `referrer`(`id`), FOREIGN KEY(`user`) REFERENCES `user`(`id`) ) ''')
		cursor.execute(''' CREATE TABLE IF NOT EXISTS "user_ip" ( "id"    INTEGER NOT NULL, "ip"    TEXT NOT NULL, PRIMARY KEY("id" AUTOINCREMENT)) ''')
		cursor.execute(''' CREATE TABLE IF NOT EXISTS "user_useragent" ( "id"    INTEGER NOT NULL, "useragent"     TEXT NOT NULL, PRIMARY KEY("id" AUTOINCREMENT) ) ''')
		cursor.execute('''  CREATE TABLE IF NOT EXISTS "user" ( "id"    INTEGER NOT NULL, "ip"    INTEGER NOT NULL, "useragent"     INTEGER NOT NULL, FOREIGN KEY("useragent") REFERENCES "user_useragent"("id"),  FOREIGN KEY("ip") REFERENCES "user_ip"("id"),  PRIMARY KEY("id" AUTOINCREMENT) ) ''')
		cursor.execute('''  CREATE TABLE IF NOT EXISTS "info" ( "key"   INTEGER, "value" TEXT ) ''')
		cursor.execute('''  CREATE INDEX request_request on request(request) ''')
		cursor.execute(''' CREATE INDEX referrer_referrer on referrer(referrer) ''')
		cursor.execute(''' CREATE UNIQUE INDEX "user_agent" ON "user_useragent" ( "useragent" ) ''')
		cursor.execute(''' CREATE UNIQUE INDEX "ip_ip" ON "user_ip" ( "ip"    ASC ) ''')
		conn.commit()

	with sqlite3.connect(args.inputdb) as old_conn, sqlite3.connect(args.outputdb) as new_conn:
		old_cursor = old_conn.cursor()
		new_cursor = new_conn.cursor()

		old_cursor.execute('SELECT DISTINCT ip FROM user')
		ip_useragent_rows = old_cursor.fetchall()

		# insert unique ip values into new user_ip table
		for row in ip_useragent_rows:
			ip = row[0]
			new_cursor.execute('INSERT OR IGNORE INTO user_ip (ip) VALUES (?)', (ip,))
		new_conn.commit()
		old_cursor.execute('SELECT DISTINCT useragent FROM user')
		ip_useragent_rows = old_cursor.fetchall()

		# insert unique ip values into new user_ip table
		for row in ip_useragent_rows:
			useragent = row[0]
			new_cursor.execute('INSERT OR IGNORE INTO user_useragent (useragent) VALUES (?)', (useragent,))
		new_conn.commit()
		old_cursor.execute('SELECT id, ip, useragent FROM user')
		for row in old_cursor.fetchall():
			user_id = row[0]
			user_ip = row[1]
			user_useragent = row[2]
			subquery = f"SELECT id FROM user_ip WHERE ip='{user_ip}'"
			new_cursor.execute(subquery)
			subquery_result = new_cursor.fetchall()
			first_row = subquery_result[0]
			ip_id = first_row[0]
			subquery = f"SELECT id FROM user_useragent WHERE useragent='{user_useragent}'"
			new_cursor.execute(subquery)
			subquery_result = new_cursor.fetchall()
			first_row = subquery_result[0]
			useragent_id = first_row[0]
			new_cursor.execute('INSERT OR IGNORE INTO user (id, ip,useragent) VALUES (?,?,?)', (user_id, ip_id, useragent_id))

		# copy data from old tables to new tables
		old_cursor.execute('SELECT id, request FROM request')
		new_cursor.executemany('INSERT INTO request (id, request) VALUES (?, ?)', old_cursor.fetchall())
		
		old_cursor.execute('SELECT id, referrer FROM referrer')
		new_cursor.executemany('INSERT INTO referrer (id, referrer) VALUES (?, ?)', old_cursor.fetchall())

		old_cursor.execute('SELECT id, hash FROM alreadyloaded')
		new_cursor.executemany('INSERT INTO alreadyloaded (id, hash) VALUES (?, ?)', old_cursor.fetchall())
		
		old_cursor.execute('SELECT id, referrer, request, visit_timestamp,user, statuscode, httpsize FROM visit')
		new_cursor.executemany('INSERT INTO visit (id, referrer, request, visit_timestamp,user, statuscode, httpsize) VALUES (?, ?,?,?,?,?,?)', old_cursor.fetchall())
		

		# commit changes and close connections
		new_conn.commit()

if __name__ == "__main__":
	main()