CREATE TABLE IF NOT EXISTS `users`
(
    id   int,
    username text,
    password text,
    wallets text,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `crypto_wallets`
(
    username text,
    name text,
    amount bigint
);

CREATE TABLE IF NOT EXISTS `start_stop_checks`
(
    username text,
    name text,
    stop BIT,
    start BIT
);


INSERT INTO users (id, username, password, wallets)
VALUES (2, "two", "twopw", "tone "),
       (3, "three", "threepw", "thone thtwo "),
       (4, "four", "fourpw", "fone ftwo fthree "),
       (6, "six", "sixpw", "sone ");

INSERT INTO crypto_wallets(username, name, amount)
VALUES ("two", "tone", 0),
       ("three", "thone", 0),
       ("three", "thtwo", 0),
       ("four", "fone", 0),
       ("four", "ftwo", 0),
       ("four", "fthree", 0),
       ("six", "sone", 0);

INSERT INTO start_stop_checks(username, name, stop, start)
VALUES ("two", "tone", 0,0),
       ("three", "thone", 0, 0),
       ("three", "thtwo", 0, 0),
       ("four", "fone", 0,0),
       ("four", "ftwo", 0,0),
       ("four", "fthree", 0,0),
       ("six", "sone", 0,0);





-- CREATE TABLE IF NOT EXISTS users 
-- (   
--     id   int,
--     username VARCHAR NOT NULL,
--     password VARCHAR NOT NULL,
--     wallets VARCHAR NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS crypto_wallets
-- (
--     username VARCHAR NOT NULL,
--     name VARCHAR NOT NULL,
--     amount BIGINT NOT NULL,
--     PRIMARY KEY (name)
-- );

-- CREATE TABLE IF NOT EXISTS start_stop_checks
-- (
--     username text,
--     name text,
--     stop BIT,
--     start BIT
-- );

-- INSERT INTO users (id, username, password, wallets)
-- VALUES (2, "two", "twopw", "tone "),
--        (3, "three", "threepw", "thone thtwo "),
--        (4, "four", "fourpw", "fone ftwo fthree "),
--        (6, "six", "sixpw", "sone ");

-- INSERT INTO crypto_wallets(username, name, amount)
-- VALUES ("two", "tone", 0),
--        ("three", "thone", 0),
--        ("three", "thtwo", 0),
--        ("four", "fone", 0),
--        ("four", "ftwo", 0),
--        ("four", "fthree", 0),
--        ("six", "sone", 0);

-- INSERT INTO start_stop_checks(username, name, stop, start)
-- VALUES ("two", "tone", 0,0),
--        ("three", "thone", 0, 0),
--        ("three", "thtwo", 0, 0),
--        ("four", "fone", 0,0),
--        ("four", "ftwo", 0,0),
--        ("four", "fthree", 0,0),
--        ("six", "sone", 0,0);

