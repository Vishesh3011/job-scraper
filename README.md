[//]: # ( ![job-scraper.png]&#40;data/job-scraper.png&#41;)
#
<div align="center">
  <img src="data/job-scraper.png" alt="Image description" width="400"/>
</div>

# Welcome to Job Scraper!
This tool is for everyone who wants to find jobs based on their interest, position, location (remote/on-site) using online job platforms like LinkedIn.
This tool can be used as both a telegram bot (access on https://t.me/jobScraperVMBot) and a CLI tool. Based on user input, users can either select a daily report for jobs as email delivery or just get a report in the moment as per their requirement. 

### Technical Requirements (if running on a CLI):
- MySQL installed

### Other Requirements:
- CSRF Token and Cookie from LinkedIn
- Knowledge of GeoIds from LinkedIn

### Project structure:
```
cmd/
├── app/
│ └── main.go
├── bot/
│ └── main.go
├── cli/
│ └── main.go
└── db/
internal/
├── cli/
│ └── main.go
├── core/
│ ├── application/
│ ├── config/
├── models/
├── repository/
├── service/
├── types/
└── utils/
```

### Setup Database
Create database:
```mysql
create database job_scraper_db;
```
Create a user in mysql:
```mysql
create user 'job_scraper'@'localhost' identified by 'job_scraper';
``` 
Grant all privileges to the user
```mysql
grant all privileges on job_scraper_db.* to 'job_scraper'@'localhost';
```
Flush privileges
```mysql
flush privileges;
```

### Load the schema and pre-requisite data
```
mysql -u job_scraper -p job_scraper_db < db/init.sql
mysql -u job_scraper -p job_scraper_db < db/load_prerequisites.sql
```

