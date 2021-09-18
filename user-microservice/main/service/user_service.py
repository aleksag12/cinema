from main import db
from main.model.user import User
from main.util.requests import created, conflict, not_found, ok, bad_request


def save_new_user(data):
    user_by_email = User.query.filter_by(email=data['email']).first()
    user_by_username = User.query.filter_by(username=data['username']).first()
    if not user_by_email and not user_by_username:
        new_user = User(
            first_name=data['first_name'],
            last_name=data['last_name'],
            username=data['username'],
            email=data['email'],
            password=data['password'],
            role=data['role'],
        )
        save_changes(new_user)
        return created('Successfully registered user.')
    else:
        return conflict('User already exists.')


def get_all_registered_users():
    return User.query.filter_by(role='ROLE_REGISTERED_USER').all()


def get_all_workers():
    return User.query.filter_by(role='ROLE_WORKER').all()


def get_all_managers():
    return User.query.filter_by(role='ROLE_MANAGER').all()


def delete_user(id):
    user = User.query.filter_by(id=int(id)).first()
    if user:
        db.session.delete(user)
        db.session.commit()
        return ok('User successfully deleted.')
    else:
        return not_found('User not found')


def save_changes(data):
    db.session.add(data)
    db.session.commit()
