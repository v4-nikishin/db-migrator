-- +db-migrator Up
CREATE TABLE test_migration (
    id              serial primary key,
    name            text,
    date            text
);
INSERT INTO test_migration (name, date) 
VALUES ('test_name', '2023-03-12');

-- +db-migrator Down
DROP TABLE test_migration;
