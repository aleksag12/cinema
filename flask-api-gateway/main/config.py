import os

basedir = os.path.abspath(os.path.dirname(__file__))

class Config:
    SECRET_KEY = os.getenv('SECRET_KEY', 'UUSRHZPjPVgcIyWyGVGPp5Rj6pFaVgSg')
    DEBUG = False


class DevelopmentConfig(Config):
    DEBUG = True
    SQLALCHEMY_DATABASE_URI = 'sqlite:///' + os.path.join(basedir, 'timesheet_main.db')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    USER_MICROSERVICE_HOST = 'localhost'
    USER_MICROSERVICE_PORT = 5001
    MOVIE_MICROSERVICE_HOST = 'localhost'
    MOVIE_MICROSERVICE_PORT = 9090
    PROJECTION_MICROSERVICE_HOST = 'localhost'
    PROJECTION_MICROSERVICE_PORT = 9091
    RATING_MICROSERVICE_HOST = 'localhost'
    RATING_MICROSERVICE_PORT = 9092
    TICKET_MICROSERVICE_HOST = 'localhost'
    TICKET_MICROSERVICE_PORT = 9093

    MAIL_SERVER = 'smtp.gmail.com'
    MAIL_PORT = 465
    MAIL_USE_SSL = True
    MAIL_USE_TLS = False
    MAIL_USERNAME = 'email@example.com'
    MAIL_PASSWORD = 'password'
    MAIL_DEFAULT_SENDER = '"MyApp" <noreply@example.com>'


class TestingConfig(Config):
    DEBUG = True
    TESTING = True
    SQLALCHEMY_DATABASE_URI = 'sqlite:///' + os.path.join(basedir, 'timesheet_test.db')
    PRESERVE_CONTEXT_ON_EXCEPTION = False
    SQLALCHEMY_TRACK_MODIFICATIONS = False


class ProductionConfig(Config):
    DEBUG = False


config_by_name = dict(
    dev=DevelopmentConfig,
    test=TestingConfig,
    prod=ProductionConfig
)

key = Config.SECRET_KEY
