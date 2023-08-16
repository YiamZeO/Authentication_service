import socket
import json
import base64

def take_tokens(request):
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client_socket.connect(('localhost', 5000))
    client_socket.send(base64.b64encode(request.encode('UTF-8')))
    data = client_socket.recv(4096)
    tokens = json.loads((base64.b64decode(data)).decode('UTF-8'))
    return tokens

user_id = ''
tokens = {'access_token': 'None',
          'refresh_token': 'None'}
while True:
    print('Menu:')
    print('1 - Request for tokens with user id')
    print('2 - Refresh request')
    print('0 - Exit')
    match (input()):
        case ('1'):
            print('Enter user id:')
            user_id = input()
            request = 'A.'+user_id
            tokens = take_tokens(request)
            print('For user id:', user_id)
            print(' New access token:', tokens['access_token'])
            print(' New refresh token:', tokens['refresh_token'])
        case ('2'):
            request = 'R.'+json.dumps(tokens)
            tokens = take_tokens(request)
            print('For user id:', user_id)
            print(' New access token:', tokens['access_token'])
            print(' New refresh token:', tokens['refresh_token'])
        case ('0'):
            break