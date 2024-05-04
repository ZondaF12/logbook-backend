CREATE TABLE IF NOT EXISTS `vehicles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `registration` VARCHAR(255) NOT NULL,
  `make` VARCHAR(255) NOT NULL,
  `model` VARCHAR(255) NOT NULL,
  `year` INT UNSIGNED NOT NULL,
  `engine_size` INT UNSIGNED NOT NULL,
  `color` VARCHAR(255) NOT NULL,
  `registered` VARCHAR(255) NOT NULL,
  `tax_date` VARCHAR(255) NOT NULL,
  `mot_date` VARCHAR(255) NOT NULL,
  `insurance_date` VARCHAR(255) DEFAULT "",
  `service_date` VARCHAR(255) DEFAULT "",
  `description` VARCHAR(255) DEFAULT "",
  `milage` INT UNSIGNED DEFAULT 0,
  `nickname` VARCHAR(255) DEFAULT "",
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES auth(id)
)