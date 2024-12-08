-- Seleziona il database
\c local_feature_flag;

CREATE TABLE IF NOT EXISTS project (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  owner_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS environment (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  public_key VARCHAR(100) NOT NULL,
  project_id INT NOT NULL,
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES project(id),
  CONSTRAINT unique_environment_per_project UNIQUE (project_id, name)
);

CREATE TABLE IF NOT EXISTS flag (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS flag_environment (
  id SERIAL PRIMARY KEY,
  environment INT NOT NULL,
  flag INT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT fk_environment FOREIGN KEY (environment) REFERENCES environment(id),
  CONSTRAINT fk_flag FOREIGN KEY (flag) REFERENCES flag(id)
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  nickname VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_project (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  project_id INT NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES project(id)
);
