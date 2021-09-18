from main import db
from main.model.user import User
from main.util.requests import bad_request, unauthorized, ok
from flask_bcrypt import Bcrypt

bcrypt = Bcrypt()

def sign_in(data):
    user = User.query.filter_by(username=data.get('username')).first()
    if not user or not bcrypt.check_password_hash(user.password.encode('utf-8'), data.get('password')):
        return unauthorized("Username or password does not match.")

    auth_token = user.encode_auth_token()
    return {
        "token": auth_token.decode(),
        "user": user
    }, 200


def sign_up(data, headers):
    if data.get('role') != 'ROLE_REGISTERED_USER':
        auth_header = headers.get('Authorization')
        if not auth_header or "Bearer " not in auth_header:
            return unauthorized("Only managers can add new workers and managers.")
        else:
            try:
                if User.decode_auth_token(auth_header.split(" ")[1]).get('role') != 'ROLE_MANAGER':
                    return unauthorized("Only managers can add new workers and managers.")
            except:
                return unauthorized("Only managers can add new workers and managers.")

    user_by_username = User.query.filter_by(username=data.get('username')).first()
    user_by_email = User.query.filter_by(email=data.get('email')).first()
    if not user_by_username and not user_by_email:
        try:
            new_user = User(
                username=data.get('username'), 
                email=data.get('email'), 
                first_name=data.get('first_name'), 
                last_name=data.get('last_name'), 
                password=bcrypt.generate_password_hash(data.get('password')).decode(),
                role=data.get('role')
            )

            db.session.add(new_user)
            db.session.commit()

            return ok('Successfully registered.')

        except Exception as e:
            return bad_request('Some error occurred. Please try again.')
    else:
        return bad_request('User already exists. Please Log in.')


def sign_out():
    return ok("Successfully logged out.")


def get_current_user(headers):
    auth_header = headers.get('Authorization')
    if not auth_header or "Bearer " not in auth_header:
        return unauthorized("Log in first.")
    else:
        try:
            user_id = User.decode_auth_token(auth_header.split(" ")[1]).get('id')
            user = User.query.filter_by(id=user_id).first()
            if not user:
                return unauthorized("Log in first.")
            return user
        except:
            return unauthorized("Log in first.")
