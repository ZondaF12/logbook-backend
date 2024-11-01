CREATE TABLE IF NOT EXISTS `logs` (
  `id` CHAR(36) NOT NULL,
  `vehicle_id` CHAR(36) NOT NULL,
  `category` INT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `date` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) DEFAULT "",
  `notes` VARCHAR(255) DEFAULT "",
  `cost` DECIMAL(10, 2) DEFAULT 0.00,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (id),
  FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
)