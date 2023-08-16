import socket
import json
import pymongo
from bson.objectid import ObjectId
import base64
import jwt
import bcrypt
import datetime

KEY_WORD = 'SECRET_WORD'
REFRESH_KEY_WORD = 'REFRESH_SECRET_WORD'

def take_tokens(request):
    m_client = pymongo.MongoClient()
    m_bd = m_client['Authentication_bd']
    c_user = m_bd['Users'].find_one(request)
    tokens = {}
    if c_user == None:
        tokens = {'access_token': 'None',
                'refresh_token': 'None'}
    else:
        jwt_access_token = jwt.encode({'user_id': str(c_user['_id']),
                                       "exp": datetime.datetime.now() + datetime.timedelta(hours=1)}, KEY_WORD, algorithm="HS256")
        jwt_refresh_token = jwt.encode({'user_id': str(c_user['_id']),
                                       "exp": datetime.datetime.now() + datetime.timedelta(days=180)}, REFRESH_KEY_WORD, algorithm="HS256")
        bcrypt_refresh_token = bcrypt.hashpw(jwt_refresh_token.encode(), bcrypt.gensalt())
        m_bd['Users'].update_one(request, {"$set": {'refresh_token': str(bcrypt_refresh_token)}})
        tokens = {'access_token': jwt_access_token,
                  'refresh_token': str(bcrypt_refresh_token)}
    m_client.close()
    return tokens

serversocket = socket.socket(
        socket.AF_INET, socket.SOCK_STREAM)
serversocket.bind(('localhost', 5000))
serversocket.listen()
print('Server running...')
while True:
    connection, address = serversocket.accept()
    print('Connection from ', address)
    data = connection.recv(4096)
    data = (base64.b64decode(data)).decode('UTF-8')
    request_type = data[0]
    request = {}
    match (request_type):
        case ('A'):
            user_id = data[2:]
            request = {'_id': ObjectId(user_id)}
        case ('R'):
            tokens = json.loads(data[2:])
            request = {'refresh_token': tokens['refresh_token']}
    tokens = take_tokens(request)
    data = base64.b64encode((json.dumps(tokens)).encode('UTF-8'))
    connection.send(data)
    print('Connection from', address, 'received new tokens')   
    connection.close()