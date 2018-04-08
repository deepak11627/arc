# ARC
Implementation of Adaptive replacement cache algorithm in Go

To use the app, you must create a database first and run the required arc.sql queries.

Database is reset everytime the app is started.

``` go run main.go```

# Logging

Zap Logger is been used for logging purpose. While the application is running the logs can be viewed in `out.log` file.

# Database

Database can be set by passing a flag named `dsn` while running the application. If the dsn flag is 
not provided then the application runs db less.

``` go run main.go -dsn="root:root@tcp(127.0.0.1:3306)/arc"```

A single table maintains the list,

```
-- Table for maintaining the Ghost lists in DB

CREATE TABLE IF NOT EXISTS `ghost_lists` (
  `list_id` varchar(2) NOT NULL,
  `ghost_key` varchar(16)  NOT NULL,
  `ghost_value` varchar(16)  NOT NULL,
  PRIMARY KEY (`list_id`, `ghost_key`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

```