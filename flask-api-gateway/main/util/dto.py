from flask_restplus import Namespace, fields


class MovieDto:
    api = Namespace('movies', description='movie related operations')
    movie = api.model('movie', {
        'id': fields.Integer(description='movie id'),
        'name': fields.String(required=True, description='movie name'),
        'description': fields.String(required=True, description='movie description'),
        'genre': fields.String(required=True, description='movie genre'),
        'length': fields.Integer(required=True, description='movie length'),
        'year': fields.Integer(required=True, description='movie year'),
        'average_rate': fields.Float(description='average rate'),
    })


class RateDto:
    api = Namespace('rates', description='rates related operations')
    rate = api.model('rate', {
        'id': fields.Integer(description='rate id'),
        'user_id': fields.Integer(description='user id'),
        'movie_id': fields.Integer(required=True, description='movie id'),
        'value': fields.Integer(required=True, description='rate value'),
    })


class CommentDto:
    api = Namespace('comments', description='comments related operations')
    comment = api.model('comment', {
        'id': fields.Integer(description='comment id'),
        'user_id': fields.Integer(description='user id'),
        'username': fields.String(description='user username'),
        'movie_id': fields.Integer(required=True, description='movie id'),
        'text': fields.String(required=True, description='comment text'),
    })


class ProjectionDto:
    api = Namespace('projections', description='projections related operations')
    projection = api.model('projection', {
        'id': fields.Integer(description='projection id'),
        'movie_id': fields.Integer(required=True, description='movie id'),
        'movie_name': fields.String(description='movie name'),
        'date_time': fields.Integer(required=True, description='date and time'),
        'price': fields.Float(required=True, description='projection price'),
    })
    seat_list = api.model('seat_list', {
        'seats': fields.List(fields.String, description='seat list'),
    })


class TicketDto:
    api = Namespace('tickets', description='tickets related operations')
    ticket = api.model('ticket', {
        'id': fields.Integer(description='ticket id'),
        'projection_id': fields.Integer(required=True, description='projection id'),
        'user_id': fields.Integer(description='user id'),
        'movie_name': fields.String(description='movie name'),
        'customer': fields.String(description='customer name'),
        'row': fields.Integer(required=True, description='seat row'),
        'column': fields.Integer(required=True, description='seat column'),
        'date_time': fields.Integer(description='date and time'),
        'price': fields.Float(description='ticket price'),
        'sold': fields.Boolean(description='is ticket sold'),
    })
    ticket_list = api.model('ticket_list', {
        'tickets': fields.List(fields.Nested(ticket), description='ticket list'),
    })


class AuthDto:
    api = Namespace('auth', description='authentification related operations')
    auth = api.model('auth', {
        'username': fields.String(required=True, description='username'),
        'password': fields.String(required=True, description='password'),
    })
    new_user = api.model('new_user', {
        'username': fields.String(required=True, description='user username'),
        'email': fields.String(required=True, description='user email address'),
        'first_name': fields.String(required=True, description='user first name'),
        'last_name': fields.String(required=True, description='user last name'),
        'password': fields.String(required=True, description='user password'),
        'role': fields.String(required=True, description='user role'),
    })
    user = api.model('user', {
        'id': fields.Integer(description='user id'),
        'username': fields.String(required=True, description='user username'),
        'email': fields.String(required=True, description='user email address'),
        'first_name': fields.String(required=True, description='user first name'),
        'last_name': fields.String(required=True, description='user last name'),
        'role': fields.String(required=True, description='user role'),
    })
