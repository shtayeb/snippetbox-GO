-- The Go tool ignores any directories called testdata, 
-- so these scripts will be ignored when compiling your application 
-- (it also ignores any directories or files which have names that begin with an _ or . character too)

CREATE TABLE snippets (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created timestamp NOT NULL,
	expires timestamp NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
	id SERIAL PRIMARY KEY ,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	hashed_password CHAR(60) NOT NULL,
	created timestamp NOT NULL
);


INSERT INTO users (name, email, hashed_password, created) VALUES (
	'Alice Jones',
	'alice@example.com',
	'$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
	'2022-01-01 10:00:00'
);
