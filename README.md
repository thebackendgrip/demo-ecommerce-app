Services
- Order service - for buying items in the system
- Catalog service - to manage inventory in the system
- User service - to manage users profiles in the system

Dependencies

Filebeat - is a lightweight log shipper, it monitors logs in the filepaths specified in its configuration and forwards them to Elasticsearch
Elasticsearch - is an open source distributed, search and analytical engine, it is used for full-text search, structured search, and analytics. Elasticsearch receives the logs forwarded from Filebeat
Kibana - is an open source visualisation and management tool designed to work with Elasticsearch, it enables users to visualise and explore managed data stored in Elasticsearch
Jaeger - is an open-source, end-to-end distributed tracing system designed to monitor and troubleshoot micro-services based applications, it captures how requests flow through a system and
Prometheus - is an open source systems monitoring and alerting toolkit designed to collect, store, and query metrics from various systems and applications
