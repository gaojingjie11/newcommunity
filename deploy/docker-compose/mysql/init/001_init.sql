CREATE DATABASE IF NOT EXISTS smart_community DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE smart_community;

CREATE TABLE IF NOT EXISTS migration_marker (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  stage VARCHAR(64) NOT NULL,
  note VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO migration_marker(stage, note)
VALUES ('microservice-skeleton', 'initial docker compose database marker');
