#!/usr/bin/python3
import sqlite3



def main():
	with sqlite3.connect('../../new_db.sqlite') as conn:
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

	with sqlite3.connect('../../db_old.db') as old_conn, sqlite3.connect('../../new_db.sqlite') as new_conn:
		old_cursor = old_conn.cursor()
		new_cursor = new_conn.cursor()

		old_cursor.execute('SELECT DISTINCT ip, useragent FROM user')
		ip_useragent_rows = old_cursor.fetchall()

		# insert unique ip values into new user_ip table
		for row in ip_useragent_rows:
			ip = row[0]
			new_cursor.execute('INSERT OR IGNORE INTO user_ip (ip) VALUES (?)', (ip,))
			
		 # insert unique useragent values into new user_useragent table
		for row in ip_useragent_rows:
			useragent = row[1]
			new_cursor.execute('INSERT OR IGNORE INTO user_useragent (useragent) VALUES (?)', (useragent,))
		"""
		# copy data from old tables to new tables
		old_cursor.execute('SELECT * FROM request')
		new_cursor.executemany('INSERT INTO request (id, request) VALUES (?, ?)', old_cursor.fetchall())

		old_cursor.execute('SELECT * FROM referrer')
		new_cursor.executemany('INSERT INTO referrer (id, referrer) VALUES (?, ?)', old_cursor.fetchall())

		old_cursor.execute('SELECT * FROM alreadyloaded')
		new_cursor.executemany('INSERT INTO alreadyloaded (id, hash) VALUES (?, ?)', old_cursor.fetchall())

		old_cursor.execute('SELECT * FROM visit')
		for row in old_cursor.fetchall():
			# get the new user_id by querying for the ip and useragent in the new tables
			ip, useragent = row[3], row[4]
			new_cursor.execute('SELECT id FROM user_ip WHERE ip = ?', (ip,))
			ip_id = new_cursor.fetchone()[0]
			new_cursor.execute('SELECT id FROM user_useragent WHERE useragent = ?', (useragent,))
			useragent_id = new_cursor.fetchone()[0]
			# insert the new visit row using the new user_id
			new_cursor.execute('INSERT INTO visit (id, referrer, request, visit_timestamp, user, statuscode, httpsize) VALUES (?, ?, ?, ?, ?, ?, ?)', (row[0], row[1], row[2], row[5], useragent_id, row[6], row[7]))
 """
		# commit changes and close connections
		new_conn.commit()

if __name__ == "__main__":
	main()