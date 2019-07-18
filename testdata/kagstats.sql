-- MySQL dump 10.13  Distrib 5.7.25, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: kagstats
-- ------------------------------------------------------
-- Server version	5.7.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE IF NOT EXISTS kagstats;
USE kagstats;

--
-- Table structure for table `events`
--

DROP TABLE IF EXISTS `events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `events` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `playerID` int(11) NOT NULL,
  `type` varchar(30) NOT NULL,
  `time` int(11) NOT NULL,
  `serverID` int(11) NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `playerID` (`playerID`),
  KEY `serverID` (`serverID`),
  CONSTRAINT `events_ibfk_1` FOREIGN KEY (`playerID`) REFERENCES `players` (`ID`),
  CONSTRAINT `events_ibfk_2` FOREIGN KEY (`serverID`) REFERENCES `servers` (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `events`
--

LOCK TABLES `events` WRITE;
/*!40000 ALTER TABLE `events` DISABLE KEYS */;
INSERT INTO `events` VALUES (1,12,'joined',1563410049,3),(2,20,'joined',1563410049,3),(3,1,'joined',1563410081,3),(4,1,'left',1563410102,3),(5,1,'joined',1563410109,3),(6,1,'left',1563410116,3),(7,1,'joined',1563410122,3),(8,1,'left',1563410125,3);
/*!40000 ALTER TABLE `events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `kills`
--

DROP TABLE IF EXISTS `kills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `kills` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `killerID` int(11) NOT NULL,
  `victimID` int(11) NOT NULL,
  `killerClass` enum('archer','builder','knight','other','none') DEFAULT 'none',
  `victimClass` enum('archer','builder','knight','other') NOT NULL DEFAULT 'archer',
  `hitter` int(11) DEFAULT '0',
  `epoch` int(11) NOT NULL,
  `serverID` int(11) NOT NULL,
  `teamKill` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `killerID` (`killerID`),
  KEY `victimID` (`victimID`),
  KEY `serverID` (`serverID`),
  CONSTRAINT `kills_ibfk_1` FOREIGN KEY (`killerID`) REFERENCES `players` (`ID`),
  CONSTRAINT `kills_ibfk_2` FOREIGN KEY (`victimID`) REFERENCES `players` (`ID`),
  CONSTRAINT `kills_ibfk_4` FOREIGN KEY (`serverID`) REFERENCES `servers` (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=317 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `kills`
--

LOCK TABLES `kills` WRITE;
/*!40000 ALTER TABLE `kills` DISABLE KEYS */;
INSERT INTO `kills` VALUES (1,1,1,'knight','knight',25,1556719870,1,0),(2,1,2,'knight','knight',14,1556720060,1,0),(3,2,1,'knight','knight',14,1556720081,1,0),(4,2,1,'knight','archer',14,1556720096,1,0),(5,2,1,'archer','knight',10,1556720117,1,0),(6,2,1,'knight','knight',14,1556720131,1,0),(7,0,1,'none','knight',28,1556720149,1,0),(8,2,1,'knight','knight',14,1556720158,1,0),(9,0,2,'none','knight',0,1556720160,1,0),(10,2,1,'knight','knight',14,1556720194,1,0),(11,1,2,'knight','knight',14,1556720219,1,0),(12,2,1,'archer','knight',23,1556720237,1,0),(13,2,1,'archer','knight',10,1556720281,1,0),(14,1,2,'knight','archer',14,1556720352,1,0),(15,2,1,'archer','archer',18,1556720366,1,0),(16,2,1,'archer','archer',18,1556720383,1,0),(17,2,3,'archer','archer',18,1556720403,1,0),(18,3,2,'archer','archer',18,1556720421,1,0),(19,3,1,'archer','archer',18,1556720443,1,0),(20,2,3,'archer','archer',19,1556720456,1,0),(21,0,3,'none','archer',0,1556720484,1,0),(22,1,2,'archer','archer',18,1556720488,1,0),(23,1,4,'archer','knight',11,1556720513,1,0),(24,0,1,'none','archer',0,1556720516,1,0),(25,3,2,'archer','archer',18,1556720530,1,0),(26,3,1,'archer','archer',18,1556720540,1,0),(27,1,3,'archer','archer',8,1556720557,1,0),(28,4,2,'knight','archer',14,1556720564,1,0),(29,2,1,'archer','archer',28,1556720567,1,0),(30,4,2,'knight','knight',14,1556720585,1,0),(31,3,1,'archer','archer',18,1556720591,1,0),(32,5,4,'knight','knight',10,1556720618,1,0),(33,2,3,'knight','archer',14,1556720619,1,0),(34,5,1,'knight','archer',14,1556720635,1,0),(35,5,4,'knight','knight',14,1556720653,1,0),(36,2,3,'archer','archer',18,1556720662,1,0),(37,3,2,'archer','archer',18,1556720682,1,0),(38,5,4,'knight','knight',14,1556720684,1,0),(39,3,5,'archer','knight',18,1556720687,1,0),(40,5,1,'knight','archer',14,1556720697,1,0),(41,2,4,'knight','knight',14,1556720699,1,0),(42,2,3,'archer','archer',10,1556720707,1,0),(43,2,5,'archer','knight',18,1556720708,1,0),(44,4,2,'knight','archer',14,1556720728,1,0),(45,5,1,'archer','archer',18,1556720728,1,0),(46,4,1,'knight','archer',14,1556720745,1,0),(47,0,3,'none','archer',29,1556720749,1,0),(48,4,5,'knight','archer',14,1556720752,1,0),(49,2,4,'knight','knight',14,1556720763,1,0),(50,1,3,'archer','archer',18,1556720783,1,0),(51,5,2,'archer','knight',18,1556720786,1,0),(52,5,1,'archer','archer',18,1556720803,1,0),(53,5,4,'archer','knight',18,1556720806,1,0),(54,2,4,'archer','knight',18,1556720826,1,0),(55,1,2,'archer','archer',18,1556720840,1,0),(56,5,1,'archer','archer',18,1556720856,1,0),(57,3,1,'archer','archer',18,1556720882,1,0),(58,5,4,'archer','knight',18,1556720883,1,0),(59,2,3,'archer','archer',18,1556720884,1,0),(60,5,2,'archer','archer',18,1556720900,1,0),(61,2,3,'archer','archer',18,1556720915,1,0),(62,2,1,'archer','archer',18,1556720920,1,0),(63,4,2,'knight','knight',14,1556720934,1,0),(64,4,3,'knight','archer',14,1556720952,1,0),(65,1,2,'archer','archer',2,1556720971,1,0),(66,5,3,'archer','archer',7,1556720978,1,0),(67,2,1,'archer','archer',19,1556720991,1,0),(68,3,4,'archer','knight',18,1556720994,1,0),(69,1,5,'archer','archer',18,1556721023,1,0),(70,2,1,'knight','archer',14,1556721036,1,0),(71,2,4,'knight','knight',14,1556721043,1,0),(72,2,3,'knight','archer',14,1556721064,1,0),(73,3,1,'archer','archer',10,1556721084,1,0),(74,5,2,'archer','archer',18,1556721092,1,0),(75,5,3,'archer','archer',18,1556721119,1,0),(76,5,1,'archer','knight',18,1556721120,1,0),(77,0,5,'none','archer',18,1556721129,1,0),(78,0,3,'none','archer',18,1556721131,1,0),(79,2,1,'archer','archer',26,1556721149,1,0),(80,2,2,'archer','archer',26,1556721149,1,0),(81,0,3,'none','archer',18,1556721157,1,0),(82,2,5,'archer','archer',18,1556721163,1,0),(83,2,1,'archer','archer',18,1556721168,1,0),(84,5,3,'archer','archer',18,1556721187,1,0),(85,6,5,'archer','archer',18,1556721189,1,0),(86,1,6,'archer','archer',18,1556721191,1,0),(87,2,1,'archer','archer',8,1556721200,1,0),(88,2,5,'archer','archer',18,1556721222,1,0),(89,2,6,'archer','archer',18,1556721230,1,0),(90,5,2,'archer','archer',18,1556721252,1,0),(91,3,1,'archer','archer',10,1556721262,1,0),(92,2,6,'archer','archer',18,1556721284,1,0),(93,5,3,'archer','archer',18,1556721289,1,0),(94,2,1,'archer','archer',18,1556721301,1,0),(95,2,5,'archer','archer',18,1556721320,1,0),(96,2,6,'archer','archer',18,1556721338,1,0),(97,2,7,'archer','archer',18,1556721340,1,0),(98,1,5,'archer','archer',18,1556721341,1,0),(99,2,1,'archer','archer',18,1556721356,1,0),(100,2,7,'archer','archer',18,1556721363,1,0),(101,2,3,'archer','archer',18,1556721367,1,0),(102,3,5,'archer','archer',7,1556721380,1,0),(103,1,7,'archer','archer',18,1556721384,1,0),(104,1,6,'archer','archer',18,1556721392,1,0),(105,1,3,'archer','archer',18,1556721402,1,0),(106,3,5,'archer','archer',18,1556721416,1,0),(107,8,6,'archer','archer',18,1556721418,1,0),(108,7,1,'archer','archer',18,1556721438,1,0),(109,2,7,'archer','archer',18,1556721441,1,0),(110,2,5,'archer','knight',18,1556721458,1,0),(111,2,6,'archer','archer',18,1556721461,1,0),(112,3,8,'archer','knight',18,1556721462,1,0),(113,2,1,'archer','archer',18,1556721465,1,0),(114,1,8,'archer','knight',18,1556721479,1,0),(115,2,1,'archer','archer',8,1556721481,1,0),(116,2,3,'archer','archer',18,1556721486,1,0),(117,5,2,'knight','archer',14,1556721493,1,0),(118,5,7,'knight','archer',15,1556721506,1,0),(119,5,6,'knight','archer',14,1556721509,1,0),(120,8,2,'knight','archer',14,1556721521,1,0),(121,3,6,'archer','archer',18,1556721523,1,0),(122,7,1,'archer','archer',18,1556721526,1,0),(123,5,5,'knight','knight',25,1556721527,1,0),(124,8,3,'knight','archer',14,1556721535,1,0),(125,7,8,'knight','knight',14,1556721559,1,0),(126,5,3,'archer','archer',18,1556721572,1,0),(127,1,2,'archer','archer',18,1556721578,1,0),(128,5,7,'archer','archer',18,1556721579,1,0),(129,8,1,'knight','archer',14,1556721581,1,0),(130,8,6,'knight','archer',14,1556721586,1,0),(131,7,2,'archer','archer',18,1556721600,1,0),(132,5,7,'archer','archer',18,1556721605,1,0),(133,3,1,'archer','archer',7,1556721606,1,0),(134,5,6,'archer','archer',18,1556721612,1,0),(135,3,8,'archer','knight',18,1556721615,1,0),(136,5,3,'archer','archer',18,1556721621,1,0),(137,5,2,'archer','archer',19,1556721635,1,0),(138,5,3,'archer','archer',19,1556721636,1,0),(139,5,5,'archer','archer',19,1556721636,1,0),(140,1,7,'archer','archer',18,1556721641,1,0),(141,6,1,'archer','archer',18,1556721641,1,0),(142,8,6,'knight','archer',14,1556721648,1,0),(143,2,1,'archer','archer',18,1556721665,1,0),(144,8,3,'knight','archer',10,1556721670,1,0),(145,3,8,'archer','knight',26,1556721671,1,0),(146,3,6,'archer','archer',26,1556721671,1,0),(147,2,7,'archer','archer',18,1556721674,1,0),(148,1,6,'archer','archer',18,1556721691,1,0),(149,3,5,'archer','archer',18,1556721700,1,0),(150,5,3,'archer','archer',18,1556721700,1,0),(151,1,8,'archer','knight',18,1556721706,1,0),(152,3,3,'archer','archer',26,1556721729,1,0),(153,7,5,'archer','archer',18,1556721740,1,0),(154,1,2,'archer','archer',18,1556721757,1,0),(155,3,5,'archer','archer',18,1556721758,1,0),(156,7,3,'archer','archer',18,1556721764,1,0),(157,7,6,'archer','archer',18,1556721786,1,0),(158,3,7,'archer','archer',7,1556721805,1,0),(159,2,5,'archer','archer',18,1556721811,1,0),(160,3,1,'archer','archer',18,1556721814,1,0),(161,3,2,'archer','archer',18,1556721829,1,0),(162,0,3,'none','archer',29,1556721831,1,0),(163,7,5,'archer','archer',18,1556721846,1,0),(164,1,6,'archer','archer',18,1556721848,1,0),(165,7,1,'archer','archer',18,1556721880,1,0),(166,1,7,'archer','archer',18,1556721880,1,0),(167,2,6,'archer','archer',7,1556721894,1,0),(168,0,5,'none','archer',18,1556721910,1,0),(169,3,7,'archer','archer',18,1556721918,1,0),(170,1,2,'archer','archer',18,1556721921,1,0),(171,5,3,'archer','archer',18,1556721945,1,0),(172,2,7,'archer','archer',18,1556721954,1,0),(173,2,5,'archer','archer',18,1556721958,1,0),(174,2,1,'archer','archer',8,1556721969,1,0),(175,5,7,'archer','archer',18,1556721984,1,0),(176,2,3,'archer','archer',18,1556721985,1,0),(177,1,2,'archer','archer',18,1556721986,1,0),(178,1,9,'archer','knight',18,1556722005,1,0),(179,5,1,'archer','archer',18,1556722024,1,0),(180,9,7,'knight','archer',14,1556722027,1,0),(181,0,5,'none','archer',29,1556722036,1,0),(182,9,7,'knight','archer',14,1556722050,1,0),(183,5,2,'knight','archer',14,1556722063,1,0),(184,5,3,'knight','archer',14,1556722089,1,0),(185,5,9,'knight','knight',14,1556722113,1,0),(186,5,9,'knight','knight',14,1556722130,1,0),(187,3,1,'archer','archer',18,1556722132,1,0),(188,2,7,'archer','archer',18,1556722144,1,0),(189,3,2,'archer','archer',18,1556722146,1,0),(190,5,3,'knight','archer',14,1556722152,1,0),(191,9,7,'knight','archer',14,1556722171,1,0),(192,2,5,'knight','knight',14,1556722172,1,0),(193,2,3,'knight','archer',14,1556722179,1,0),(194,2,5,'knight','knight',14,1556722194,1,0),(195,2,9,'knight','knight',14,1556722194,1,0),(196,7,1,'archer','archer',18,1556722195,1,0),(197,2,1,'knight','archer',14,1556722211,1,0),(198,2,9,'knight','knight',10,1556722211,1,0),(199,5,10,'knight','knight',14,1556722212,1,0),(200,9,1,'knight','archer',14,1556722233,1,0),(201,3,9,'knight','knight',14,1556722240,1,0),(202,2,3,'archer','knight',11,1556722250,1,0),(203,2,9,'archer','knight',23,1556722267,1,0),(204,0,2,'none','archer',29,1556722268,1,0),(205,3,1,'archer','archer',18,1556722284,1,0),(206,0,3,'none','archer',29,1556722294,1,0),(207,2,7,'archer','archer',18,1556722315,1,0),(208,2,3,'archer','archer',18,1556722317,1,0),(209,2,3,'knight','archer',14,1556722335,1,0),(210,9,2,'knight','knight',14,1556722336,1,0),(211,9,7,'knight','knight',14,1556722352,1,0),(212,0,9,'none','knight',29,1556722364,1,0),(213,2,3,'archer','archer',18,1556722370,1,0),(214,2,7,'archer','archer',18,1556722399,1,0),(215,9,2,'knight','archer',14,1556722411,1,0),(216,7,2,'archer','archer',18,1556722435,1,0),(217,0,1,'none','knight',0,1556724559,1,0),(218,1,12,'knight','knight',14,1556724843,1,0),(219,1,12,'knight','knight',14,1556724889,1,0),(220,1,12,'knight','knight',14,1556724906,1,0),(221,14,13,'archer','archer',18,1556728335,2,0),(222,13,15,'archer','archer',18,1556728355,2,0),(223,14,13,'archer','archer',18,1556728366,2,0),(224,15,13,'archer','archer',18,1556728399,2,0),(225,15,13,'knight','archer',14,1556728425,2,0),(226,15,14,'knight','archer',14,1556728427,2,0),(227,13,15,'knight','knight',18,1556728446,2,0),(228,0,13,'none','knight',0,1556728449,2,0),(229,0,13,'none','archer',28,1556728482,2,0),(230,15,14,'knight','knight',14,1556728497,2,0),(231,15,14,'knight','knight',14,1556728513,2,0),(232,15,13,'knight','knight',14,1556728529,2,0),(233,14,13,'knight','archer',14,1556728567,2,0),(234,15,14,'knight','knight',27,1556728596,2,0),(235,14,15,'knight','knight',14,1556728623,2,0),(236,13,14,'archer','knight',18,1556728639,2,0),(237,14,13,'knight','archer',14,1556728664,2,0),(238,15,13,'knight','knight',14,1556728695,2,0),(239,0,15,'none','knight',0,1556728698,2,0),(240,14,13,'knight','archer',14,1556728728,2,0),(241,15,14,'knight','knight',14,1556728741,2,0),(242,13,15,'archer','knight',10,1556728796,2,0),(243,14,15,'knight','knight',14,1556728819,2,0),(244,14,13,'knight','knight',14,1556728862,2,0),(245,15,13,'knight','knight',14,1556728881,2,0),(246,15,14,'knight','knight',14,1556728908,2,0),(247,15,13,'knight','knight',14,1556728916,2,0),(248,13,14,'knight','knight',14,1556728936,2,0),(249,14,13,'knight','knight',14,1556728960,2,0),(250,14,15,'knight','knight',14,1556728971,2,0),(251,13,15,'archer','knight',19,1556728989,2,0),(252,15,14,'knight','knight',28,1556728989,2,0),(253,0,14,'none','knight',29,1556729012,2,0),(254,14,13,'knight','knight',14,1556729092,2,0),(255,15,14,'knight','knight',14,1556729111,2,0),(256,14,13,'knight','archer',14,1556729172,2,0),(257,14,15,'knight','knight',14,1556729177,2,0),(258,14,15,'knight','knight',14,1556729202,2,0),(259,14,1,'knight','archer',14,1556729247,2,0),(260,14,13,'knight','knight',10,1556729263,2,0),(261,14,1,'knight','archer',14,1556729269,2,0),(262,14,13,'knight','knight',10,1556729320,2,0),(263,16,14,'knight','knight',14,1556729334,2,0),(264,16,15,'knight','knight',14,1556729348,2,0),(265,15,14,'knight','knight',14,1556729377,2,0),(266,14,16,'knight','knight',14,1556729415,2,0),(267,13,14,'knight','knight',14,1556729418,2,0),(268,13,15,'knight','knight',14,1556729423,2,0),(269,17,14,'knight','knight',14,1556729435,2,0),(270,17,13,'knight','knight',14,1556729435,2,0),(271,16,15,'knight','knight',14,1556729448,2,0),(272,17,16,'knight','knight',14,1556729457,2,0),(273,13,14,'knight','knight',14,1556729469,2,0),(274,13,17,'knight','knight',14,1556729469,2,0),(275,14,15,'knight','knight',14,1556729492,2,0),(276,14,13,'knight','knight',14,1556729503,2,0),(277,17,14,'knight','knight',14,1556729518,2,0),(278,17,15,'knight','knight',14,1556729521,2,0),(279,10,17,'knight','knight',14,1556729539,2,0),(280,14,13,'knight','knight',16,1556729546,2,0),(281,15,10,'knight','knight',14,1556729552,2,0),(282,14,10,'knight','knight',14,1556729574,2,0),(283,17,15,'knight','knight',14,1556729579,2,0),(284,17,14,'knight','knight',14,1556729587,2,0),(285,13,14,'knight','knight',14,1556729657,2,0),(286,17,14,'knight','knight',14,1556729706,2,0),(287,14,13,'knight','knight',14,1556729726,2,0),(288,14,17,'knight','knight',14,1556729738,2,0),(289,14,13,'knight','knight',14,1556729763,2,0),(290,14,17,'knight','knight',14,1556729771,2,0),(291,17,14,'knight','knight',14,1556729791,2,0),(292,14,17,'knight','knight',14,1556729818,2,0),(293,17,13,'knight','knight',10,1556729831,2,0),(294,14,17,'knight','knight',14,1556729845,2,0),(295,14,17,'knight','knight',14,1556729860,2,0),(296,17,14,'knight','knight',14,1556729904,2,0),(297,14,17,'knight','knight',16,1556729908,2,0),(298,17,14,'knight','knight',14,1556729917,2,0),(299,14,13,'knight','knight',14,1556729952,2,0),(300,13,14,'knight','knight',14,1556729982,2,0),(301,17,18,'archer','archer',18,1556729999,2,0),(302,14,17,'knight','knight',14,1556730023,2,0),(303,14,18,'knight','archer',14,1556730030,2,0),(304,18,17,'knight','knight',14,1556730051,2,0),(305,13,14,'knight','knight',14,1556730052,2,0),(306,13,17,'knight','knight',14,1556730068,2,0),(307,13,18,'knight','knight',14,1556730074,2,0),(308,13,14,'knight','knight',14,1556730089,2,0),(309,17,18,'knight','knight',14,1556730100,2,0),(310,14,18,'knight','knight',14,1556730123,2,0),(311,13,14,'knight','knight',14,1556730130,2,0),(312,13,17,'knight','knight',14,1556730131,2,0),(313,17,18,'knight','archer',14,1556730151,2,0),(314,17,18,'knight','archer',14,1556730165,2,0),(315,18,17,'archer','archer',18,1556730206,2,0),(316,17,18,'knight','knight',14,1556730243,2,0);
/*!40000 ALTER TABLE `kills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `players`
--

DROP TABLE IF EXISTS `players`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `players` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) COLLATE utf8mb4_bin NOT NULL,
  `charactername` varchar(30) COLLATE utf8mb4_bin NOT NULL,
  `clantag` varchar(30) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `players`
--

LOCK TABLES `players` WRITE;
/*!40000 ALTER TABLE `players` DISABLE KEYS */;
INSERT INTO `players` VALUES (1,'Verrazano','Witch',''),(2,'Turtlebutt','yellow powerranger','[NI]'),(3,'Asu','Asu','HOMEK'),(4,'Nalas','Day','HOMEK'),(5,'Bunnie','Bunnie','HOMEK'),(6,'Matixar','Matixar',''),(7,'BattleCat','The Cat Wizard','HOMEK'),(8,'otaviobia','otaviobia',''),(9,'killer_bug','xXIlluminatiPartyXx','INFINITY'),(10,'Ferrezinhre','Ferrezinhre',''),(11,'Egas','Egas',''),(12,'Henry','Henry',''),(13,'AlexSelister2002','Alex_Selister',''),(14,'fyri14','viv-ha','SD'),(15,'KloynCrasty','KloynCrasty',''),(16,'Artis','Artis',''),(17,'mav3rick2','mav3rick2',''),(18,'alexshow66','MiKKiLiN','');
/*!40000 ALTER TABLE `players` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `servers`
--

DROP TABLE IF EXISTS `servers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `servers` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `gamemode` varchar(30) NOT NULL,
  `tags` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `servers`
--

LOCK TABLES `servers` WRITE;
/*!40000 ALTER TABLE `servers` DISABLE KEYS */;
INSERT INTO `servers` VALUES (1,'Verra\'s TDM with stats - US East','...','TDM','US,TDM,Test'),(2,'Verra\'s TDM with stats - EU','...','TDM','EU,TDM,Test'),(3,'Verra\'s Testing Server','...','TDM','');
/*!40000 ALTER TABLE `servers` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-07-17 19:35:59
