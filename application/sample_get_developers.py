from kubernetes import client, config


def read_crd_in_namespace(namespace, crd_plural):
    # Automatically loads the configuration from the Kubernetes environment
    config.load_incluster_config()

    # Create a CustomObjectsApi object to interact with CRDs
    custom_api = client.CustomObjectsApi()

    try:
        # Use the API to list CRD resources in the specified namespace
        resources = custom_api.list_namespaced_custom_object(
            group="dmaganto.infra",
            version="v1alpha1",  # Replace with the appropriate version of your CRD
            namespace=namespace,
            plural=crd_plural,
        )

        # Print the found resources
        for resource in resources["items"]:
            print(resource)

    except client.ApiException as e:
        print(f"Error listing CRD resources: {e}")


if __name__ == "__main__":
    # Specify the namespace name and the CRD plural you want to read
    namespace = "default"
    crd_plural = "developers"  # Replace with the plural of your CRD

    read_crd_in_namespace(namespace, crd_plural)
