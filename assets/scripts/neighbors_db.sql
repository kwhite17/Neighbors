CREATE TABLE IF NOT EXISTS shelters (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    City VARCHAR(100) NOT NULL,
    Country VARCHAR(100) NOT NULL,
    Name VARCHAR(100) NOT NULL,
    PostalCode VARCHAR(100) NOT NULL,
    State VARCHAR(100) NOT NULL,
    Street VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Category VARCHAR(100) NOT NULL,
    Gender VARCHAR(100) NOT NULL,
    Quantity TINYINT NOT NULL,
    Size VARCHAR(100) NOT NULL,
    Status VARCHAR(100) NOT NULL,
    ShelterID VARCHAR(100) NOT NULL
);
