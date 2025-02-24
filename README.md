# My-Microservice-Project
My main project I am working on using Go and other technologies, updating...
For now 4 services are going to be done:
1. Room Service (Mostly finished)
2. User Service (Mostly finished)
3. Reservation Service (Working on)
4. Notification Service (0%)

- PostgreSQL, Redis, S3 for the data storage.
- Very simple Diagram will look like this:
                                +----------------------+
                                |   User Service       |
                                |  (Auth & Roles)      |
                                |  (JWT Generation)    |
                                +----------------------+
                                          |
                        ------------------|------------------
                       |                                    |
              +----------------------+          +----------------------+
              |   Room Service      |           | Reservation Service  |
              |  (CRUD Rooms)       |-----------| (CRUD Reservations)  |
              |  (Owner Access)     |           |   (Guest Access)     |
              +----------------------+          +----------------------+
                                |                     |
                                |                     |
                              +-------------------------+
                              |   Notification Service  |
                              |  (Sends Notifications)  |
                              +-------------------------+
 +-------+
 | Kafka |    (Kafka will be added as well)
 +-------+

