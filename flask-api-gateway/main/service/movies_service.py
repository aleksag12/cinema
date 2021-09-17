from main.util.helper import get_path, guard_check
from main.util.requests import ok, unauthorized
import requests


def save_new_movie(data, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.post(get_path('MOVIE_MICROSERVICE') + '/api/movies', json = data, headers = headers)
    return None, r.status_code


def update_movie(id, data, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.put(get_path('MOVIE_MICROSERVICE') + '/api/movies/' + str(id), data = data, headers = headers)
    return None, r.status_code


def get_all_movies(headers):
    r = requests.get(get_path('MOVIE_MICROSERVICE') + '/api/movies', headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code


def get_movie(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('MOVIE_MICROSERVICE') + '/api/movies/' + id, headers = headers)
    return r.json(), r.status_code


def delete_movie(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_MANAGER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.delete(get_path('MOVIE_MICROSERVICE') + '/api/movies/' + str(id), headers = headers)
    return None, r.status_code
