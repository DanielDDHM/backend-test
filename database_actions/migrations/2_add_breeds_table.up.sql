CREATE TABLE breeds (
  id INT AUTO_INCREMENT PRIMARY KEY,
  pet_size VARCHAR(50) NOT NULL,
  pet_name VARCHAR(255) NOT NULL,
  species VARCHAR(50) NOT NULL,
  average_male_adult_weight INT NOT NULL,
  average_female_adult_weight INT NOT NULL
);