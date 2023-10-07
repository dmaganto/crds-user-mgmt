import random
import string
from jinja2 import Template

metadata_name = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10))
team = random.sample(['claims', 'front', 'back'], k=2)

template = Template('''
apiVersion: dmaganto.infra/v1alpha1
kind: Application
metadata:
  name: {{metadata_name}}-service
  namespace: default
spec:
  team: {{team[0]}}
  slackChannel: {{metadata_name}}-service
''')

print(template.render(metadata_name=metadata_name,
                      team=team))

