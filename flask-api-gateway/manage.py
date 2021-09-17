import os
import unittest

from flask_restplus import Api
from flask import Blueprint
from flask_script import Manager

from main import create_app

from main.controller.movies_controller import api as movies_ns
from main.controller.rates_controller import api as rates_ns
from main.controller.comments_controller import api as comments_ns
from main.controller.auth_controller import api as auth_ns
from main.controller.projections_controller import api as projection_ns
from main.controller.tickets_controller import api as ticket_ns


blueprint = Blueprint('api', __name__)

api = Api(blueprint,
          title='GO-CINEMA API GATEWAY',
          version='1.0',
          description='an api-gateway for go-cinema web service'
          )

api.add_namespace(movies_ns, path='/api/movies')
api.add_namespace(rates_ns, path='/api/rates')
api.add_namespace(comments_ns, path='/api/comments')
api.add_namespace(projection_ns, path='/api/projections')
api.add_namespace(ticket_ns, path='/api/tickets')
api.add_namespace(auth_ns, path='/auth')

app = create_app(os.getenv('BOILERPLATE_ENV') or 'dev')
app.register_blueprint(blueprint)

app.app_context().push()

manager = Manager(app)

@manager.command
def run():
    app.run()

@manager.command
def test():
    tests = unittest.TestLoader().discover('flask-api-gateway/test', pattern='test*.py')
    result = unittest.TextTestRunner(verbosity=2).run(tests)
    if result.wasSuccessful():
        return 0
    return 1

if __name__ == '__main__':
    manager.run()
