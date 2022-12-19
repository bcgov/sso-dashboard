import os
from starlette.config import Config
from sqlalchemy.engine.url import URL

env = Config('.env')

def getEnv(env_name: str, default=''):
    return os.environ.get(env_name, env.get(env_name, default=default))


config = dict(
    db_hostname=getEnv("DB_HOSTNAME"),
    db_port=getEnv("DB_PORT", 5432),
    db_database=getEnv("DB_DATABASE"),
    db_username=getEnv("DB_USERNAME"),
    db_password=getEnv("DB_PASSWORD"),
)

config["database_url"] = URL.create(
    "postgresql",
    username=config.get('db_username'),
    password=config.get('db_password'),
    host=config.get('db_hostname'),
    database=config.get('db_database'),
    port=config.get('db_port'),
).render_as_string(hide_password=False)
