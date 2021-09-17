import requests
from main.util.helper import get_path, guard_check
from main.util.requests import unauthorized


def save_new_projection(projection, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.post(get_path('PROJECTION_MICROSERVICE') + '/api/projections', json = projection, headers = headers)
    return None, r.status_code


def get_all(headers):
    r = requests.get(get_path('PROJECTION_MICROSERVICE') + '/api/projections', headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code


def get_reserved_seats(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER', 'ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.get(get_path('PROJECTION_MICROSERVICE') + '/api/seats/' + str(id), headers = headers)
    return r.json(), r.status_code


def delete_projection(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.delete(get_path('PROJECTION_MICROSERVICE') + '/api/projections/' + str(id), headers = headers)
    return None, r.status_code
