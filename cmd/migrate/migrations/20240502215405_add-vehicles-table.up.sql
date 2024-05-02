CREATE TABLE IF NOT EXISTS `vehicles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `registration` VARCHAR(255) NOT NULL,
  `make` VARCHAR(255) NOT NULL,
  `model` VARCHAR(255) NOT NULL,
  `year` INT UNSIGNED NOT NULL,
  `engine_size` INT UNSIGNED NOT NULL,
  `color` VARCHAR(255) NOT NULL,
  `registered` DATE NOT NULL,
  `tax_date` VARCHAR(255) NOT NULL,
  `mot_date` DATE NOT NULL,
  `insurance_date` DATE,
  `service_date` DATE,
  `description` VARCHAR(255),
  `milage` INT UNSIGNED,


  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES auth(id)
)