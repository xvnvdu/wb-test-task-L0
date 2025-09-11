CREATE DATABASE orders_db;
CREATE USER orders_user WITH PASSWORD '12345' ;

GRANT ALL PRIVILEGES ON DATABASE orders_db TO orders_user ;
\c orders_db;

GRANT ALL PRIVILEGES ON SCHEMA public TO orders_user;

ALTER DATABASE orders_db OWNER TO orders_user;
ALTER SCHEMA public OWNER TO orders_user;
