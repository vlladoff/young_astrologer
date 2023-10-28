CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    image BYTEA NOT NULL,
    hd_image BYTEA NOT NULL
);

CREATE TABLE data (
    id SERIAL PRIMARY KEY,
    datetime DATE NOT NULL,
    title VARCHAR NOT NULL,
    explanation VARCHAR NOT NULL,
    url VARCHAR UNIQUE NOT NULL,
    hd_url VARCHAR UNIQUE NOT NULL,
    images_id INT,
    FOREIGN KEY (images_id) REFERENCES images(id)
);
CREATE INDEX data_url ON data(url);
CREATE INDEX data_hd_url ON data(hd_url);