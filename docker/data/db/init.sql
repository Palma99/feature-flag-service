-- Seleziona il database
\c local_feature_flag;

-- Creazione della tabella project
CREATE TABLE IF NOT EXISTS project (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

-- Creazione della tabella environment
CREATE TABLE IF NOT EXISTS environment (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  public_key VARCHAR(100) NOT NULL,
  private_key VARCHAR(100) NOT NULL,
  project_id INT NOT NULL,
  CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES project(id)
);

CREATE TABLE IF NOT EXISTS flag (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS flag_environment (
  id SERIAL PRIMARY KEY,
  environment INT NOT NULL,
  flag INT NOT NULL,
  CONSTRAINT fk_environment FOREIGN KEY (environment) REFERENCES environment(id),
  CONSTRAINT fk_flag FOREIGN KEY (flag) REFERENCES flag(id)
)