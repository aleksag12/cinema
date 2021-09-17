from flask import request
from flask_restplus import Resource

from ..service.auth_service import delete_user, get_all_managers, get_all_workers, sign_in, sign_up, get_current_user, sign_out, get_all_registered_users
from ..util.dto import AuthDto

api = AuthDto.api
_auth = AuthDto.auth
_user = AuthDto.user
_new_user = AuthDto.new_user


@api.route('/<id>')
@api.param('id', 'The User identifier')
@api.response(404, 'User not found.')
class User(Resource):
    @api.doc('delete user')
    def delete(self, id):
        headers = request.headers
        return delete_user(id, headers)

@api.route('/sign-in')
class SignIn(Resource):
    @api.doc('sign in')
    @api.expect(_auth, validate=True)
    def post(self):
        data = request.json
        return sign_in(data)

@api.route('/sign-up')
class SignUp(Resource):
    @api.doc('sign up')
    @api.expect(_new_user, validate=True)
    def post(self):
        data = request.json
        headers = request.headers
        return sign_up(data, headers)


@api.route('/sign-out')
class SignOut(Resource):
    @api.doc('sign out')
    def get(self):
        headers = request.headers
        return sign_out(headers)


@api.route('/current-user')
class CurrentUser(Resource):
    @api.doc('get current user')
    @api.marshal_with(_user)
    def get(self):
        headers = request.headers
        return get_current_user(headers)


@api.route('/registered-users')
class RegisteredUsersList(Resource):
    @api.doc('list_of_all_registered_users')
    @api.marshal_list_with(_user)
    def get(self):
        headers = request.headers
        return get_all_registered_users(headers)


@api.route('/workers')
class WorkersList(Resource):
    @api.doc('list_of_all_workers')
    @api.marshal_list_with(_user)
    def get(self):
        headers = request.headers
        return get_all_workers(headers)


@api.route('/managers')
class ManagersList(Resource):
    @api.doc('list_of_all_managers')
    @api.marshal_list_with(_user)
    def get(self):
        headers = request.headers
        return get_all_managers(headers)
