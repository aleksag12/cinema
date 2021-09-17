from flask import Flask
from flask_cors import CORS

from .config import config_by_name


def create_app(config_name):
    app = Flask(__name__)
    app.config.from_object(config_by_name[config_name])
    CORS(app)

    @app.after_request
    def after_request(response):
        response.headers.add('Access-Control-Allow-Headers', 'Content-Type, x-requested-with')
        response.headers.add('Access-Control-Allow-Credentials', 'true')
        return response

    return app
