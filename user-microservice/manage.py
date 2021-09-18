import os

from flask_restplus import Api
from flask import Blueprint
from flask_migrate import Migrate, MigrateCommand
from flask_script import Manager

from main import create_app, db

from main.controller.user_controller import api as user_ns
from main.controller.auth_controller import api as auth_ns
from main.model.user import User
from flask_bcrypt import Bcrypt

bcrypt = Bcrypt()


blueprint = Blueprint('api', __name__)

api = Api(blueprint,
          title='USER MICROSERVICE API',
          version='1.0',
          description='an api for user microservice'
        )

api.add_namespace(user_ns, path='/users')
api.add_namespace(auth_ns, path='/auth')

app = create_app(os.getenv('BOILERPLATE_ENV') or 'dev')
app.register_blueprint(blueprint)

app.app_context().push()

manager = Manager(app)

migrate = Migrate(app, db)

manager.add_command('db', MigrateCommand)

@manager.command
def run():
    app.run()

if __name__ == '__main__':
    try:
        db.session.add(User(username='manager', email='manager@gmail.com', first_name="James", last_name="Harden", role="ROLE_MANAGER", password=bcrypt.generate_password_hash('sifra123').decode()))
        db.session.add(User(username='worker', email='worker@gmail.com', first_name="LeBron", last_name="James", role="ROLE_WORKER", password=bcrypt.generate_password_hash('sifra123').decode()))
        db.session.commit()
    except:
        pass
    manager.run()
