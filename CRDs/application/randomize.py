import random
import string
from jinja2 import Template


keywords = ["claim", "document", "payment", "notification", "security", "user", "data", "report", "analytics", "inventory",
            "workflow", "customer", "authentication", "search", "booking", "shipping", "inventory", "notification", "integration", "monitoring"]

metadata_name = [random.choice(keywords) + "-service" for _ in range(20)]
team = random.sample(['claims', 'front', 'back'], k=2)

template = Template('''
apiVersion: dmaganto.infra/v1alpha1
kind: Application
metadata:
  name: {{metadata_name[0]}}
  namespace: default
spec:
  team: {{team[0]}}
  slackChannel: {{metadata_name[0]}}
''')

print(template.render(metadata_name=metadata_name,
                      team=team))

