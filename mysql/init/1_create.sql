CREATE DATABASE IF NOT EXISTS newRelease;

USE newRelease;

CREATE TABLE Artist (
    artistId    VARCHAR(100) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    url         VARCHAR(200) NOT NULL,
    iconUrl     VARCHAR(200) NOT NULL,
    PRIMARY KEY (artistId)
);

CREATE TABLE User (
    userId         VARCHAR(100) NOT NULL,
    accessToken    VARCHAR(300) NOT NULL,
    tokenType      VARCHAR(20) NOT NULL,
    refreshToken   VARCHAR(200) NOT NULL,
    expiry         DATETIME     NOT NULL,
    playlistId     VARCHAR(100) NOT NULL,
    ifRemixAdd     BOOLEAN  DEFAULT TRUE,
    ifAcousticAdd  BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (userId)
);

CREATE TABLE ListenTo (
    listenId    INT             NOT NULL AUTO_INCREMENT,
    userId      VARCHAR(100)    NOT NULL,
    artistId    VARCHAR(100)    NOT NULL,
    listenCount INT             NOT NULL,
    timestamp   DATETIME        NOT NULL,
    ifFollowing BOOLEAN NOT NULL,
    PRIMARY KEY (listenId),
    FOREIGN KEY (userId) REFERENCES User(userId) ON DELETE CASCADE,
    FOREIGN KEY (artistId) REFERENCES Artist(artistId) ON DELETE CASCADE
);