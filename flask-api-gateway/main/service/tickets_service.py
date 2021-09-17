import requests
from main.util.requests import ok, unauthorized
from main.util.helper import get_path, guard_check


def get_all(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('TICKET_MICROSERVICE') + '/api/tickets/personal', headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code


def delete_ticket(id, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER', 'ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")
    r = requests.delete(get_path('TICKET_MICROSERVICE') + '/api/tickets/cancel/' + str(id), headers = headers)
    return None, r.status_code


def add_new_tickets(tickets, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER', 'ROLE_REGISTERED_USER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.post(get_path('TICKET_MICROSERVICE') + '/api/tickets', json = tickets, headers = headers)
    return None, r.status_code


def get_all_sold_tickets(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.get(get_path('TICKET_MICROSERVICE') + '/api/tickets/sold', headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code


def get_all_reserved_tickets(headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER']):
        return unauthorized("You are not allowed to complete this aciton.")

    r = requests.get(get_path('TICKET_MICROSERVICE') + '/api/tickets/reserved', headers = headers)
    if r.json() == None:
        return [], r.status_code
    return r.json(), r.status_code


def update_ticket(ticket, headers):
    auth_header = headers.get('Authorization')
    if not guard_check(auth_header, ['ROLE_WORKER']):
        return unauthorized("You are not allowed to complete this aciton.")
    
    r = requests.put(get_path('TICKET_MICROSERVICE') + '/api/tickets/' + str(ticket["id"]), json = ticket, headers = headers)
    return None, r.status_code
