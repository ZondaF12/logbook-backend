CREATE TABLE IF NOT EXISTS `profiles` (
  `id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `username` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `bio` VARCHAR(255) DEFAULT "",
  `avatar` VARCHAR(255) DEFAULT "",
  `public` BOOLEAN NOT NULL DEFAULT TRUE,

  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES auth(id),
  UNIQUE KEY (username)
);