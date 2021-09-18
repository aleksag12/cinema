def bad_request(message):
    response_object = {
        'status': 'fail',
        'message': message
    }
    return response_object, 400


def created(message):
    response_object = {
        'status': 'success',
        'message': message
    }
    return response_object, 201


def conflict(message):
    response_object = {
        'status': 'fail',
        'message': message
    }
    return response_object, 409


def ok(message):
    response_object = {
        'status': 'success',
        'message': message
    }
    return response_object, 200


def not_found(message):
    response_object = {
        'status': 'fail',
        'message': message
    }
    return response_object, 404


def unauthorized(message):
    response_object = {
        'status': 'fail',
        'message': message
    }
    return response_object, 401
