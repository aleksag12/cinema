from flask import request
from flask_restplus import Resource

from ..util.dto import ProjectionDto
from ..service.projections_service import get_all, get_reserved_seats, save_new_projection, delete_projection

api = ProjectionDto.api
_projection = ProjectionDto.projection
_seat_list = ProjectionDto.seat_list


@api.route('/')
class ProjectionList(Resource):
    @api.doc('list_of_all_projections')
    @api.marshal_list_with(_projection)
    def get(self):
        headers = request.headers
        return get_all(headers)

    @api.response(201, 'Projection successfully created.')
    @api.doc('create a new projection')
    @api.expect(_projection, validate=True)
    def post(self):
        data = request.json
        headers = request.headers
        return save_new_projection(data, headers)


@api.route('/reserved-seats/<id>')
@api.param('id', 'The Projection identifier')
@api.response(404, 'Projection not found.')
class ReservedSeats(Resource):
    @api.doc('get reserved seats')
    @api.marshal_with(_seat_list)
    def get(self, id):
        headers = request.headers
        return get_reserved_seats(id, headers)


@api.route('/<id>')
@api.param('id', 'The Projection identifier')
@api.response(404, 'Projection not found.')
class Projection(Resource):
    @api.doc('delete projection')
    def delete(self, id):
        headers = request.headers
        return delete_projection(id, headers)
