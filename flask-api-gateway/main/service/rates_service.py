import requests
from main.util.requests import bad_request, unauthorized
from main.util.helper import get_path, guard_check


def get_rate(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.get(get_path('RATING_MICROSERVICE') + '/api/rates/' + str(id), headers = headers)
    return r.json(), r.status_code


def rate_movie(id, value, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    try:
        r = requests.post(get_path('RATING_MICROSERVICE') + '/api/rates', json = { 'movie_id': int(id), 'value': int(value) }, headers = headers)
        return r.json(), r.status_code
    except:
        return bad_request("Rate value must be a number")
