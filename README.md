# tasker

<p align='center'.>
  <img src="./docs/assets/robotasks.png" width=25% height=25%>
</p>

Create tasks and save them in the database.  
Will replace the todo list in vim/notepadqq that is used for daily work.

## Setup

### Create the database

The following script will create the database, user and table for the development.

```bash
./scripts/setup.sh create 'myPASShere'
```

## Developing

Once the database is running, export the environment variables for the app.  
For example:

```bash
 export DB_PASS='taskerPWD22'
 export DB_USER='tasker'
 export DB_HOST='localhost'
 export DB_DATABASE='tasks'
 export DB_PORT='5432
```