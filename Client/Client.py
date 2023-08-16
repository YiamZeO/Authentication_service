import requests

if __name__ == '__main__':
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
                response = requests.post('http://127.0.0.1:5000/user/authentication', json={'user_id': user_id})
                tokens = response.json()
                print('For user id:', user_id)
                print(' New access token:', tokens['access_token'])
                print(' New refresh token:', tokens['refresh_token'])
            case ('2'):
                response = requests.post('http://127.0.0.1:5000/user/refresh', json=tokens)
                tokens = response.json()
                print('For user id:', user_id)
                print(' New access token:', tokens['access_token'])
                print(' New refresh token:', tokens['refresh_token'])
            case ('0'):
                break