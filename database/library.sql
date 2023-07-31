DROP TABLE IF EXISTS book;
DROP TABLE IF EXISTS author;
CREATE TABLE IF NOT EXISTS author
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS book
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    title     VARCHAR(255) NOT NULL,
    author_id INT,
    FOREIGN KEY (author_id) REFERENCES author (id),
    UNIQUE KEY unique_id (title, author_id)
);


INSERT INTO author (name)
VALUES ('Eric Freeman'),
       ('Elizabeth Freeman'),
       ('Robert C. Martin');


INSERT INTO book (title, author_id)
VALUES ('Head First Design Patterns', (SELECT id FROM author WHERE author.name = 'Eric Freeman')),
       ('Head First Design Patterns', (SELECT id FROM author WHERE author.name = 'Elizabeth Freeman')),
       ('Clean Code', (SELECT id FROM author WHERE author.name = 'Robert C. Martin')),
       ('Head First JavaScript Programming: A Brain-Friendly Guide', (SELECT id FROM author WHERE author.name = 'Eric Freeman'));
