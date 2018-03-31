-- Table for maintaining the Ghost lists in DB

CREATE TABLE IF NOT EXISTS `ghost_lists` (
  `list_id` varchar(2) NOT NULL,
  `ghost_key` varchar(16)  NOT NULL,
  `ghost_value` varchar(16)  NOT NULL,
  PRIMARY KEY (`list_id`, `ghost_key`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
