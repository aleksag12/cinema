from flask_restplus import Namespace, fields


class UserDto:
    api = Namespace('user', description='user related operations')
    new_user = api.model('new_user', {
        'username': fields.String(required=True, description='user username'),
        'email': fields.String(required=True, description='user email address'),
        'first_name': fields.String(required=True, description='user first name'),
        'last_name': fields.String(required=True, description='user last name'),
        'password': fields.String(required=True, description='user password'),
    })
    user = api.model('user', {
        'id': fields.Integer(description='user id'),
        'username': fields.String(required=True, description='user username'),
        'email': fields.String(required=True, description='user email address'),
        'first_name': fields.String(required=True, description='user first name'),
        'last_name': fields.String(required=True, description='user last name'),
        'role': fields.String(required=True, description='user role'),
    })
    token = api.model('token', {
        'token': fields.String(description='token'),
        'user': fields.Nested(user, description='user info'),
    })


class AuthDto:
    api = Namespace('auth', description='authentification related operations')
    auth = api.model('auth', {
        'username': fields.String(required=True, description='username'),
        'password': fields.String(required=True, description='password'),
    })
