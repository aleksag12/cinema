from flask import request
from flask_restplus import Resource

from ..util.dto import UserDto
from ..service.user_service import get_all_registered_users, get_all_workers, get_all_managers, delete_user

api = UserDto.api
_user = UserDto.user


@api.route('/registered-users')
class RegisteredUserList(Resource):
    @api.doc('list_of_registered_users')
    @api.marshal_list_with(_user)
    def get(self):
        return get_all_registered_users()


@api.route('/workers')
class WorkersList(Resource):
    @api.doc('list_of_all_workers')
    @api.marshal_list_with(_user)
    def get(self):
        return get_all_workers()


@api.route('/managers')
class ManagersList(Resource):
    @api.doc('list_of_all_managers')
    @api.marshal_list_with(_user)
    def get(self):
        return get_all_managers()


@api.route('/<id>')
@api.param('id', 'The User identifier')
@api.response(404, 'User not found.')
class User(Resource):
    @api.doc('delete user')
    def delete(self, id):
        return delete_user(id)