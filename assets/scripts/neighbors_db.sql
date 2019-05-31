CREATE TABLE IF NOT EXISTS shelters (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name VARCHAR(100) NOT NULL,
    City VARCHAR(100) NOT NULL,
    Country VARCHAR(100) NOT NULL,
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
    ShelterID INTEGER NOT NULL,
    FOREIGN KEY(ShelterID) REFERENCES shelters(ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS shelterSessions (
    SessionKey VARCHAR(50) PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Password VARCHAR(100) NOT NULL,
    ShelterID INTEGER NOT NULL,
    LoginTime BIGINT NOT NULL,
    LastSeenTime BIGINT NOT NULL,
    FOREIGN KEY(ShelterID) REFERENCES shelters(ID) ON DELETE CASCADE
);
