from flask import request
from flask_restplus import Resource

from ..util.dto import TicketDto
from ..service.tickets_service import delete_ticket, get_all, add_new_tickets, get_all_reserved_tickets, get_all_sold_tickets, update_ticket

api = TicketDto.api
_ticket = TicketDto.ticket
_ticket_list = TicketDto.ticket_list


@api.route('/')
class TicketList(Resource):
    @api.doc('list_of_all_tickets')
    @api.marshal_list_with(_ticket)
    def get(self):
        headers = request.headers
        return get_all(headers)

    @api.response(200, 'Tickets successfully added.')
    @api.response(400, 'Seats are already taken.')
    @api.doc('add new tickets')
    @api.expect(_ticket_list, validate=True)
    def post(self):
        data = request.json
        headers = request.headers
        return add_new_tickets(data, headers)

    @api.response(200, 'Ticket successfully updated.')
    @api.doc('update ticket')
    @api.expect(_ticket)
    def put(self):
        data = request.json
        headers = request.headers
        return update_ticket(data, headers)


@api.route('/<id>')
@api.param('id', 'The Ticket identifier')
@api.response(404, 'Ticket not found.')
@api.response(400, 'Ticket can\'t be deleted.')
class Ticket(Resource):
    @api.doc('delete ticket')
    def delete(self, id):
        headers = request.headers
        return delete_ticket(id, headers)


@api.route('/sold/')
class TicketList(Resource):
    @api.doc('list_of_all_sold_tickets')
    @api.marshal_list_with(_ticket)
    def get(self):
        headers = request.headers
        return get_all_sold_tickets(headers)


@api.route('/reserved/')
class TicketList(Resource):
    @api.doc('list_of_all_reserved_tickets')
    @api.marshal_list_with(_ticket)
    def get(self):
        headers = request.headers
        return get_all_reserved_tickets(headers)
