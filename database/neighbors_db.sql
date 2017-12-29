create database if not exists neighbors;

use neighbors;

create table if not exists neighbors (
    NeighborID int not null primary key auto_increment,
    Username varchar(100) not null, 
    Password varchar(100) not null,
    Email varchar(100), 
    Phone varchar(100), 
    Location varchar(100) not null);

create table if not exists samaritans (
    SamaritanID int not null primary key auto_increment,
    Username varchar(100) not null,
    Password varchar(100) not null,
    Email varchar(100), 
    Phone varchar(100), 
    Location varchar(100) not null);

create table if not exists items (
    ItemID int not null primary key auto_increment, 
    Type varchar(100) not null, 
    Gender varchar(100), 
    Size varchar(100) not null, 
    Quantity int not null, 
    DropoffLocation varchar(100) not null, 
    Requestor int not null, 
    Fulfiller int, 
    FOREIGN KEY Requestor (Requestor) 
    REFERENCES neighbors(NeighborID) 
    ON DELETE CASCADE 
    ON UPDATE NO ACTION, 
    FOREIGN KEY Fulfiller (Fulfiller) 
    REFERENCES samaritans(SamaritanID) 
    ON DELETE NO ACTION 
    ON UPDATE NO ACTION);

create user 'neighbors_dba'@'localhost' IDENTIFIED BY 'neighbors_dba';

GRANT ALL PRIVILEGES ON neighbors.* TO 'neighbors_dba'@'localhost';

FLUSH PRIVILEGES;