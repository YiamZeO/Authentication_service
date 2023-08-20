import pymongo
from bson.objectid import ObjectId
import jwt
import bcrypt
import datetime
import flask
import base64

KEY_WORD = 'SECRET_WORD'
REFRESH_KEY_WORD = 'REFRESH_SECRET_WORD'
app = flask.Flask(__name__)

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

@app.route("/user/authentication", methods=['GET'])
def user_authentication():
    if flask.request.method == 'GET':
        user_id = flask.request.args.get('user_id')
        request = {'_id': ObjectId(user_id)}
        tokens = take_tokens(request)
        return tokens

@app.route("/user/refresh", methods=['POST'])
def user_refresh():
    if flask.request.method == 'POST':
        request = {'refresh_token': (base64.b64decode((flask.request.get_data()))).decode('UTF-8')}
        tokens = take_tokens(request)
        return tokens

if __name__ == '__main__':
    app.run()
    