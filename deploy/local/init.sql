-- creating env for sentry
CREATE DATABASE "sentry";
CREATE ROLE "sentry" WITH LOGIN PASSWORD 'sentry';
GRANT ALL PRIVILEGES ON DATABASE "sentry" to "sentry";
ALTER USER "sentry" WITH SUPERUSER;

-- creating env for keycloak
CREATE DATABASE "keycloak";
CREATE ROLE "keycloak" WITH LOGIN PASSWORD 'keycloak';
GRANT ALL PRIVILEGES ON DATABASE "keycloak" to "keycloak";
