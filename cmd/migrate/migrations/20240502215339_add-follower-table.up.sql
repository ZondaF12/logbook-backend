CREATE TABLE IF NOT EXISTS `followers` (
  `id` CHAR(36) NOT NULL,
  `follower_id` CHAR(36) NOT NULL,
  `following_id` CHAR(36) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  FOREIGN KEY (follower_id) REFERENCES auth(id),
  FOREIGN KEY (following_id) REFERENCES auth(id)
)