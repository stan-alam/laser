package postgres

// @note: separating all queries into a separate file gives us easier control
//        in tracking down and modifying queries.
//
// @note: we may want to connect the POC on a service to our Users in a future
//        iteration.  This makes more sense if we assume all POC's are also
//        users.
//
// @note: future iterations of tokens table may include scope, matching how
//        something like github allows you to associate select permissions and
//        have multiple tokens registered, and then we could adjust the login
//        to look for a specfic type of token, allowing multiple tokens to
//        exist per user in the table.
//
//        we could also restrict tokens used for login operations by adding a
//        unique key clause across the new token type field.
const (
	SQLCreateExtensionPGCrypto string = `CREATE EXTENSION IF NOT EXISTS pgcrypto;`

	SQLCreateTableUsers = `CREATE TABLE IF NOT EXISTS users (
		name text primary key not null,
		email text not null,
		password text not null);`
	SQLCreateTableTokens = `CREATE TABLE IF NOT EXISTS tokens (
		id uuid PRIMARY KEY not null DEFAULT gen_random_uuid(),
		username text REFERENCES users (name),
		created timestamp DEFAULT transaction_timestamp());`
	SQLCreateTableServices = `create table if not exists services (
		name text primary key not null,
		technology text not null,
		poc text not null);`

	SQLUserLogin   = `SELECT * FROM users WHERE email = $1 AND password = crypt($2, password);`
	SQLUserInsert  = `INSERT INTO users (name, email, password) VALUES ($1, $2, crypt($3, gen_salt('bf')));`
	SQLUserSelect  = `SELECT * FROM users WHERE email = $1;`
	SQLUserDelete  = `DELETE FROM users WHERE name = $1;`
	SQLUserUpdate  = `UPDATE users SET email = $2, password = $3 WHERE name = $1;`
	SQLUsersSelect = `SELECT * FROM users;`

	SQLServiceInsert  = `INSERT INTO services (name, technology, poc) VALUES ($1, $2, $3);`
	SQLServiceSelect  = `SELECT * FROM services WHERE name = $1;`
	SQLServiceDelete  = `DELETE FROM services WHERE name = $1;`
	SQLServiceUpdate  = `UPDATE services SET technology = $2, poc = $3 WHERE name = $1;`
	SQLServicesSelect = `SELECT * FROM services;`

	SQLTokenInsert = `INSERT INTO tokens (username) VALUES ($1) RETURNING id;`
	SQLTokenSelect = `SELECT id, username, created FROM tokens WHERE id = $1;`
	SQLTokenDelete = `DELETE FROM tokens WHERE id = $1;`
)
