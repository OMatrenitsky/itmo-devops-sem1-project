#!/bin/bash
#!/bin/bash
set -e

go mod tidy

psql -h localhost -U validator -d project-sem-1 <<EOF
CREATE TABLE IF NOT EXISTS prices (
    id INTEGER,
    created_at DATE,
    name TEXT,
    category TEXT,
    price NUMERIC
);
EOF