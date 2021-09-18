from flask import request
from flask_restplus import Resource

from ..service.auth_service import sign_in, get_current_user, sign_out, sign_up
from ..util.dto import AuthDto, UserDto

api = AuthDto.api
_auth = AuthDto.auth
_token = UserDto.token
_user = UserDto.user
_new_user = UserDto.new_user


@api.route('/sign-in')
class SignIn(Resource):
    @api.doc('sign in')
    @api.expect(_auth, validate=True)
    @api.marshal_with(_token)
    def post(self):
        data = request.json
        return sign_in(data=data)


@api.route('/sign-up')
class SignUp(Resource):
    @api.doc('sign up')
    @api.expect(_new_user, validate=True)
    def post(self):
        data = request.json
        return sign_up(data, request.headers)


@api.route('/sign-out')
class SignOut(Resource):
    @api.doc('sign out')
    def get(self):
        return sign_out()


@api.route('/current-user')
class CurrentUser(Resource):
    @api.doc('get current user')
    @api.marshal_with(_user)
    def get(self):
        headers = request.headers
        return get_current_user(headers)
