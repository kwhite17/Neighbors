CREATE DATABASE IF NOT EXISTS neighbors;

USE neighbors;

CREATE TABLE IF NOT EXISTS users (
    ID INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    Username VARCHAR(100) NOT NULL UNIQUE KEY,
    Password VARCHAR(100) NOT NULL,
    Email VARCHAR(100),
    Phone VARCHAR(100),
    Location VARCHAR(100) NOT NULL,
    Role ENUM('NEIGHBOR', 'SAMARITAN') NOT NULL
);


CREATE TABLE userSession (
 SessionKey VARCHAR(100) PRIMARY KEY,
 UserID INT NOT NULL,
 LoginTime BIGINT not null,
 LastSeenTime BIGINT not null,
 FOREIGN KEY UserID (UserID)
 REFERENCES users(ID)
 ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items (
    ID INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    Category ENUM('SOCKS', 'UNDERWEAR') NOT NULL,
    Gender ENUM('MALE', 'FEMALE'),
    Size VARCHAR(100) NOT NULL,
    Quantity INT NOT NULL,
    DropoffLocation VARCHAR(100) NOT NULL,
    Requestor INT NOT NULL,
    Fulfiller INT,
    OrderStatus ENUM('REQUESTED', 'ASSIGNED', 'PURCHASED', 'DELIVERED'),
    FOREIGN KEY Requestor (Requestor)
    REFERENCES users(ID)
    ON DELETE CASCADE,
    FOREIGN KEY Fulfiller (Fulfiller)
    REFERENCES users(ID)
    ON DELETE CASCADE
);

CREATE USER IF NOT EXISTS 'neighbors_dba'@'localhost' IDENTIFIED BY 'neighbors_dba';
GRANT ALL PRIVILEGES ON neighbors.* TO 'neighbors_dba'@'localhost';
FLUSH PRIVILEGES;
