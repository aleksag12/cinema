from flask import request
from flask_restplus import Resource

from ..util.dto import CommentDto
from ..service.comments_service import add_comment, delete_comment, get_comments

api = CommentDto.api
_comment = CommentDto.comment

@api.route('/<id>')
@api.param('id', 'The Movie identifier')
@api.response(404, 'Movie not found.')
class GetComments(Resource):
    @api.doc('get comments for movie')
    @api.marshal_list_with(_comment)
    def get(self, id):
        headers = request.headers
        return get_comments(id, headers)

    @api.doc('get comments for movie')
    @api.marshal_with(_comment)
    def delete(self, id):
        headers = request.headers
        return delete_comment(id, headers)

@api.route('/')
class Comment(Resource):
    @api.response(200, 'Comment successfully added.')
    @api.doc('add a new comment')
    @api.expect(_comment, validate=True)
    @api.marshal_with(_comment)
    def post(self):
        data = request.json
        headers = request.headers
        return add_comment(data, headers)
