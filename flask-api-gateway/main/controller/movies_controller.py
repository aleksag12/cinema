from flask import request
from requests.api import head
from flask_restplus import Resource

from ..util.dto import MovieDto
from ..service.movies_service import get_movie, save_new_movie, get_all_movies, delete_movie, update_movie

api = MovieDto.api
_movie = MovieDto.movie


@api.route('/')
class MovieList(Resource):
    @api.doc('list_of_all_movies') 
    @api.marshal_list_with(_movie)
    def get(self):
        headers = request.headers
        return get_all_movies(headers)


    @api.response(201, 'Movie successfully created.')
    @api.doc('create a new movie')
    @api.expect(_movie, validate=True)
    def post(self):
        data = request.json
        headers = request.headers
        return save_new_movie(data, headers)


@api.route('/<id>')
@api.param('id', 'The Movie identifier')
@api.response(404, 'Movie not found.')
class Movie(Resource):
    @api.doc('get movie')
    @api.marshal_with(_movie)
    def get(self, id):
        headers = request.headers
        return get_movie(id, headers)

    @api.doc('delete movie')
    def delete(self, id):
        headers = request.headers
        return delete_movie(id, headers)

    @api.doc('update movie')
    @api.expect(_movie, validate=True)
    def put(self, id):
        headers = request.headers
        data = request.data
        return update_movie(id, data, headers)
