import requests
from main.util.helper import get_path, guard_check
from main.util.requests import ok, unauthorized


def get_comments(id, headers):
    r = requests.get(get_path('RATING_MICROSERVICE') + '/api/comments/' + str(id), headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code

def add_comment(comment, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.post(get_path('RATING_MICROSERVICE') + '/api/comments', json = comment, headers = headers)
    return r.json(), r.status_code

def delete_comment(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.delete(get_path('RATING_MICROSERVICE') + '/api/comments/' + str(id), headers = headers)
    return None, r.status_code