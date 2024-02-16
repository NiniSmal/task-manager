CREATE TABLE tasks(
                      id BIGSERIAL,
                      name TEXT NOT NULL,
                      status TEXT DEFAULT 'not done'
);