import random
import string
from jinja2 import Template

names = ["James", "Emma", "William", "Olivia", "Liam", "Ava", "Henry", "Isabella", "Alexander", "Sophia", "Benjamin", "Mia", "Michael", "Charlotte", "Daniel", "Amelia", "Matthew", "Evelyn", "David", "Abigail",
         "Joseph", "Harper", "Andrew", "Grace", "John", "Elizabeth", "Samuel", "Lily", "Christopher", "Chloe", "Anthony", "Ella", "Robert", "Scarlett", "William", "Zoe", "Nicholas", "Natalie", "Jonathan", "Avery"]
surnames = ["Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia", "Rodriguez", "Martinez",
            "Hernandez", "Lopez", "Gonzalez", "Perez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson"]

random_name = random.choice(names)
random_surname = random.choice(surnames)

metadata_name = f"{random_name.lower()}.{random_surname.lower()}"
spec_fullname = f"{random_name} {random_surname}"
spec_email = f"{random_name.lower()}.{random_surname.lower()}" + '@example.com'
role_type = random.choice(['devops', 'developer', 'productowner', 'agilecoach', 'federateddevops'])
teams = random.sample(['claims', 'front', 'back'], k=2)

template = Template('''
apiVersion: dmaganto.infra/v1alpha1
kind: Developer
metadata:
  name: {{ metadata_name }}
  namespace: default
spec:
  fullName: {{ spec_fullname }}
  email: {{ spec_email }}
  roleType: {{ role_type }}
  teams:
    - {{ teams[0] }}
    - {{ teams[1] }}
''')

print(template.render(metadata_name=metadata_name,
                      spec_fullname=spec_fullname,
                      spec_email=spec_email,
                      role_type=role_type,
                      teams=teams))
