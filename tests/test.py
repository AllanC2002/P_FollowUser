import requests

BASE_URL_FOLLOW = "http://54.163.143.185:8080"
#BASE_URL_FOLLOW = "http://localhost:8080"  
BASE_URL_LOGIN = "http://52.203.72.116:8080"   

login_data = {
    "User_mail": "allancorrea",
    "password": "1234"
}

login_response = requests.post(f"{BASE_URL_LOGIN}/login", json=login_data)

if login_response.status_code != 200:
    print("Error en login:", login_response.status_code, login_response.json())
    exit()

token = login_response.json().get("token")
print("Token obtenido:", token)

data_follow = {
    "id_following": 1
}

# Headers con token para el microservicio Go
headers = {
    "Authorization": f"Bearer {token}"
}

# Hacer petición para seguir usuario
follow_response = requests.post(f"{BASE_URL_FOLLOW}/follow", json=data_follow, headers=headers)

print("Respuesta follow:")
print(follow_response.status_code, follow_response.json())
