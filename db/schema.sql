-- Latest Migration: 20241027170349
CREATE TABLE goose_db_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TIMESTAMP DEFAULT (datetime('now'))
	);
CREATE TABLE `users` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`email` text NOT NULL,`password_hash` text NOT NULL,CONSTRAINT `uni_users_email` UNIQUE (`email`));
CREATE INDEX `idx_users_deleted_at` ON `users`(`deleted_at`);
CREATE TABLE `sessions` (`token` varchar(43),`data` blob,`expiry` datetime,PRIMARY KEY (`token`));
CREATE INDEX `idx_sessions_expiry` ON `sessions`(`expiry`);