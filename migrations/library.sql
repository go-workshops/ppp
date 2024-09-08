CREATE TABLE authors1
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE books1
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID         NOT NULL,
    title     VARCHAR(100) NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors1 (id)
);

CREATE TABLE book_copies1
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID        NOT NULL,
    status  VARCHAR(20) NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books1 (id)
);

CREATE TABLE authors2
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE books2
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID         NOT NULL,
    title     VARCHAR(100) NOT NULL,
    FOREIGN KEY (author_id) REFERENCES authors2 (id)
);

CREATE TABLE book_copies2
(
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID        NOT NULL,
    status  VARCHAR(20) NOT NULL,
    FOREIGN KEY (book_id) REFERENCES books2 (id)
);
