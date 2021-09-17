from main.util.requests import unauthorized
import requests
from main.util.helper import get_path, guard_check


def sign_in(data):
    r = requests.post(get_path('USER_MICROSERVICE') + '/auth/sign-in', json = data)
    return r.json(), r.status_code


def sign_up(data, headers):
    r = requests.post(get_path('USER_MICROSERVICE') + '/auth/sign-up', json = data, headers = headers)
    return r.json(), r.status_code


def get_current_user(headers):
    r = requests.get(get_path('USER_MICROSERVICE') + '/auth/current-user', headers = headers)
    return r.json(), r.status_code


def sign_out(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER', 'ROLE_REGISTERED_USER', 'ROLE_WORKER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('USER_MICROSERVICE') + '/auth/sign-out')
    return r.json(), r.status_code


def get_all_registered_users(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('USER_MICROSERVICE') + '/users/registered-users')
    return r.json(), r.status_code


def get_all_workers(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('USER_MICROSERVICE') + '/users/workers')
    return r.json(), r.status_code


def get_all_managers(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('USER_MICROSERVICE') + '/users/managers')
    return r.json(), r.status_code


def delete_user(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.delete(get_path('USER_MICROSERVICE') + '/users/' + str(id))
    return r.json(), r.status_code