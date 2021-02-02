# Some prerequisites to run an application

Create all necessary database entities

```sql
  
   CREATE DATABASE intern
   
   CREATE TABLE sellers(
    id INT PRIMARY_KEY
   );
   
   CREATE TABLE goods(
    offer_id INT PRIMARY KEY,
    name     VARCHAR(70) NOT NULL,
    price    DECIMAL NOT NULL,
    quantity INT NOT NULL,
    available BOOLEAN NOT NULL,
    seller_id INT NOT NULL,
    CONSTRAINT fk_seller
      FOREIGN KEY(seller_id)
        REFERENCES sellers(id)
    );
```

Be sure to export the following environmental variables before running docker-compose. If you are on linux, run the following commands:

```bash
  export POSTGRES_USER=<your_postgres_user>
  export POSTGRES_PASSWORD=<your_postgres_password>
  export POSTGRES_DB=<database_name>
  export ADDR=<0.0.0.0:<some_port>>
  
  sudo docker-compose up
```

