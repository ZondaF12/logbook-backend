CREATE TABLE IF NOT EXISTS `images` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,  -- Unique identifier for each image record
  `filename` VARCHAR(255) NOT NULL,  -- Name of the image file
  `file_type` VARCHAR(100) NOT NULL,  -- MIME type of the file
  `s3_location` VARCHAR(500) NOT NULL,  -- S3 URL or path to the image
  `uploaded_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the image was uploaded
  
  `user_id` INT UNSIGNED,
  `vehicle_id` INT UNSIGNED,

  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES auth(id),
  FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
)