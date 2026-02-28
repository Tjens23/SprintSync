from diagrams import Diagram, Cluster, Edge
from diagrams.onprem.client import User
from diagrams.onprem.queue import RabbitMQ
from diagrams.onprem.compute import Server
from diagrams.saas.identity import Auth0

with Diagram("authDiagram", show=False, direction="LR"):
    user = User("User")
    
    with Cluster("Auth Service (Port 3001)"):
        auth_server = Server("Fiber Server")
        auth_routes = Server("/login\n/callback\n/logout\n/user")
    
    auth0 = Auth0("Auth0\n(OAuth2/OIDC)")
    
    rabbitmq = RabbitMQ("RabbitMQ\n(Topic Exchange)")
    
    with Cluster("User Service (Port 3002)"):
        user_server = Server("Fiber Server")
        user_consumer = Server("Message\nConsumer")
    
    # User flow
    user >> Edge(label="1. Login Request") >> auth_server
    auth_server >> auth_routes
    auth_routes >> Edge(label="2. OAuth2 Redirect") >> auth0
    auth0 >> Edge(label="3. Callback + Token") >> auth_routes
    
    # Event publishing
    auth_routes >> Edge(label="4. Publish Events\n(user.authenticated\nuser.logged.out)", color="orange") >> rabbitmq
    
    # Event consuming
    rabbitmq >> Edge(label="5. Consume Events", color="green") >> user_consumer
    user_consumer >> user_server