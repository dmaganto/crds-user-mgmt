import random
import string
from jinja2 import Template

keywords = ["claim", "document", "payment", "notification", "security", "user", "data", "report", "analytics", "inventory",
            "workflow", "customer", "authentication", "search", "booking", "shipping", "inventory", "notification", "integration", "monitoring"]

metadata_name = [random.choice(keywords) for _ in range(20)]
app_names = [random.choice(keywords) + "-service" for _ in range(20)]

template = Template('''
apiVersion: dmaganto.infra/v1alpha1
kind: Team
metadata:
  name: {{metadata_name[0]}}-team
  namespace: default
spec:
  applications: 
    - {{apps[0]}}
    - {{apps[1]}}
  slackChannel: {{metadata_name[0]}}-channel
''')

print(template.render(metadata_name=metadata_name,
                      apps=app_names))

