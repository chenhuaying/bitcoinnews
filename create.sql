CREATE TABLE IF NOT EXISTS `cc_news` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL UNIQUE,
  `media` varchar(30) DEFAULT NULL,
  `url` varchar(255) NOT NULL,
  `ctime` datetime NOT NULL,
  `time` datetime NOT NULL,
  `content` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
