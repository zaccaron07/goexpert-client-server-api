CREATE TABLE IF NOT EXISTS exchange_rate (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    codein TEXT NOT NULL,
    name TEXT NOT NULL,
    high TEXT NOT NULL,
    low TEXT NOT NULL,
    var_bid TEXT NOT NULL,
    pct_change TEXT NOT NULL,
    bid TEXT NOT NULL,
    ask TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    create_date TEXT NOT NULL
);