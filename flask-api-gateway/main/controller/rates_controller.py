from flask import request
from flask_restplus import Resource

from ..util.dto import RateDto
from ..service.rates_service import get_rate, rate_movie

api = RateDto.api
_rate = RateDto.rate

@api.route('/<id>')
@api.param('id', 'The Movie identifier')
@api.response(404, 'Movie not found.')
class GetRate(Resource):
    @api.doc('get rate for movie')
    @api.marshal_with(_rate)
    def get(self, id):
        headers = request.headers
        return get_rate(id, headers)

@api.route('/rate-movie/<id>/<value>')
@api.param('id', 'The Movie identifier')
@api.response(404, 'Movie not found.')
class RateMovie(Resource):
    @api.doc('rate movie')
    def get(self, id, value):
        headers = request.headers
        return rate_movie(id, value, headers)