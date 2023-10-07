import random
import string
from jinja2 import Template

metadata_name = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10))
spec_fullname = ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))
spec_email = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10)) + '@example.com'
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
