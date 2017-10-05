CREATE TABLE vpn_users (
  id         BIGSERIAL PRIMARY KEY,
  name       VARCHAR(255)                NOT NULL,
  email      VARCHAR(255)                NOT NULL,
  password   VARCHAR(255)                NOT NULL,
  details    TEXT                        NOT NULL,
  online     BOOLEAN                     NOT NULL DEFAULT FALSE,
  enable     BOOLEAN                     NOT NULL DEFAULT FALSE,
  _begin     DATE                        NOT NULL DEFAULT current_date,
  _end       DATE                        NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_vpn_users_email
  ON vpn_users (email);
CREATE INDEX idx_vpn_users_name
  ON vpn_users (name);

CREATE TABLE vpn_logs (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT                      NOT NULL REFERENCES vpn_users,
  trusted_ip   VARCHAR(45),
  trusted_port SMALLINT,
  remote_ip    VARCHAR(45),
  remote_port  SMALLINT,
  _begin       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  _end         TIMESTAMP WITHOUT TIME ZONE,
  received     FLOAT                       NOT NULL DEFAULT '0.0',
  send         FLOAT                       NOT NULL DEFAULT '0.0'
);
