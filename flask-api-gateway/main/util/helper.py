from flask import current_app
import jwt


def get_path(service):
    return 'http://' + current_app.config.get(service + '_HOST') + ':' + str(current_app.config.get(service + '_PORT'))


def decode_auth_token(auth_token):
    try:
        payload = jwt.decode(auth_token, current_app.config.get('SECRET_KEY'))
        return payload['sub']
    except jwt.ExpiredSignatureError:
        return 'Signature expired. Please log in again.'
    except jwt.InvalidTokenError:
        return 'Invalid token. Please log in again.'


def guard_check(auth_header, roles):
    if not auth_header or "Bearer " not in auth_header:
        return False
    else:
        try:
            if decode_auth_token(auth_header.split(" ")[1]).get('role') not in roles:
                return False
        except:
            return False
    return True